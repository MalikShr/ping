package engine

import (
	"math"
)

var CastlePerm = [64]int{
	13, 15, 15, 15, 12, 15, 15, 14,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	7, 15, 15, 15, 3, 15, 15, 11,
}

func SQOFFBOARD(sq int) bool {
	return sq > 63 || sq < 0
}

func ValidPawnDelta1(sq int, direction int) bool {
	return math.Abs(float64(FileOf(sq)-FileOf(sq+(7*direction)))) == 1
}

func ValidPawnDelta2(sq int, direction int) bool {
	return math.Abs(float64(FileOf(sq)-FileOf(sq+(9*direction)))) == 1
}

func GenerateAllMoves(pos *BoardStruct, list *MoveList, genQuiet bool) {
	list.Count = 0

	if genQuiet {
		if pos.SideToMove == White {
			// castling
			if pos.CastlePerm&wKCastle != 0 && pos.CastlePerm&wQCastle != 0 {
				if pos.Squares[F1] == Empty && pos.Squares[G1] == Empty {
					if !SqAttacked(E1, pos, Black) && !SqAttacked(F1, pos, Black) {
						list.AddQuietMove(pos, MOVE(E1, G1, Empty, Empty, MFLAGCA))
					}
				}

				if pos.Squares[D1] == Empty && pos.Squares[C1] == Empty && pos.Squares[B1] == Empty {
					if !SqAttacked(E1, pos, Black) && !SqAttacked(D1, pos, Black) {
						list.AddQuietMove(pos, MOVE(E1, C1, Empty, Empty, MFLAGCA))
					}
				}
			}
		} else {
			// castling
			if pos.CastlePerm&bKCastle != 0 && pos.CastlePerm&bQCastle != 0 {
				if pos.Squares[F8] == Empty && pos.Squares[G8] == Empty {
					if !SqAttacked(E8, pos, White) && !SqAttacked(F8, pos, White) {
						list.AddQuietMove(pos, MOVE(E8, G8, Empty, Empty, MFLAGCA))
					}
				}

				if pos.Squares[D8] == Empty && pos.Squares[C8] == Empty && pos.Squares[B8] == Empty {
					if !SqAttacked(E8, pos, White) && !SqAttacked(D8, pos, White) {
						list.AddQuietMove(pos, MOVE(E8, C8, Empty, Empty, MFLAGCA))
					}
				}
			}
		}
	}

	// Copy bitboard of all pieces and loop over them
	allPieces := pos.Sides[pos.SideToMove]

	for allPieces != 0 {
		sq := allPieces.PopBit()
		piece := pos.Squares[sq]

		if PiecePawn[piece] {
			GenPawnMoves(sq, pos, list, genQuiet)
		}

		if IsKn(piece) {
			GenKnightAttacks(sq, pos, list, genQuiet)
		}

		if IsBQ(piece) {
			GenBishopAttacks(sq, pos, list, genQuiet)
		}

		if IsRQ(piece) {
			GenRookAttacks(sq, pos, list, genQuiet)
		}

		if IsKi(piece) {
			GenKingAttacks(sq, pos, list, genQuiet)
		}
	}
}

func GenPawnMoves(sq int, pos *BoardStruct, list *MoveList, genQuiet bool) {
	attacks := PawnAttacks[pos.SideToMove^1][sq] & pos.Sides[pos.SideToMove^1]
	for attacks != 0 {
		targetSq := attacks.PopBit()

		list.AddPawnCapMove(pos, sq, targetSq, pos.Squares[targetSq], pos.SideToMove)
	}

	if pos.EnPas != NoSq {
		enPasAttacks := PawnAttacks[pos.SideToMove^1][sq] & (1 << pos.EnPas)

		if enPasAttacks != 0 {
			targetSq := enPasAttacks.PopBit()

			move := MOVE(sq, targetSq, Empty, Empty, MFLAGEP)
			list.AddEnPassantMove(pos, move)
		}
	}

	if genQuiet {
		dir := 1
		fistPawnRank := R2

		if pos.SideToMove == Black {
			dir = -1
			fistPawnRank = R7
		}

		onePawnPush := sq + (8 * dir)
		twoPawnPush := sq + (16 * dir)

		if !SQOFFBOARD(onePawnPush) && pos.Sides[Both]&(1<<onePawnPush) == 0 {
			list.AddPawnMove(pos, sq, onePawnPush, pos.SideToMove)

			if RankOf(sq) == fistPawnRank && pos.Sides[Both]&(1<<twoPawnPush) == 0 {
				move := MOVE(sq, twoPawnPush, Empty, Empty, MFLAGPS)
				list.AddQuietMove(pos, move)
			}
		}
	}

}

func GenerateAllPawnCaptureMoves(sq int, pos *BoardStruct, list *MoveList, side uint8) {
	direction := 1

	if side == Black {
		direction = -1
	}

	if !SQOFFBOARD(sq+(7*direction)) && PieceCol[pos.Squares[sq+(7*direction)]] == pos.SideToMove^1 && ValidPawnDelta1(sq, direction) {
		list.AddPawnCapMove(pos, sq, sq+(7*direction), pos.Squares[sq+(7*direction)], side)
	}

	if !SQOFFBOARD(sq+(9*direction)) && PieceCol[pos.Squares[sq+(9*direction)]] == pos.SideToMove^1 && ValidPawnDelta2(sq, direction) {
		list.AddPawnCapMove(pos, sq, sq+(9*direction), pos.Squares[sq+(9*direction)], side)
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

func GenKnightAttacks(sq int, pos *BoardStruct, list *MoveList, genQuiet bool) {
	attacks := KnightAttacks[sq] & ^pos.Sides[pos.SideToMove]

	for attacks != 0 {
		targetSq := attacks.PopBit()

		if pos.Sides[pos.SideToMove^1]&(1<<targetSq) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Squares[targetSq], Empty, 0))
			continue
		}
		if genQuiet {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}
}

func GenBishopAttacks(sq int, pos *BoardStruct, list *MoveList, genQuiet bool) {
	attacks := GetBishopAttacks(sq, pos.Sides[Both]) & ^pos.Sides[pos.SideToMove]

	for attacks != 0 {
		targetSq := attacks.PopBit()

		if pos.Sides[pos.SideToMove^1]&(1<<targetSq) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Squares[targetSq], Empty, 0))
			continue
		}
		if genQuiet {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}
}

func GenRookAttacks(sq int, pos *BoardStruct, list *MoveList, genQuiet bool) {
	attacks := GetRookAttacks(sq, pos.Sides[Both]) & ^pos.Sides[pos.SideToMove]

	for attacks != 0 {
		targetSq := attacks.PopBit()

		if pos.Sides[pos.SideToMove^1]&(1<<targetSq) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Squares[targetSq], Empty, 0))
			continue
		}
		if genQuiet {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}
}

func GenKingAttacks(square int, pos *BoardStruct, list *MoveList, genQuiet bool) {
	attacks := KingAttacks[square] & ^pos.Sides[pos.SideToMove]

	for attacks != 0 {
		sq := attacks.PopBit()

		if pos.Sides[pos.SideToMove^1]&(1<<sq) != 0 {
			list.AddCaptureMove(pos, MOVE(square, sq, pos.Squares[sq], Empty, 0))
			continue
		}
		if genQuiet {
			list.AddQuietMove(pos, MOVE(square, sq, Empty, Empty, 0))
		}
	}
}

func SqAttacked(targetSq int, pos *BoardStruct, side uint8) bool {
	if side == White {
		if PawnAttacks[side][targetSq]&pos.Pieces[wPawn] != 0 {
			return true
		}

		if KnightAttacks[targetSq]&pos.Pieces[wKnight] != 0 {
			return true
		}

		if GetBishopAttacks(targetSq, pos.Sides[Both])&pos.Pieces[wBishop] != 0 || GetBishopAttacks(targetSq, pos.Sides[Both])&pos.Pieces[wQueen] != 0 {
			return true
		}

		if GetRookAttacks(targetSq, pos.Sides[Both])&pos.Pieces[wRook] != 0 || GetRookAttacks(targetSq, pos.Sides[Both])&pos.Pieces[wQueen] != 0 {
			return true
		}

		if KingAttacks[targetSq]&pos.Pieces[wKing] != 0 {
			return true
		}
	} else {
		if PawnAttacks[side][targetSq]&pos.Pieces[bPawn] != 0 {
			return true
		}

		if KnightAttacks[targetSq]&pos.Pieces[bKnight] != 0 {
			return true
		}

		if GetBishopAttacks(targetSq, pos.Sides[Both])&pos.Pieces[bBishop] != 0 || GetBishopAttacks(targetSq, pos.Sides[Both])&pos.Pieces[bQueen] != 0 {

			return true
		}

		if GetRookAttacks(targetSq, pos.Sides[Both])&pos.Pieces[bRook] != 0 || GetRookAttacks(targetSq, pos.Sides[Both])&pos.Pieces[bQueen] != 0 {
			return true
		}

		if KingAttacks[targetSq]&pos.Pieces[bKing] != 0 {
			return true
		}
	}

	return false
}
