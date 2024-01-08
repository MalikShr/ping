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

var PawnPassed = [8]int{0, 5, 10, 20, 35, 60, 100, 200}

var PawnTable = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 0, -10, -10, 0, 10, 10,
	5, 0, 0, 5, 5, 0, 0, 5,
	0, 0, 10, 20, 20, 10, 0, 0,
	5, 5, 5, 10, 10, 5, 5, 5,
	10, 10, 10, 20, 20, 10, 10, 10,
	20, 20, 20, 30, 30, 20, 20, 20,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var KnightTable = [64]int{
	0, -10, 0, 0, 0, 0, -10, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 10, 10, 10, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 5, 0,
	5, 10, 15, 20, 20, 15, 10, 5,
	5, 10, 10, 20, 20, 10, 10, 5,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var BishopTable = [64]int{
	0, 0, -10, 0, 0, -10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 10, 15, 15, 10, 0, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 0, 10, 15, 15, 10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var RookTable = [64]int{
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	25, 25, 25, 25, 25, 25, 25, 25,
	0, 0, 5, 10, 10, 5, 0, 0,
}

var KingE = [64]int{
	-50, -10, 0, 0, 0, 0, -10, -50,
	-10, 0, 10, 10, 10, 10, 0, -10,
	0, 10, 20, 20, 20, 20, 10, 0,
	0, 10, 20, 40, 40, 20, 10, 0,
	0, 10, 20, 40, 40, 20, 10, 0,
	0, 10, 20, 20, 20, 20, 10, 0,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-50, -10, 0, 0, 0, 0, -10, -50,
}

var KingO = [64]int{
	0, 5, 5, -10, -10, 0, 10, 5,
	-30, -30, -30, -30, -30, -30, -30, -30,
	-50, -50, -50, -50, -50, -50, -50, -50,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
}

func MIRROR(sq int) int {
	return Mirror64[sq]
}

func calculatePieceMobility(piece uint8, sq int, pos *BoardStruct) int {
	var list MoveList

	if IsKn(piece) && PieceCol[piece] == pos.SideToMove {
		GenKnightMoves(sq, pos, &list, false)
		return list.Count / 2
	}

	if IsBQ(piece) && PieceCol[piece] == pos.SideToMove {
		GenBishopMoves(sq, pos, &list, false)
		return list.Count / 2
	}

	if IsRQ(piece) && PieceCol[piece] == pos.SideToMove {
		GenRookMoves(sq, pos, &list, false)
		return list.Count / 3
	}

	return 0
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

		score += PawnTable[sq]

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

		score -= PawnTable[MIRROR(sq)]

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
		score += KnightTable[sq]
		score += calculatePieceMobility(piece, sq, pos)
	}

	piece = bKnight
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score -= KnightTable[MIRROR(sq)]
		score -= calculatePieceMobility(piece, sq, pos)
	}

	piece = wBishop
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score += BishopTable[sq]
		score += calculatePieceMobility(piece, sq, pos)
	}

	piece = bBishop
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score -= BishopTable[MIRROR(sq)]
		score -= calculatePieceMobility(piece, sq, pos)
	}

	piece = wRook
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score += RookTable[sq]
		score += calculatePieceMobility(piece, sq, pos)

		if (pos.Pawns[Both] & FileBBMask[FileOf(sq)]) == 0 {
			score += RookOpenFile
		} else if (pos.Pawns[White] & FileBBMask[FileOf(sq)]) == 0 {
			score += RookSemiOpenFile
		}
	}

	piece = bRook
	for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
		sq = pos.PieceList[piece][pieceNum]
		score -= RookTable[MIRROR(sq)]
		score -= calculatePieceMobility(piece, sq, pos)

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
		score += KingE[sq]
	} else {
		score += KingO[sq]
	}

	piece = bKing
	sq = pos.PieceList[piece][0]

	if pos.Material[White] <= EndGameMaterial() {
		score -= KingE[MIRROR(sq)]
	} else {
		score -= KingO[MIRROR(sq)]
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
