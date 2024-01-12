package engine

var FileBBMask [8]Bitboard
var RankBBMask [8]Bitboard

var ClearFile [8]Bitboard
var ClearRank [8]Bitboard

var KnightAttacks [64]Bitboard
var KingAttacks [64]Bitboard
var PawnAttacks [2][64]Bitboard
var BishopRelevantBits [64]Bitboard

const (
	FA int = iota
	FB
	FC
	FD
	FE
	FF
	FG
	FH
	FNone
)

const (
	R1 int = iota
	R2
	R3
	R4
	R5
	R6
	R7
	R8
	RNone
)

func InitTables() {
	for sq := 0; sq < 8; sq++ {
		FileBBMask[sq] = 0
		RankBBMask[sq] = 0
		ClearFile[sq] = FullBB
		ClearRank[sq] = FullBB
	}

	for r := R8; r >= R1; r-- {
		for f := FA; f <= FH; f++ {
			sq := FR2SQ(f, r)
			FileBBMask[f] |= (1 << sq)
			RankBBMask[r] |= (1 << sq)
			ClearFile[f].ClearBit(sq)
			ClearRank[r].ClearBit(sq)
		}
	}

	for sq := 0; sq < 64; sq++ {
		PawnAttacks[White][sq] = calcPawnAttacks(White, sq)
		PawnAttacks[Black][sq] = calcPawnAttacks(Black, sq)
		KnightAttacks[sq] = calcKnightAttacks(sq)
		KingAttacks[sq] = calcKingAttacks(sq)
	}

}

func calcPawnAttacks(side uint8, sq int) Bitboard {
	// result attacks bitboard
	attacks := Bitboard(0)

	// piece bitboard
	bitboard := Bitboard(0)

	// set piece on board
	bitboard.SetBit(sq)

	// white pawns
	if side == White {
		// generate pawn attacks
		if (bitboard>>7)&ClearFile[FA] != 0 {
			attacks |= (bitboard >> 7)
		}
		if (bitboard>>9)&ClearFile[FH] != 0 {
			attacks |= (bitboard >> 9)
		}
	} else if side == Black {
		// generate pawn attacks
		if (bitboard<<7)&ClearFile[FH] != 0 {
			attacks |= (bitboard << 7)
		}
		if (bitboard<<9)&ClearFile[FA] != 0 {
			attacks |= (bitboard << 9)
		}
	}

	// return attack map
	return attacks
}

func calcKnightAttacks(sq int) Bitboard {
	attacks := Bitboard(0)
	bitboard := Bitboard(0)

	bitboard.SetBit(sq)

	// Generate Knight Attacks 17, 16, 10, 6
	if (bitboard>>17)&ClearFile[FH] != 0 {
		attacks |= (bitboard >> 17)
	}

	if (bitboard>>15)&ClearFile[FA] != 0 {
		attacks |= (bitboard >> 15)
	}

	if (bitboard>>10)&(ClearFile[FG]&ClearFile[FH]) != 0 {
		attacks |= (bitboard >> 10)
	}

	if (bitboard>>6)&(ClearFile[FA]&ClearFile[FB]) != 0 {
		attacks |= (bitboard >> 6)
	}

	if (bitboard<<17)&ClearFile[FA] != 0 {
		attacks |= (bitboard << 17)
	}

	if (bitboard<<15)&ClearFile[FH] != 0 {
		attacks |= (bitboard << 15)
	}

	if (bitboard<<10)&(ClearFile[FA]&ClearFile[FB]) != 0 {
		attacks |= (bitboard << 10)
	}

	if (bitboard<<6)&(ClearFile[FG]&ClearFile[FH]) != 0 {
		attacks |= (bitboard << 6)
	}

	return attacks
}

func calcKingAttacks(sq int) Bitboard {
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

	if (bitboard>>9)&ClearFile[FH] != 0 {
		attacks |= (bitboard >> 9)
	}

	if (bitboard>>7)&ClearFile[FA] != 0 {
		attacks |= (bitboard >> 7)
	}

	if (bitboard>>1)&ClearFile[FH] != 0 {
		attacks |= (bitboard >> 1)
	}
	if bitboard<<8 != 0 {
		attacks |= (bitboard << 8)
	}
	if (bitboard<<9)&ClearFile[FA] != 0 {
		attacks |= (bitboard << 9)
	}
	if (bitboard<<7)&ClearFile[FH] != 0 {
		attacks |= (bitboard << 7)
	}
	if (bitboard<<1)&ClearFile[FA] != 0 {
		attacks |= (bitboard << 1)
	}

	// return attack map
	return attacks
}
