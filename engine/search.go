package engine

import (
	"fmt"
	"math"
	"time"
)

const (
	MaxDepth = 64
	INFINITE = 30000

	// A constant representing the score of the principal variation
	// move from the transposition table.
	PVMoveScore uint16 = 120

	// A constant representing the score offsets of the killer moves.
	FirstKillerMoveScore  uint16 = 10
	SecondKillerMoveScore uint16 = 20

	// A constant to offset the score of the pv and MVV-LVA move higher
	// than killers and history heuristic moves.
	MvvLvaOffset uint16 = math.MaxUint16 - 500
)

type Search struct {
	Starttime int64
	Stoptime  int64
	Depth     int
	Timeset   bool

	Nodes int64

	Stopped bool

	TT      map[uint64]int
	PvArray [MaxDepth]int

	Fh  float32
	Fhf float32
}

func GetPvLine(depth int, pos *BoardStruct, search *Search) int {
	move := search.TT[pos.Hash]
	count := 0

	for move != NoMove && count < depth {

		if pos.MoveExists(move) {
			pos.DoMove(move)
			search.PvArray[count] = move
			count++
		} else {
			break
		}
		move = search.TT[pos.Hash]
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
		if list.Moves[i].Score > bestScore {
			bestScore = list.Moves[i].Score
			bestNum = i
		}
	}

	tempMove = list.Moves[moveNum]
	list.Moves[moveNum] = list.Moves[bestNum]
	list.Moves[bestNum] = tempMove
}

func Quiescence(alpha int, beta int, pos *BoardStruct, info *Search) int {
	if (info.Nodes & 2047) == 0 {
		checkUp(info)
	}

	info.Nodes++

	if isRepetition(pos) || pos.Rule50 >= 100 {
		return 0
	}

	if pos.Ply > MaxDepth-1 {
		return EvalPosition(pos)
	}

	Score := EvalPosition(pos)

	if Score >= beta {
		return beta
	}

	if Score > alpha {
		alpha = Score
	}

	var list MoveList
	GenerateAllMoves(pos, &list, false)

	Legal := 0
	OldAlpha := alpha
	BestMove := NoMove
	Score = -INFINITE

	for MoveNum := 0; MoveNum < list.Count; MoveNum++ {

		PickNextMove(MoveNum, &list)

		if !pos.DoMove(list.Moves[MoveNum].Move) {
			continue
		}

		Legal++
		Score = -Quiescence(-beta, -alpha, pos, info)
		pos.UndoMove()

		if info.Stopped {
			return 0
		}

		if Score > alpha {
			if Score >= beta {
				if Legal == 1 {
					info.Fhf++
				}
				info.Fh++

				return beta
			}
			alpha = Score
			BestMove = list.Moves[MoveNum].Move
		}
	}

	if alpha != OldAlpha {
		info.TT[pos.Hash] = BestMove
	}

	return alpha
}

func AlphaBeta(alpha int, beta int, depth int, pos *BoardStruct, info *Search) int {
	if depth == 0 {
		return Quiescence(alpha, beta, pos, info)
	}

	if (info.Nodes & 2047) == 0 {
		checkUp(info)
	}

	info.Nodes++

	if (isRepetition(pos) || pos.Rule50 >= 100) && pos.Ply != 0 {
		return 0
	}

	if pos.Ply > MaxDepth-1 {
		return EvalPosition(pos)
	}

	inCheck := SqAttacked(pos.KingSq[pos.SideToMove], pos, pos.SideToMove^1)

	var list MoveList
	GenerateAllMoves(pos, &list, true)

	Legal := 0
	OldAlpha := alpha
	BestMove := NoMove
	Score := -INFINITE
	PvMove := info.TT[pos.Hash]

	if PvMove != NoMove {
		for MoveNum := 0; MoveNum < list.Count; MoveNum++ {
			if list.Moves[MoveNum].Move == PvMove {
				list.Moves[MoveNum].Score = MvvLvaOffset + PVMoveScore
			}

		}
	}

	for MoveNum := 0; MoveNum < list.Count; MoveNum++ {
		PickNextMove(MoveNum, &list)

		if !pos.DoMove(list.Moves[MoveNum].Move) {
			continue
		}

		Legal++
		Score = -AlphaBeta(-beta, -alpha, depth-1, pos, info)
		pos.UndoMove()

		if info.Stopped {
			return 0
		}

		if Score > alpha {
			if Score >= beta {
				if Legal == 1 {
					info.Fhf++
				}
				info.Fh++

				if list.Moves[MoveNum].Move&MFLAGCAP != 0 {
					pos.SearchKillers[1][pos.Ply] = pos.SearchKillers[0][pos.Ply]
					pos.SearchKillers[0][pos.Ply] = list.Moves[MoveNum].Move
				}

				return beta
			}
			alpha = Score
			BestMove = list.Moves[MoveNum].Move

			if list.Moves[MoveNum].Move&MFLAGCAP != 0 {
				pos.SearchHistory[pos.Squares[FromSq(BestMove)]][ToSq(BestMove)] += uint16(depth)
			}
		}
	}

	if Legal == 0 {
		if inCheck {
			return -INFINITE + pos.Ply
		} else {
			return 0
		}
	}

	if alpha != OldAlpha {
		info.TT[pos.Hash] = BestMove
	}

	return alpha
}

func SearchPosition(pos *BoardStruct, info *Search) {
	bestMove := NoMove
	bestScore := -INFINITE
	pvMoves := 0

	clearForSearch(pos, info)

	for currentDepth := 1; currentDepth <= info.Depth; currentDepth++ {
		bestScore = AlphaBeta(-INFINITE, INFINITE, currentDepth, pos, info)

		if info.Stopped {
			break
		}

		pvMoves = GetPvLine(currentDepth, pos, info)
		bestMove = info.PvArray[0]

		fmt.Printf("\ninfo score cp %d depth %d nodes %d time %d pv",
			bestScore, currentDepth, info.Nodes, time.Now().UnixMilli()-int64(info.Starttime))

		for pvNum := 0; pvNum < pvMoves; pvNum++ {
			fmt.Printf(" %s", PrMove(info.PvArray[pvNum]))
		}
		fmt.Println()
		fmt.Printf(" Ordering: %f\n", info.Fhf/info.Fh)
	}

	fmt.Printf("bestmove %s\n", PrMove(bestMove))
}

func clearForSearch(pos *BoardStruct, info *Search) {
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

	info.TT = make(map[uint64]int)
	pos.Ply = 0

	info.Stopped = false
	info.Nodes = 0
	info.Fh = 0
	info.Fhf = 0
}

func isRepetition(pos *BoardStruct) bool {
	for i := pos.HistoryPly - pos.Rule50; i < pos.HistoryPly-1; i++ {
		if pos.Hash == pos.History[i].Hash {
			return true
		}
	}

	return false
}

func checkUp(info *Search) {
	if info.Timeset && time.Now().UnixMilli() > int64(info.Stoptime) {
		info.Stopped = true
	}
}
