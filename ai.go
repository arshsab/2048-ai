package main

import (
	"fmt"
	// "math/rand"
	"time"
)

type Stats struct {
	leaves int
	dp     map[uint64]DpEntry
}

type DpEntry struct {
	depth int
	score float64
}

func ChooseBestMove(board uint64, d time.Duration) (best Move) {
	end := time.Now().Add(d)
	best = NONE

	i := 1

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from timeout at depth: ", i)
		}
	}()

	// iterative deepening
	for ; ; i++ {
		stats := Stats{leaves: 0, dp: make(map[uint64]DpEntry)}

		move, score := max(board, 0, i, &stats, end)
		best = move

		fmt.Println("Finished Search at depth", i, "with score:", score, "and move:", move)
		fmt.Println("Leaves: ", stats.leaves, " DpEntries: ", len(stats.dp))
	}

	return
}

func max(board uint64, depth int, limit int, stats *Stats, end time.Time) (Move, float64) {
	if time.Now().After(end) {
		panic("Timed out!")
	}

	if depth >= limit {
		stats.leaves++

		return NONE, float64(heuristic(board))
	}

	if dpVal, ok := stats.dp[board]; ok && dpVal.depth <= depth {
		return NONE, dpVal.score
	}

	bestScore := 0.0
	bestMove := NONE

	for _, move := range Moves {
		moved := MakeMove(board, move)

		if moved == board {
			continue
		}

		score := expectation(moved, depth, limit, stats, end)

		if score >= bestScore {
			bestScore = score
			bestMove = move
		}
	}

	stats.dp[board] = DpEntry{depth: depth, score: bestScore}

	return bestMove, bestScore
}

func expectation(board uint64, depth int, limit int, stats *Stats, end time.Time) float64 {
	empty := 0
	total := 0.0

	for i := uint(0); i < 16; i++ {
		if (board & (0xf << (i * 4))) == 0 {
			board1 := board ^ (0x1 << (i * 4))
			board2 := board ^ (0x2 << (i * 4))

			_, score1 := max(board1, depth+1, limit, stats, end)
			_, score2 := max(board2, depth+2, limit, stats, end)

			total += (score1 * .9) + (score2 * .1)

			empty++
		}
	}

	return total / float64(empty)
}

func heuristic(board uint64) int {
	count := 0

	for {
		score := -10000.0
		moved := board

		for _, move := range Moves {
			potential := MakeMove(board, move)

			if potential == board {
				continue
			}

			potentialScore := simpleHeuristic(potential)

			if potentialScore < score {
				score = potentialScore
				moved = potential
			}
		}

		if moved == board {
			break
		}

		count++
		board = InsertRandomTile(moved)
	}

	return count
}

func spaces(board uint64) float64 {
	mask := uint64(0xf)
	empties := 0

	for i := 0; i < 16; i++ {
		if (board & mask) == 0 {
			empties++
		}

		mask <<= 4
	}

	return float64(empties)
}

func simpleHeuristic(board uint64) float64 {
	mask := uint64(0xf)
	empties := 0

	for i := 0; i < 16; i++ {
		if board&mask == 0 {
			empties++
		}

		mask <<= 4
	}

	ret := gradeRow(board>>48) + gradeRow((board>>32)&0xffff) + gradeRow((board>>16)&0xffff) + gradeRow(board&0xffff)

	return float64(ret) + (2.0 * float64(empties))
}

func gradeRow(row uint64) float64 {
	switches := 0
	magnitude := 0.0
	stability := 4

	pos := sgn((row>>12)&0xf - (row>>8)&0xf)
	for i := 3; i > 0; i-- {
		one := (row >> uint((i * 4))) & 0xf
		two := (row >> uint(((i - 1) * 4))) & 0xf

		if one == 0 {
			stability--
			continue
		}

		sign := sgn((one - two))

		if sign != pos {
			switches++
		}

		magnitude += float64(one)
	}

	if (row & 0xf) == 0 {
		stability--
	}

	magnitude += float64((row & 0xf))

	return (float64((stability - switches)) * magnitude) + (2.0 * (magnitude / float64((stability + 1))))
}

func stability(board uint64) float64 {
	afters := [4]uint64{ShiftUp(board), ShiftDown(board), ShiftLeft(board), ShiftRight(board)}
	mask := uint64(0xf)

	total := uint64(0)

	for i := 0; i < 16; i++ {
		tile := (board & mask)
		val := tile >> uint((4 * i))

		for _, moved := range afters {
			if (moved & mask) == tile {
				total += val * val // square for added effect
			}
		}

		mask <<= 4
	}

	return float64(total)
}

func sgn(s uint64) int {
	if s < 0 {
		return -1
	}

	return 1
}
