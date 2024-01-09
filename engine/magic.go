package engine

type MoveGen struct {
	pos  *BoardStruct
	list *MoveList
	side uint8
	quit bool
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
