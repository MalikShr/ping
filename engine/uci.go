package engine

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type UCIInterface struct {
	OpeningBook   map[uint64][]PolyEntry
	OptionUseBook bool
}

func (inter *UCIInterface) ParsePosition(cmd string, pos *BoardStruct) {

	cmd = strings.TrimPrefix(cmd, "position")
	cmd = strings.TrimPrefix(cmd, " ")

	parts := strings.Split(cmd, "moves")

	if len(cmd) == 0 || len(parts) > 2 {
		err := fmt.Errorf("%v wrong length=%v", parts, len(parts))
		fmt.Println("info string Error", fmt.Sprint(err))
		return
	}

	alt := strings.Split(parts[0], " ")
	alt[0] = strings.TrimSpace(alt[0])
	if alt[0] == "startpos" {
		parts[0] = FENStart
	} else if alt[0] == "fen" {
		parts[0] = strings.TrimSpace(strings.TrimPrefix(parts[0], "fen"))
	} else {
		err := fmt.Errorf("%#v must be %#v or %#v", alt[0], "fen", "startpos")
		fmt.Println("info string Error", err.Error())
		return
	}

	pos.ParseFen(parts[0])

	if len(parts) == 2 {
		parts[1] = strings.ToLower(strings.TrimSpace(parts[1]))

		mvs := strings.Fields(strings.ToLower(parts[1]))

		for i := 0; i < len(mvs); i++ {
			move := ParseMove(mvs[i], pos)

			if move == NoMove {
				break
			}

			pos.DoMove(move)
			pos.Ply = 0
		}
	}
}

func Uci() {
	quit := false
	reader := bufio.NewReader(os.Stdin)

	var pos BoardStruct
	var info Search
	inter := UCIInterface{}
	info.TT = make(map[uint64]int)

	inter.handleUci()

	//
	pos.ParseFen("8/8/8/2pP4/8/8/8/8 w - c6 0 1")
	var list MoveList
	GenPawnMoves(D5, &pos, &list, true)
	//

	var err error
	inter.OpeningBook, err = LoadPolyglotFile("../book/baron30.bin")

	if err == nil {
		inter.OptionUseBook = true
	} else {
		fmt.Printf("ERROR \n")
	}

	for !quit {
		cmd, _ := reader.ReadString('\n')

		words := strings.Fields(cmd)

		switch words[0] {
		case "uci":
			inter.handleUci()
		case "isready":
			fmt.Println("readyok")
		case "position":
			inter.ParsePosition(cmd, &pos)
		case "setoption":
			inter.HandleSetOption(cmd)
		case "ucinewgame":
			inter.ParsePosition("position startpos\n", &pos)
		case "go":
			inter.handleGo(cmd, &info, &pos)
		case "perft":
			if len(words) >= 2 {
				depth, _ := strconv.Atoi(words[1])
				PerftTest(depth, &pos)
			}
		case "print":
			fmt.Println(pos.String())
		case "poly":
			fmt.Printf("Poly Key: %d\n", PolyKeyFromBoard(&pos))
			fmt.Printf("Poly Key HEX: %x\n", PolyKeyFromBoard(&pos))
		case "quit":
			quit = true
		default:
			fmt.Println("info string unknown cmd ", cmd)
		}

	}
}

func (inter *UCIInterface) handleUci() {
	fmt.Println("id name Ping")
	fmt.Println("id author MalikShr")

	fmt.Println("option name Hash type spin default 128 min 16 max 1024")
	fmt.Println("option name Threads type spin default 1 min 1 max 16")
	fmt.Print("option name UseBook type check default false\n")
	fmt.Println("uciok")
}

func (inter *UCIInterface) handleGo(cmd string, info *Search, pos *BoardStruct) {
	// If Opening book is enabled a random move will be logged instead of searching

	if inter.OptionUseBook {
		if inter.OpeningBook[PolyKeyFromBoard(pos)] != nil {
			entries := inter.OpeningBook[PolyKeyFromBoard(pos)]

			bestMove := entries[rand.Intn(len(entries))].Move
			moveDelay := time.Duration(rand.Intn(2500-500) + 500)

			if ParseMove(bestMove, pos) != NoMove {
				time.Sleep(moveDelay * time.Millisecond)
				fmt.Printf("bestmove %s\n", bestMove)
				return
			}
		}
	}

	cmd = strings.TrimPrefix(cmd, "go")
	cmd = strings.TrimPrefix(cmd, " ")
	words := strings.Fields(cmd)

	depth := -1
	movestogo := 30
	movetime := -1
	gameTime := -1
	inc := 0
	info.Timeset = false

	for i := 0; i < len(words)-1; i++ {
		switch words[i] {

		case "infinite":
		case "binc":
			if pos.SideToMove == Black {
				inc, _ = strconv.Atoi(words[i+1])
			}
		case "winc":
			if pos.SideToMove == White {
				inc, _ = strconv.Atoi(words[i+1])
			}
		case "wtime":
			if pos.SideToMove == White {
				gameTime, _ = strconv.Atoi(words[i+1])
			}
		case "btime":
			if pos.SideToMove == Black {
				gameTime, _ = strconv.Atoi(words[i+1])
			}
		case "movestogo":
			movestogo, _ = strconv.Atoi(words[i+1])
		case "movetime":
			movetime, _ = strconv.Atoi(words[i+1])
		case "depth":
			depth, _ = strconv.Atoi(words[i+1])
		}

		if movetime != -1 {
			gameTime = movetime
			movestogo = 1
		}
	}

	info.Starttime = time.Now().UnixMilli()
	info.Depth = depth

	if gameTime != -1 {
		info.Timeset = true
		gameTime /= movestogo
		gameTime -= 50
		info.Stoptime = info.Starttime + int64(gameTime) + int64(inc)
	}

	if depth == -1 {
		info.Depth = MaxDepth
	}

	fmt.Printf("time:%d start:%d stop:%d depth:%d timeset:%t\n", gameTime, info.Starttime, info.Stoptime, info.Depth, info.Timeset)
	SearchPosition(pos, info)
}

func (inter *UCIInterface) HandleSetOption(cmd string) {
	fields := strings.Fields(cmd)
	var option, value string
	parsingWhat := ""

	for _, field := range fields {
		if field == "name" {
			parsingWhat = "name"
		} else if field == "value" {
			parsingWhat = "value"
		} else if parsingWhat == "name" {
			option += field + " "
		} else if parsingWhat == "value" {
			value += field + " "
		}
	}

	option = strings.TrimSuffix(option, " ")
	value = strings.TrimSuffix(value, " ")

	switch option {
	case "UseBook":
		if value == "true" {
			inter.OptionUseBook = true
		} else if value == "false" {
			inter.OptionUseBook = false
		}
	case "BookPath":
		var err error
		inter.OpeningBook, err = LoadPolyglotFile(value)

		if err == nil {
			fmt.Println("Opening book loaded...")
		} else {
			fmt.Println("Failed to load opening book...")
		}
	}
}
