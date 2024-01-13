package engine

import "fmt"

// Thanks for many ideas to the Blunder chess engine https://github.com/algerbrex/blunder

const (
	// Constants represeting the four possible move types.
	Quiet     uint8 = 0
	Attack    uint8 = 1
	Castle    uint8 = 2
	Promotion uint8 = 3

	// Constants representing move flags indicating what kind of promotion
	// is occuring.
	KnightPromotion uint8 = 0
	BishopPromotion uint8 = 1
	RookPromotion   uint8 = 2
	QueenPromotion  uint8 = 3

	// A constant representing a move flag indicating an attack is an en passant
	// attack.
	AttackEP uint8 = 1

	// A constant representing a null flag
	NoFlag uint8 = 0
)

type Move uint32

func NewMove(from int, to int, moveType uint8, flag uint8) Move {
	return Move(uint32(from)<<26 | uint32(to)<<20 | uint32(moveType)<<18 | uint32(flag)<<16)
}

// Get the from square of the move.
func (move Move) FromSq() int {
	return int((move & 0xfc000000) >> 26)
}

// Get the to square of the move.
func (move Move) ToSq() int {
	return int((move & 0x3f00000) >> 20)
}

// Get the type of the move.
func (move Move) MoveType() uint8 {
	return uint8((move & 0xc0000) >> 18)
}

// Get the flag of the move.
func (move Move) Flag() uint8 {
	return uint8((move & 0x30000) >> 16)
}

// Get the score of a move.
func (move Move) Score() uint16 {
	return uint16(move & 0xffff)
}

// Add a score to the move for move ordering.
func (move *Move) AddScore(score uint16) {
	(*move) &= 0xffff0000
	(*move) |= Move(score)
}

func (move Move) Equals(m Move) bool {
	return (move & 0xffff0000) == (m & 0xffff0000)
}

func (move Move) String() string {
	ff := FileOf(move.FromSq())
	rf := RankOf(move.FromSq())
	ft := FileOf(move.ToSq())
	rt := RankOf(move.ToSq())

	if move.MoveType() == Promotion {
		pchar := 'q'

		switch move.Flag() {
		case KnightPromotion:
			pchar = 'n'
		case BishopPromotion:
			pchar = 'b'
		case RookPromotion:
			pchar = 'r'
		}
		return fmt.Sprintf("%c%c%c%c%c", ('a' + ff), ('1' + rf), ('a' + ft), ('1' + rt), pchar)
	}

	return fmt.Sprintf("%c%c%c%c", ('a' + ff), ('1' + rf), ('a' + ft), ('1' + rt))
}
