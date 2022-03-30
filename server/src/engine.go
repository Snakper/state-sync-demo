package src

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

func (g *Game) OpenForecast() {
	for _, c := range g.clients {
		c.SetForecast(true)
	}
}

func (g *Game) OpenReconciliation() {
	for _, c := range g.clients {
		c.SetReconciliation(true)
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
	c.player.Image = eimg
	g.clients[c.player.Id] = c.DeepCopyClient(c)
}

func (g *Game) Update() error {
	// 根据网络消息更新位置
	for _, c := range g.clients {
		msg := c.GetMessage()
		// 不开启预测，直接采用服务端位置
		if !c.forecast {
			if msg != nil {
				c.player.Pos = msg.Pos
			}
			continue
		}
		// 开启预测，但未开启对账
		if c.forecast && !c.reconciliation {
			if msg != nil {
				c.player.Target = msg.Target
				c.player.Pos = msg.Pos
			}
			res := ProcessOne(c.player, 60)
			c.player.Pos = res.Pos
		}
		// 预测及对账
		if c.forecast && c.reconciliation {
			if msg != nil && msg.Index != 0 {
				buf, ok := c.ControlBuffer[msg.Index]
				if ok {
					// 对账失败，强制同步位置
					if !(buf.Target == msg.Pos) {
						c.player.Pos = msg.Pos
						c.ControlBuffer = map[int]ControlMsg{}
						continue
					}
					// 删除缓存
					delete(c.ControlBuffer, msg.Index)
				}
				// 缓存不存在，对账失败，强制同步位置
				if !ok {
					c.player.Pos = msg.Pos
				}
			}
			res := ProcessOne(c.player, 60)
			c.player.Pos = res.Pos
		}
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
			target := p.Pos
			target.X = float64(x)
			target.Y = float64(y)
			p.Target = target
			c.Move(target)
		}
	}
	// 渲染
	for _, c := range g.clients {
		p := c.player
		sizeX, sizeY := p.Image.Size()
		op.GeoM.Translate(p.Pos.X-float64(sizeX/2), p.Pos.Y-float64(sizeY/2))
		screen.DrawImage(p.Image, op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Hello,ebiten!\nTPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720 //窗口分辨率
}
