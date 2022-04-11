package src

import (
	"math"
)

type Client struct {
	player         *Player
	playerChan     chan []ControlMsg
	forecast       bool
	reconciliation bool
	interpolation  bool
	ControlBuffer  map[int]ControlMsg
	index          int
}

func NewClient() *Client {
	client := &Client{
		playerChan:    make(chan []ControlMsg, 100),
		ControlBuffer: map[int]ControlMsg{},
		index:         1,
	}
	return client
}

func (c *Client) NewPlayer() *Player {
	p := &Player{
		Id:          "player1",
		Speed:       150,
		Pos:         Vec{0, 0},
		Status:      STOP,
		Destination: Vec{0, 0},
		Target:      Vec{0, 0},
		Client:      c,
	}
	c.player = p
	return p
}

func (c *Client) DeepCopyPlayer(p *Player) *Player {
	newP := &Player{
		Id:          p.Id,
		Speed:       p.Speed,
		Pos:         p.Pos,
		Status:      p.Status,
		Destination: p.Destination,
		Target:      p.Target,
		Image:       p.Image,
		Client:      p.Client,
	}
	return newP
}

func (c *Client) DeepCopyClient(cli *Client) *Client {
	client := &Client{
		player:        cli.DeepCopyPlayer(cli.player),
		playerChan:    cli.playerChan,
		forecast:      cli.forecast,
		ControlBuffer: cli.ControlBuffer,
		index:         cli.index,
	}
	return client
}

func (c *Client) GetMessage() []*ControlMsg {
	for {
		select {
		case m := <-c.playerChan:
			ms := make([]*ControlMsg, 0)
			for idx := range m {
				ms = append(ms, &m[idx])
			}
			return ms
		default:
			return nil
		}
	}
}

func (c *Client) SetForecast(open bool) {
	c.forecast = open
}

func (c *Client) SetReconciliation(open bool) {
	c.reconciliation = open
}

func (c *Client) SetInterpolation(open bool) {
	c.interpolation = open
}

func (c *Client) Connect(s *Server) {
	s.Clients[c.player.Id] = c.DeepCopyClient(c)
}

func (c *Client) Move(target Vec) {
	msg := &ControlMsg{
		Id:     c.player.Id,
		Target: target,
		Index:  c.index,
	}
	c.ControlBuffer[c.index] = *msg
	if c.index > 100000 {
		c.index = 0
	}
	c.index++
	c.player.sendToServer(msg)
}

func ProcessOne(p *Player, frame float64) *ControlMsg {
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
	res := &ControlMsg{
		Id:     p.Id,
		Pos:    p.Pos,
		Target: p.Target,
	}
	return res
}
