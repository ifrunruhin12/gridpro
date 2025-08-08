package backend

import (
	"math/rand"
)

const (
	MaxDepth = 10 // Increased depth for stronger play
	winScore = 1000000
)

// Opening book for first few moves
var openingBook = map[string]int{
	// Empty board - always start in center
	"000000000000000000000000000000000000000000": 3,
	// Common responses after center start
	"000000000000000000000000000000020000000000": 3, // Opponent plays center
	"000000000000000000000000000000000200000000": 3, // Opponent plays adjacent
	"000000000000000000000000000000000020000000": 2, // Opponent plays one further
}

// GetAIMove finds the best move for the AI
func GetAIMove(board *Board) int {
	if board == nil {
		return -1
	}

	// Check for immediate win
	if move := findWinningMove(board, AI); move != -1 {
		return move
	}

	// Check for opponent's immediate threat
	if move := findWinningMove(board, Player); move != -1 {
		return move
	}

	// Check opening book
	if move, found := lookupOpeningBook(board); found {
		return move
	}

	// For endgame, use perfect play if possible
	if countPieces(board) >= 30 { // Late game
		if move := findForcedWin(board); move != -1 {
			return move
		}
	}

	// Default to minimax
	return getBestMove(board)
}

func lookupOpeningBook(board *Board) (int, bool) {
	key := ""
	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			key += string('0' + byte(board.Grid[row][col]))
		}
	}
	move, exists := openingBook[key]
	return move, exists
}

func countPieces(board *Board) int {
	count := 0
	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			if board.Grid[row][col] != Empty {
				count++
			}
		}
	}
	return count
}

func findForcedWin(board *Board) int {
	// Try all possible moves to see if any lead to a forced win
	for _, col := range board.ValidMoves() {
		sim := board.Clone()
		sim.Drop(col)
		if sim.CheckWin() == AI {
			return col
		}
		// If opponent can't win in their next move
		if !canOpponentWin(sim) {
			return col
		}
	}
	return -1
}

func canOpponentWin(board *Board) bool {
	for _, col := range board.ValidMoves() {
		sim := board.Clone()
		sim.CurrentTurn = Player
		if sim.Drop(col) && sim.CheckWin() == Player {
			return true
		}
	}
	return false
}

func findWinningMove(board *Board, player int) int {
	for _, col := range board.ValidMoves() {
		sim := board.Clone()
		sim.CurrentTurn = player
		if sim.Drop(col) && sim.CheckWin() == player {
			return col
		}
	}
	return -1
}

func getBestMove(board *Board) int {
	bestScore := -1000000
	bestMoves := []int{}
	alpha := -1000000
	beta := 1000000

	// Try center first (best initial move)
	if board.IsValidMove(3) {
		return 3
	}

	// Then try other columns
	for _, col := range board.ValidMoves() {
		sim := board.Clone()
		sim.Drop(col)

		score := minimax(sim, MaxDepth-1, false, alpha, beta)

		if score > bestScore {
			bestScore = score
			bestMoves = []int{col}
		} else if score == bestScore {
			bestMoves = append(bestMoves, col)
		}

		// Update alpha
		if score > alpha {
			alpha = score
		}

		// Prune if possible
		if alpha >= beta {
			break
		}
	}

	// Return random move among best moves
	if len(bestMoves) > 0 {
		return bestMoves[rand.Intn(len(bestMoves))]
	}

	// Fallback: return first valid move
	moves := board.ValidMoves()
	if len(moves) > 0 {
		return moves[0]
	}

	return -1
}

func minimax(board *Board, depth int, isMaximizing bool, alpha, beta int) int {
	// Check terminal states
	winner := board.CheckWin()
	if winner == AI {
		return winScore + depth // Prefer faster wins
	}
	if winner == Player {
		return -winScore - depth // Prefer slower losses
	}
	if board.IsDraw() || depth == 0 {
		return evaluateBoard(board)
	}

	if isMaximizing {
		maxEval := -1000000
		for _, col := range board.ValidMoves() {
			sim := board.Clone()
			sim.Drop(col)
			score := minimax(sim, depth-1, false, alpha, beta)
			if score > maxEval {
				maxEval = score
			}
			// Update alpha
			if score > alpha {
				alpha = score
			}
			// Prune if possible
			if beta <= alpha {
				break
			}
		}
		return maxEval
	} else {
		minEval := 1000000
		for _, col := range board.ValidMoves() {
			sim := board.Clone()
			sim.CurrentTurn = Player
			sim.Drop(col)
			score := minimax(sim, depth-1, true, alpha, beta)
			if score < minEval {
				minEval = score
			}
			// Update beta
			if score < beta {
				beta = score
			}
			// Prune if possible
			if beta <= alpha {
				break
			}
		}
		return minEval
	}
}

func evaluateBoard(board *Board) int {
	if board == nil {
		return 0
	}

	score := 0

	// Check for immediate wins/losses
	switch board.CheckWin() {
	case AI:
		return winScore
	case Player:
		return -winScore
	}

	// Evaluate based on potential threats and opportunities
	score += evaluateThreats(board, AI)*10 - evaluateThreats(board, Player)*10

	// Enhanced center control (higher weight for center column)
	centerWeights := []int{1, 2, 3, 4, 3, 2, 1}
	for col := 0; col < Cols; col++ {
		if col < len(centerWeights) {
			for row := 0; row < Rows; row++ {
				switch board.Grid[row][col] {
				case AI:
					score += centerWeights[col]
				case Player:
					score -= centerWeights[col]
				}
			}
		}
	}

	return score
}

func evaluateThreats(board *Board, player int) int {
	threats := 0

	// Check horizontal threats
	for row := 0; row < Rows; row++ {
		for col := 0; col <= Cols-4; col++ {
			count := 0
			empty := 0
			for i := 0; i < 4; i++ {
				if board.Grid[row][col+i] == player {
					count++
				} else if board.Grid[row][col+i] == Empty {
					empty++
				} else {
					count = 0
					break
				}
			}
			if count >= 3 && empty == 1 {
				threats++
			}
		}
	}

	// Check vertical threats
	for row := 0; row <= Rows-4; row++ {
		for col := 0; col < Cols; col++ {
			count := 0
			empty := 0
			for i := 0; i < 4; i++ {
				if board.Grid[row+i][col] == player {
					count++
				} else if board.Grid[row+i][col] == Empty {
					empty++
				} else {
					count = 0
					break
				}
			}
			if count >= 3 && empty == 1 {
				threats++
			}
		}
	}

	// Check diagonal down-right (\ direction)
	for row := 0; row <= Rows-4; row++ {
		for col := 0; col <= Cols-4; col++ {
			count := 0
			empty := 0
			for i := 0; i < 4; i++ {
				if board.Grid[row+i][col+i] == player {
					count++
				} else if board.Grid[row+i][col+i] == Empty {
					empty++
				} else {
					count = 0
					break
				}
			}
			if count >= 3 && empty == 1 {
				threats++
			}
		}
	}

	// Check diagonal up-right (/ direction)
	for row := 3; row < Rows; row++ {
		for col := 0; col <= Cols-4; col++ {
			count := 0
			empty := 0
			for i := 0; i < 4; i++ {
				if board.Grid[row-i][col+i] == player {
					count++
				} else if board.Grid[row-i][col+i] == Empty {
					empty++
				} else {
					count = 0
					break
				}
			}
			if count >= 3 && empty == 1 {
				threats++
			}
		}
	}

	return threats
}
