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

// Preferred column order for move ordering (center-out)
var preferredCols = []int{3, 2, 4, 1, 5, 0, 6}

// Simple transposition table for minimax caching
var transposition = map[string]int{}

// GetAIMove finds the best move for the AI
func GetAIMove(board *Board) int {
	if board == nil {
		return -1
	}

	// If the board is empty and AI starts, force center
	if countPieces(board) == 0 {
		return 3
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

	// Search using preferred move ordering (center-out)
	for _, col := range orderedMoves(board) {
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
	// Transposition table lookup
	if val, ok := transposition[boardKey(board, depth, isMaximizing)]; ok {
		return val
	}

	// Check terminal states
	winner := board.CheckWin()
	if winner == AI {
		return winScore + depth // Prefer faster wins
	}
	if winner == Player {
		return -winScore - depth // Prefer slower losses
	}
	if board.IsDraw() || depth == 0 {
		val := evaluateBoard(board)
		transposition[boardKey(board, depth, isMaximizing)] = val
		return val
	}

	if isMaximizing {
		maxEval := -1000000
		for _, col := range orderedMoves(board) {
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
		transposition[boardKey(board, depth, isMaximizing)] = maxEval
		return maxEval
	} else {
		minEval := 1000000
		for _, col := range orderedMoves(board) {
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
		transposition[boardKey(board, depth, isMaximizing)] = minEval
		return minEval
	}
}

func evaluateBoard(board *Board) int {
	if board == nil {
		return 0
	}

	score := 0

	// Immediate wins/losses
	switch board.CheckWin() {
	case AI:
		return winScore
	case Player:
		return -winScore
	}

	// Center column control
	centerWeights := []int{1, 2, 3, 5, 3, 2, 1}
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

	// Window-based scoring (classic Connect 4 heuristic)
	scoreWindow := func(a, b, c, d int) int {
		ai, opp, emp := 0, 0, 0
		vals := [4]int{a, b, c, d}
		for _, v := range vals {
			switch v {
			case AI:
				ai++
			case Player:
				opp++
			default:
				emp++
			}
		}
		// Offensive
		if ai == 4 {
			return 100000
		}
		if ai == 3 && emp == 1 {
			return 200
		}
		if ai == 2 && emp == 2 {
			return 20
		}
		// Defensive (slightly higher to prioritize blocking)
		if opp == 3 && emp == 1 {
			return -300
		}
		if opp == 2 && emp == 2 {
			return -30
		}
		return 0
	}

	// Horizontal
	for row := 0; row < Rows; row++ {
		for col := 0; col <= Cols-4; col++ {
			score += scoreWindow(
				board.Grid[row][col],
				board.Grid[row][col+1],
				board.Grid[row][col+2],
				board.Grid[row][col+3],
			)
		}
	}
	// Vertical
	for row := 0; row <= Rows-4; row++ {
		for col := 0; col < Cols; col++ {
			score += scoreWindow(
				board.Grid[row][col],
				board.Grid[row+1][col],
				board.Grid[row+2][col],
				board.Grid[row+3][col],
			)
		}
	}
	// Diagonal \
	for row := 0; row <= Rows-4; row++ {
		for col := 0; col <= Cols-4; col++ {
			score += scoreWindow(
				board.Grid[row][col],
				board.Grid[row+1][col+1],
				board.Grid[row+2][col+2],
				board.Grid[row+3][col+3],
			)
		}
	}
	// Diagonal /
	for row := 3; row < Rows; row++ {
		for col := 0; col <= Cols-4; col++ {
			score += scoreWindow(
				board.Grid[row][col],
				board.Grid[row-1][col+1],
				board.Grid[row-2][col+2],
				board.Grid[row-3][col+3],
			)
		}
	}

	return score
}

// Preferred move ordering filtered by current valid moves
func orderedMoves(board *Board) []int {
	valid := board.ValidMoves()
	present := make(map[int]bool, len(valid))
	for _, c := range valid {
		present[c] = true
	}
	res := make([]int, 0, len(valid))
	for _, c := range preferredCols {
		if present[c] {
			res = append(res, c)
		}
	}
	return res
}

// Board key for transposition table
func boardKey(b *Board, depth int, isMax bool) string {
	buf := make([]byte, 0, Rows*Cols+3)
	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			buf = append(buf, byte('0'+b.Grid[row][col]))
		}
	}
	buf = append(buf, byte('0'+b.CurrentTurn))
	buf = append(buf, byte(depth%10+'0'))
	if isMax {
		buf = append(buf, 'M')
	} else {
		buf = append(buf, 'm')
	}
	return string(buf)
}
