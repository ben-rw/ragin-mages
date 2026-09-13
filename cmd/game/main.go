package main

import (
	_ "image/png"
	"log"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/scenes/lobby"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/ws"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update(messages []protocol.Message) error
	Draw(screen *ebiten.Image)
}

type Game struct {
	Conn  *ws.Connection
	Scene Scene
}

func (g *Game) Update() error {
	msgs := g.Conn.Check()
	for _, msg := range msgs {
		if msg.Type == protocol.SceneChange {
			sceneChangeData, err := msg.UnmarshalMessageData()
			if err != nil {
				log.Printf("couldn't start new scene: %v", err)
				continue
			}
			data, ok := sceneChangeData.(*protocol.SceneChangeData)
			if !ok {
				log.Println("failed type assertion")
				continue
			}
			g.Scene = StartNewScene(data.SceneType, g.Conn)
		}
	}

	return g.Scene.Update(msgs)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(shared.ScreenWidth), int(shared.ScreenHeight)
}

func main() {
	conn, err := ws.ConnectToWebsocket()
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		Conn:  conn,
		Scene: lobby.NewLobby(conn),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
