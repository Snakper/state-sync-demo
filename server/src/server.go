package src

import (
	"fmt"
	"math"
	"time"
)

type Server struct {
	Msg     []*ControlMsg
	tempMsg []ControlMsg
	Clients map[string]*Client
}

func NewServer() *Server {
	return &Server{
		Msg:     []*ControlMsg{},
		tempMsg: []ControlMsg{},
		Clients: map[string]*Client{},
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

func (s *Server) sendToClient(msg []ControlMsg) {
	go func() {
		time.Sleep(lag)
		SendToNetwork(msg)
		for _, cli := range s.Clients {
			cli.playerChan <- msg
		}
		s.tempMsg = s.tempMsg[:0]
	}()
}

func (s *Server) Run() {
	go func() {
		hz := time.Duration(1000 / ServerFrame)
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
	clientMsg := map[string]*ControlMsg{} // 控制帧校验
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
		speed := 1 / float64(ServerFrame) * p.Speed
		per := speed / total
		if per > 1.0 {
			per = 1.0
		}
		speedX := per * disX
		speedY := per * disY
		p.Destination.X = speedX
		p.Destination.Y = speedY
		p.Status = MOVING
		if msg.Index != 0 {
			clientMsg[p.Id] = msg
		}
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
		msg := ControlMsg{
			Id:     p.Id,
			Pos:    p.Pos,
			Target: p.Target,
		}
		// 对账
		if buf, ok := clientMsg[p.Id]; ok {
			msg.Index = buf.Index
		}
		s.tempMsg = append(s.tempMsg, msg)
		fmt.Println("Server Player Move:", msg)
	}
	s.sendToClient(s.tempMsg)
	s.Msg = s.Msg[:0]
}
