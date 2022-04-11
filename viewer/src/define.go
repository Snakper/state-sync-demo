package src

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type Vec struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Status int

var msg = make([]ControlMsg, 0)
var lock = &sync.Mutex{}

const (
	STOP   = Status(0)
	MOVING = Status(1)
)

type Player struct {
	id          string
	speed       float64 // 速度 = 更号（x² + y²）
	pos         Vec     // 当前位置
	status      Status
	destination Vec // 速度分量
	target      Vec // 目的位置
	image       *ebiten.Image
}

type ControlMsg struct {
	Id     string `json:"id"`
	Pos    Vec    `json:"pos"`
	Target Vec    `json:"target"`
	Index  int    `json:"index"`
}
