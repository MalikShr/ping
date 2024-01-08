package engine

import (
	"fmt"
	"strings"
)

const (
	Empty uint8 = 0

	wPawn   uint8 = 1
	wKnight uint8 = 2
	wBishop uint8 = 3
	wRook   uint8 = 4
	wQueen  uint8 = 5
	wKing   uint8 = 6
	bPawn   uint8 = 7
	bKnight uint8 = 8
	bBishop uint8 = 9
	bRook   uint8 = 10
	bQueen  uint8 = 11
	bKing   uint8 = 12

	White uint8 = 0
	Black uint8 = 1
	Both  uint8 = 2

	// Constants mapping each board coordinate to its square
	A1, B1, C1, D1, E1, F1, G1, H1 = 0, 1, 2, 3, 4, 5, 6, 7
	A2, B2, C2, D2, E2, F2, G2, H2 = 8, 9, 10, 11, 12, 13, 14, 15
	A3, B3, C3, D3, E3, F3, G3, H3 = 16, 17, 18, 19, 20, 21, 22, 23
	A4, B4, C4, D4, E4, F4, G4, H4 = 24, 25, 26, 27, 28, 29, 30, 31
	A5, B5, C5, D5, E5, F5, G5, H5 = 32, 33, 34, 35, 36, 37, 38, 39
	A6, B6, C6, D6, E6, F6, G6, H6 = 40, 41, 42, 43, 44, 45, 46, 47
	A7, B7, C7, D7, E7, F7, G7, H7 = 48, 49, 50, 51, 52, 53, 54, 55
	A8, B8, C8, D8, E8, F8, G8, H8 = 56, 57, 58, 59, 60, 61, 62, 63

	// A constant representing no square
	NoSq = 64

	wKCastle = 1
	wQCastle = 2
	bKCastle = 4
	bQCastle = 8

	FENStart = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 "

	NoMove = 0

	maxGameMoves     = 2048
	maxPositionMoves = 256
)

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

type BoardStruct struct {
	Hash uint64

	Pieces [64]uint8
	Pawns  [3]uint64

	KingSq [2]int

	SideToMove uint8
	EnPas      int
	Rule50     int

	Ply        int
	HistoryPly int

	CastlePerm int

	PieceNum [13]int
	Material [2]int

	History [maxGameMoves]State

	PieceList [13][10]int

	SearchHistory [13][64]int
	SearchKillers [2][MaxDepth]int
}

var SetMask [64]uint64
var ClearMask [64]uint64

var PieceKeys [13][64]uint64
var SideKey uint64
var CastleKeys [16]uint64

var FileBBMask [8]uint64
var RankBBMask [8]uint64

var BlackPassedMask [64]uint64
var WhitePassedMask [64]uint64
var IsolatedMask [64]uint64

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

func CLEARBIT(bb *uint64, sq int) {
	*bb &= ClearMask[sq]
}

func SETBIT(bb *uint64, sq int) {
	*bb |= SetMask[sq]
}

func (pos *BoardStruct) ResetBoard() {
	for i := 0; i < 64; i++ {
		pos.Pieces[i] = Empty
	}

	for i := 0; i < 2; i++ {
		pos.Material[i] = 0
		pos.Pawns[i] = 0
	}

	for i := 0; i < 3; i++ {
		pos.Pawns[i] = 0
	}

	for i := 0; i < 13; i++ {
		pos.PieceNum[i] = 0
	}

	pos.KingSq[White] = NoSq
	pos.KingSq[Black] = NoSq

	pos.SideToMove = Both
	pos.EnPas = NoSq
	pos.Rule50 = 0

	pos.Ply = 0
	pos.HistoryPly = 0

	pos.CastlePerm = 0

	pos.Hash = 0
}

func (pos *BoardStruct) UpdateListsMaterial() {

	for i := 0; i < 64; i++ {
		sq := i
		piece := pos.Pieces[i]

		if piece != Empty {
			color := PieceCol[piece]

			pos.Material[color] += PieceVal[piece]

			pos.PieceList[piece][pos.PieceNum[piece]] = sq
			pos.PieceNum[piece]++

			if piece == wKing {
				pos.KingSq[White] = sq
			}

			if piece == bKing {
				pos.KingSq[Black] = sq
			}

			if piece == wPawn {
				SETBIT(&pos.Pawns[White], sq)
				SETBIT(&pos.Pawns[Both], sq)
			}
			if piece == bPawn {
				SETBIT(&pos.Pawns[Black], sq)
				SETBIT(&pos.Pawns[Both], sq)
			}
		}
	}
}

func (pos *BoardStruct) ParseFen(fen string) {

	pos.ResetBoard()

	// Load in each field of the FEN string.
	fields := strings.Fields(fen)
	pieces := fields[0]
	color := fields[1]
	castling := fields[2]
	ep := fields[3]

	// Loop over each square of the board, rank by rank, from left to right,
	// loading in pieces at squares described by the FEN string.
	for index, sq := 0, uint8(56); index < len(pieces); index++ {
		char := pieces[index]
		switch char {
		case 'p', 'n', 'b', 'r', 'q', 'k', 'P', 'N', 'B', 'R', 'Q', 'K':
			piece := CharToPiece[char]
			pos.Pieces[sq] = piece
			sq++
		case '/':
			sq -= 16
		case '1', '2', '3', '4', '5', '6', '7', '8':
			sq += pieces[index] - '0'
		}
	}

	// Set the side to move for the position.
	pos.SideToMove = Black
	if color == "w" {
		pos.SideToMove = White
	}

	if ep != "-" {
		pos.EnPas = ParseMove(ep, pos)
	}

	for _, char := range castling {
		switch char {
		case 'K':
			pos.CastlePerm |= wKCastle
		case 'Q':
			pos.CastlePerm |= wQCastle
		case 'k':
			pos.CastlePerm |= bKCastle
		case 'q':
			pos.CastlePerm |= bQCastle
		}
	}

	pos.Hash = GeneratePosKey(pos)
	pos.UpdateListsMaterial()
}

func (pos *BoardStruct) String() string {
	boardStr := "\nGame Board:\n\n"

	for rank := R8; rank >= R1; rank-- {
		boardStr += fmt.Sprintf("%d  ", rank+1)
		for file := FA; file <= FH; file++ {
			sq := FR2SQ(file, rank)
			piece := pos.Pieces[sq]
			boardStr += fmt.Sprintf("%3c", PceChar[piece])
		}
		boardStr += "\n"
	}

	boardStr += "\n   "
	for file := FA; file <= FH; file++ {
		boardStr += fmt.Sprintf("%3c", 'a'+file)
	}
	boardStr += "\n"
	boardStr += fmt.Sprintf("side:%c\n", SideChar[pos.SideToMove])

	ep := "None"

	if pos.EnPas != NoSq {
		ep = fmt.Sprintf("enPas:%d\n", pos.EnPas)
	}

	boardStr += ep

	castleWK := "-"
	castleWQ := "-"
	castleBK := "-"
	castleBQ := "-"

	if pos.CastlePerm&wKCastle != 0 {
		castleWK = "K"
	}
	if pos.CastlePerm&wQCastle != 0 {
		castleWQ = "Q"
	}
	if pos.CastlePerm&bKCastle != 0 {
		castleBK = "k"
	}
	if pos.CastlePerm&bQCastle != 0 {
		castleBQ = "q"
	}

	pos.Hash = GeneratePosKey(pos)

	boardStr += fmt.Sprintf("castle:%s%s%s%s\n", castleWK, castleWQ, castleBK, castleBQ)
	boardStr += fmt.Sprintf("PosKey:%d\n", pos.Hash)

	return boardStr
}

func InitEvalMasks() {

	sq := NoSq

	for sq = 0; sq < 8; sq++ {
		FileBBMask[sq] = 0
		RankBBMask[sq] = 0
	}

	for r := R8; r >= R1; r-- {
		for f := FA; f <= FH; f++ {
			sq = r*8 + f
			FileBBMask[f] |= (1 << sq)
			RankBBMask[r] |= (1 << sq)
		}
	}

	for r := R8; r >= R1; r-- {
		for f := FA; f <= FH; f++ {
			sq = r*8 + f
			FileBBMask[f] |= (1 << sq)
			RankBBMask[r] |= (1 << sq)
		}
	}

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

func InitBitMasks() {
	for i := 0; i < 64; i++ {
		SetMask[i] = 0
		ClearMask[i] = 0
	}

	for i := 0; i < 64; i++ {
		SetMask[i] = (1 << i)
	}

	for i, value := range SetMask {
		ClearMask[i] = ^value
	}
}
