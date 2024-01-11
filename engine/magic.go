package engine

var KnightAttacks [64]Bitboard
var KingAttacks [64]Bitboard
var PawnAttacks [2][64]Bitboard

// not A file constant
var notAFile Bitboard = 18374403900871474942

// not H file constant
var notHFile Bitboard = 9187201950435737471

// not HG file constant
var notHGFile Bitboard = 4557430888798830399

// not AB file constant
var notABFile Bitboard = 18229723555195321596

func InitAttacks() {
	for sq := 0; sq < 64; sq++ {
		KnightAttacks[sq] = MaskKnightAttacks(sq)
		KingAttacks[sq] = MaskKingAttacks(sq)
		PawnAttacks[White][sq] = MaskPawnAttacks(White, sq)
		PawnAttacks[Black][sq] = MaskPawnAttacks(Black, sq)
	}
}

func MaskPawnAttacks(side uint8, sq int) Bitboard {
	// result attacks bitboard
	attacks := Bitboard(0)

	// piece bitboard
	bitboard := Bitboard(0)

	// set piece on board
	bitboard.SetBit(sq)

	// white pawns
	if side == Black {
		// generate pawn attacks
		if (bitboard>>7)&notAFile != 0 {
			attacks |= (bitboard >> 7)
		}
		if (bitboard>>9)&notHFile != 0 {
			attacks |= (bitboard >> 9)
		}
	} else if side == White {
		// generate pawn attacks
		if (bitboard<<7)&notHFile != 0 {
			attacks |= (bitboard << 7)
		}
		if (bitboard<<9)&notAFile != 0 {
			attacks |= (bitboard << 9)
		}
	}

	// return attack map
	return attacks
}

func MaskKnightAttacks(sq int) Bitboard {
	attacks := Bitboard(0)
	bitboard := Bitboard(0)

	bitboard.SetBit(sq)

	// Generate Knight Attacks 17, 16, 10, 6
	if (bitboard>>17)&notHFile != 0 {
		attacks |= (bitboard >> 17)
	}

	if (bitboard>>15)&notAFile != 0 {
		attacks |= (bitboard >> 15)
	}

	if (bitboard>>10)&notHGFile != 0 {
		attacks |= (bitboard >> 10)
	}

	if (bitboard>>6)&notABFile != 0 {
		attacks |= (bitboard >> 6)
	}

	if (bitboard<<17)&notAFile != 0 {
		attacks |= (bitboard << 17)
	}

	if (bitboard<<15)&notHFile != 0 {
		attacks |= (bitboard << 15)
	}

	if (bitboard<<10)&notABFile != 0 {
		attacks |= (bitboard << 10)
	}

	if (bitboard<<6)&notHGFile != 0 {
		attacks |= (bitboard << 6)
	}

	return attacks
}

func MaskKingAttacks(sq int) Bitboard {
	// result attacks bitboard
	attacks := Bitboard(0)

	// piece bitboard
	bitboard := Bitboard(0)

	// set piece on board
	bitboard.SetBit(sq)

	// generate king attacks
	if bitboard>>8 != 0 {
		attacks |= (bitboard >> 8)
	}

	if (bitboard>>9)&notHFile != 0 {
		attacks |= (bitboard >> 9)
	}

	if (bitboard>>7)&notAFile != 0 {
		attacks |= (bitboard >> 7)
	}

	if (bitboard>>1)&notHFile != 0 {
		attacks |= (bitboard >> 1)
	}
	if bitboard<<8 != 0 {
		attacks |= (bitboard << 8)
	}
	if (bitboard<<9)&notAFile != 0 {
		attacks |= (bitboard << 9)
	}
	if (bitboard<<7)&notHFile != 0 {
		attacks |= (bitboard << 7)
	}
	if (bitboard<<1)&notAFile != 0 {
		attacks |= (bitboard << 1)
	}

	// return attack map
	return attacks
}

func GenKnightAttacks(sq int, ownPieces Bitboard, block Bitboard, pos *BoardStruct, list *MoveList, quiet bool) {
	attacks := KnightAttacks[sq] & ^ownPieces

	for attacks != 0 {
		targetSq := attacks.PopBit()

		if block&(1<<targetSq) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Squares[targetSq], Empty, 0))
			continue
		}
		if quiet {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}
}

func GenKingAttacks(square int, ownPieces Bitboard, block Bitboard, pos *BoardStruct, list *MoveList, quiet bool) {
	attacks := KingAttacks[square] & ^ownPieces

	for attacks != 0 {
		sq := attacks.PopBit()

		if block&(1<<sq) != 0 {
			list.AddCaptureMove(pos, MOVE(square, sq, pos.Squares[sq], Empty, 0))
			continue
		}
		if quiet {
			list.AddQuietMove(pos, MOVE(square, sq, Empty, Empty, 0))
		}
	}
}

func GenBishopAttacks(sq int, ownPieces Bitboard, block Bitboard, pos *BoardStruct, list *MoveList, quiet bool) {
	attacks := Bitboard(0)

	rank := RankOf(sq)
	file := FileOf(sq)

	for f, r := file+1, rank+1; f < 8 && r < 8; f, r = f+1, r+1 {
		if ownPieces&(1<<(r*8+f)) != 0 {
			break
		}

		attacks |= 1 << (r*8 + f)

		if block&(1<<(r*8+f)) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, r*8+f, pos.Squares[r*8+f], Empty, 0))
			break
		}
		if quiet {
			list.AddQuietMove(pos, MOVE(sq, r*8+f, Empty, Empty, 0))
		}
	}

	for f, r := file-1, rank+1; f >= 0 && r < 8; f, r = f-1, r+1 {
		if ownPieces&(1<<(r*8+f)) != 0 {
			break
		}

		attacks |= 1 << (r*8 + f)
		if block&(1<<(r*8+f)) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, r*8+f, pos.Squares[r*8+f], Empty, 0))
			break
		}
		if quiet {
			list.AddQuietMove(pos, MOVE(sq, r*8+f, Empty, Empty, 0))
		}
	}

	for f, r := file+1, rank-1; f < 8 && r >= 0; f, r = f+1, r-1 {
		if ownPieces&(1<<(r*8+f)) != 0 {
			break
		}

		attacks |= 1 << (r*8 + f)
		if block&(1<<(r*8+f)) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, r*8+f, pos.Squares[r*8+f], Empty, 0))
			break
		}
		if quiet {
			list.AddQuietMove(pos, MOVE(sq, r*8+f, Empty, Empty, 0))
		}
	}

	for f, r := file-1, rank-1; f >= 0 && r >= 0; f, r = f-1, r-1 {
		if ownPieces&(1<<(r*8+f)) != 0 {
			break
		}

		attacks |= 1 << (r*8 + f)
		if block&(1<<(r*8+f)) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, r*8+f, pos.Squares[r*8+f], Empty, 0))
			break
		}
		if quiet {
			list.AddQuietMove(pos, MOVE(sq, r*8+f, Empty, Empty, 0))
		}
	}
}

func GenRookAttacks(sq int, ownPieces Bitboard, block Bitboard, pos *BoardStruct, list *MoveList, quiet bool) {
	attacks := Bitboard(0)

	rank := RankOf(sq)
	file := FileOf(sq)

	for f := file + 1; f < 8; f++ {
		if ownPieces&(1<<(rank*8+f)) != 0 {
			break
		}
		attacks |= 1 << (rank*8 + f)

		if block&(1<<(rank*8+f)) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, rank*8+f, pos.Squares[rank*8+f], Empty, 0))
			break
		}
		if quiet {
			list.AddQuietMove(pos, MOVE(sq, rank*8+f, Empty, Empty, 0))
		}
	}

	for f := file - 1; f >= 0; f-- {
		if ownPieces&(1<<(rank*8+f)) != 0 {
			break
		}
		attacks |= 1 << (rank*8 + f)

		if block&(1<<(rank*8+f)) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, rank*8+f, pos.Squares[rank*8+f], Empty, 0))
			break
		}
		if quiet {
			list.AddQuietMove(pos, MOVE(sq, rank*8+f, Empty, Empty, 0))
		}
	}

	for r := rank + 1; r < 8; r++ {
		if ownPieces&(1<<(r*8+file)) != 0 {
			break
		}
		attacks |= 1 << (r*8 + file)

		if block&(1<<(r*8+file)) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, r*8+file, pos.Squares[r*8+file], Empty, 0))
			break
		}
		if quiet {
			list.AddQuietMove(pos, MOVE(sq, r*8+file, Empty, Empty, 0))
		}
	}

	for r := rank - 1; r >= 0; r-- {
		if ownPieces&(1<<(r*8+file)) != 0 {
			break
		}
		attacks |= 1 << (r*8 + file)

		if block&(1<<(r*8+file)) != 0 {
			list.AddCaptureMove(pos, MOVE(sq, r*8+file, pos.Squares[r*8+file], Empty, 0))
			break
		}
		if quiet {
			list.AddQuietMove(pos, MOVE(sq, r*8+file, Empty, Empty, 0))
		}
	}
}
