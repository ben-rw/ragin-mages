package room

import (
// "math/rand"
)

// create index to choose from 8 player sprites randomly
func NewPlayerSpriteIndex() *[]int {
	playerSpriteIndex := make([]int, 8)
	for i := 0; i < 8; i++ {
		playerSpriteIndex[i] = i
	}

	return &playerSpriteIndex
}

// returns sprite index, remove index from playerSpriteIndex to avoid duplicates
// returns sprites based on join order
func (r *Room) AssignPlayerSprite() int {
	r.Mu.Lock()
	spriteIndex := (*r.PlayerSpriteIndex)[0]
	*r.PlayerSpriteIndex = (*r.PlayerSpriteIndex)[1:]
	r.Mu.Unlock()

	return spriteIndex
}

// returns random sprites
// func (r *Room) AssignPlayerSprite() int {
// 	r.Mu.Lock()
// 	i := rand.Intn(len(*r.PlayerSpriteIndex))
// 	s := *r.PlayerSpriteIndex
// 	spriteIndex := s[i]
// 	s[i] = s[len(s)-1]
// 	*r.PlayerSpriteIndex = s[:len(s)-1]
// 	r.Mu.Unlock()
//
// 	return spriteIndex
// }
