package main

import (
	"ping/engine"
)

func Init() {
	engine.InitBitMasks()
	engine.InitHashKeys()
	engine.InitTables()
	engine.InitEvalMasks()
	engine.InitMvvLva()
	engine.InitMagic()

}

func main() {
	Init()

	engine.Uci()
}
