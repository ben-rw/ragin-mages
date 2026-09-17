package lobby

import (
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

const backgroundPath = "assets/images/center_background.png"
const songPath = "assets/audio/music/amalfi-coast-loop.ogg"

type Lobby struct {
	shared.Roster
	Conn          *ws.Connection
	Sprites       []*shared.Sprite
	Background    *ebiten.Image
	audioPlayer   *audio.Player
	sceneChanging bool
	lobbyImage    *ebiten.Image
}

func NewLobby(c *ws.Connection) *Lobby {
	bg, _, err := ebitenutil.NewImageFromFileSystem(shared.AssetsFS, backgroundPath)
	if err != nil {
		log.Printf("couldn't load background: %v")
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
		Conn:          c,
		Sprites:       []*shared.Sprite{},
		Background:    bg,
		audioPlayer:   audioPlayer,
		sceneChanging: false,
		lobbyImage:    lobbyImage,
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

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) &&
		l.Player.Data.Host == true {
		l.sceneChanging = true
	}

	var fadeFinished = false
	if l.sceneChanging == true {
		fadeFinished = sound.FadeOut(l.audioPlayer)
	} else if !l.audioPlayer.IsPlaying() {
		l.audioPlayer.SetVolume(0.2)
		l.audioPlayer.SetBufferSize(8192)
		l.audioPlayer.Play()
	}

	if fadeFinished {
		l.Conn.WriteMsg(protocol.SceneChange, protocol.SceneChangeData{
			SceneType: protocol.RandomScene,
		})
	}

	return nil
}

func (l *Lobby) Draw(screen *ebiten.Image) {
	l.lobbyImage.Clear()
	l.lobbyImage.Fill(shared.BackgroundColor)

	opts := ebiten.DrawImageOptions{}

	// s := l.Background.Bounds().Size()
	// scaleX := shared.ScreenWidth / float64(s.X)
	// scaleY := shared.ScreenHeight / float64(s.Y)
	// opts.GeoM.Scale(scaleX, scaleY)
	// screen.DrawImage(l.Background, &opts)
	//
	// opts.GeoM.Reset()

	for _, player := range l.Players {
		opts.GeoM.Translate(player.X, player.Y)

		l.lobbyImage.DrawImage(
			player.Img.SubImage(
				player.SpriteSheet.Rect(player.ActiveAnimation.Frame()),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	for _, player := range l.Players {
		textOpts := text.DrawOptions{
			LayoutOptions: player.NameTag.LayoutOptions,
		}
		textOpts.GeoM.Translate(player.NameTag.X, player.NameTag.Y)
		text.Draw(l.lobbyImage, player.Data.Name, player.NameTag.Face, &textOpts)

		textOpts.GeoM.Reset()
	}

	controlsText := "Controls: 'Left Click' to Attack, 'Right Click' to Reflect, 'Q' to Repel"
	textOpts := text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: 2,
		},
	}
	tX, tY := shared.TopRight()
	textOpts.GeoM.Translate(tX/2, tY/2)
	text.Draw(l.lobbyImage, controlsText, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)

	textOpts.GeoM.Reset()

	waitText := "Waiting for host..."
	if l.Player.Data.Host == true {
		waitText = "Press ENTER to start!"
	}
	textOpts = text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: 2,
		},
	}
	tX, tY = shared.BottomRight()
	textOpts.GeoM.Translate(tX/2, tY/2)
	text.Draw(l.lobbyImage, waitText, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)

	textOpts.GeoM.Reset()

	opts.GeoM.Scale(2, 2)
	screen.DrawImage(l.lobbyImage, &opts)

	opts.GeoM.Reset()
}
