package src

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	lag         = time.Millisecond * 200
	ClientFrame = float64(60)
	ServerFrame = float64(10)
)

type Vec struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Status int

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
	Client      *Client
}

type ControlMsg struct {
	Id     string `json:"id"`
	Pos    Vec    `json:"pos"`
	Target Vec    `json:"target"`
	Index  int    `json:"index"`
}

var serverChan = make(chan *ControlMsg, 100)

func (p *Player) sendToServer(msg *ControlMsg) {
	go func() {
		time.Sleep(lag)
		serverChan <- msg
	}()
}
