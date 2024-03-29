package engine

import "fmt"

const (
	maxPositionMoves = 256
)

type MoveList struct {
	Moves [maxPositionMoves]Move
	Count int
}

var VictimScore = [13]uint16{0, 100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600}
var MvvLvaScores [13][13]uint16

func InitMvvLva() {
	for Attacker := wPawn; Attacker <= bKing; Attacker++ {
		for Victim := wPawn; Victim <= bKing; Victim++ {
			MvvLvaScores[Victim][Attacker] = VictimScore[Victim] + 6 - (VictimScore[Attacker] / 100)
		}
	}
}

func (list *MoveList) AddQuietMove(pos *BoardStruct, move Move) {
	list.Moves[list.Count] = move

	if pos.SearchKillers[0][pos.Ply] == move {
		list.Moves[list.Count].AddScore(MvvLvaOffset - FirstKillerMoveScore)
	} else if pos.SearchKillers[1][pos.Ply] == move {
		list.Moves[list.Count].AddScore(MvvLvaOffset - SecondKillerMoveScore)
	} else {
		list.Moves[list.Count].AddScore(pos.SearchHistory[pos.Squares[move.FromSq()]][move.ToSq()])
	}

	list.Count++
}

func (list *MoveList) AddCaptureMove(pos *BoardStruct, move Move) {
	captured := pos.Squares[move.ToSq()]
	attacker := pos.Squares[move.FromSq()]

	list.Moves[list.Count] = move
	list.Moves[list.Count].AddScore(MvvLvaOffset + MvvLvaScores[captured][attacker])
	list.Count++
}

func (list *MoveList) AddEnPassantMove(pos *BoardStruct, move Move) {
	list.Moves[list.Count] = move
	list.Moves[list.Count].AddScore(MvvLvaOffset + 15)
	list.Count++
}

func (list *MoveList) AddPawnCapMove(pos *BoardStruct, from int, to int, cap uint8, side uint8) {
	beforePromRank := R7

	if side == Black {
		beforePromRank = R2
	}

	if RankOf(from) == beforePromRank {
		list.AddCaptureMove(pos, NewMove(from, to, Promotion, KnightPromotion))
		list.AddCaptureMove(pos, NewMove(from, to, Promotion, BishopPromotion))
		list.AddCaptureMove(pos, NewMove(from, to, Promotion, RookPromotion))
		list.AddCaptureMove(pos, NewMove(from, to, Promotion, QueenPromotion))
	} else {
		list.AddCaptureMove(pos, NewMove(from, to, Quiet, NoFlag))
	}
}

func (list *MoveList) AddPawnMove(pos *BoardStruct, from int, to int, side uint8) {
	beforePromRank := R7

	if side == Black {
		beforePromRank = R2
	}

	if RankOf(from) == beforePromRank {
		list.AddQuietMove(pos, NewMove(from, to, Promotion, KnightPromotion))
		list.AddQuietMove(pos, NewMove(from, to, Promotion, BishopPromotion))
		list.AddQuietMove(pos, NewMove(from, to, Promotion, RookPromotion))
		list.AddQuietMove(pos, NewMove(from, to, Promotion, QueenPromotion))
	} else {
		list.AddQuietMove(pos, NewMove(from, to, Quiet, NoFlag))
	}
}

func (list *MoveList) PickNextMove(moveIndex int) {
	var tempMove Move
	bestScore := uint16(0)
	bestNum := moveIndex

	for i := moveIndex; i < list.Count; i++ {
		if list.Moves[i].Score() > bestScore {
			bestScore = list.Moves[i].Score()
			bestNum = i
		}
	}

	tempMove = list.Moves[moveIndex]
	list.Moves[moveIndex] = list.Moves[bestNum]
	list.Moves[bestNum] = tempMove
}

func (list *MoveList) String() string {
	moveListStr := "MoveList:\n"

	for i := 0; i < list.Count; i++ {
		move := list.Moves[i]
		score := move.Score()

		moveListStr += fmt.Sprintf("Move:%d > %s (score:%d)\n", i+1, move.String(), score)
	}
	moveListStr += fmt.Sprintf("MoveList Total %d Moves:\n\n", list.Count)

	return moveListStr
}
