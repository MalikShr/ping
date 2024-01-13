package engine

import "fmt"

type MoveList struct {
	Moves [maxPositionMoves]Move
	Count int
}

var VictimScore = [13]uint16{0, 10, 20, 30, 40, 50, 60, 10, 20, 30, 40, 50, 60}
var MvvLvaScores [13][13]uint16

func InitMvvLva() {
	for Attacker := wPawn; Attacker <= bKing; Attacker++ {
		for Victim := wPawn; Victim <= bKing; Victim++ {
			MvvLvaScores[Victim][Attacker] = VictimScore[Victim] + 6 - (VictimScore[Attacker] / 100)
		}
	}
}

func (list *MoveList) AddQuietMove(pos *BoardStruct, move int) {
	list.Moves[list.Count].Move = move

	if pos.SearchKillers[0][pos.Ply] == move {
		list.Moves[list.Count].Score = MvvLvaOffset - FirstKillerMoveScore
	} else if pos.SearchKillers[1][pos.Ply] == move {
		list.Moves[list.Count].Score = MvvLvaOffset - SecondKillerMoveScore
	} else {

		list.Moves[list.Count].Score = pos.SearchHistory[pos.Squares[FromSq(move)]][ToSq(move)]
	}

	list.Count++
}

func (list *MoveList) AddCaptureMove(pos *BoardStruct, move int) {
	list.Moves[list.Count].Move = move
	list.Moves[list.Count].Score = MvvLvaOffset + MvvLvaScores[Captured(move)][pos.Squares[FromSq(move)]]
	list.Count++
}

func (list *MoveList) AddEnPassantMove(pos *BoardStruct, move int) {
	list.Moves[list.Count].Move = move
	list.Moves[list.Count].Score = MvvLvaOffset + 15
	list.Count++
}

func (list *MoveList) AddPawnCapMove(pos *BoardStruct, from int, to int, cap uint8, side uint8) {
	possibleProms := [4]uint8{wQueen, wRook, wBishop, wKnight}
	beforePromRank := R7

	if side == Black {
		possibleProms = [4]uint8{bQueen, bRook, bBishop, bKnight}
		beforePromRank = R2
	}

	if RankOf(from) == beforePromRank {
		list.AddCaptureMove(pos, NewMove(from, to, cap, possibleProms[0], 0))
		list.AddCaptureMove(pos, NewMove(from, to, cap, possibleProms[1], 0))
		list.AddCaptureMove(pos, NewMove(from, to, cap, possibleProms[2], 0))
		list.AddCaptureMove(pos, NewMove(from, to, cap, possibleProms[3], 0))
	} else {
		list.AddCaptureMove(pos, NewMove(from, to, cap, Empty, 0))
	}
}

func (list *MoveList) AddPawnMove(pos *BoardStruct, from int, to int, side uint8) {
	possibleProms := [4]uint8{wQueen, wRook, wBishop, wKnight}
	beforePromRank := R7

	if side == Black {
		possibleProms = [4]uint8{bQueen, bRook, bBishop, bKnight}
		beforePromRank = R2
	}

	if RankOf(from) == beforePromRank {
		list.AddQuietMove(pos, NewMove(from, to, Empty, possibleProms[0], 0))
		list.AddQuietMove(pos, NewMove(from, to, Empty, possibleProms[1], 0))
		list.AddQuietMove(pos, NewMove(from, to, Empty, possibleProms[2], 0))
		list.AddQuietMove(pos, NewMove(from, to, Empty, possibleProms[3], 0))
	} else {
		list.AddQuietMove(pos, NewMove(from, to, Empty, Empty, 0))
	}
}

func (list *MoveList) String() string {
	moveListStr := "MoveList:\n"

	for i := 0; i < list.Count; i++ {
		move := list.Moves[i].Move
		score := list.Moves[i].Score

		moveListStr += fmt.Sprintf("Move:%d > %s (score:%d)\n", i+1, PrMove(move), score)
	}
	moveListStr += fmt.Sprintf("MoveList Total %d Moves:\n\n", list.Count)

	return moveListStr
}
