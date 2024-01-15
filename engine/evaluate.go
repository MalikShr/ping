package engine

import (
	"math"
)

const (
	IsolatedPawnPenalty = -10
	DoublePawnPenalty   = -10
	RookOpenFile        = 10
	RookSemiOpenFile    = 5
	QueenOpenFile       = 5
	QueenSemiOpenFile   = 3
	BishopPair          = 30

	EndGameMaterialScore = 21340
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

func materialDraw(pos *BoardStruct) bool {
	whiteKnights := pos.Pieces[wKnight].CountBits()
	blackKnights := pos.Pieces[bKnight].CountBits()

	whiteBishops := pos.Pieces[wBishop].CountBits()
	blackBishops := pos.Pieces[bBishop].CountBits()

	whiteRooks := pos.Pieces[wRook].CountBits()
	blackRooks := pos.Pieces[bRook].CountBits()

	whiteQueens := pos.Pieces[wQueen].CountBits()
	blackQueens := pos.Pieces[bQueen].CountBits()

	if whiteRooks == 0 && blackRooks == 0 &&
		whiteQueens == 0 && blackQueens == 0 {

		if blackBishops == 0 && whiteBishops == 0 {
			if whiteKnights < 3 && blackKnights < 3 {
				return true
			}
		} else if whiteKnights == 0 && blackKnights == 0 {
			if math.Abs(float64(whiteBishops-blackBishops)) < 2 {
				return true
			}
		} else if (whiteKnights < 3 && whiteBishops == 0) ||
			(whiteBishops == 1 && whiteKnights == 0) {

			if (blackKnights < 3 && blackBishops == 0) ||
				(blackBishops == 1 && blackKnights == 0) {
				return true
			}
		}
	} else if whiteQueens == 0 && blackQueens == 0 {
		if whiteRooks == 1 && blackRooks == 1 {
			if (whiteKnights+whiteBishops) < 2 &&
				(blackKnights+blackBishops) < 2 {
				return true
			}
		} else if whiteRooks == 1 && blackRooks == 0 {
			if (whiteKnights+whiteBishops == 0) &&
				(((blackKnights + blackBishops) == 1) ||
					((blackKnights + blackBishops) == 2)) {
				return true
			}
		} else if blackRooks == 1 && whiteRooks == 0 {
			if (blackKnights+blackBishops == 0) &&
				(((whiteKnights + whiteBishops) == 1) ||
					((whiteKnights + whiteBishops) == 2)) {
				return true
			}
		}
	}
	return false
}

func getMaterial(pos *BoardStruct, side uint8) int {
	materialVal := 0

	for piece := AllPieces[side][Pawn]; piece <= AllPieces[side][King]; piece++ {
		materialVal += pos.Pieces[piece].CountBits() * PieceVal[piece]
	}

	return materialVal
}

func EvalPosition(pos *BoardStruct) int {
	whiteMaterial := getMaterial(pos, White)
	blackMaterial := getMaterial(pos, Black)

	score := whiteMaterial - blackMaterial

	if pos.Pieces[wPawn].CountBits() != 0 && pos.Pieces[bPawn].CountBits() != 0 && materialDraw(pos) {
		return 0
	}

	for piece := wPawn; piece <= bKing; piece++ {
		pieceBB := pos.Pieces[piece]

		for pieceBB != 0 {
			sq := pieceBB.PopBit()

			switch piece {
			case wPawn:
				score += PstPawn[sq]

				if (IsolatedMask[sq] & pos.Pawns[White]) == 0 {
					score += IsolatedPawnPenalty
				}

				if (FileBBMask[FileOf(sq)] & pos.Pawns[White]).CountBits() >= 2 {
					score += DoublePawnPenalty
				}

				if (WhitePassedMask[sq] & pos.Pawns[Black]) == 0 {
					score += PawnPassed[RankOf(sq)]
				}
			case bPawn:
				score -= PstPawn[Mirror(sq)]

				if (IsolatedMask[sq] & pos.Pawns[Black]) == 0 {
					score -= IsolatedPawnPenalty
				}

				if (FileBBMask[FileOf(sq)] & pos.Pawns[Black]).CountBits() >= 2 {
					score -= DoublePawnPenalty
				}

				if (BlackPassedMask[sq] & pos.Pawns[White]) == 0 {
					score -= PawnPassed[7-RankOf(sq)]
				}
			case wKnight:
				score += PstKnight[sq]
			case bKnight:
				score -= PstKnight[Mirror(sq)]
			case wBishop:
				score += PstBishop[sq]
			case bBishop:
				score -= PstBishop[Mirror(sq)]
			case wRook:
				score += PstRook[sq]

				if (pos.Pawns[Both] & FileBBMask[FileOf(sq)]) == 0 {
					score += RookOpenFile
				} else if (pos.Pawns[White] & FileBBMask[FileOf(sq)]) == 0 {
					score += RookSemiOpenFile
				}
			case bRook:
				score -= PstRook[Mirror(sq)]

				if (pos.Pawns[Both] & FileBBMask[FileOf(sq)]) == 0 {
					score -= RookOpenFile
				} else if (pos.Pawns[Black] & FileBBMask[FileOf(sq)]) == 0 {
					score -= RookSemiOpenFile
				}
			case wQueen:
				if (pos.Pawns[Both] & FileBBMask[FileOf(sq)]) == 0 {
					score += QueenOpenFile
				} else if (pos.Pawns[White] & FileBBMask[FileOf(sq)]) == 0 {
					score += QueenSemiOpenFile
				}
			case bQueen:
				if (pos.Pawns[Both] & FileBBMask[FileOf(sq)]) == 0 {
					score -= QueenOpenFile
				} else if (pos.Pawns[Black] & FileBBMask[FileOf(sq)]) == 0 {
					score -= QueenSemiOpenFile
				}
			case wKing:
				if blackMaterial <= EndGameMaterialScore {
					score += PstKingEG[sq]
				} else {
					score += PstKingMG[sq]
				}
			case bKing:
				if whiteMaterial <= EndGameMaterialScore {
					score -= PstKingEG[Mirror(sq)]
				} else {
					score -= PstKingMG[Mirror(sq)]
				}
			}
		}
	}

	if pos.Pieces[wBishop].CountBits() >= 2 {
		score += BishopPair
	}
	if pos.Pieces[bBishop].CountBits() >= 2 {
		score -= BishopPair
	}

	if pos.SideToMove == White {
		return score
	} else {
		return -score
	}
}
