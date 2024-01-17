package engine

import "fmt"

const Version = "1.0.0"

const (
	Pawn   uint8 = 1
	Knight uint8 = 2
	Bishop uint8 = 3
	Rook   uint8 = 4
	Queen  uint8 = 5
	King   uint8 = 6

	PceChar  string = ".PNBRQKpnbrqk"
	SideChar string = "wb-"
	RankChar string = "12345678"
	FileChar string = "abcdefgh"
)

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
)

var AllPieces = [2][7]uint8{
	{Empty, wPawn, wKnight, wBishop, wRook, wQueen, wKing},
	{Empty, bPawn, bKnight, bBishop, bRook, bQueen, bKing},
}

var Mirror64 = [64]int{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

// A constant mapping piece characters to Piece objects.
var CharToPiece = map[byte]uint8{
	'P': wPawn,
	'N': wKnight,
	'B': wBishop,
	'R': wRook,
	'Q': wQueen,
	'K': wKing,
	'p': bPawn,
	'n': bKnight,
	'b': bBishop,
	'r': bRook,
	'q': bQueen,
	'k': bKing,
}

func FR2SQ(f int, r int) int {
	return f + r*8
}

func FileOf(sq int) int {
	return sq % 8
}

func RankOf(sq int) int {
	return sq / 8
}

func Mirror(sq int) int {
	return Mirror64[sq]
}

func PrSq(sq int) string {
	file := FileOf(sq)
	rank := RankOf(sq)

	return fmt.Sprintf("%c%c", ('a' + file), ('1' + rank))
}
