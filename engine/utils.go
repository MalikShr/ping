package engine

import "fmt"

const PceChar string = ".PNBRQKpnbrqk"
const SideChar string = "wb-"
const RankChar string = "12345678"
const FileChar string = "abcdefgh"

var (
	PieceBig = [13]bool{false, false, true, true, true, true, true, false, true, true, true, true, true}
	PieceMaj = [13]bool{false, false, false, false, true, true, true, false, false, false, true, true, true}
	PieceMin = [13]bool{false, false, true, true, false, false, false, false, true, true, false, false, false}
	PieceCol = [13]uint8{Both, White, White, White, White, White, White, Black, Black, Black, Black, Black, Black}

	PiecePawn        = [13]bool{false, true, false, false, false, false, false, true, false, false, false, false, false}
	PieceKnight      = [13]bool{false, false, true, false, false, false, false, false, true, false, false, false, false}
	PieceKing        = [13]bool{false, false, false, false, false, false, true, false, false, false, false, false, true}
	PieceRookQueen   = [13]bool{false, false, false, false, true, true, false, false, false, false, true, true, false}
	PieceBishopQueen = [13]bool{false, false, false, true, false, true, false, false, false, true, false, true, false}
	PieceSlides      = [13]bool{false, false, false, true, true, true, false, false, false, true, true, true, false}
	NonPawnPieces    = [10]uint8{wKnight, wBishop, wRook, wQueen, wKing, bKnight, bBishop, bRook, bQueen, bKing}

	Mirror64 = [64]int{
		56, 57, 58, 59, 60, 61, 62, 63,
		48, 49, 50, 51, 52, 53, 54, 55,
		40, 41, 42, 43, 44, 45, 46, 47,
		32, 33, 34, 35, 36, 37, 38, 39,
		24, 25, 26, 27, 28, 29, 30, 31,
		16, 17, 18, 19, 20, 21, 22, 23,
		8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7,
	}
)

// Given a board square, return it's file.
func FileOf(sq int) int {
	return sq % 8
}

// Given a board square, return it's rank.
func RankOf(sq int) int {
	return sq / 8
}

func PrMove(move int) string {

	var MvStr string

	ff := FileOf(FROMSQ(move))
	rf := RankOf(FROMSQ(move))
	ft := FileOf(TOSQ(move))
	rt := RankOf(TOSQ(move))

	promoted := PROMOTED(move)

	if promoted != 0 {
		pchar := 'q'
		if IsKn(promoted) {
			pchar = 'n'
		} else if IsRQ(promoted) && !IsBQ(promoted) {
			pchar = 'r'
		} else if !IsRQ(promoted) && IsBQ(promoted) {
			pchar = 'b'
		}
		MvStr = fmt.Sprintf("%c%c%c%c%c", ('a' + ff), ('1' + rf), ('a' + ft), ('1' + rt), pchar)
	} else {
		MvStr = fmt.Sprintf("%c%c%c%c", ('a' + ff), ('1' + rf), ('a' + ft), ('1' + rt))
	}

	return MvStr
}

func PrSq(sq int) string {

	var SqStr string

	file := FileOf(sq)
	rank := RankOf(sq)

	SqStr = fmt.Sprintf("%c%c", ('a' + file), ('1' + rank))

	return SqStr
}

func ParseMove(ptrChar string, pos *BoardStruct) int {
	if ptrChar[1] > '8' || ptrChar[1] < '1' {
		return NoMove
	}
	if ptrChar[3] > '8' || ptrChar[3] < '1' {
		return NoMove
	}
	if ptrChar[0] > 'h' || ptrChar[0] < 'a' {
		return NoMove
	}
	if ptrChar[2] > 'h' || ptrChar[2] < 'a' {
		return NoMove
	}

	from := FR2SQ(int(ptrChar[0]-'a'), int(ptrChar[1]-'1'))
	to := FR2SQ(int(ptrChar[2]-'a'), int(ptrChar[3]-'1'))

	var list MoveList
	GenerateAllMoves(pos, &list, true)
	promPiece := Empty

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		move := list.Moves[moveNum].Move
		if FROMSQ(move) == from && TOSQ(move) == to {
			promPiece = PROMOTED(move)
			if promPiece != Empty {
				if IsRQ(promPiece) && !IsBQ(promPiece) && ptrChar[4] == 'r' {
					return move
				} else if !IsRQ(promPiece) && IsBQ(promPiece) && ptrChar[4] == 'b' {
					return move
				} else if IsRQ(promPiece) && IsBQ(promPiece) && ptrChar[4] == 'q' {
					return move
				} else if IsKn(promPiece) && ptrChar[4] == 'n' {
					return move
				}
				continue
			}
			return move
		}
	}

	return NoMove
}

func MIRROR(sq int) int {
	return Mirror64[sq]
}

func CheckBoard(pos *BoardStruct) bool {
	var tPieceNum [13]int
	var tBigPce [2]int
	var tMajPce [2]int
	var tMinPce [2]int
	var tMaterial [2]int

	var pcount int

	var tPawns [3]Bitboard
	tPawns[White] = pos.Pawns[White]
	tPawns[Black] = pos.Pawns[Black]
	tPawns[Both] = pos.Pawns[Both]

	// Check piece lists
	for tPiece := wPawn; tPiece <= bKing; tPiece++ {
		for tPieceNum := 0; tPieceNum < pos.PieceNum[tPiece]; tPieceNum++ {
			sq := pos.PieceList[tPiece][tPieceNum]
			if pos.Pieces[sq] != tPiece {
				fmt.Println("Error: Piece mismatch at square", sq)
				return false
			}
		}
	}

	// Check piece count and other counters
	for sq := 0; sq < 64; sq++ {
		tPiece := pos.Pieces[sq]
		tPieceNum[tPiece]++
		color := PieceCol[tPiece]
		if PieceBig[tPiece] {
			tBigPce[color]++
		}
		if PieceMin[tPiece] {
			tMinPce[color]++
		}
		if PieceMaj[tPiece] {
			tMajPce[color]++
		}

		if color < 2 {
			tMaterial[color] += PieceVal[tPiece]
		}
	}

	for tPiece := wPawn; tPiece <= bKing; tPiece++ {
		if tPieceNum[tPiece] != pos.PieceNum[tPiece] {
			fmt.Println("Error: Piece count mismatch for piece", tPiece)
			return false
		}
	}

	// Check bitboards count
	pcount = tPawns[White].CountBits()
	if pcount != pos.PieceNum[wPawn] {
		fmt.Println("Error: White pawn count mismatch")
		return false
	}
	pcount = tPawns[Black].CountBits()
	if pcount != pos.PieceNum[bPawn] {
		fmt.Println("Error: Black pawn count mismatch")
		return false
	}
	pcount = tPawns[Both].CountBits()
	if pcount != pos.PieceNum[bPawn]+pos.PieceNum[wPawn] {
		fmt.Println("Error: Combined pawn count mismatch")
		return false
	}

	// Check bitboards squares
	for tPawns[White] != 0 {
		sq := tPawns[White].PopBit()
		if pos.Pieces[sq] != wPawn {
			fmt.Println("Error: White pawn bitboard mismatch at square", sq)
			return false
		}
	}

	for tPawns[Black] != 0 {
		sq := tPawns[Black].PopBit()
		if pos.Pieces[sq] != bPawn {
			fmt.Println("Error: Black pawn bitboard mismatch at square", sq)
			return false
		}
	}

	for tPawns[Both] != 0 {
		sq := tPawns[Both].PopBit()
		if pos.Pieces[sq] != bPawn && pos.Pieces[sq] != wPawn {
			fmt.Println("Error: Combined pawn bitboard mismatch at square", sq)
			return false
		}
	}

	if tMaterial[White] != pos.Material[White] || tMaterial[Black] != pos.Material[Black] {
		fmt.Println("Error: Material mismatch")
		return false
	}

	if pos.SideToMove != White && pos.SideToMove != Black {
		fmt.Println("Error: Invalid side")
		return false
	}

	if GeneratePosKey(pos) != pos.Hash {
		fmt.Println("Error: Position key mismatch")
		return false
	}

	if pos.EnPas != NoSq && (RankOf(pos.EnPas) != R6 || (pos.SideToMove == Black && RankOf(pos.EnPas) != R3)) {
		fmt.Println("Error: Invalid en passant square")
		return false
	}

	if pos.Pieces[pos.KingSq[White]] != wKing {
		fmt.Println("Error: White king position mismatch")
		return false
	}

	if pos.Pieces[pos.KingSq[Black]] != bKing {
		fmt.Println("Error: Black king position mismatch")
		return false
	}

	if pos.CastlePerm < 0 || pos.CastlePerm > 15 {
		fmt.Println("Error: Invalid castle permissions")
		return false
	}

	return true
}
