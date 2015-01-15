package main

import (
	"fmt"
	"time"
)

func main() {
	// todo: random board
	var board uint64 = 0x22

	PrintBoard(board)
	for !Dead(board) {
		move := ChooseBestMove(board, 200*time.Millisecond)

		board = MakeMove(board, move)

		fmt.Println("Made Move: ", move)
		PrintBoard(board)

		board = InsertRandomTile(board)

		fmt.Println("Inserted random tile:")
		PrintBoard(board)
	}
}
