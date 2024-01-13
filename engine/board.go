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

type BoardStruct struct {
	Hash uint64

	Squares [64]uint8
	Pieces  [13]Bitboard
	Pawns   [3]Bitboard
	Sides   [3]Bitboard

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

	SearchHistory [13][64]uint16
	SearchKillers [2][MaxDepth]int
}

type State struct {
	Hash uint64

	Move       int
	CastlePerm int
	EnPas      int
	Rule50     int
}

func (pos *BoardStruct) ResetBoard() {

	for i := 0; i < 64; i++ {
		pos.Squares[i] = Empty
	}

	for i := 0; i < 2; i++ {
		pos.Material[i] = 0
		pos.Pawns[i] = 0
		pos.KingSq[i] = NoSq
	}

	for i := 0; i < 13; i++ {
		pos.Pieces[i] = 0
	}

	for i := 0; i < 3; i++ {
		pos.Pawns[i] = 0
		pos.Sides[i] = 0
	}

	for i := 0; i < 13; i++ {
		pos.PieceNum[i] = 0
	}

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
		piece := pos.Squares[i]

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

			pos.Sides[color].SetBit(sq)
			pos.Sides[Both].SetBit(sq)
			pos.Pieces[piece].SetBit(sq)

			if piece == wPawn {
				pos.Pawns[White].SetBit(sq)
				pos.Pawns[Both].SetBit(sq)
			}
			if piece == bPawn {
				pos.Pawns[Black].SetBit(sq)
				pos.Pawns[Both].SetBit(sq)
			}
		}
	}
}

func (pos *BoardStruct) DoMove(move int) bool {
	from := FromSq(move)
	to := ToSq(move)

	side := pos.SideToMove

	state := State{
		Hash:       pos.Hash,
		Move:       move,
		Rule50:     pos.Rule50,
		EnPas:      pos.EnPas,
		CastlePerm: pos.CastlePerm,
	}

	if move&MFLAGEP != 0 {
		if side == White {
			pos.ClearPiece(to - 8)
		} else {
			pos.ClearPiece(to + 8)
		}
	} else if move&MFLAGCA != 0 {
		switch to {
		case C1:
			pos.MovePiece(A1, D1)
		case C8:
			pos.MovePiece(A8, D8)
		case G1:
			pos.MovePiece(H1, F1)
		case G8:
			pos.MovePiece(H8, F8)
		default:
			return false
		}
	}

	if pos.EnPas != NoSq {
		pos.HASHEP()
	}
	pos.HASHCASTLE()

	pos.History[pos.HistoryPly] = state

	pos.CastlePerm &= CastlePerm[from]
	pos.CastlePerm &= CastlePerm[to]
	pos.EnPas = NoSq

	pos.HASHCASTLE()

	captured := Captured(move)
	pos.Rule50++

	if captured != Empty {
		pos.ClearPiece(to)
		pos.Rule50 = 0
	}

	pos.HistoryPly++
	pos.Ply++

	if PiecePawn[pos.Squares[from]] {
		pos.Rule50 = 0
		if move&MFLAGPS != 0 {
			if side == White {
				pos.EnPas = from + 8
			} else {
				pos.EnPas = from - 8
			}
			pos.HASHEP()
		}
	}

	pos.MovePiece(from, to)

	promotedPiece := Promoted(move)
	if promotedPiece != Empty {
		pos.ClearPiece(to)
		pos.AddPiece(to, promotedPiece)
	}

	if PieceKing[pos.Squares[to]] {
		pos.KingSq[pos.SideToMove] = to
	}

	pos.SideToMove ^= 1
	pos.HASHSIDE()

	if SqAttacked(pos.KingSq[side], pos, pos.SideToMove) {
		pos.UndoMove()

		return false
	}

	return true
}

func (pos *BoardStruct) UndoMove() {
	pos.HistoryPly--
	pos.Ply--

	move := pos.History[pos.HistoryPly].Move
	from := FromSq(move)
	to := ToSq(move)

	if pos.EnPas != NoSq {
		pos.HASHEP()
	}

	pos.HASHCASTLE()

	pos.CastlePerm = pos.History[pos.HistoryPly].CastlePerm
	pos.Rule50 = pos.History[pos.HistoryPly].Rule50
	pos.EnPas = pos.History[pos.HistoryPly].EnPas

	if pos.EnPas != NoSq {
		pos.HASHEP()
	}
	pos.HASHCASTLE()

	pos.SideToMove ^= 1
	pos.HASHSIDE()

	if MFLAGEP&move != 0 {
		if pos.SideToMove == White {
			pos.AddPiece(to-8, bPawn)
		} else {
			pos.AddPiece(to+8, wPawn)
		}
	} else if MFLAGCA&move != 0 {
		switch to {
		case C1:
			pos.MovePiece(D1, A1)
		case C8:
			pos.MovePiece(D8, A8)
		case G1:
			pos.MovePiece(F1, H1)
		case G8:
			pos.MovePiece(F8, H8)
		}
	}

	pos.MovePiece(to, from)

	if PieceKing[pos.Squares[from]] {
		pos.KingSq[pos.SideToMove] = from
	}

	captured := Captured(move)
	if captured != Empty {
		pos.AddPiece(to, captured)
	}

	if Promoted(move) != Empty {
		pos.ClearPiece(from)

		pawn := wPawn

		if PieceCol[Promoted(move)] == Black {
			pawn = bPawn
		}

		pos.AddPiece(from, pawn)
	}
}

func (pos *BoardStruct) MoveExists(move int) bool {
	var list MoveList
	GenerateAllMoves(pos, &list, true)

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		if !pos.DoMove(list.Moves[moveNum].Move) {
			continue
		}
		pos.UndoMove()
		if list.Moves[moveNum].Move == move {
			return true
		}
	}

	return false
}

func (pos *BoardStruct) MovePiece(from int, to int) {
	piece := pos.Squares[from]
	col := PieceCol[piece]

	pos.Pieces[piece].ClearBit(from)
	pos.Sides[col].ClearBit(from)
	pos.Sides[Both].ClearBit(from)
	pos.HASHPIECE(piece, from)
	pos.Squares[from] = Empty

	pos.Pieces[piece].SetBit(to)
	pos.Sides[col].SetBit(to)
	pos.Sides[Both].SetBit(to)
	pos.HASHPIECE(piece, to)
	pos.Squares[to] = piece

	if !PieceBig[piece] {
		pos.Pawns[col].ClearBit(from)
		pos.Pawns[Both].ClearBit(from)
		pos.Pawns[col].SetBit(to)
		pos.Pawns[Both].SetBit(to)
	}

	for index := 0; index < pos.PieceNum[piece]; index++ {
		if pos.PieceList[piece][index] == from {
			pos.PieceList[piece][index] = to
			break
		}
	}
}

func (pos *BoardStruct) AddPiece(sq int, piece uint8) {
	col := PieceCol[piece]

	pos.HASHPIECE(piece, sq)

	pos.Squares[sq] = piece
	pos.Sides[col].SetBit(sq)
	pos.Sides[Both].SetBit(sq)
	pos.Pieces[piece].SetBit(sq)

	if !PieceBig[piece] {
		pos.Pawns[col].SetBit(sq)
		pos.Pawns[Both].SetBit(sq)
	}

	pos.Material[col] += PieceVal[piece]
	pos.PieceList[piece][pos.PieceNum[piece]] = sq

	pos.PieceNum[piece]++
}

func (pos *BoardStruct) ClearPiece(sq int) {
	piece := pos.Squares[sq]

	col := PieceCol[piece]
	tPieceNum := -1

	pos.HASHPIECE(piece, sq)

	pos.Sides[col].ClearBit(sq)
	pos.Sides[Both].ClearBit(sq)
	pos.Pieces[piece].ClearBit(sq)
	pos.Squares[sq] = Empty
	pos.Material[col] -= PieceVal[piece]

	if !PieceBig[piece] {
		pos.Pawns[col].ClearBit(sq)
		pos.Pawns[Both].ClearBit(sq)
	}

	for i := 0; i < pos.PieceNum[piece]; i++ {
		if pos.PieceList[piece][i] == sq {
			tPieceNum = i
			break
		}
	}

	pos.PieceNum[piece]--

	pos.PieceList[piece][tPieceNum] = pos.PieceList[piece][pos.PieceNum[piece]]
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
			pos.Squares[sq] = piece
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
		file := ep[0] - 'a'
		rank := int(ep[1]-'0') - 1

		pos.EnPas = FR2SQ(int(file), rank)
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
			piece := pos.Squares[sq]
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

	ep := "enPas: none\n"

	if pos.EnPas != NoSq {
		ep = fmt.Sprintf("enPas: %s\n", PrSq(pos.EnPas))
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
