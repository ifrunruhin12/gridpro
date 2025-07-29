package backend

import (
	"math/rand"
)

func GetRandomMove(g *GameState) int {
	moves := g.ValidMoves()
	if len(moves) == 0 {
		return -1
	}
	return moves[rand.Intn(len(moves))]
}
