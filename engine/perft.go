package engine

import (
	"fmt"
	"time"
)

var leafNodes int64

func Perft(depth int, pos *BoardStruct) {
	if depth == 0 {
		leafNodes++
		return
	}

	var list MoveList

	GenerateAllMoves(pos, &list, true)

	for moveNum := 0; moveNum < list.Count; moveNum++ {

		if !pos.DoMove(list.Moves[moveNum]) {
			continue
		}
		Perft(depth-1, pos)
		pos.UndoMove()
	}
}

func PerftTest(depth int, pos *BoardStruct) {
	fmt.Printf("\nStarting Test To Depth:%d\n", depth)
	leafNodes = 0

	start := time.Now().UnixMilli()

	var list MoveList
	GenerateAllMoves(pos, &list, true)

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		move := list.Moves[moveNum]
		if !pos.DoMove(move) {
			continue
		}
		var cumnodes int64 = leafNodes
		Perft(depth-1, pos)
		pos.UndoMove()
		var oldnodes int64 = leafNodes - cumnodes
		fmt.Printf("move %d > %s : %d\n", moveNum+1, move.String(), oldnodes)
	}

	totalTimeInMs := time.Now().UnixMilli() - start
	totalTimeInS := float64(totalTimeInMs) / 1000

	fmt.Println()
	fmt.Println("===========================")
	fmt.Printf("Total time (ms) : %d\n", totalTimeInMs)
	fmt.Printf("Nodes searched : %d\n", leafNodes)
	fmt.Printf("Nodes/second : %d\n\n", int64(float64(leafNodes)/totalTimeInS))
}
