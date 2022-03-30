package src

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player map[string]*Player
}

func NewGameEngine() *Game {
	return &Game{
		player: map[string]*Player{},
	}
}

//func (g *Game) OpenForecast() {
//	for _, c := range g.clients {
//		c.SetForecast(true)
//	}
//}
//
//func (g *Game) OpenReconciliation() {
//	for _, c := range g.clients {
//		c.SetReconciliation(true)
//	}
//}

func (g *Game) AddPlayer(p *Player) {
	f, err := os.Open("img.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	eimg := ebiten.NewImageFromImage(img)
	p.image = eimg
	g.player[p.id] = p
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// 渲染
	lock.Lock()
	for _, m := range msg {
		p, ok := g.player[m.Id]
		if !ok {
			p = &Player{
				id:  m.Id,
				pos: m.Pos,
			}
			g.AddPlayer(p)
		}
		p.pos = m.Pos
	}
	msg = msg[:0]
	lock.Unlock()
	for _, p := range g.player {
		sizeX, sizeY := p.image.Size()
		op.GeoM.Translate(p.pos.X-float64(sizeX/2), p.pos.Y-float64(sizeY/2))
		screen.DrawImage(p.image, op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Hello,ebiten!\nTPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720 //窗口分辨率
}
