package main

import (
	"fmt"
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

		return NONE, float64(Heuristic(board))
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
