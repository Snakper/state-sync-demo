package main

import (
	"fmt"
	"math"
	"time"
)

type Server struct {
	Msg       []*ControlMsg
	Frame     float64
	Clients   map[string]*Client
	MsgBuffer map[string]ControlMsg
}

func NewServer(frame float64) *Server {
	return &Server{
		Msg:       []*ControlMsg{},
		Frame:     frame,
		Clients:   map[string]*Client{},
		MsgBuffer: map[string]ControlMsg{},
	}
}

func (s *Server) GetMessage() {
	for {
		select {
		case m := <-serverChan:
			s.Msg = append(s.Msg, m)
		default:
			return
		}
	}
}

func (s *Server) sendToClient(c *Client, msg *ControlMsg) {
	go func() {
		time.Sleep(lag)
		c.playerChan <- msg
	}()
}

func (s *Server) Run() {
	go func() {
		hz := time.Duration(1000 / s.Frame)
		timer := time.NewTimer(time.Millisecond * hz)
		for {
			select {
			case <-timer.C:
				timer.Reset(time.Millisecond * hz)
				s.process()
			default:
			}
		}
	}()
}

func (s *Server) process() {
	s.GetMessage()
	for _, msg := range s.Msg {
		c, ok := s.Clients[msg.id]
		if !ok {
			continue
		}
		p := c.player
		p.target = msg.target
		disX := p.target.X - p.pos.X
		disY := p.target.Y - p.pos.Y
		total := math.Sqrt(math.Pow(disX, 2) + math.Pow(disY, 2))
		speed := 1 / float64(s.Frame) * p.speed
		per := speed / total
		if per > 1.0 {
			per = 1.0
		}
		speedX := per * disX
		speedY := per * disY
		p.destination.X = speedX
		p.destination.Y = speedY
		p.status = MOVING
		s.MsgBuffer[p.id] = *msg
	}

	for _, c := range s.Clients {
		p := c.player
		if p.status == STOP {
			continue
		}
		p.pos.X += p.destination.X
		p.pos.Y += p.destination.Y

		// 会有误差
		nextX := p.pos.X + p.destination.X
		targetX := p.target.X - p.pos.X
		nextTargetX := p.target.X - nextX
		nextY := p.pos.Y + p.destination.Y
		targetY := p.target.Y - p.pos.Y
		nextTargetY := p.target.Y - nextY
		if targetX == 0 && targetY == 0 {
			p.pos.X = p.target.X
			p.pos.Y = p.target.Y
			p.status = STOP
		}
		if (targetX > 0 && nextTargetX < 0) || (targetX < 0 && nextTargetX > 0) {
			p.pos.X = p.target.X
			p.pos.Y = p.target.Y
			p.status = STOP
		}
		if (targetY > 0 && nextTargetY < 0) || (targetY < 0 && nextTargetY > 0) {
			p.pos.X = p.target.X
			p.pos.Y = p.target.Y
			p.status = STOP
		}
		msg := &ControlMsg{
			id:     p.id,
			pos:    p.pos,
			target: p.target,
		}
		// 移动完毕，对账
		if p.status == STOP {
			if buf, ok := s.MsgBuffer[p.id]; ok {
				msg.index = buf.index
				delete(s.MsgBuffer, p.id)
			}
		}
		s.sendToClient(c, msg)
		fmt.Println("Server Player Move:", msg)
	}

	s.Msg = s.Msg[:0]
}
