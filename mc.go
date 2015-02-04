package main

import (
	// "container/heap"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const MAX_INT = int(^uint(0) >> 1)
const SMOOTHING_FACTOR = 0.3

var topTile uint64
var maxTable map[uint64]Entry
var expectationTable map[uint64]Entry
var removeQueue *PosHeap

func init() {
	maxTable = make(map[uint64]Entry)
	expectationTable = make(map[uint64]Entry)
	topTile = 0
	// removeQueue = &PosHeap{}

	// heap.Init(removeQueue)
}

type Entry struct {
	visitations int
	score       float64
}

func MonteCarloChooseMove(board uint64) Move {
	fmt.Println("----")

	start := time.Now()

	newTop := extractMax(board)

	if newTop > topTile {
		topTile = newTop
		maxTable = make(map[uint64]Entry)
		expectationTable = make(map[uint64]Entry)
	}

	duration := time.Now().Sub(start)
	fmt.Println("Deletions took:", duration)

	for i := 0; i < 1000; i++ {
		maxStep(board, maxTable, expectationTable)
	}

	best := NONE
	bScore := 0.0

	for _, move := range Moves {
		moved := MakeMove(board, move)

		if moved == board {
			continue
		}

		entry, ok := expectationTable[moved]

		if !ok {
			panic("!!")
		}

		realScore := entry.score / float64(entry.visitations)

		if realScore > bScore {
			bScore = realScore
			best = move
		}
	}

	fmt.Println(bScore, best)

	duration = time.Now().Sub(start)

	fmt.Println("Move took: ", duration)

	return best
}

func extractMax(board uint64) uint64 {
	max := uint64(0)

	for i := 0; i < 16; i++ {
		if (board & 0xf) > max {
			max = board & 0xf
		}

		board >>= 4
	}

	return max
}

func maxStep(board uint64, maxTable map[uint64]Entry, expectationTable map[uint64]Entry) int {
	if Dead(board) {
		return 0
	}

	// Random start to avoid bias
	start := rand.Int31n(4)

	best := 0.0
	bPos := board

	for i := start; i < (4 + start); i++ {
		move := Moves[(i & 0x3)]
		moved := MakeMove(board, move)

		if moved == board {
			continue
		}

		if entry, ok := expectationTable[moved]; ok {
			// Value in [0, 1)
			score := math.Exp(-float64((entry.visitations * entry.visitations)) / entry.score)

			if score > best {
				best = score
				bPos = moved
			}
		} else {
			best = 1
			bPos = moved
		}
	}

	val := expectationStep(bPos, maxTable, expectationTable) + 1

	entry, ok := maxTable[board]

	if !ok {
		entry = Entry{visitations: 1, score: float64(val)}

		// heap.Push(removeQueue, board)
	} else {
		entry.visitations++
		entry.score = (float64(val) * SMOOTHING_FACTOR) + (float64(entry.score) * (1 - SMOOTHING_FACTOR))
	}

	maxTable[board] = entry

	return val
}

func expectationStep(board uint64, maxTable map[uint64]Entry, expectationTable map[uint64]Entry) int {
	min := MAX_INT
	minPos := uint64(0)
	mask := uint64(0xf)

	// Random tile and random starting position to avoid bias.
	tile := uint64(0x1)

	if rand.Float64() < 0.1 {
		tile = uint64(0x2)
	}

	start := rand.Int31n(16)

	for i := start; i < (start + 16); i++ {
		shift := uint((i & 0xf) << 2)

		_mask := mask << shift
		_tile := tile << shift

		if (board & _mask) == 0 {
			placed := board ^ _tile

			if entry, ok := maxTable[placed]; ok {
				if entry.visitations < min {
					min = entry.visitations
					minPos = placed
				}
			} else {
				min = 0
				minPos = placed
			}
		}
	}

	val := maxStep(minPos, maxTable, expectationTable) + 1

	entry, ok := expectationTable[board]

	if !ok {
		entry = Entry{visitations: 1, score: float64(val)}

		// heap.Push(removeQueue, board)
	} else {
		entry.visitations++
		entry.score += float64(val)
	}

	expectationTable[board] = entry

	return val
}

type PosHeap []uint64

func (p PosHeap) Len() int {
	return len(p)
}

func (p PosHeap) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PosHeap) Less(i, j int) bool {
	a := p[i]
	b := p[j]

	return value(a) < value(b)
}

func value(board uint64) uint64 {
	mask := uint64(0xf)

	total := uint64(0)
	for i := 0; i < 16; i++ {
		shift := board & mask

		total += 1 << (shift - 1)

		board >>= 4
	}

	return total
}

func (p *PosHeap) Push(x interface{}) {
	*p = append(*p, x.(uint64))
}

func (p *PosHeap) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}
