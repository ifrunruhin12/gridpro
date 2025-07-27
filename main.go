package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	game := &GameState{
		CurrentPlayer: 1,
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nCurrent Board:")
		game.Print()

		if game.CheckWin() != 0 {
			fmt.Printf("\n🎉 Player %d wins!\n", -game.CurrentPlayer) // because turn already flipped
			break
		}

		if game.IsDraw() {
			fmt.Println("\n🤝 It's a draw!")
			break
		}

		fmt.Printf("Player %d, enter column (0-6): ", game.CurrentPlayer)

		scanner.Scan()
		input := scanner.Text()

		col, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || col < 0 || col > 6 {
			fmt.Println("Invalid input. Enter a number from 0 to 6.")
			continue
		}

		if !game.Drop(col) {
			fmt.Println("That column is full. Try another.")
		}
	}
}

