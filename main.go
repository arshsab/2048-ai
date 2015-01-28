package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	// todo: random board
	var board uint64 = 0x11

	PrintBoard(board)
	for !Dead(board) {
		move := ChooseBestMove(board, 1*time.Second)

		board = MakeMove(board, move)

		fmt.Println("Made Move: ", move)
		PrintBoard(board)

		board = InsertRandomTile(board)

		fmt.Println("Inserted random tile:")
		PrintBoard(board)

		fmt.Println("Max Score: ", MAX_SCORE)
	}
}
