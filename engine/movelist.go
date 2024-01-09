package engine

import "fmt"

type MoveList struct {
	Moves [maxPositionMoves]Move
	Count int
}

func (list *MoveList) AddQuietMove(pos *BoardStruct, move int) {
	list.Moves[list.Count].Move = move

	if pos.SearchKillers[0][pos.Ply] == move {
		list.Moves[list.Count].Score = 900_000
	} else if pos.SearchKillers[1][pos.Ply] == move {
		list.Moves[list.Count].Score = 800_000
	} else {

		list.Moves[list.Count].Score = pos.SearchHistory[pos.Pieces[FROMSQ(move)]][TOSQ(move)]
	}

	list.Count++
}

func (list *MoveList) AddCaptureMove(pos *BoardStruct, move int) {
	list.Moves[list.Count].Move = move
	list.Moves[list.Count].Score = MvvLvaScores[CAPTURED(move)][pos.Pieces[FROMSQ(move)]] + 1_000_000
	list.Count++
}

func (list *MoveList) AddEnPassantMove(pos *BoardStruct, move int) {
	list.Moves[list.Count].Move = move
	list.Moves[list.Count].Score = 105 + 1_000_000
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
		list.AddCaptureMove(pos, MOVE(from, to, cap, possibleProms[0], 0))
		list.AddCaptureMove(pos, MOVE(from, to, cap, possibleProms[1], 0))
		list.AddCaptureMove(pos, MOVE(from, to, cap, possibleProms[2], 0))
		list.AddCaptureMove(pos, MOVE(from, to, cap, possibleProms[3], 0))
	} else {
		list.AddCaptureMove(pos, MOVE(from, to, cap, Empty, 0))
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
		list.AddQuietMove(pos, MOVE(from, to, Empty, possibleProms[0], 0))
		list.AddQuietMove(pos, MOVE(from, to, Empty, possibleProms[1], 0))
		list.AddQuietMove(pos, MOVE(from, to, Empty, possibleProms[2], 0))
		list.AddQuietMove(pos, MOVE(from, to, Empty, possibleProms[3], 0))
	} else {
		list.AddQuietMove(pos, MOVE(from, to, Empty, Empty, 0))
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
