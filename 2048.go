package main

import (
	"fmt"
	"math/rand"
)

type Move int

const (
	UP Move = iota
	DOWN
	LEFT
	RIGHT
	NONE
)

var Moves = [4]Move{UP, DOWN, LEFT, RIGHT}

var shiftr_table [65536]uint64
var shiftl_table [65536]uint64

func init() {
	for i := range shiftr_table {
		shiftr_table[i] = shiftRowRight(uint64(i))
		shiftl_table[i] = reverse(shiftRowRight(reverse(uint64(i))))
	}

	fmt.Println("Setup shift tables.")
}

func MakeMove(board uint64, move Move) uint64 {
	switch move {
	case UP:
		return ShiftUp(board)
	case DOWN:
		return ShiftDown(board)
	case LEFT:
		return ShiftLeft(board)
	case RIGHT:
		return ShiftRight(board)
	default:
		return board
	}
}

func ShiftUp(board uint64) uint64 {
	return Transpose(ShiftLeft(Transpose(board)))
}

func ShiftDown(board uint64) uint64 {
	return Transpose(ShiftRight(Transpose(board)))
}

func ShiftLeft(board uint64) uint64 {
	row1 := (board >> 48) & 0xffff
	row2 := (board >> 32) & 0xffff
	row3 := (board >> 16) & 0xffff
	row4 := board & 0xffff

	return (shiftl_table[row1] << 48) ^ (shiftl_table[row2] << 32) ^ (shiftl_table[row3] << 16) ^ shiftl_table[row4]
}

func ShiftRight(board uint64) uint64 {
	row1 := (board >> 48) & 0xffff
	row2 := (board >> 32) & 0xffff
	row3 := (board >> 16) & 0xffff
	row4 := board & 0xffff

	return (shiftr_table[row1] << 48) ^ (shiftr_table[row2] << 32) ^ (shiftr_table[row3] << 16) ^ shiftr_table[row4]
}

func Transpose(board uint64) uint64 {
	// Diagonals have the same shift, so isolate diagonals.

	diag := board & 0xf0000f0000f0000f

	top1 := (board & 0x0f0000f0000f0000) >> 12
	top2 := (board & 0x00f0000f00000000) >> 24
	top3 := (board & 0x000f000000000000) >> 36

	low1 := (board & 0x0000f0000f0000f0) << 12
	low2 := (board & 0x00000000f0000f00) << 24
	low3 := (board & 0x000000000000f000) << 36

	return diag ^ top1 ^ top2 ^ top3 ^ low1 ^ low2 ^ low3
}

func Dead(board uint64) bool {
	for _, move := range Moves {
		if MakeMove(board, move) != board {
			return false
		}
	}

	return true
}

func DecodeBoard(board uint64) [4][4]uint {
	var ret [4][4]uint

	ret[0] = decode((board >> 48) & 0xffff)
	ret[1] = decode((board >> 32) & 0xffff)
	ret[2] = decode((board >> 16) & 0xffff)
	ret[3] = decode(board & 0xffff)

	return ret
}

func PrintBoard(board uint64) {
	decoded := DecodeBoard(board)

	fmt.Println("Board:")

	for _, val := range decoded {
		fmt.Println("    ", val)
	}
}

func InsertRandomTile(board uint64) uint64 {
	empty := 0
	mask := uint64(0xf)

	for i := 0; i < 16; i++ {
		if (board & mask) == 0 {
			empty++
		}

		mask <<= 4
	}

	pos := rand.Intn(empty)
	til := 1

	if rand.Float64() < 0.1 {
		til = 2
	}

	mask = 0xf

	for i := uint(0); i < 16; i++ {
		if (board & mask) == 0 {
			empty--
		}

		if empty == pos {
			til <<= (i * 4)

			return board ^ uint64(til)
		}

		mask <<= 4
	}

	panic("Should have returned by now.")
}

func reverse(row uint64) uint64 {
	var ret uint64 = 0

	ret ^= (row & 0xf) << 12
	ret ^= (row & 0xf0) << 4
	ret ^= (row & 0xf00) >> 4
	ret ^= (row & 0xf000) >> 12

	return ret
}

func shiftRowRight(row uint64) uint64 {
	return uint64(encode(pushRight(combineRight(pushRight(decode(row))))))
}

func encode(row [4]uint) uint64 {
	var ret uint64 = 0

	for _, value := range row {
		ret <<= 4

		ret = ret ^ uint64(value)
	}

	return ret
}

func decode(row uint64) [4]uint {
	return [4]uint{uint(row >> 12), uint((row >> 8) & 0xf), uint((row >> 4) & 0xf), uint(row & 0xf)}
}

func pushRight(arr [4]uint) [4]uint {
	swapInd := 3

	for i := 3; i >= 0; i-- {
		if arr[i] != 0 {
			temp := arr[i]
			arr[i] = arr[swapInd]
			arr[swapInd] = temp

			swapInd--
		}
	}

	return arr
}

func combineRight(arr [4]uint) [4]uint {
	for i := 3; i > 0; i-- {
		if arr[i] != 0 && arr[i] == arr[i-1] {
			arr[i] = arr[i] + 1
			arr[i-1] = 0
		}
	}

	return arr
}
