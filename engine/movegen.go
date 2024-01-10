package engine

func GenerateAllMoves(pos *BoardStruct, list *MoveList, quiet bool) {
	list.Count = 0

	if pos.SideToMove == White {
		for pieceNum := 0; pieceNum < pos.PieceNum[wPawn]; pieceNum++ {
			sq := pos.PieceList[wPawn][pieceNum]

			if pos.Pieces[sq+8] == Empty && quiet {
				list.AddPawnMove(pos, sq, sq+8, White)

				if RankOf(sq) == R2 && pos.Pieces[sq+16] == Empty {
					move := MOVE(sq, (sq + 16), Empty, Empty, MFLAGPS)
					list.AddQuietMove(pos, move)
				}
			}

			GenPawnAttacks(sq, pos.Sides[pos.SideToMove], pos.Sides[pos.SideToMove^1], pos, list, pos.SideToMove, false)
		}

		if quiet {
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

			if pos.Pieces[sq-8] == Empty && quiet {
				list.AddPawnMove(pos, sq, sq-8, Black)
				if RankOf(sq) == R7 && pos.Pieces[sq-16] == Empty {
					list.AddQuietMove(pos, MOVE(sq, (sq-16), Empty, Empty, MFLAGPS))
				}
			}

			GenPawnAttacks(sq, pos.Sides[pos.SideToMove], pos.Sides[pos.SideToMove^1], pos, list, pos.SideToMove, false)
		}

		if quiet {
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

	var pieceIndex uint8 = 0

	for pieceIndex < 10 {
		piece := NonPawnPieces[pieceIndex]

		for pieceNum := 0; pieceNum < pos.PieceNum[piece]; pieceNum++ {
			sq := pos.PieceList[piece][pieceNum]

			if IsKn(piece) && PieceCol[piece] == pos.SideToMove {
				GenKnightAttacks(sq, pos.Sides[pos.SideToMove], pos.Sides[pos.SideToMove^1], pos, list, quiet)
			}

			if IsBQ(piece) && PieceCol[piece] == pos.SideToMove {
				GenBishopAttacks(sq, pos.Sides[pos.SideToMove], pos.Sides[pos.SideToMove^1], pos, list, quiet)
			}

			if IsRQ(piece) && PieceCol[piece] == pos.SideToMove {
				GenRookAttacks(sq, pos.Sides[pos.SideToMove], pos.Sides[pos.SideToMove^1], pos, list, quiet)
			}

			if IsKi(piece) && PieceCol[piece] == pos.SideToMove {
				GenKingAttacks(sq, pos.Sides[pos.SideToMove], pos.Sides[pos.SideToMove^1], pos, list, quiet)
			}
		}
		pieceIndex++
	}

}
