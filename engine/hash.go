package engine

import "math/rand"

var PieceKeys [13][64]uint64
var SideKey uint64
var CastleKeys [16]uint64

func Rand64() uint64 {
	return (uint64(rand.Int()) |
		uint64(rand.Int())<<15 |
		uint64(rand.Int())<<30 |
		uint64(rand.Int())<<45 |
		(uint64(rand.Int())&0xf)<<60)
}

func InitHashKeys() {
	for i := 0; i < 13; i++ {
		for j := 0; j < 64; j++ {
			PieceKeys[i][j] = Rand64()
		}
	}
	SideKey = Rand64()

	for i := 0; i < 16; i++ {
		CastleKeys[i] = Rand64()
	}

}

func GeneratePosKey(pos *BoardStruct) uint64 {
	var piece uint8
	var finalKey uint64

	// Pieces
	for sq := 0; sq < 64; sq++ {
		piece = pos.Squares[sq]

		if piece != Empty {
			// Validate piece
			finalKey ^= PieceKeys[piece][sq]
		}
	}

	// Side
	if pos.SideToMove == White {
		finalKey ^= SideKey
	}

	// En Passant square
	if pos.EnPas != NoSq {
		// Validate en passant square
		finalKey ^= PieceKeys[Empty][pos.EnPas]
	}

	finalKey ^= CastleKeys[pos.CastlePerm]

	return finalKey
}
