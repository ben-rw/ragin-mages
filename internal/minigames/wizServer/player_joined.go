package wizServer

import (
	"context"
	"log"
	"time"

	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
	"github.com/coder/websocket/wsjson"
)

func (w *WizardServer) PlayerJoined(player *room.Player, playerDataList []*protocol.PlayerData) {
	w.Mu.Lock()

	log.Printf("joined: %v\n", player)

	if _, ok := w.wizards[player.Username]; !ok {
		// create new wizard
		w.wizards[player.Username] = NewWizard(player.Username)
		w.Mu.Unlock()

		// send other clients info about the new wizard
		playerUpdateMsg := &protocol.PlayerUpdateData{
			PlayerData: player.PlayerToPlayerData(),
		}

		pumsg, err := protocol.MarshalToMessage(protocol.PlayerUpdate, playerUpdateMsg)
		if err != nil {
			log.Printf("couldn't connect player %v mid-game, failed to marshal player update message: %v\n", player.Username, err)
		}
		w.BroadcastExcludeSender(pumsg, player.Username)
	} else {
		w.Mu.Unlock()
	}

	// send the connecting wizard's client the list of the room's players and identify the connecting wizard to theirself
	// joinResponseMessage := &protocol.LobbyStateSnapshotData{
	// 	PlayerData: player.PlayerToPlayerData(),
	// 	PlayerList: playerDataList,
	// }
	//
	// player.Mu.RLock()
	// client := player.Client
	// player.Mu.RUnlock()
	//
	// if client != nil {
	// 	err := wsjson.Write(context.Background(), client.Conn, joinResponseMessage)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }

}

func (w *WizardServer) SendStateSnapshot(c *room.Client) {
	w.Mu.RLock()

	playerData := &protocol.PlayerData{
		Score:       c.Player.Score,
		SpriteIndex: c.Player.SpriteIndex,
		Host:        c.Player.Host,
		X:           c.Player.X,
		Y:           c.Player.Y,
		Username:    c.Player.Username,
	}

	wizardsData := make([]*protocol.FullWizardData, 0, len(w.wizards))
	for _, player := range w.Players {
		wizardsData = append(wizardsData, &protocol.FullWizardData{
			X:               w.wizards[player.Username].X,
			Y:               w.wizards[player.Username].Y,
			Dx:              w.wizards[player.Username].Dx,
			Dy:              w.wizards[player.Username].Dy,
			KnockbackDx:     w.wizards[player.Username].KnockbackDx,
			KnockbackDy:     w.wizards[player.Username].KnockbackDy,
			Noclip:          w.wizards[player.Username].Noclip,
			AttackCooldown:  w.wizards[player.Username].AttackCooldown,
			ProjectileSpeed: w.wizards[player.Username].ProjectileSpeed,
			ProjectileScale: w.wizards[player.Username].ProjectileScale,
			Knockback:       w.wizards[player.Username].Knockback,
			Health:          w.wizards[player.Username].Health,
			Username:        player.Username,
			Score:           player.Score,
			SpriteIndex:     player.SpriteIndex,
			Host:            player.Host,
		})
	}

	enemiesData := make([]*protocol.FullEnemyData, 0, len(w.enemies))
	for _, enemy := range w.enemies {
		enemiesData = append(enemiesData, &protocol.FullEnemyData{
			X:             enemy.X,
			Y:             enemy.Y,
			Dx:            enemy.Dx,
			Dy:            enemy.Dy,
			KnockbackDx:   enemy.KnockbackDx,
			KnockbackDy:   enemy.KnockbackDy,
			Noclip:        enemy.Noclip,
			Health:        enemy.Health,
			MoveSpeed:     enemy.MoveSpeed,
			ID:            enemy.ID,
			FollowsPlayer: enemy.FollowsPlayer,
		})
	}

	chestsData := make([]*protocol.ChestUpdateData, 0, len(w.openedChests))
	for key := range w.openedChests {
		chestsData = append(chestsData, &protocol.ChestUpdateData{
			X: key.X,
			Y: key.Y,
		})
	}

	msg, err := protocol.MarshalToMessage(protocol.WizStateSnapshot, &protocol.WizStateSnapshotData{
		PlayerData:  playerData,
		WizardsData: wizardsData,
		EnemiesData: enemiesData,
		ChestsData:  chestsData,
		PhaseTimer:  int(time.Until(w.phaseTimer).Milliseconds()),
		Phase:       int(w.phase),
		Round:       w.round,
	})

	w.Mu.RUnlock()

	if err != nil {
		log.Printf("couldn't marshal message: %v", err)
	}

	c.Player.Mu.RLock()
	if c.Player.Client == c {
		c.Player.Mu.RUnlock()
		err := wsjson.Write(context.Background(), c.Conn, msg)
		if err != nil {
			log.Println(err)
		}
	} else {
		c.Player.Mu.RUnlock()
	}
}
