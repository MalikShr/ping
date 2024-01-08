package engine

import "fmt"

var VictimScore = [13]int{0, 100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600}
var MvvLvaScores [13][13]int

func InitMvvLva() {
	for Attacker := wPawn; Attacker <= bKing; Attacker++ {
		for Victim := wPawn; Victim <= bKing; Victim++ {
			MvvLvaScores[Victim][Attacker] = VictimScore[Victim] + 6 - (VictimScore[Attacker] / 100)
		}
	}
}

func MoveExists(pos *BoardStruct, move int) bool {
	var list MoveList
	GenerateAllMoves(pos, &list, true)

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		if !pos.MakeMove(list.Moves[moveNum].Move) {
			continue
		}
		pos.TakeMove()
		if list.Moves[moveNum].Move == move {
			return true
		}
	}

	return false
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

func (list *MoveList) GenerateAllPawnCaptureMoves(sq int, pos *BoardStruct, side uint8) {
	direction := 1

	if side == Black {
		direction = -1
	}

	if !SQOFFBOARD(sq+(7*direction)) && PieceCol[pos.Pieces[sq+(7*direction)]] == pos.SideToMove^1 && ValidPawnDelta1(sq, direction) {
		list.AddPawnCapMove(pos, sq, sq+(7*direction), pos.Pieces[sq+(7*direction)], side)
	}

	if !SQOFFBOARD(sq+(9*direction)) && PieceCol[pos.Pieces[sq+(9*direction)]] == pos.SideToMove^1 && ValidPawnDelta2(sq, direction) {
		list.AddPawnCapMove(pos, sq, sq+(9*direction), pos.Pieces[sq+(9*direction)], side)
	}

	if pos.EnPas != NoSq {
		if sq+(7*direction) == pos.EnPas && ValidPawnDelta1(sq, direction) {
			move := MOVE(sq, sq+(7*direction), Empty, Empty, MFLAGEP)
			list.AddEnPassantMove(pos, move)
		}
		if sq+(9*direction) == pos.EnPas && ValidPawnDelta2(sq, direction) {
			move := MOVE(sq, sq+(9*direction), Empty, Empty, MFLAGEP)
			list.AddEnPassantMove(pos, move)
		}
	}
}

func GenerateAllMoves(pos *BoardStruct, list *MoveList, quit bool) {
	list.Count = 0

	if pos.SideToMove == White {
		for pieceNum := 0; pieceNum < pos.PieceNum[wPawn]; pieceNum++ {
			sq := pos.PieceList[wPawn][pieceNum]

			if pos.Pieces[sq+8] == Empty && quit {
				list.AddPawnMove(pos, sq, sq+8, White)

				if RankOf(sq) == R2 && pos.Pieces[sq+16] == Empty {
					move := MOVE(sq, (sq + 16), Empty, Empty, MFLAGPS)
					list.AddQuietMove(pos, move)
				}
			}

			list.GenerateAllPawnCaptureMoves(sq, pos, White)
		}

		if quit {
			// castling
			if pos.CastlePerm&wKCastle != 0 && pos.CastlePerm&wQCastle != 0 {
				if pos.Pieces[F1] == Empty && pos.Pieces[G1] == Empty {
					if !pos.SqAttacked(E1, Black) && !pos.SqAttacked(F1, Black) {
						list.AddQuietMove(pos, MOVE(E1, G1, Empty, Empty, MFLAGCA))
					}
				}

				if pos.Pieces[D1] == Empty && pos.Pieces[C1] == Empty && pos.Pieces[B1] == Empty {
					if !pos.SqAttacked(E1, Black) && !pos.SqAttacked(D1, Black) {
						list.AddQuietMove(pos, MOVE(E1, C1, Empty, Empty, MFLAGCA))
					}
				}
			}

		}

	} else {

		for pieceNum := 0; pieceNum < pos.PieceNum[bPawn]; pieceNum++ {
			sq := pos.PieceList[bPawn][pieceNum]

			if pos.Pieces[sq-8] == Empty && quit {
				list.AddPawnMove(pos, sq, sq-8, Black)
				if RankOf(sq) == R7 && pos.Pieces[sq-16] == Empty {
					list.AddQuietMove(pos, MOVE(sq, (sq-16), Empty, Empty, MFLAGPS))
				}
			}

			list.GenerateAllPawnCaptureMoves(sq, pos, Black)
		}

		if quit {
			// castling
			if pos.CastlePerm&bKCastle != 0 && pos.CastlePerm&bQCastle != 0 {
				if pos.Pieces[F8] == Empty && pos.Pieces[G8] == Empty {
					if !pos.SqAttacked(E8, White) && !pos.SqAttacked(F8, White) {
						list.AddQuietMove(pos, MOVE(E8, G8, Empty, Empty, MFLAGCA))
					}
				}

				if pos.Pieces[D8] == Empty && pos.Pieces[C8] == Empty && pos.Pieces[B8] == Empty {
					if !pos.SqAttacked(E8, White) && !pos.SqAttacked(D8, White) {
						list.AddQuietMove(pos, MOVE(E8, C8, Empty, Empty, MFLAGCA))
					}
				}
			}
		}
	}

	for sq := 0; sq < 64; sq++ {
		piece := pos.Pieces[sq]

		if piece == Empty || piece == wPawn || piece == bPawn {
			continue
		}

		if IsKn(piece) && PieceCol[piece] == pos.SideToMove {
			GenKnightMoves(sq, pos, list, quit)
		}

		if IsBQ(piece) && PieceCol[piece] == pos.SideToMove {
			GenBishopMoves(sq, pos, list, quit)
		}

		if IsRQ(piece) && PieceCol[piece] == pos.SideToMove {
			GenRookMoves(sq, pos, list, quit)
		}

		if IsKi(piece) && PieceCol[piece] == pos.SideToMove {
			GenKingMoves(sq, pos, list, quit)
		}
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
