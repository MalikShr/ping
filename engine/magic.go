package engine

func GenBishopMoves(sq int, pos *BoardStruct, list *MoveList, quit bool) {
	rank := RankOf(sq)
	file := FileOf(sq)

	for i := 1; i < 8; i++ {
		f := file - i
		r := rank - i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if pos.Pieces[targetSq] != Empty {
			if PieceCol[pos.Pieces[targetSq]] == pos.SideToMove^1 {
				list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if quit {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	for i := 1; i < 8; i++ {
		f := file + i
		r := rank - i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if pos.Pieces[targetSq] != Empty {
			if PieceCol[pos.Pieces[targetSq]] == pos.SideToMove^1 {
				list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if quit {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	for i := 1; i < 8; i++ {
		f := file - i
		r := rank + i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if pos.Pieces[targetSq] != Empty {
			if PieceCol[pos.Pieces[targetSq]] == pos.SideToMove^1 {
				list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if quit {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	for i := 1; i < 8; i++ {
		f := file + i
		r := rank + i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if pos.Pieces[targetSq] != Empty {
			if PieceCol[pos.Pieces[targetSq]] == pos.SideToMove^1 {
				list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if quit {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}
}

func GenRookMoves(sq int, pos *BoardStruct, list *MoveList, quit bool) {
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

		if pos.Pieces[targetSq] != Empty {
			if PieceCol[pos.Pieces[targetSq]] == pos.SideToMove^1 {
				list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if quit {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	for i := 1; i < 8; i++ {
		f := file - i
		r := rank

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if pos.Pieces[targetSq] != Empty {
			if PieceCol[pos.Pieces[targetSq]] == pos.SideToMove^1 {
				list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if quit {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
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

		if pos.Pieces[targetSq] != Empty {
			if PieceCol[pos.Pieces[targetSq]] == pos.SideToMove^1 {
				list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if quit {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}

	for i := 1; i < 8; i++ {
		f := file
		r := rank - i

		targetSq := FR2SQ(f, r)

		if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
			break
		}

		if pos.Pieces[targetSq] != Empty {
			if PieceCol[pos.Pieces[targetSq]] == pos.SideToMove^1 {
				list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Pieces[targetSq], Empty, 0))
			}
			break
		}

		if quit {
			list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
		}
	}
}

func GenKnightMoves(sq int, pos *BoardStruct, list *MoveList, quit bool) {
	file := FileOf(sq)
	rank := RankOf(sq)

	// up-left
	AddMove(file-1, rank-2, sq, pos, list, quit)

	//up-right
	AddMove(file+1, rank-2, sq, pos, list, quit)

	// down-left
	AddMove(file-1, rank+2, sq, pos, list, quit)

	// down-right
	AddMove(file+1, rank+2, sq, pos, list, quit)

	//left-up
	AddMove(file-2, rank-1, sq, pos, list, quit)

	//left-down
	AddMove(file-2, rank+1, sq, pos, list, quit)

	//right-up
	AddMove(file+2, rank-1, sq, pos, list, quit)

	//right-down
	AddMove(file+2, rank+1, sq, pos, list, quit)
}

func GenKingMoves(sq int, pos *BoardStruct, list *MoveList, quit bool) {
	file := FileOf(sq)
	rank := RankOf(sq)

	// up
	AddMove(file, rank-1, sq, pos, list, quit)

	// left
	AddMove(file-1, rank, sq, pos, list, quit)

	//down
	AddMove(file, rank+1, sq, pos, list, quit)

	//right
	AddMove(file+1, rank, sq, pos, list, quit)

	// up-left
	AddMove(file-1, rank-1, sq, pos, list, quit)

	// up-right
	AddMove(file+1, rank-1, sq, pos, list, quit)

	// down-left
	AddMove(file-1, rank+1, sq, pos, list, quit)

	//down-right
	AddMove(file+1, rank+1, sq, pos, list, quit)
}

func AddMove(f int, r int, sq int, pos *BoardStruct, list *MoveList, quit bool) {
	targetSq := FR2SQ(f, r)

	if SQOFFBOARD(targetSq) || LastBordIndex(f, r) {
		return
	}

	if pos.Pieces[targetSq] != Empty {
		if PieceCol[pos.Pieces[targetSq]] == pos.SideToMove^1 {
			list.AddCaptureMove(pos, MOVE(sq, targetSq, pos.Pieces[targetSq], Empty, 0))
		}
		return
	}

	if quit {
		list.AddQuietMove(pos, MOVE(sq, targetSq, Empty, Empty, 0))
	}
}

func LastBordIndex(f int, r int) bool {
	return r < 0 || r > 7 || f < 0 || f > 7
}
