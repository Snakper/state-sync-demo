package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Vec struct {
	X float64
	Y float64
}

type Status int

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
	client      *Client
}

type ControlMsg struct {
	id     string
	pos    Vec
	target Vec
}

var serverChan = make(chan *ControlMsg, 100)

var lag = time.Millisecond * 1

func (p *Player) sendToServer(msg *ControlMsg) {
	go func() {
		time.Sleep(lag)
		serverChan <- msg
	}()
}
