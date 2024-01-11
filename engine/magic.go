package engine

// not A file constant
var notAFile Bitboard = 18374403900871474942

// not H file constant
var notHFile Bitboard = 9187201950435737471

// not HG file constant
var notHGFile Bitboard = 4557430888798830399

// not AB file constant
var notABFile Bitboard = 18229723555195321596

var KnightAttacks [64]Bitboard
var KingAttacks [64]Bitboard
var PawnAttacks [2][64]Bitboard
var BishopRelevantBits [64]Bitboard

var BishopMasks [64]Bitboard
var RookMasks [64]Bitboard

var BishopAttacks [64][512]Bitboard
var RookAttacks [64][4096]Bitboard

// bishop relevant occupancy bit count for every square on board
var RelevantBishopBits = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

// rook relevant occupancy bit count for every square on board
var RelevantRookBits = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

func InitTables() {
	for sq := 0; sq < 64; sq++ {
		bishopAttackMask := MaskBishopAttacks(sq)
		rookAttackMask := MaskRookAttacks(sq)

		PawnAttacks[White][sq] = MaskPawnAttacks(White, sq)
		PawnAttacks[Black][sq] = MaskPawnAttacks(Black, sq)
		KnightAttacks[sq] = MaskKnightAttacks(sq)
		KingAttacks[sq] = MaskKingAttacks(sq)
		BishopMasks[sq] = bishopAttackMask
		RookMasks[sq] = rookAttackMask

		bishopRelevantBitCount := bishopAttackMask.CountBits()
		rookRelevantBitsCount := rookAttackMask.CountBits()

		bishopOccupancyIndicies := (1 << bishopRelevantBitCount)
		rookOccupancyIndicies := (1 << rookRelevantBitsCount)

		for i := 0; i < bishopOccupancyIndicies; i++ {
			occupancy := SetOccupancies(i, bishopRelevantBitCount, bishopAttackMask)
			magicIndex := (occupancy * MagicB[sq]) >> (64 - RelevantBishopBits[sq])
			BishopAttacks[sq][magicIndex] = GenBishopAttacksFly(sq, occupancy)
		}

		for i := 0; i < rookOccupancyIndicies; i++ {
			occupancy := SetOccupancies(i, rookRelevantBitsCount, rookAttackMask)

			magicIndex := (occupancy * MagicR[sq]) >> (64 - RelevantRookBits[sq])
			RookAttacks[sq][magicIndex] = GenRookAttacksFly(sq, occupancy)
		}
	}
}

func SetOccupancies(index int, bitsInMask int, attackMask Bitboard) Bitboard {
	// occupancy map
	occupancy := Bitboard(0)

	// Loop over range of bits within attack mask
	for count := 0; count < bitsInMask; count++ {
		sq := attackMask.PopBit()

		// make sure occupancy is on board
		if index&(1<<count) != 0 {
			// update occupancy map
			occupancy |= (1 << sq)
		}
	}

	return occupancy
}

func MaskPawnAttacks(side uint8, sq int) Bitboard {
	// result attacks bitboard
	attacks := Bitboard(0)

	// piece bitboard
	bitboard := Bitboard(0)

	// set piece on board
	bitboard.SetBit(sq)

	// white pawns
	if side == White {
		// generate pawn attacks
		if (bitboard>>7)&notAFile != 0 {
			attacks |= (bitboard >> 7)
		}
		if (bitboard>>9)&notHFile != 0 {
			attacks |= (bitboard >> 9)
		}
	} else if side == Black {
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

func MaskBishopAttacks(sq int) Bitboard {
	attacks := Bitboard(0)

	rank := RankOf(sq)
	file := FileOf(sq)

	for f, r := file+1, rank+1; f < 7 && r < 7; f, r = f+1, r+1 {
		attacks |= 1 << (r*8 + f)
	}

	for f, r := file-1, rank+1; f >= 1 && r < 7; f, r = f-1, r+1 {
		attacks |= 1 << (r*8 + f)
	}

	for f, r := file+1, rank-1; f < 7 && r >= 1; f, r = f+1, r-1 {
		attacks |= 1 << (r*8 + f)
	}

	for f, r := file-1, rank-1; f >= 1 && r >= 1; f, r = f-1, r-1 {
		attacks |= 1 << (r*8 + f)
	}

	return attacks
}

func MaskRookAttacks(sq int) Bitboard {
	attacks := Bitboard(0)

	rank := RankOf(sq)
	file := FileOf(sq)

	for f := file + 1; f < 7; f++ {
		attacks |= 1 << (rank*8 + f)
	}

	for f := file - 1; f >= 1; f-- {
		attacks |= 1 << (rank*8 + f)
	}

	for r := rank + 1; r < 7; r++ {
		attacks |= 1 << (r*8 + file)
	}

	for r := rank - 1; r >= 1; r-- {
		attacks |= 1 << (r*8 + file)
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

func GenBishopAttacksFly(sq int, block Bitboard) Bitboard {
	attacks := Bitboard(0)

	rank := RankOf(sq)
	file := FileOf(sq)

	for f, r := file+1, rank+1; f < 8 && r < 8; f, r = f+1, r+1 {
		attacks |= 1 << (r*8 + f)

		if block&(1<<(r*8+f)) != 0 {
			break
		}
	}

	for f, r := file-1, rank+1; f >= 0 && r < 8; f, r = f-1, r+1 {
		attacks |= 1 << (r*8 + f)

		if block&(1<<(r*8+f)) != 0 {
			break
		}
	}

	for f, r := file+1, rank-1; f < 8 && r >= 0; f, r = f+1, r-1 {
		attacks |= 1 << (r*8 + f)

		if block&(1<<(r*8+f)) != 0 {
			break
		}
	}

	for f, r := file-1, rank-1; f >= 0 && r >= 0; f, r = f-1, r-1 {

		attacks |= 1 << (r*8 + f)

		if block&(1<<(r*8+f)) != 0 {
			break
		}
	}

	return attacks
}

func GenRookAttacksFly(sq int, block Bitboard) Bitboard {
	attacks := Bitboard(0)

	rank := RankOf(sq)
	file := FileOf(sq)

	for f := file + 1; f < 8; f++ {
		attacks |= 1 << (rank*8 + f)

		if block&(1<<(rank*8+f)) != 0 {
			break
		}
	}

	for f := file - 1; f >= 0; f-- {
		attacks |= 1 << (rank*8 + f)

		if block&(1<<(rank*8+f)) != 0 {
			break
		}
	}

	for r := rank + 1; r < 8; r++ {
		attacks |= 1 << (r*8 + file)

		if block&(1<<(r*8+file)) != 0 {
			break
		}
	}

	for r := rank - 1; r >= 0; r-- {
		attacks |= 1 << (r*8 + file)

		if block&(1<<(r*8+file)) != 0 {
			break
		}
	}

	return attacks
}

func GetBishopAttacks(sq int, occupancy Bitboard) Bitboard {
	occupancy &= BishopMasks[sq]
	occupancy *= MagicB[sq]
	occupancy >>= 64 - RelevantBishopBits[sq]

	return BishopAttacks[sq][occupancy]
}

func GetRookAttacks(sq int, occupancy Bitboard) Bitboard {
	occupancy &= RookMasks[sq]
	occupancy *= MagicR[sq]
	occupancy >>= 64 - RelevantRookBits[sq]

	return RookAttacks[sq][occupancy]
}

// rook magic numbers
var MagicR = [64]Bitboard{
	0x8a80104000800020,
	0x140002000100040,
	0x2801880a0017001,
	0x100081001000420,
	0x200020010080420,
	0x3001c0002010008,
	0x8480008002000100,
	0x2080088004402900,
	0x800098204000,
	0x2024401000200040,
	0x100802000801000,
	0x120800800801000,
	0x208808088000400,
	0x2802200800400,
	0x2200800100020080,
	0x801000060821100,
	0x80044006422000,
	0x100808020004000,
	0x12108a0010204200,
	0x140848010000802,
	0x481828014002800,
	0x8094004002004100,
	0x4010040010010802,
	0x20008806104,
	0x100400080208000,
	0x2040002120081000,
	0x21200680100081,
	0x20100080080080,
	0x2000a00200410,
	0x20080800400,
	0x80088400100102,
	0x80004600042881,
	0x4040008040800020,
	0x440003000200801,
	0x4200011004500,
	0x188020010100100,
	0x14800401802800,
	0x2080040080800200,
	0x124080204001001,
	0x200046502000484,
	0x480400080088020,
	0x1000422010034000,
	0x30200100110040,
	0x100021010009,
	0x2002080100110004,
	0x202008004008002,
	0x20020004010100,
	0x2048440040820001,
	0x101002200408200,
	0x40802000401080,
	0x4008142004410100,
	0x2060820c0120200,
	0x1001004080100,
	0x20c020080040080,
	0x2935610830022400,
	0x44440041009200,
	0x280001040802101,
	0x2100190040002085,
	0x80c0084100102001,
	0x4024081001000421,
	0x20030a0244872,
	0x12001008414402,
	0x2006104900a0804,
	0x1004081002402,
}

// bishop magic numbers
var MagicB = [64]Bitboard{
	0x40040844404084,
	0x2004208a004208,
	0x10190041080202,
	0x108060845042010,
	0x581104180800210,
	0x2112080446200010,
	0x1080820820060210,
	0x3c0808410220200,
	0x4050404440404,
	0x21001420088,
	0x24d0080801082102,
	0x1020a0a020400,
	0x40308200402,
	0x4011002100800,
	0x401484104104005,
	0x801010402020200,
	0x400210c3880100,
	0x404022024108200,
	0x810018200204102,
	0x4002801a02003,
	0x85040820080400,
	0x810102c808880400,
	0xe900410884800,
	0x8002020480840102,
	0x220200865090201,
	0x2010100a02021202,
	0x152048408022401,
	0x20080002081110,
	0x4001001021004000,
	0x800040400a011002,
	0xe4004081011002,
	0x1c004001012080,
	0x8004200962a00220,
	0x8422100208500202,
	0x2000402200300c08,
	0x8646020080080080,
	0x80020a0200100808,
	0x2010004880111000,
	0x623000a080011400,
	0x42008c0340209202,
	0x209188240001000,
	0x400408a884001800,
	0x110400a6080400,
	0x1840060a44020800,
	0x90080104000041,
	0x201011000808101,
	0x1a2208080504f080,
	0x8012020600211212,
	0x500861011240000,
	0x180806108200800,
	0x4000020e01040044,
	0x300000261044000a,
	0x802241102020002,
	0x20906061210001,
	0x5a84841004010310,
	0x4010801011c04,
	0xa010109502200,
	0x4a02012000,
	0x500201010098b028,
	0x8040002811040900,
	0x28000010020204,
	0x6000020202d0240,
	0x8918844842082200,
	0x4010011029020020,
}
