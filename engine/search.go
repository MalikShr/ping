package engine

import (
	"fmt"
	"time"
)

const (
	MaxDepth = 64
	INFINITE = 30000

	// A constant representing the score of the principal variation
	// move from the transposition table.
	PVMoveScore uint16 = 10_000

	// A constant representing the score offsets of the killer moves.
	FirstKillerMoveScore  uint16 = 1000
	SecondKillerMoveScore uint16 = 2000

	// A constant to offset the score of the pv and MVV-LVA move higher
	// than killers and history heuristic moves.
	MvvLvaOffset uint16 = 10_000
)

type Search struct {
	Starttime int64
	Stoptime  int64
	Depth     int
	Timeset   bool

	Nodes int64

	Stopped bool

	TT      TranspositionTable
	PvArray [MaxDepth]Move

	Fh  float32
	Fhf float32
}

func GetPvLine(depth int, pos *BoardStruct, search *Search) int {
	move := search.TT.Probe(pos.Hash).Best
	count := 0

	for move != NoMove && count < depth {
		if pos.MoveExists(move) {
			pos.DoMove(move)
			search.PvArray[count] = move
			count++
		} else {
			break
		}
		move = search.TT.Probe(pos.Hash).Best
	}

	for pos.Ply > 0 {
		pos.UndoMove()
	}

	return count
}

func PickNextMove(moveNum int, list *MoveList) {
	var tempMove Move
	bestScore := uint16(0)
	bestNum := moveNum

	for i := moveNum; i < list.Count; i++ {
		if list.Moves[i].Score() > bestScore {
			bestScore = list.Moves[i].Score()
			bestNum = i
		}
	}

	tempMove = list.Moves[moveNum]
	list.Moves[moveNum] = list.Moves[bestNum]
	list.Moves[bestNum] = tempMove
}

func Quiescence(alpha int, beta int, pos *BoardStruct, search *Search) int {
	if (search.Nodes & 2047) == 0 {
		checkUp(search)
	}

	search.Nodes++

	if isRepetition(pos) || pos.Rule50 >= 100 {
		return 0
	}

	if pos.Ply > MaxDepth-1 {
		return EvalPosition(pos)
	}

	score := EvalPosition(pos)

	if score >= beta {
		return beta
	}

	if score > alpha {
		alpha = score
	}

	var list MoveList
	GenerateAllMoves(pos, &list, false)

	legal := 0
	oldAlpha := alpha
	bestMove := NoMove
	score = -INFINITE

	for moveNum := 0; moveNum < list.Count; moveNum++ {

		PickNextMove(moveNum, &list)

		if !pos.DoMove(list.Moves[moveNum]) {
			continue
		}

		legal++
		score = -Quiescence(-beta, -alpha, pos, search)
		pos.UndoMove()

		if search.Stopped {
			return 0
		}

		if score > alpha {
			if score >= beta {
				if legal == 1 {
					search.Fhf++
				}
				search.Fh++

				return beta
			}
			alpha = score
			bestMove = list.Moves[moveNum]
		}
	}

	if alpha != oldAlpha {
		search.TT.Store(pos.Hash, bestMove)
	}

	return alpha
}

func AlphaBeta(alpha int, beta int, depth int, pos *BoardStruct, search *Search) int {
	if depth == 0 {
		return Quiescence(alpha, beta, pos, search)
	}

	if (search.Nodes & 2047) == 0 {
		checkUp(search)
	}

	search.Nodes++

	if (isRepetition(pos) || pos.Rule50 >= 100) && pos.Ply != 0 {
		return 0
	}

	if pos.Ply > MaxDepth-1 {
		return EvalPosition(pos)
	}

	kingBB := pos.Pieces[AllPieces[pos.SideToMove][King]]
	inCheck := SqAttacked(kingBB.Msb(), pos, pos.SideToMove^1)

	var list MoveList
	GenerateAllMoves(pos, &list, true)

	legal := 0
	oldAlpha := alpha
	bestMove := NoMove
	score := -INFINITE
	pvMove := search.TT.Probe(pos.Hash).Best

	if pvMove != NoMove {
		for moveNum := 0; moveNum < list.Count; moveNum++ {
			if list.Moves[moveNum].Equals(pvMove) {
				list.Moves[moveNum].AddScore(MvvLvaOffset + PVMoveScore)
			}

		}
	}

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		PickNextMove(moveNum, &list)

		if !pos.DoMove(list.Moves[moveNum]) {
			continue
		}

		legal++
		score = -AlphaBeta(-beta, -alpha, depth-1, pos, search)
		pos.UndoMove()

		if search.Stopped {
			return 0
		}

		if score > alpha {
			if score >= beta {
				if legal == 1 {
					search.Fhf++
				}
				search.Fh++

				if list.Moves[moveNum].MoveType() == Attack {
					pos.SearchKillers[1][pos.Ply] = pos.SearchKillers[0][pos.Ply]
					pos.SearchKillers[0][pos.Ply] = list.Moves[moveNum]
				}

				return beta
			}
			alpha = score
			bestMove = list.Moves[moveNum]

			if list.Moves[moveNum].MoveType() == Attack {
				pos.SearchHistory[pos.Squares[bestMove.FromSq()]][bestMove.ToSq()] += uint16(depth)
			}
		}
	}

	if legal == 0 {
		if inCheck {
			return -INFINITE + pos.Ply
		} else {
			return 0
		}
	}

	if alpha != oldAlpha {
		search.TT.Store(pos.Hash, bestMove)
	}

	return alpha
}

func SearchPosition(pos *BoardStruct, search *Search) {
	bestMove := NoMove
	bestScore := -INFINITE
	pvMoves := 0

	clearForSearch(pos, search)

	for currentDepth := 1; currentDepth <= search.Depth; currentDepth++ {
		bestScore = AlphaBeta(-INFINITE, INFINITE, currentDepth, pos, search)

		if search.Stopped {
			break
		}

		pvMoves = GetPvLine(currentDepth, pos, search)
		bestMove = search.PvArray[0]

		fmt.Printf("\ninfo score cp %d depth %d nodes %d time %d pv",
			bestScore, currentDepth, search.Nodes, time.Now().UnixMilli()-int64(search.Starttime))

		for pvNum := 0; pvNum < pvMoves; pvNum++ {
			fmt.Printf(" %s", search.PvArray[pvNum].String())
		}
		fmt.Println()
		//fmt.Printf(" Ordering: %f\n", info.Fhf/info.Fh)
	}

	fmt.Printf("bestmove %s\n", bestMove.String())
}

func clearForSearch(pos *BoardStruct, search *Search) {
	for index := 0; index < 13; index++ {
		for index2 := 0; index2 < 64; index2++ {
			pos.SearchHistory[index][index2] = 0
		}
	}

	for index := 0; index < 2; index++ {
		for index2 := 0; index2 < MaxDepth; index2++ {
			pos.SearchKillers[index][index2] = 0
		}
	}

	pos.Ply = 0

	search.Stopped = false
	search.Nodes = 0
	search.Fh = 0
	search.Fhf = 0
}

func isRepetition(pos *BoardStruct) bool {
	for i := pos.HistoryPly - pos.Rule50; i < pos.HistoryPly-1; i++ {
		if pos.Hash == pos.History[i].Hash {
			return true
		}
	}

	return false
}

func checkUp(search *Search) {
	if search.Timeset && time.Now().UnixMilli() > int64(search.Stoptime) {
		search.Stopped = true
	}
}
