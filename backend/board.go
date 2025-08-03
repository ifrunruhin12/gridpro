package backend

const (
	Empty  = 0
	Player = 1
	AI     = 2
	Rows   = 6
	Cols   = 7
)

type Board struct {
	Grid         [Rows][Cols]int `json:"grid"`
	CurrentTurn  int             `json:"current_turn"` // 1 = human, 2 = AI
	LastMoveRow  int             `json:"last_move_row"`
	LastMoveCol  int             `json:"last_move_col"`
}

// Drop places a piece in the selected column.
// Returns true if the move is valid, false otherwise.
func (b *Board) Drop(col int) bool {
	if col < 0 || col >= Cols {
		return false
	}
	for row := Rows - 1; row >= 0; row-- {
		if b.Grid[row][col] == Empty {
			b.Grid[row][col] = b.CurrentTurn
			b.LastMoveRow = row
			b.LastMoveCol = col
			b.toggleTurn()
			return true
		}
	}
	return false
}

func (b *Board) toggleTurn() {
	if b.CurrentTurn == Player {
		b.CurrentTurn = AI
	} else {
		b.CurrentTurn = Player
	}
}

// CheckWin checks if the last move resulted in a win.
// Returns 0 if no win, 1 if Player wins, 2 if AI wins.
func (b *Board) CheckWin() int {
	lastRow := b.LastMoveRow
	lastCol := b.LastMoveCol
	if lastRow == -1 || lastCol == -1 {
		return 0
	}
	player := b.Grid[lastRow][lastCol]

	directions := [][2]int{
		{0, 1},  // horizontal
		{1, 0},  // vertical
		{1, 1},  // diagonal /
		{1, -1}, // diagonal \
	}

	for _, d := range directions {
		count := 1
		count += b.countDirection(lastRow, lastCol, d[0], d[1], player)
		count += b.countDirection(lastRow, lastCol, -d[0], -d[1], player)
		if count >= 4 {
			return player
		}
	}
	return 0
}

// countDirection counts the number of consecutive same-colored discs in one direction
func (b *Board) countDirection(row, col, dr, dc, player int) int {
	count := 0
	for {
		row += dr
		col += dc
		if row < 0 || row >= Rows || col < 0 || col >= Cols || b.Grid[row][col] != player {
			break
		}
		count++
	}
	return count
}

// IsDraw returns true if the board is full and no winner
func (b *Board) IsDraw() bool {
	for c := 0; c < Cols; c++ {
		if b.Grid[0][c] == Empty {
			return false
		}
	}
	return b.CheckWin() == 0
}

// Clone creates a deep copy of the board
func (b *Board) Clone() *Board {
	newBoard := *b
	return &newBoard
}

func (b *Board) IsValidMove(col int) bool {
	return col >= 0 && col < Cols && b.Grid[0][col] == Empty
}

func (b *Board) ValidMoves() []int {
	moves := []int{}
	for col := 0; col < Cols; col++ {
		if b.IsValidMove(col) {
			moves = append(moves, col)
		}
	}
	return moves
}

