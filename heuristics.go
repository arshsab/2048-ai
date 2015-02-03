package main

import (
	"math"
	"math/rand"
)

// MetaHeuristics

const INITIAL_TEMPERATURE int = 10000
const COOLING_RATE float64 = 0.03

type Weights struct {
	spaces     float64
	stability  float64
	uniformity float64
}

var DEFAULT_WEIGHTS Weights = Weights{spaces: 0.5, stability: 0.5, uniformity: 0.5}

func Anneal(eval func(Weights) float64) Weights {
	// keep spaces constant because the ratios are the only thing that matters
	current := Weights{spaces: 0.5, stability: rand.Float64(), uniformity: rand.Float64()}
	best := current

	currentScore := eval(current)
	bestScore := currentScore

	temp := float64(INITIAL_TEMPERATURE)

	for temp > 1.0 {
		newPosition := Weights{
			spaces:     0.5,
			stability:  current.stability + ((rand.Float64() / 10) - .05),
			uniformity: current.uniformity + ((rand.Float64() / 10) - .05),
		}

		newScore := eval(newPosition)

		if newScore > currentScore || math.Exp(((currentScore-newScore)/temp)) < rand.Float64() {
			current = newPosition
			currentScore = newScore
		}

		if currentScore > bestScore {
			best = newPosition
			bestScore = currentScore
		}

		temp *= 1 - COOLING_RATE
	}

	return best
}

// Game time heuristics:

var MAX_SCORE int = 0

func Heuristic(board uint64, weights Weights) int {
	count := 0

	for {
		score := -10000.0
		moved := board

		for _, move := range Moves {
			potential := MakeMove(board, move)

			if potential == board {
				continue
			}

			potentialScore := weighted(potential, weights)

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

	if count > MAX_SCORE {
		MAX_SCORE = count
	}

	return count
}

func weighted(board uint64, weights Weights) (score float64) {
	score = spaces(board) * weights.spaces * maxTile(board)
	score += stability(board) * weights.stability
	score -= (nonuniformity(board) + nonuniformity(Transpose(board))) * weights.uniformity

	return
}

func maxTile(board uint64) float64 {
	max := board & 0xf

	for board != 0 {
		board >>= 1

		pot := board & 0xf

		if pot > max {
			max = pot
		}
	}

	return float64(max)
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
