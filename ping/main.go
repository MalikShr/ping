package main

import (
	"ping/engine"
)

func Init() {
	engine.InitBitMasks()
	engine.InitHashKeys()
	engine.InitEvalMasks()
	engine.InitMvvLva()
	engine.InitAttacks()
}

func main() {
	Init()

	engine.Uci()

}
