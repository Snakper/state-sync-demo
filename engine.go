package main

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	clients    map[string]*Client
	mainPlayer string
}

func NewGameEngine(main string) *Game {
	return &Game{
		clients:    map[string]*Client{},
		mainPlayer: main,
	}
}

func (g *Game) AddClient(c *Client) {
	f, err := os.Open("img.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	eimg := ebiten.NewImageFromImage(img)
	c.player.image = eimg
	g.clients[c.player.id] = c
	g.clients[c.player.id].player = &Player{
		id:          c.player.id,
		speed:       c.player.speed,
		pos:         c.player.pos,
		status:      c.player.status,
		destination: c.player.destination,
		target:      c.player.target,
		image:       c.player.image,
		client:      c.player.client,
	}
}

func (g *Game) Update() error {
	// 根据网络消息更新位置
	for _, c := range g.clients {
		msg := c.GetMessage()
		// 不开启预测，直接采用服务端位置
		if !c.forecast {
			if msg != nil {
				c.player.pos = msg.pos
			}
			continue
		}
		res := ProcessOne(msg, c.player, 60)
		c.player.pos = res.pos
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// 发送控制信息
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		c, ok := g.clients[g.mainPlayer]
		if ok {
			p := c.player
			pos := p.pos
			pos.X = float64(x)
			pos.Y = float64(y)
			c.Move(pos)
		}
	}
	// 渲染
	for _, c := range g.clients {
		p := c.player
		sizeX, sizeY := p.image.Size()
		op.GeoM.Translate(p.pos.X-float64(sizeX/2), p.pos.Y-float64(sizeY/2))
		screen.DrawImage(p.image, op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Hello,ebiten!\nTPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720 //窗口分辨率
}
