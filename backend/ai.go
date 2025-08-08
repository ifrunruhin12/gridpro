package backend

import (
	"math/rand"
)

const (
	MaxDepth = 6  // Increased depth for better AI with pruning
)

// GetAIMove finds the best move for the AI using Minimax with Alpha-Beta pruning
func GetAIMove(board *Board) int {
	bestScore := -1000000
	bestMoves := []int{}
	alpha := -1000000
	beta := 1000000

	// Get all valid moves
	validMoves := board.ValidMoves()
	
	// If only one move is available, return it immediately
	if len(validMoves) == 1 {
		return validMoves[0]
	}

	// Try each valid move
	for _, col := range validMoves {
		sim := board.Clone()
		sim.Drop(col)

		// Start the minimax search
		score := minimax(sim, MaxDepth-1, false, alpha, beta)

		// Update best move if this score is better
		if score > bestScore {
			bestScore = score
			bestMoves = []int{col}
			// Update alpha (best already explored option for max player)
			if score > alpha {
				alpha = score
			}
		} else if score == bestScore {
			bestMoves = append(bestMoves, col)
		}

		// Alpha-Beta Pruning
		if alpha >= beta {
			break
		}
	}

	// If multiple moves have the same score, choose one at random
	if len(bestMoves) > 0 {
		return bestMoves[rand.Intn(len(bestMoves))]
	}

	// Fallback: return first valid move if no best move found
	if len(validMoves) > 0 {
		return validMoves[0]
	}

	return -1
}

// minimax implements the minimax algorithm with alpha-beta pruning
func minimax(board *Board, depth int, isMaximizing bool, alpha, beta int) int {
	// Check terminal states
	winner := board.CheckWin()
	if winner == AI {
		return 1000 + depth // Prefer faster wins
	}
	if winner == Player {
		return -1000 - depth // Prefer slower losses
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

// evaluateBoard provides a simple evaluation of the board position
// Positive score is good for AI, negative is good for player
func evaluateBoard(board *Board) int {
	score := 0
	
	// Check center control (columns 2,3,4 are more valuable)
	centerCols := []int{2, 3, 4}
	for _, col := range centerCols {
		for row := 0; row < Rows; row++ {
			if board.Grid[row][col] == AI {
				score += 2
			} else if board.Grid[row][col] == Player {
				score -= 2
			}
		}
	}
	
	// Check for potential 3-in-a-row
	score += countPotentialWins(board, AI) * 5
	score -= countPotentialWins(board, Player) * 5
	
	return score
}

// countPotentialWins counts the number of potential winning positions for a player
func countPotentialWins(board *Board, player int) int {
	count := 0
	// Check all possible 4-in-a-row positions
	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			// Skip if this position is already occupied by the opponent
			if board.Grid[row][col] != Empty && board.Grid[row][col] != player {
				continue
			}
			
			// Check horizontal
			if col <= Cols-4 {
				potential := true
				for i := 0; i < 4; i++ {
					if board.Grid[row][col+i] != Empty && board.Grid[row][col+i] != player {
						potential = false
						break
					}
				}
				if potential {
					count++
				}
			}
			
			// Check vertical
			if row <= Rows-4 {
				potential := true
				for i := 0; i < 4; i++ {
					if board.Grid[row+i][col] != Empty && board.Grid[row+i][col] != player {
						potential = false
						break
					}
				}
				if potential {
					count++
				}
			}
			
			// Check diagonal down-right
			if row <= Rows-4 && col <= Cols-4 {
				potential := true
				for i := 0; i < 4; i++ {
					if board.Grid[row+i][col+i] != Empty && board.Grid[row+i][col+i] != player {
						potential = false
						break
					}
				}
				if potential {
					count++
				}
			}
			
			// Check diagonal up-right
			if row >= 3 && col <= Cols-4 {
				potential := true
				for i := 0; i < 4; i++ {
					if board.Grid[row-i][col+i] != Empty && board.Grid[row-i][col+i] != player {
						potential = false
						break
					}
				}
				if potential {
					count++
				}
			}
		}
	}
	return count
}
