package main

import (
	"log"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/scenes/lobby"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/scenes/memory"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/scenes/wizarena"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/ws"
	"github.com/ben-rw/ragin-mages/internal/protocol"
)

func StartNewScene(sceneType protocol.SceneType, c *ws.Connection) Scene {
	switch sceneType {
	case protocol.LobbyScene:
		return lobby.NewLobby(c)
	case protocol.MemoryScene:
		return memory.NewMemory(c)
	case protocol.WizardsScene:
		return wizarena.NewWizArena(c)
	default:
		log.Println("invalid scene name")
		return nil
	}
}
