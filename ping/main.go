package main

import (
	"ping/engine"
)

func Init() {
	engine.InitBitMasks()
	engine.InitHashKeys()
	engine.InitEvalMasks()
	engine.InitMvvLva()
	engine.InitTables()
}

func main() {
	Init()

	engine.Uci()
}
