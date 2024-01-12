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
	SliderPieces     = []uint8{wBishop, wRook, wQueen, bBishop, bRook, bQueen}

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

func FileOf(sq int) int {
	return sq % 8
}

func RankOf(sq int) int {
	return sq / 8
}

func FR2SQ(f int, r int) int {
	return f + r*8
}

func IsBQ(p uint8) bool {
	return PieceBishopQueen[p]
}

func IsRQ(p uint8) bool {
	return PieceRookQueen[p]
}

func IsKn(p uint8) bool {
	return PieceKnight[p]
}

func IsKi(p uint8) bool {
	return PieceKing[p]
}

func FromSq(m int) int {
	return m & 0x7F
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
		if FromSq(move) == from && ToSq(move) == to {
			promPiece = Promoted(move)
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
