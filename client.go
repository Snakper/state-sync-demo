package main

import "math"

type Client struct {
	player         *Player
	playerChan     chan *ControlMsg
	forecast       bool
	reconciliation bool
	ControlBuffer  map[int]ControlMsg
	index          int
}

func NewClient() *Client {
	client := &Client{
		playerChan:    make(chan *ControlMsg, 100),
		ControlBuffer: map[int]ControlMsg{},
		index:         1,
	}
	return client
}

func (c *Client) NewPlayer() *Player {
	p := &Player{
		id:          "player1",
		speed:       150,
		pos:         Vec{0, 0},
		status:      STOP,
		destination: Vec{0, 0},
		target:      Vec{0, 0},
		client:      c,
	}
	c.player = p
	return p
}

func (c *Client) DeepCopyPlayer(p *Player) *Player {
	newP := &Player{
		id:          p.id,
		speed:       p.speed,
		pos:         p.pos,
		status:      p.status,
		destination: p.destination,
		target:      p.target,
		image:       p.image,
		client:      p.client,
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

func (c *Client) GetMessage() *ControlMsg {
	for {
		select {
		case m := <-c.playerChan:
			return m
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

func (c *Client) Connect(s *Server) {
	s.Clients[c.player.id] = c.DeepCopyClient(c)
}

func (c *Client) Move(target Vec) {
	msg := &ControlMsg{
		id:     c.player.id,
		target: target,
		index:  c.index,
	}
	c.ControlBuffer[c.index] = *msg
	if c.index > 100000 {
		c.index = 0
	}
	c.index++
	c.player.sendToServer(msg)
}

func ProcessOne(p *Player, frame float64) *ControlMsg {
	disX := p.target.X - p.pos.X
	disY := p.target.Y - p.pos.Y
	total := math.Sqrt(math.Pow(disX, 2) + math.Pow(disY, 2))
	speed := 1 / frame * p.speed
	per := speed / total
	if per > 1.0 {
		per = 1.0
	}
	speedX := per * disX
	speedY := per * disY
	p.destination.X = speedX
	p.destination.Y = speedY
	p.status = MOVING

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
	res := &ControlMsg{
		id:     p.id,
		pos:    p.pos,
		target: p.target,
	}
	return res
}
