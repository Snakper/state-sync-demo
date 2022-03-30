package main

import (
	"log"

	"server/src"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	src.Listen()

	svr := src.NewServer()
	client := src.NewClient()
	p1 := client.NewPlayer()
	g := src.NewGameEngine(p1.Id)
	g.AddClient(client)
	client.Connect(svr)
	svr.Run()

	// 客户端预测
	g.OpenForecast()
	// 客户端、服务端对账
	g.OpenReconciliation()
	// 客户端插值
	g.OpenInterpolation()

	// 游戏引擎
	ebiten.SetWindowSize(1280, 720)        //窗口大小
	ebiten.SetWindowTitle("Hello, World!") //窗口标题
	ebiten.SetMaxTPS(int(src.ClientFrame))
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
