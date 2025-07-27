package main

import (
	"fmt"
)

type GameState struct {
	Board [6][7]int
	CurrentPlayer int
	NextPlayer int
	LastMove int
}

func (g *GameState) Drop(col int) bool {

	if col < 0 || col > 6 {
		return false
	}

	for i := 5; i >= 0; i-- {
		if g.Board[i][col] == 0 {
			g.Board[i][col] = g.CurrentPlayer
			g.LastMove = col
			g.NextPlayer = g.CurrentPlayer
			g.CurrentPlayer = -g.CurrentPlayer
			return true
		}
	}
	return false
}


func (g *GameState) CheckWin() int {
    board := g.Board

    // Horizontal
    for i := 0; i < 6; i++ {
        for j := 0; j < 4; j++ {
            if board[i][j] != 0 &&
                board[i][j] == board[i][j+1] &&
                board[i][j] == board[i][j+2] &&
                board[i][j] == board[i][j+3] {
                return board[i][j]
            }
        }
    }

    // Vertical
    for i := 0; i < 3; i++ {
        for j := 0; j < 7; j++ {
            if board[i][j] != 0 &&
                board[i][j] == board[i+1][j] &&
                board[i][j] == board[i+2][j] &&
                board[i][j] == board[i+3][j] {
                return board[i][j]
            }
        }
    }

    // Diagonal ↘
    for i := 0; i < 3; i++ {
        for j := 0; j < 4; j++ {
            if board[i][j] != 0 &&
                board[i][j] == board[i+1][j+1] &&
                board[i][j] == board[i+2][j+2] &&
                board[i][j] == board[i+3][j+3] {
                return board[i][j]
            }
        }
    }

    // Diagonal ↗
    for i := 3; i < 6; i++ {
        for j := 0; j < 4; j++ {
            if board[i][j] != 0 &&
                board[i][j] == board[i-1][j+1] &&
                board[i][j] == board[i-2][j+2] &&
                board[i][j] == board[i-3][j+3] {
                return board[i][j]
            }
        }
    }

    return 0 // no winner yet
}


func (g *GameState) IsDraw() bool {
    for _, row := range g.Board {
        for _, cell := range row {
            if cell == 0 {
                return false
            }
        }
    }
    return g.CheckWin() == 0
}


func (g *GameState) Clone() *GameState {
	newBoard := g.Board //Arrays are value-copied in Go

	return &GameState{
		Board:          newBoard,
		CurrentPlayer:  g.CurrentPlayer,
		NextPlayer:     g.NextPlayer,
		LastMove:       g.LastMove,
	}
}

func (g *GameState) ValidMoves() []int {
    moves := []int{}
    for col := 0; col < 7; col++ {
        if g.Board[0][col] == 0 {
            moves = append(moves, col)
        }
    }
    return moves
}


func (g *GameState) Print() {
    for i := range g.Board {
        for j := range g.Board[i] {
            switch g.Board[i][j] {
            case 1:
                fmt.Print("X")
            case -1:
                fmt.Print("O")
            default:
                fmt.Print(".")
            }
        }
        fmt.Println()
    }
}

