package engine

const (
	// These masks help determine whether or not the squares between
	// the king and it's rooks are clear for castling
	F1_G1, B1_C1_D1 = 0x60, 0xe
	F8_G8, B8_C8_D8 = 0x6000000000000000, 0xe00000000000000
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

func GenerateAllMoves(pos *BoardStruct, list *MoveList, genQuiet bool) {
	// Copy bitboard of all pieces and loop over them
	allPieces := pos.Sides[pos.SideToMove]

	for allPieces != 0 {
		sq := allPieces.PopBit()
		piece := pos.Squares[sq]

		if PiecePawn[piece] {
			genPawnMoves(sq, pos, list, genQuiet)
		}

		if IsKn(piece) {
			knighAttacks := KnightAttacks[sq] & ^pos.Sides[pos.SideToMove]
			genMovesFromBB(sq, knighAttacks, pos, list, genQuiet)
		}

		if IsBQ(piece) {
			bishopAttacks := genBishopMoves(sq, pos.Sides[Both]) & ^pos.Sides[pos.SideToMove]
			genMovesFromBB(sq, bishopAttacks, pos, list, genQuiet)
		}

		if IsRQ(piece) {
			rookAttacks := genRookMoves(sq, pos.Sides[Both]) & ^pos.Sides[pos.SideToMove]
			genMovesFromBB(sq, rookAttacks, pos, list, genQuiet)
		}

		if IsKi(piece) {
			kingAttacks := KingAttacks[sq] & ^pos.Sides[pos.SideToMove]
			genMovesFromBB(sq, kingAttacks, pos, list, genQuiet)
		}
	}

	if genQuiet {
		genCastlingMoves(pos, list)
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

		if genBishopMoves(targetSq, pos.Sides[Both])&pos.Pieces[wBishop] != 0 || genBishopMoves(targetSq, pos.Sides[Both])&pos.Pieces[wQueen] != 0 {
			return true
		}

		if genRookMoves(targetSq, pos.Sides[Both])&pos.Pieces[wRook] != 0 || genRookMoves(targetSq, pos.Sides[Both])&pos.Pieces[wQueen] != 0 {
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

		if genBishopMoves(targetSq, pos.Sides[Both])&pos.Pieces[bBishop] != 0 || genBishopMoves(targetSq, pos.Sides[Both])&pos.Pieces[bQueen] != 0 {

			return true
		}

		if genRookMoves(targetSq, pos.Sides[Both])&pos.Pieces[bRook] != 0 || genRookMoves(targetSq, pos.Sides[Both])&pos.Pieces[bQueen] != 0 {
			return true
		}

		if KingAttacks[targetSq]&pos.Pieces[bKing] != 0 {
			return true
		}
	}

	return false
}

func genPawnMoves(sq int, pos *BoardStruct, list *MoveList, genQuiet bool) {
	attacks := PawnAttacks[pos.SideToMove^1][sq] & pos.Sides[pos.SideToMove^1]
	for attacks != 0 {
		targetSq := attacks.PopBit()

		list.AddPawnCapMove(pos, sq, targetSq, pos.Squares[targetSq], pos.SideToMove)
	}

	if pos.EnPas != NoSq {
		enPasAttacks := PawnAttacks[pos.SideToMove^1][sq] & (1 << pos.EnPas)

		if enPasAttacks != 0 {
			targetSq := enPasAttacks.PopBit()

			move := NewMove(sq, targetSq, Attack, AttackEP)
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
				move := NewMove(sq, twoPawnPush, Quiet, NoFlag)
				list.AddQuietMove(pos, move)
			}
		}
	}

}

func genCastlingMoves(pos *BoardStruct, list *MoveList) {
	if pos.SideToMove == White {
		if pos.CastlePerm&wKCastle != 0 && (pos.Sides[Both]&F1_G1) == 0 && !SqAttacked(E1, pos, Black) {
			if !SqAttacked(F1, pos, Black) && !SqAttacked(G1, pos, Black) {
				list.AddQuietMove(pos, NewMove(E1, G1, Castle, NoFlag))
			}
		}

		if pos.CastlePerm&wQCastle != 0 && (pos.Sides[Both]&B1_C1_D1) == 0 && !SqAttacked(E1, pos, Black) {
			if !SqAttacked(C1, pos, Black) && !SqAttacked(D1, pos, Black) {
				list.AddQuietMove(pos, NewMove(E1, C1, Castle, NoFlag))
			}
		}
	} else {
		if pos.CastlePerm&bKCastle != 0 && (pos.Sides[Both]&F8_G8) == 0 && !SqAttacked(E8, pos, White) {
			if !SqAttacked(F8, pos, White) && !SqAttacked(G8, pos, White) {
				list.AddQuietMove(pos, NewMove(E8, G8, Castle, NoFlag))
			}
		}

		if pos.CastlePerm&bQCastle != 0 && (pos.Sides[Both]&B8_C8_D8) == 0 && !SqAttacked(E8, pos, White) {
			if !SqAttacked(C8, pos, White) && !SqAttacked(D8, pos, White) {
				list.AddQuietMove(pos, NewMove(E8, C8, Castle, NoFlag))
			}
		}
	}

}

func genMovesFromBB(sq int, attacks Bitboard, pos *BoardStruct, list *MoveList, genQuiet bool) {
	for attacks != 0 {
		targetSq := attacks.PopBit()

		if pos.Sides[pos.SideToMove^1]&(1<<targetSq) != 0 {
			list.AddCaptureMove(pos, NewMove(sq, targetSq, Attack, NoFlag))
			continue
		}
		if genQuiet {
			list.AddQuietMove(pos, NewMove(sq, targetSq, Quiet, NoFlag))
		}
	}
}

func genBishopMoves(sq int, occupancy Bitboard) Bitboard {
	occupancy &= BishopMasks[sq]
	occupancy *= MagicB[sq]
	occupancy >>= 64 - RelevantBishopBits[sq]

	return BishopAttacks[sq][occupancy]
}

func genRookMoves(sq int, occupancy Bitboard) Bitboard {
	occupancy &= RookMasks[sq]
	occupancy *= MagicR[sq]
	occupancy >>= 64 - RelevantRookBits[sq]

	return RookAttacks[sq][occupancy]
}
