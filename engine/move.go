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

	NoMove = Move(0)
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

func ParseMove(ptrChar string, pos *BoardStruct) Move {
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

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		move := list.Moves[moveNum]

		if move.FromSq() == from && move.ToSq() == to {
			if move.MoveType() == Promotion {
				moveFlag := move.Flag()

				if moveFlag == KnightPromotion && ptrChar[4] == 'n' {
					return move
				} else if moveFlag == BishopPromotion && ptrChar[4] == 'b' {
					return move
				} else if moveFlag == RookPromotion && ptrChar[4] == 'r' {
					return move
				} else if moveFlag == QueenPromotion && ptrChar[4] == 'q' {
					return move
				}
				continue
			}
			return move
		}
	}

	return NoMove
}
