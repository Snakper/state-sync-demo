package src

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
		SendToNetwork(*msg)
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
		c, ok := s.Clients[msg.Id]
		if !ok {
			continue
		}
		p := c.player
		p.Target = msg.Target
		disX := p.Target.X - p.Pos.X
		disY := p.Target.Y - p.Pos.Y
		total := math.Sqrt(math.Pow(disX, 2) + math.Pow(disY, 2))
		speed := 1 / float64(s.Frame) * p.Speed
		per := speed / total
		if per > 1.0 {
			per = 1.0
		}
		speedX := per * disX
		speedY := per * disY
		p.Destination.X = speedX
		p.Destination.Y = speedY
		p.Status = MOVING
		s.MsgBuffer[p.Id] = *msg
	}

	for _, c := range s.Clients {
		p := c.player
		if p.Status == STOP {
			continue
		}
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
		msg := &ControlMsg{
			Id:     p.Id,
			Pos:    p.Pos,
			Target: p.Target,
		}
		// 移动完毕，对账
		if p.Status == STOP {
			if buf, ok := s.MsgBuffer[p.Id]; ok {
				msg.Index = buf.Index
				delete(s.MsgBuffer, p.Id)
			}
		}
		s.sendToClient(c, msg)
		fmt.Println("Server Player Move:", msg)
	}

	s.Msg = s.Msg[:0]
}
