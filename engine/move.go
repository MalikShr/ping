package engine

import (
	"math"
)

type Move struct {
	Move  int
	Score int
}

type MoveList struct {
	Moves [maxPositionMoves]Move
	Count int
}

type State struct {
	Hash uint64

	Move       int
	CastlePerm int
	EnPas      int
	Rule50     int
}

var MFLAGEP int = 0x40000
var MFLAGPS int = 0x80000
var MFLAGCA int = 0x1000000

var MFLAGCAP int = 0x7C000
var MFLAGPROM int = 0xF00000

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

func FROMSQ(m int) int {
	return m & 0x7F
}

func TOSQ(m int) int {
	return (m >> 7) & 0x7F
}

func CAPTURED(m int) uint8 {
	return uint8((m >> 14)) & 0xF
}

func PROMOTED(m int) uint8 {
	return uint8((m >> 20) & 0xF)
}

func (pos *BoardStruct) HASHPIECE(piece uint8, sq int) {
	pos.Hash ^= PieceKeys[piece][sq]
}

func (pos *BoardStruct) HASHCASTLE() {
	pos.Hash ^= CastleKeys[pos.CastlePerm]
}

func (pos *BoardStruct) HASHSIDE() {
	pos.Hash ^= SideKey
}

func (pos *BoardStruct) HASHEP() {
	pos.Hash ^= PieceKeys[Empty][pos.EnPas]
}

func MOVE(from int, to int, capture uint8, promotion uint8, flag int) int {
	return from | (to << 7) | (int(capture) << 14) | (int(promotion) << 20) | flag
}

func SQOFFBOARD(sq int) bool {
	return sq > 63 || sq < 0
}

var CastlePerm = [64]int{
	13, 15, 15, 15, 12, 15, 15, 14,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	7, 15, 15, 15, 3, 15, 15, 11,
}

func (pos *BoardStruct) ClearPiece(sq int) {
	piece := pos.Pieces[sq]

	col := PieceCol[piece]
	tPieceNum := -1

	pos.HASHPIECE(piece, sq)

	pos.Pieces[sq] = Empty
	pos.Material[col] -= PieceVal[piece]

	if !PieceBig[piece] {
		CLEARBIT(&pos.Pawns[col], sq)
		CLEARBIT(&pos.Pawns[Both], sq)
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

func ValidPawnDelta1(sq int, direction int) bool {
	return math.Abs(float64(FileOf(sq)-FileOf(sq+(7*direction)))) == 1
}

func ValidPawnDelta2(sq int, direction int) bool {
	return math.Abs(float64(FileOf(sq)-FileOf(sq+(9*direction)))) == 1
}

func (pos *BoardStruct) SqAttacked(sq int, side uint8) bool {

	var list MoveList

	mg := MoveGen{
		pos:  pos,
		list: &list,
		side: side,
		quit: false,
	}

	// pawns
	if side == White {
		if sq-9 > 0 && ((pos.Pieces[sq-9] == wPawn && ValidPawnDelta2(sq, -1)) || (pos.Pieces[sq-7] == wPawn && ValidPawnDelta1(sq, -1))) {
			return true
		}
	} else {
		if sq+9 < 64 && ((pos.Pieces[sq+9] == bPawn && ValidPawnDelta2(sq, 1)) || (pos.Pieces[sq+7] == bPawn && ValidPawnDelta1(sq, 1))) {
			return true
		}
	}

	for sq := 0; sq < 64; sq++ {
		piece := pos.Pieces[sq]

		if piece == Empty || piece == wPawn || piece == bPawn {
			continue
		}

		if IsKn(piece) && PieceCol[piece] == side {
			mg.GenKnightMoves(sq)
		}

		if IsBQ(piece) && PieceCol[piece] == side {
			mg.GenBishopMoves(sq)
		}

		if IsRQ(piece) && PieceCol[piece] == side {
			mg.GenRookMoves(sq)
		}

		if IsKi(piece) && PieceCol[piece] == side {
			mg.GenKingMoves(sq)
		}
	}

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		if TOSQ(list.Moves[moveNum].Move) == sq {
			return true
		}
	}

	return false
}

func (pos *BoardStruct) AddPiece(sq int, piece uint8) {
	col := PieceCol[piece]

	pos.HASHPIECE(piece, sq)

	pos.Pieces[sq] = piece

	if !PieceBig[piece] {
		SETBIT(&pos.Pawns[col], sq)
		SETBIT(&pos.Pawns[Both], sq)
	}

	pos.Material[col] += PieceVal[piece]
	pos.PieceList[piece][pos.PieceNum[piece]] = sq

	pos.PieceNum[piece]++
}

func (pos *BoardStruct) MovePiece(from int, to int) {
	piece := pos.Pieces[from]
	col := PieceCol[piece]

	pos.HASHPIECE(piece, from)
	pos.Pieces[from] = Empty

	pos.HASHPIECE(piece, to)
	pos.Pieces[to] = piece

	if !PieceBig[piece] {
		CLEARBIT(&pos.Pawns[col], from)
		CLEARBIT(&pos.Pawns[Both], from)
		SETBIT(&pos.Pawns[col], to)
		SETBIT(&pos.Pawns[Both], to)
	}

	for index := 0; index < pos.PieceNum[piece]; index++ {
		if pos.PieceList[piece][index] == from {
			pos.PieceList[piece][index] = to
			break
		}
	}
}

func (pos *BoardStruct) MakeMove(move int) bool {
	from := FROMSQ(move)
	to := TOSQ(move)

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

	captured := CAPTURED(move)
	pos.Rule50++

	if captured != Empty {
		pos.ClearPiece(to)
		pos.Rule50 = 0
	}

	pos.HistoryPly++
	pos.Ply++

	if PiecePawn[pos.Pieces[from]] {
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

	promotedPiece := PROMOTED(move)
	if promotedPiece != Empty {
		pos.ClearPiece(to)
		pos.AddPiece(to, promotedPiece)
	}

	if PieceKing[pos.Pieces[to]] {
		pos.KingSq[pos.SideToMove] = to
	}

	pos.SideToMove ^= 1
	pos.HASHSIDE()

	if pos.SqAttacked(pos.KingSq[side], pos.SideToMove) {
		pos.TakeMove()

		return false
	}

	return true
}

func (pos *BoardStruct) TakeMove() {
	pos.HistoryPly--
	pos.Ply--

	move := pos.History[pos.HistoryPly].Move
	from := FROMSQ(move)
	to := TOSQ(move)

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

	if PieceKing[pos.Pieces[from]] {
		pos.KingSq[pos.SideToMove] = from
	}

	captured := CAPTURED(move)
	if captured != Empty {
		pos.AddPiece(to, captured)
	}

	if PROMOTED(move) != Empty {
		pos.ClearPiece(from)

		pawn := wPawn

		if PieceCol[PROMOTED(move)] == Black {
			pawn = bPawn
		}

		pos.AddPiece(from, pawn)
	}
}
