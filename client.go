package main

type Client struct {
	player     *Player
	playerChan chan *ControlMsg
	forecast   bool
}

func NewClient() *Client {
	client := &Client{
		playerChan: make(chan *ControlMsg, 100),
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

func (c *Client) Connect(s *Server) {
	s.Clients[c.player.id] = c
}

func (c *Client) Move(target Vec) {
	msg := &ControlMsg{
		id:     c.player.id,
		pos:    Vec{},
		target: target,
	}
	c.player.sendToServer(msg)
}
