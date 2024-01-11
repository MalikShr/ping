package engine

type Move struct {
	Move  int
	Score int
}

var MFLAGEP int = 0x40000
var MFLAGPS int = 0x80000
var MFLAGCA int = 0x1000000

var MFLAGCAP int = 0x7C000
var MFLAGPROM int = 0xF00000

func ToSq(m int) int {
	return (m >> 7) & 0x7F
}

func Captured(m int) uint8 {
	return uint8((m >> 14)) & 0xF
}

func Promoted(m int) uint8 {
	return uint8((m >> 20) & 0xF)
}

func MOVE(from int, to int, capture uint8, promotion uint8, flag int) int {
	return from | (to << 7) | (int(capture) << 14) | (int(promotion) << 20) | flag
}
