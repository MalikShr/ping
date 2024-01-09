package engine

type MoveGen struct {
	pos  *BoardStruct
	list *MoveList
	side uint8
	quit bool
}

func GenerateAllMoves(pos *BoardStruct, list *MoveList, quit bool) {
	list.Count = 0

	if pos.SideToMove == White {
		for pieceNum := 0; pieceNum < pos.PieceNum[wPawn]; pieceNum++ {
			sq := pos.PieceList[wPawn][pieceNum]

			if pos.Pieces[sq+8] == Empty && quit {
				list.AddPawnMove(pos, sq, sq+8, White)

				if RankOf(sq) == R2 && pos.Pieces[sq+16] == Empty {
					move := MOVE(sq, (sq + 16), Empty, Empty, MFLAGPS)
					list.AddQuietMove(pos, move)
				}
			}

			GenerateAllPawnCaptureMoves(sq, pos, list, White)
		}

		if quit {
			// castling
			if pos.CastlePerm&wKCastle != 0 && pos.CastlePerm&wQCastle != 0 {
				if pos.Pieces[F1] == Empty && pos.Pieces[G1] == Empty {
					if !pos.SqAttacked(E1, Black) && !pos.SqAttacked(F1, Black) {
						list.AddQuietMove(pos, MOVE(E1, G1, Empty, Empty, MFLAGCA))
					}
				}

				if pos.Pieces[D1] == Empty && pos.Pieces[C1] == Empty && pos.Pieces[B1] == Empty {
					if !pos.SqAttacked(E1, Black) && !pos.SqAttacked(D1, Black) {
						list.AddQuietMove(pos, MOVE(E1, C1, Empty, Empty, MFLAGCA))
					}
				}
			}

		}

	} else {

		for pieceNum := 0; pieceNum < pos.PieceNum[bPawn]; pieceNum++ {
			sq := pos.PieceList[bPawn][pieceNum]

			if pos.Pieces[sq-8] == Empty && quit {
				list.AddPawnMove(pos, sq, sq-8, Black)
				if RankOf(sq) == R7 && pos.Pieces[sq-16] == Empty {
					list.AddQuietMove(pos, MOVE(sq, (sq-16), Empty, Empty, MFLAGPS))
				}
			}

			GenerateAllPawnCaptureMoves(sq, pos, list, Black)
		}

		if quit {
			// castling
			if pos.CastlePerm&bKCastle != 0 && pos.CastlePerm&bQCastle != 0 {
				if pos.Pieces[F8] == Empty && pos.Pieces[G8] == Empty {
					if !pos.SqAttacked(E8, White) && !pos.SqAttacked(F8, White) {
						list.AddQuietMove(pos, MOVE(E8, G8, Empty, Empty, MFLAGCA))
					}
				}

				if pos.Pieces[D8] == Empty && pos.Pieces[C8] == Empty && pos.Pieces[B8] == Empty {
					if !pos.SqAttacked(E8, White) && !pos.SqAttacked(D8, White) {
						list.AddQuietMove(pos, MOVE(E8, C8, Empty, Empty, MFLAGCA))
					}
				}
			}
		}
	}

	mg := MoveGen{
		pos:  pos,
		list: list,
		side: pos.SideToMove,
		quit: quit,
	}

	for sq := 0; sq < 64; sq++ {
		piece := pos.Pieces[sq]

		if piece == Empty || piece == wPawn || piece == bPawn {
			continue
		}

		if IsKn(piece) && PieceCol[piece] == pos.SideToMove {
			mg.GenKnightMoves(sq)
		}

		if IsBQ(piece) && PieceCol[piece] == pos.SideToMove {
			mg.GenBishopMoves(sq)
		}

		if IsRQ(piece) && PieceCol[piece] == pos.SideToMove {
			mg.GenRookMoves(sq)
		}

		if IsKi(piece) && PieceCol[piece] == pos.SideToMove {
			mg.GenKingMoves(sq)
		}
	}
}

func GenerateAllPawnCaptureMoves(sq int, pos *BoardStruct, list *MoveList, side uint8) {
	direction := 1

	if side == Black {
		direction = -1
	}

	if !SQOFFBOARD(sq+(7*direction)) && PieceCol[pos.Pieces[sq+(7*direction)]] == pos.SideToMove^1 && ValidPawnDelta1(sq, direction) {
		list.AddPawnCapMove(pos, sq, sq+(7*direction), pos.Pieces[sq+(7*direction)], side)
	}

	if !SQOFFBOARD(sq+(9*direction)) && PieceCol[pos.Pieces[sq+(9*direction)]] == pos.SideToMove^1 && ValidPawnDelta2(sq, direction) {
		list.AddPawnCapMove(pos, sq, sq+(9*direction), pos.Pieces[sq+(9*direction)], side)
	}

	if pos.EnPas != NoSq {
		if sq+(7*direction) == pos.EnPas && ValidPawnDelta1(sq, direction) {
			move := MOVE(sq, sq+(7*direction), Empty, Empty, MFLAGEP)
			list.AddEnPassantMove(pos, move)
		}
		if sq+(9*direction) == pos.EnPas && ValidPawnDelta2(sq, direction) {
			move := MOVE(sq, sq+(9*direction), Empty, Empty, MFLAGEP)
			list.AddEnPassantMove(pos, move)
		}
	}
}

func (mg MoveGen) GenBishopMoves(sq int) {
	rank := RankOf(sq)
	file := FileOf(sq)

	for i := 1; i < 8; i++ {
		f := file - i
		r := rank - i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if mg.pos.Pieces[targetSq] != Empty {
			if PieceCol[mg.pos.Pieces[targetSq]] == mg.side^1 {
				mg.list.AddCaptureMove(mg.pos, MOVE(sq, targetSq, mg.pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if mg.quit {
			mg.list.AddQuietMove(mg.pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	for i := 1; i < 8; i++ {
		f := file + i
		r := rank - i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if mg.pos.Pieces[targetSq] != Empty {
			if PieceCol[mg.pos.Pieces[targetSq]] == mg.side^1 {
				mg.list.AddCaptureMove(mg.pos, MOVE(sq, targetSq, mg.pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if mg.quit {
			mg.list.AddQuietMove(mg.pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	for i := 1; i < 8; i++ {
		f := file - i
		r := rank + i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if mg.pos.Pieces[targetSq] != Empty {
			if PieceCol[mg.pos.Pieces[targetSq]] == mg.side^1 {
				mg.list.AddCaptureMove(mg.pos, MOVE(sq, targetSq, mg.pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if mg.quit {
			mg.list.AddQuietMove(mg.pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	for i := 1; i < 8; i++ {
		f := file + i
		r := rank + i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if mg.pos.Pieces[targetSq] != Empty {
			if PieceCol[mg.pos.Pieces[targetSq]] == mg.side^1 {
				mg.list.AddCaptureMove(mg.pos, MOVE(sq, targetSq, mg.pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if mg.quit {
			mg.list.AddQuietMove(mg.pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}
}

func (mg MoveGen) GenRookMoves(sq int) {
	file := FileOf(sq)
	rank := RankOf(sq)

	// Check for horizontal moves
	for i := 1; i < 8; i++ {
		f := file + i
		r := rank

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if mg.pos.Pieces[targetSq] != Empty {
			if PieceCol[mg.pos.Pieces[targetSq]] == mg.side^1 {
				mg.list.AddCaptureMove(mg.pos, MOVE(sq, targetSq, mg.pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if mg.quit {
			mg.list.AddQuietMove(mg.pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	for i := 1; i < 8; i++ {
		f := file - i
		r := rank

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if mg.pos.Pieces[targetSq] != Empty {
			if PieceCol[mg.pos.Pieces[targetSq]] == mg.side^1 {
				mg.list.AddCaptureMove(mg.pos, MOVE(sq, targetSq, mg.pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if mg.quit {
			mg.list.AddQuietMove(mg.pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	// Check for vertical moves
	for i := 1; i < 8; i++ {
		f := file
		r := rank + i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if mg.pos.Pieces[targetSq] != Empty {
			if PieceCol[mg.pos.Pieces[targetSq]] == mg.side^1 {
				mg.list.AddCaptureMove(mg.pos, MOVE(sq, targetSq, mg.pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if mg.quit {
			mg.list.AddQuietMove(mg.pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	for i := 1; i < 8; i++ {
		f := file
		r := rank - i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if mg.pos.Pieces[targetSq] != Empty {
			if PieceCol[mg.pos.Pieces[targetSq]] == mg.side^1 {
				mg.list.AddCaptureMove(mg.pos, MOVE(sq, targetSq, mg.pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if mg.quit {
			mg.list.AddQuietMove(mg.pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}
}

func (mg MoveGen) GenKnightMoves(sq int) {
	file := FileOf(sq)
	rank := RankOf(sq)

	// up-left
	mg.AddMove(sq, file-1, rank-2)

	//up-right
	mg.AddMove(sq, file+1, rank-2)

	// down-left
	mg.AddMove(sq, file-1, rank+2)

	// down-right
	mg.AddMove(sq, file+1, rank+2)

	//left-up
	mg.AddMove(sq, file-2, rank-1)

	//left-down
	mg.AddMove(sq, file-2, rank+1)

	//right-up
	mg.AddMove(sq, file+2, rank-1)

	//right-down
	mg.AddMove(sq, file+2, rank+1)
}

func (mg MoveGen) GenKingMoves(sq int) {
	file := FileOf(sq)
	rank := RankOf(sq)

	// up
	mg.AddMove(sq, file, rank-1)

	// left
	mg.AddMove(sq, file-1, rank)

	//down
	mg.AddMove(sq, file, rank+1)

	//right
	mg.AddMove(sq, file+1, rank)

	// up-left
	mg.AddMove(sq, file-1, rank-1)

	// up-right
	mg.AddMove(sq, file+1, rank-1)

	// down-left
	mg.AddMove(sq, file-1, rank+1)

	//down-right
	mg.AddMove(sq, file+1, rank+1)
}

func (mg MoveGen) AddMove(sq int, f int, r int) {
	targetSq := FR2SQ(f, r)

	if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
		return
	}

	if mg.pos.Pieces[targetSq] != Empty {
		if PieceCol[mg.pos.Pieces[targetSq]] == mg.side^1 {
			mg.list.AddCaptureMove(mg.pos, MOVE(sq, targetSq, mg.pos.Pieces[targetSq], Empty, 0))
		}
		return
	}

	if mg.quit {
		mg.list.AddQuietMove(mg.pos, MOVE(sq, targetSq, Empty, Empty, 0))
	}
}

func LastBordIndex(f int, r int) bool {
	return r < 0 || r > 7 || f < 0 || f > 7
}
