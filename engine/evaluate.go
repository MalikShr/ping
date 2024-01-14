package engine

import "math"

const (
	PawnIsolated      = -10
	RookOpenFile      = 10
	RookSemiOpenFile  = 5
	QueenOpenFile     = 5
	QueenSemiOpenFile = 3
	BishopPair        = 30
)

var BlackPassedMask [64]Bitboard
var WhitePassedMask [64]Bitboard
var IsolatedMask [64]Bitboard

var PawnPassed = [8]int{0, 5, 10, 20, 35, 60, 100, 200}

var PstPawn = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, -20, -20, 10, 10, 5,
	5, -5, -10, 0, 0, -10, -5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, 5, 10, 25, 25, 10, 5, 5,
	10, 10, 20, 30, 30, 20, 10, 10,
	50, 50, 50, 50, 50, 50, 50, 50,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PstKnight = [64]int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-30, 5, 10, 15, 15, 10, 5, -30,
	-30, 0, 15, 20, 20, 15, 0, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-40, -20, 0, 0, 0, 0, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}

var PstBishop = [64]int{
	-20, -10, -10, -10, -10, -10, -10, -20,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-10, 5, 5, 10, 10, 5, 5, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-20, -10, -10, -10, -10, -10, -10, -20,
}

var PstRook = [64]int{
	0, 0, 0, 5, 5, 0, 0, 0,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	5, 10, 10, 10, 10, 10, 10,
	5, 0, 0, 0, 0, 0, 0, 0, 0,
}

var PstKingMG = [64]int{
	20, 30, 10, 0, 0, 10, 30, 20,
	20, 20, 0, 0, 0, 0, 20, 20,
	-10, -20, -20, -20, -20, -20, -20, -10,
	-20, -30, -30, -40, -40, -30, -30, -20,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
}

var PstKingEG = [64]int{
	-50, -30, -30, -30, -30, -30, -30, -50,
	-30, -30, 0, 0, 0, 0, -30, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -20, -10, 0, 0, -10, -20, -30,
	-50, -40, -30, -20, -20, -30, -40, -50,
}

var PieceVal = [13]int{0, 100, 320, 330, 500, 900, 20_000, 100, 320, 330, 500, 900, 20_000}

func InitEvalMasks() {

	for sq := 0; sq < 64; sq++ {
		IsolatedMask[sq] = 0
		WhitePassedMask[sq] = 0
		BlackPassedMask[sq] = 0
	}

	for sq := 0; sq < 64; sq++ {
		tsq := sq + 8

		for tsq < 64 {
			WhitePassedMask[sq] |= (1 << tsq)
			tsq += 8
		}

		tsq = sq - 8
		for tsq >= 0 {
			BlackPassedMask[sq] |= (1 << tsq)
			tsq -= 8
		}

		if FileOf(sq) > FA {
			IsolatedMask[sq] |= FileBBMask[FileOf(sq)-1]

			tsq = sq + 7
			for tsq < 64 {
				WhitePassedMask[sq] |= (1 << tsq)
				tsq += 8
			}

			tsq = sq - 9
			for tsq >= 0 {
				BlackPassedMask[sq] |= (1 << tsq)
				tsq -= 8
			}
		}

		if FileOf(sq) < FH {
			IsolatedMask[sq] |= FileBBMask[FileOf(sq)+1]

			tsq = sq + 9
			for tsq < 64 {
				WhitePassedMask[sq] |= (1 << tsq)
				tsq += 8
			}

			tsq = sq - 7
			for tsq >= 0 {
				BlackPassedMask[sq] |= (1 << tsq)
				tsq -= 8
			}
		}
	}
}

func MaterialDraw(pos *BoardStruct) bool {
	if pos.PieceNum[wRook] == 0 && pos.PieceNum[bRook] == 0 &&
		pos.PieceNum[wQueen] == 0 && pos.PieceNum[bQueen] == 0 {

		if pos.PieceNum[bBishop] == 0 && pos.PieceNum[wBishop] == 0 {
			if pos.PieceNum[wKnight] < 3 && pos.PieceNum[bKnight] < 3 {
				return true
			}
		} else if pos.PieceNum[wKnight] == 0 && pos.PieceNum[bKnight] == 0 {
			if math.Abs(float64(pos.PieceNum[wBishop]-pos.PieceNum[bBishop])) < 2 {
				return true
			}
		} else if (pos.PieceNum[wKnight] < 3 && pos.PieceNum[wBishop] == 0) ||
			(pos.PieceNum[wBishop] == 1 && pos.PieceNum[wKnight] == 0) {

			if (pos.PieceNum[bKnight] < 3 && pos.PieceNum[bBishop] == 0) ||
				(pos.PieceNum[bBishop] == 1 && pos.PieceNum[bKnight] == 0) {
				return true
			}
		}
	} else if pos.PieceNum[wQueen] == 0 && pos.PieceNum[bQueen] == 0 {
		if pos.PieceNum[wRook] == 1 && pos.PieceNum[bRook] == 1 {
			if (pos.PieceNum[wKnight]+pos.PieceNum[wBishop]) < 2 &&
				(pos.PieceNum[bKnight]+pos.PieceNum[bBishop]) < 2 {
				return true
			}
		} else if pos.PieceNum[wRook] == 1 && pos.PieceNum[bRook] == 0 {
			if (pos.PieceNum[wKnight]+pos.PieceNum[wBishop] == 0) &&
				(((pos.PieceNum[bKnight] + pos.PieceNum[bBishop]) == 1) ||
					((pos.PieceNum[bKnight] + pos.PieceNum[bBishop]) == 2)) {
				return true
			}
		} else if pos.PieceNum[bRook] == 1 && pos.PieceNum[wRook] == 0 {
			if (pos.PieceNum[bKnight]+pos.PieceNum[bBishop] == 0) &&
				(((pos.PieceNum[wKnight] + pos.PieceNum[wBishop]) == 1) ||
					((pos.PieceNum[wKnight] + pos.PieceNum[wBishop]) == 2)) {
				return true
			}
		}
	}
	return false
}

func EndGameMaterial() int {
	return (1*PieceVal[wRook] + 2*PieceVal[wKnight] + 2*PieceVal[wPawn] + PieceVal[wKing])
}

func EvalPosition(pos *BoardStruct) int {
	var piece uint8
	var sq int
	score := pos.Material[White] - pos.Material[Black]

	if pos.PieceNum[wPawn] != 0 && pos.PieceNum[bPawn] != 0 && MaterialDraw(pos) {
		return 0
	}

	piece = wPawn
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]

		score += PstPawn[sq]

		if (IsolatedMask[sq] & pos.Pawns[White]) == 0 {
			score += PawnIsolated
		}

		if (WhitePassedMask[sq] & pos.Pawns[Black]) == 0 {
			score += PawnPassed[RankOf(sq)]
		}

	}

	piece = bPawn
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]

		score -= PstPawn[Mirror(sq)]

		if (IsolatedMask[sq] & pos.Pawns[Black]) == 0 {
			score -= PawnIsolated
		}

		if (BlackPassedMask[sq] & pos.Pawns[White]) == 0 {
			score -= PawnPassed[7-RankOf(sq)]
		}
	}

	piece = wKnight
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score += PstKnight[sq]
	}

	piece = bKnight
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score -= PstKnight[Mirror(sq)]
	}

	piece = wBishop
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score += PstBishop[sq]
	}

	piece = bBishop
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score -= PstBishop[Mirror(sq)]
	}

	piece = wRook
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score += PstRook[sq]

		if (pos.Pawns[Both] & FileBBMask[FileOf(sq)]) == 0 {
			score += RookOpenFile
		} else if (pos.Pawns[White] & FileBBMask[FileOf(sq)]) == 0 {
			score += RookSemiOpenFile
		}
	}

	piece = bRook
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score -= PstRook[Mirror(sq)]

		if (pos.Pawns[Both] & FileBBMask[FileOf(sq)]) == 0 {
			score -= RookOpenFile
		} else if (pos.Pawns[Black] & FileBBMask[FileOf(sq)]) == 0 {
			score -= RookSemiOpenFile
		}
	}

	piece = wQueen
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]

		if (pos.Pawns[Both] & FileBBMask[FileOf(sq)]) == 0 {
			score += QueenOpenFile
		} else if (pos.Pawns[White] & FileBBMask[FileOf(sq)]) == 0 {
			score += QueenSemiOpenFile
		}
	}

	piece = bQueen
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]

		if (pos.Pawns[Both] & FileBBMask[FileOf(sq)]) == 0 {
			score -= QueenOpenFile
		} else if (pos.Pawns[Black] & FileBBMask[FileOf(sq)]) == 0 {
			score -= QueenSemiOpenFile
		}
	}
	//8/p6k/6p1/5p2/P4K2/8/5pB1/8 b - - 2 62
	piece = wKing
	sq = pos.PieceList[piece][0]

	if pos.Material[Black] <= EndGameMaterial() {
		score += PstKingEG[sq]
	} else {
		score += PstKingMG[sq]
	}

	piece = bKing
	sq = pos.PieceList[piece][0]

	if pos.Material[White] <= EndGameMaterial() {
		score -= PstKingEG[Mirror(sq)]
	} else {
		score -= PstKingMG[Mirror(sq)]
	}

	if pos.PieceNum[wBishop] >= 2 {
		score += BishopPair
	}
	if pos.PieceNum[bBishop] >= 2 {
		score -= BishopPair
	}

	if pos.SideToMove == White {
		return score
	} else {
		return -score
	}
}
