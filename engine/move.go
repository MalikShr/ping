package engine

import "fmt"

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

func PrMove(move int) string {
	var MvStr string

	ff := FileOf(FromSq(move))
	rf := RankOf(FromSq(move))
	ft := FileOf(ToSq(move))
	rt := RankOf(ToSq(move))

	promoted := Promoted(move)

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
