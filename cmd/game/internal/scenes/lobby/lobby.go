package lobby

import (
	"fmt"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/sound"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/ws"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"log"
)

const (
	backgroundPath = "assets/images/ninja_adventure/Backgrounds/ninja_actor_background.png"
	mouseLeftPath  = "assets/images/ninja_adventure/Ui/Input/Mouse/MouseButtonLeft.png"
	mouseRightPath = "assets/images/ninja_adventure/Ui/Input/Mouse/MouseButtonRight.png"
	songPath       = "assets/audio/music/amalfi-coast-loop.ogg"
)

var (
	waitText     = "Waiting for host..."
	hostText     = "Press ENTER to start!"
	titleText    = "RAGIN'\n  MAGES"
	controlsText = fmt.Sprintf("Attack\nReflect")
)

type Background struct {
	img   *ebiten.Image
	alpha float32
}

type Lobby struct {
	shared.Roster
	Conn               *ws.Connection
	Sprites            []*shared.Sprite
	Background         *Background
	lobbyImage         *ebiten.Image
	mouseLeftImage     *ebiten.Image
	mouseRightImage    *ebiten.Image
	audioPlayer        *audio.Player
	foregroundAlpha    float32
	sceneChanging      bool
	musicFadeFinished  bool
	screenFadeFinished bool
}

func NewLobby(c *ws.Connection) *Lobby {
	bgImg, _, err := ebitenutil.NewImageFromFileSystem(shared.AssetsFS, backgroundPath)
	if err != nil {
		log.Printf("couldn't load background: %v")
	}
	bg := &Background{
		img:   bgImg,
		alpha: 1.0,
	}

	mlImg, _, err := ebitenutil.NewImageFromFileSystem(shared.AssetsFS, mouseLeftPath)
	if err != nil {
		log.Printf("couldn't load image: %v")
	}

	mrImg, _, err := ebitenutil.NewImageFromFileSystem(shared.AssetsFS, mouseRightPath)
	if err != nil {
		log.Printf("couldn't load image: %v")
	}

	audioPlayer, err := sound.NewAudioPlayer(songPath, true, 0)
	if err != nil {
		log.Printf("couldn't create audio player: %v", err)
	}

	lobbyImage := ebiten.NewImage(int(shared.ScreenWidth/2), int(shared.ScreenHeight/2))

	return &Lobby{
		Roster: shared.Roster{
			Players: map[string]*shared.Player{},
			Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		},
		Conn:               c,
		Sprites:            []*shared.Sprite{},
		Background:         bg,
		lobbyImage:         lobbyImage,
		mouseLeftImage:     mlImg,
		mouseRightImage:    mrImg,
		audioPlayer:        audioPlayer,
		foregroundAlpha:    0,
		sceneChanging:      false,
		musicFadeFinished:  false,
		screenFadeFinished: false,
	}
}

func (l *Lobby) Update(messages []protocol.Message) error {
	for _, message := range messages {
		switch message.Type {
		case protocol.JoinResponse:
			err := l.HandleJoinResponse(message)
			if err != nil {
				log.Println(err)
				continue
			}

			playerUpdateData := protocol.PlayerUpdateData{
				PlayerData: l.Player.Data,
			}

			l.Conn.WriteMsg(protocol.PlayerUpdate, playerUpdateData)

		case protocol.PlayerUpdate:
			err := l.HandlePlayerUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}

		default:
		}
	}

	for _, player := range l.Players {
		player.ActiveAnimation = player.GetActiveAnimation()
		player.ActiveAnimation.Update()
	}

	// change scene when host presses enter
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) &&
		l.Player.Data.Host == true {
		l.sceneChanging = true
	}

	// fade music before scene change
	if l.sceneChanging == true && l.musicFadeFinished == false {
		// check that both music and images have finished fading
		l.musicFadeFinished = sound.FadeOut(l.audioPlayer)
	} else if !l.audioPlayer.IsPlaying() && !l.sceneChanging {
		l.audioPlayer.SetVolume(0.2)
		l.audioPlayer.SetBufferSize(8192)
		l.audioPlayer.Play()
	}

	// fade screen out before scene change
	if l.sceneChanging {
		l.screenFadeFinished = func() bool {
			if l.foregroundAlpha <= 0.5 &&
				l.Background.alpha > 0 {
				l.Background.alpha -= 0.02
				l.foregroundAlpha -= 0.02
				return false
			} else if l.foregroundAlpha > 0 {
				l.foregroundAlpha -= 0.02
				return false
			}
			return true
		}()
	} else {
		// fade background after join
		if l.Background.alpha > 0.3 {
			l.Background.alpha -= 0.005
		}

		// fade everything else in after the background fades out
		if l.Background.alpha <= 0.3 && l.foregroundAlpha < 1.0 {
			l.foregroundAlpha += 0.01
		}
	}

	// wait for music to fade to avoid song overlap
	if l.musicFadeFinished && l.screenFadeFinished {
		l.Conn.WriteMsg(protocol.SceneChange, protocol.SceneChangeData{
			SceneType: protocol.RandomScene,
		})
	}

	return nil
}

func (l *Lobby) Draw(screen *ebiten.Image) {
	l.lobbyImage.Clear()
	// l.lobbyImage.Fill(shared.BackgroundColor)

	opts := ebiten.DrawImageOptions{}
	textOpts := text.DrawOptions{}

	opts.ColorScale.ScaleAlpha(l.Background.alpha)
	l.lobbyImage.DrawImage(l.Background.img, &opts)
	opts.ColorScale.Reset()

	if l.foregroundAlpha < 1.0 {
		opts.ColorScale.ScaleAlpha(l.foregroundAlpha)
		textOpts.ColorScale.ScaleAlpha(l.foregroundAlpha)
	}

	for _, player := range l.Players {
		opts.GeoM.Translate(player.X, player.Y)

		player.ActiveAnimation = player.GetActiveAnimation()
		l.lobbyImage.DrawImage(
			player.Img.SubImage(
				player.SpriteSheet.Rect(player.ActiveAnimation.Frame()),
			).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}

	x, y := shared.TopCenterRight()
	opts.GeoM.Translate(x/2-20, y/2-9)
	l.lobbyImage.DrawImage(l.mouseLeftImage, &opts)
	opts.GeoM.Reset()

	opts.GeoM.Translate(x/2-20, y/2+15)
	l.lobbyImage.DrawImage(l.mouseRightImage, &opts)
	opts.GeoM.Reset()

	for _, player := range l.Players {
		textOpts.LayoutOptions = player.NameTag.LayoutOptions
		textOpts.GeoM.Translate(player.NameTag.X, player.NameTag.Y)
		text.Draw(l.lobbyImage, player.Data.Name, player.NameTag.Face, &textOpts)
		textOpts.GeoM.Reset()
	}

	textOpts.PrimaryAlign = text.AlignStart
	textOpts.LineSpacing = 24
	textOpts.GeoM.Translate(x/2, y/2-4)
	text.Draw(l.lobbyImage, controlsText, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)
	textOpts.GeoM.Reset()

	if l.Player.Data.Host == true {
		waitText = hostText
	}
	textOpts.PrimaryAlign = text.AlignEnd
	x, y = shared.BottomRight()
	textOpts.GeoM.Translate(x/2, y/2)
	text.Draw(l.lobbyImage, waitText, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)
	textOpts.GeoM.Reset()

	x, y = shared.Center()
	textOpts.PrimaryAlign = text.AlignCenter
	textOpts.LineSpacing = 24
	textOpts.GeoM.Translate(x/2, y/2-20)
	text.Draw(l.lobbyImage, titleText, &text.GoTextFace{Source: shared.FontSrc, Size: 18}, &textOpts)
	textOpts.GeoM.Reset()

	opts.ColorScale.Reset()
	textOpts.ColorScale.Reset()

	opts.GeoM.Scale(2, 2)
	screen.DrawImage(l.lobbyImage, &opts)

	opts.GeoM.Reset()
}
