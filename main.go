package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"sort"
	// "time"
)

func main() {
	runAi()
}

func testPositions() {
	fmt.Println("Loaded.")

	boards := loadPositions()

	fmt.Println("Loaded.")

	for _, board := range boards {
		total := 0

		arr := []int{}

		for i := 0; i < 10; i++ {
			score := randomPlay(board)
			total += score

			arr = append(arr, score)

			fmt.Println("Score", (float64(total) / float64((i + 1))), "iteration", i)
		}

		sort.Ints(arr)

		fmt.Println(arr)
	}
}

func randomPlay(board uint64) int {
	ret := 0

	for !Dead(board) {
		move := Moves[rand.Int31n(4)]
		moved := MakeMove(board, move)

		if moved == board {
			continue
		}

		board = InsertRandomTile(moved)

		// fmt.Println(board, ret)
		ret++
	}

	return ret
}

func runAi() {
	var board uint64 = 0

	board = InsertRandomTile(board)
	board = InsertRandomTile(board)

	boards := []uint64{}

	PrintBoard(board)
	for !Dead(board) {
		boards = append(boards, board)

		move := MonteCarloChooseMove(board)

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
