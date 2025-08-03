package backend

import (
	"math/rand"
)

const (
	MaxDepth = 4
)

func GetAIMove(board *Board) int {
	bestScore := -1000000
	bestMoves := []int{}

	for _, col := range board.ValidMoves() {
		sim := board.Clone()
		sim.Drop(col)

		score := minimax(sim, MaxDepth-1, false)

		if score > bestScore {
			bestScore = score
			bestMoves = []int{col}
		} else if score == bestScore {
			bestMoves = append(bestMoves, col)
		}
	}

	if len(bestMoves) > 0 {
		return bestMoves[rand.Intn(len(bestMoves))]
	}

	return -1
}

func minimax(board *Board, depth int, isMaximizing bool) int {
	winner := board.CheckWin()
	if winner == AI {
		return 1000
	}
	if winner == Player {
		return -1000
	}
	if board.IsDraw() || depth == 0 {
		return 0
	}

	if isMaximizing {
		maxEval := -1000000
		for _, col := range board.ValidMoves() {
			sim := board.Clone()
			sim.Drop(col)
			score := minimax(sim, depth-1, false)
			if score > maxEval {
				maxEval = score
			}
		}
		return maxEval
	} else {
		minEval := 1000000
		for _, col := range board.ValidMoves() {
			sim := board.Clone()
			sim.CurrentTurn = Player
			sim.Drop(col)
			score := minimax(sim, depth-1, true)
			if score < minEval {
				minEval = score
			}
		}
		return minEval
	}
}

