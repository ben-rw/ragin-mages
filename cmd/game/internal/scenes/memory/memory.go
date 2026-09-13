package memory

import (
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/ws"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	//	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type Memory struct {
	shared.Roster
	Conn    *ws.Connection
	Sprites []*shared.Sprite
}

func NewMemory(c *ws.Connection) *Memory {
	log.Println("scene changed to Memory")
	return &Memory{
		Roster: shared.Roster{
			Players: map[string]*shared.Player{},
			Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		},
		Conn:    c,
		Sprites: []*shared.Sprite{},
	}
}

func (m *Memory) Update(messages []protocol.Message) error {
	for _, message := range messages {
		switch message.Type {
		case protocol.JoinResponse:
			err := m.HandleJoinResponse(message)
			if err != nil {
				log.Println(err)
				continue
			}

			playerUpdateData := protocol.PlayerUpdateData{
				PlayerData: m.Player.Data,
			}

			getInLine(m.Player, len(m.Players), 0)
			m.Conn.WriteMsg(protocol.PlayerUpdate, playerUpdateData)

		case protocol.PlayerUpdate:
			err := m.HandlePlayerUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}

		default:
		}
	}

	for _, player := range m.Players {
		player.ActiveAnimation = player.GetActiveAnimation()
		player.ActiveAnimation.Update()
	}

	return nil
}

func (m *Memory) Draw(screen *ebiten.Image) {
	screen.Fill(shared.BackgroundColor)

	opts := ebiten.DrawImageOptions{}
	// opts.GeoM.Translate(m.Player.X, m.Player.Y)
	//
	// screen.DrawImage(
	// 	m.Player.Img.SubImage(
	// 		image.Rect(0, 0, 16, 16),
	// 	).(*ebiten.Image),
	// 	&opts,
	// )
	//
	// opts.GeoM.Reset()

	for _, player := range m.Players {
		opts.GeoM.Translate(player.X, player.Y)

		player.ActiveAnimation = player.GetActiveAnimation()
		screen.DrawImage(
			player.Img.SubImage(
				player.SpriteSheet.Rect(player.ActiveAnimation.Frame()),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}
	for _, player := range m.Players {
		textOpts := text.DrawOptions{
			LayoutOptions: player.NameTag.LayoutOptions,
		}
		textOpts.GeoM.Translate(player.NameTag.X, player.NameTag.Y)
		text.Draw(screen, player.Data.Name, player.NameTag.Face, &textOpts)

		textOpts.GeoM.Reset()
	}

	waitText := "Hit the drums in the right order!"
	// if m.Player.Data.Host == true {
	// 	waitText = "Press ENTER to start!"
	// }
	textOpts := text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: 2,
		},
	}
	textOpts.GeoM.Translate(shared.BottomRight())
	text.Draw(screen, waitText, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)
}
