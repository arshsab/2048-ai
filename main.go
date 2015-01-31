package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

func mian() {
	var board uint64 = 0

	board = InsertRandomTile(board)
	board = InsertRandomTile(board)

	boards := []uint64{}

	PrintBoard(board)
	for !Dead(board) {
		boards = append(boards, board)

		move := ChooseBestMove(board, (50 * time.Millisecond), DEFAULT_WEIGHTS)

		board = MakeMove(board, move)

		fmt.Println("Made Move: ", move)
		PrintBoard(board)

		board = InsertRandomTile(board)

		fmt.Println("Inserted random tile:")
		PrintBoard(board)
	}

	fmt.Println("Saving boards:")
	fmt.Println(len(boards))

	savePositions(boards)
}

func loadPositions() []uint64 {
	positions := []uint64{}
	f, err := os.Open("positions")

	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return nil
	}

	r := bufio.NewReader(f)
	dec := gob.NewDecoder(r)

	dec.Decode(&positions)

	return positions
}

func savePositions(positions []uint64) {
	os.Remove("positions")
	f, err := os.Create("positions")

	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)

	enc.Encode(positions)

	w.Flush()
}
