package main

var SPACES_WEIGHT float64 = 1.0
var STABILITY_WEIGHT float64 = 1.0
var UNIFORMITY_WEIGHT float64 = 1.0

var MAX_SCORE int = 100000

func OptimizeHeuristics() {
	// todo
}

func SetHeuristicValues(spaces float64, stability float64, uniformity float64, maxScore int) {
	SPACES_WEIGHT = spaces
	STABILITY_WEIGHT = stability
	UNIFORMITY_WEIGHT = uniformity
	MAX_SCORE = maxScore
}

func Heuristic(board uint64) int {
	count := 0

	for count < MAX_SCORE {
		score := -10000.0
		moved := board

		for _, move := range Moves {
			potential := MakeMove(board, move)

			if potential == board {
				continue
			}

			potentialScore := weighted(potential)

			if potentialScore > score {
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

func weighted(board uint64) (score float64) {
	score = spaces(board)
	score += stability(board)
	score -= nonuniformity(board) + nonuniformity(Transpose(board))

	return
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

func nonuniformity(board uint64) float64 {
	total := uint64(0)

	for i := 0; i < 4; i++ {
		row := (board & 0xffff)

		for j := 0; j < 3; j++ {
			diff := (row & 0xf) - ((row >> 1) & 0xf)

			total += diff * diff

			row >>= 4
		}

		board >>= 16
	}

	return float64(total)
}
