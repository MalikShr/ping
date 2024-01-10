package engine

import "fmt"

type Bitboard uint64

var BitTable = [64]int{
	63, 30, 3, 32, 25, 41, 22, 33, 15, 50, 42, 13, 11, 53, 19, 34, 61, 29, 2,
	51, 21, 43, 45, 10, 18, 47, 1, 54, 9, 57, 0, 35, 62, 31, 40, 4, 49, 5, 52,
	26, 60, 6, 23, 44, 46, 27, 56, 16, 7, 39, 48, 24, 59, 14, 12, 55, 38, 28,
	58, 20, 37, 17, 36, 8,
}

func (bb *Bitboard) PopBit() int {
	b := *bb ^ (*bb - 1)
	fold := uint32((b & 0xffffffff) ^ (b >> 32))
	*bb &= (*bb - 1)
	return BitTable[(fold*0x783a9b23)>>26]
}

func (b Bitboard) CountBits() int {
	r := 0
	for ; b != 0; r++ {
		b &= b - 1
	}
	return r
}

func (bb *Bitboard) CLEARBIT(sq int) {
	*bb &= ClearMask[sq]
}

func (bb *Bitboard) SETBIT(sq int) {
	*bb |= SetMask[sq]
}

func (bb Bitboard) String() {
	var shiftMe Bitboard = 1

	var f, sq int

	fmt.Println()

	for rank := R8; rank >= R1; rank-- {
		for f = FA; f <= FH; f++ {
			sq = FR2SQ(f, rank)

			if (shiftMe<<uint(sq))&bb != 0 {
				fmt.Print("X")
			} else {
				fmt.Print("-")
			}
		}
		fmt.Println()
	}

	fmt.Printf("\n\n")
}
