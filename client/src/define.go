package src

import (
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type Vec struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var ClientFrame = float64(60)

type Status int

var msg = make([]*ControlMsg, 0)
var lock = &sync.Mutex{}

const (
	STOP   = Status(0)
	MOVING = Status(1)
)

type Player struct {
	Id          string
	Speed       float64 // 速度 = 更号（x² + y²）
	Pos         Vec     // 当前位置
	Status      Status
	Destination Vec // 速度分量
	Target      Vec // 目的位置
	Image       *ebiten.Image
}

type ControlMsg struct {
	Id     string `json:"id"`
	Pos    Vec    `json:"pos"`
	Target Vec    `json:"target"`
	Index  int    `json:"index"`
}

func NewPlayer() *Player {
	p := &Player{
		Id:          "player1",
		Speed:       150,
		Pos:         Vec{0, 0},
		Status:      STOP,
		Destination: Vec{0, 0},
		Target:      Vec{0, 0},
	}
	return p
}

func ProcessOne(p *Player, frame float64) {
	disX := p.Target.X - p.Pos.X
	disY := p.Target.Y - p.Pos.Y
	total := math.Sqrt(math.Pow(disX, 2) + math.Pow(disY, 2))
	speed := 1 / frame * p.Speed
	per := speed / total
	if per > 1.0 {
		per = 1.0
	}
	speedX := per * disX
	speedY := per * disY
	p.Destination.X = speedX
	p.Destination.Y = speedY
	p.Status = MOVING

	p.Pos.X += p.Destination.X
	p.Pos.Y += p.Destination.Y
	// 会有误差
	nextX := p.Pos.X + p.Destination.X
	targetX := p.Target.X - p.Pos.X
	nextTargetX := p.Target.X - nextX
	nextY := p.Pos.Y + p.Destination.Y
	targetY := p.Target.Y - p.Pos.Y
	nextTargetY := p.Target.Y - nextY
	if targetX == 0 && targetY == 0 {
		p.Pos.X = p.Target.X
		p.Pos.Y = p.Target.Y
		p.Status = STOP
	}
	if (targetX > 0 && nextTargetX < 0) || (targetX < 0 && nextTargetX > 0) {
		p.Pos.X = p.Target.X
		p.Pos.Y = p.Target.Y
		p.Status = STOP
	}
	if (targetY > 0 && nextTargetY < 0) || (targetY < 0 && nextTargetY > 0) {
		p.Pos.X = p.Target.X
		p.Pos.Y = p.Target.Y
		p.Status = STOP
	}
}
