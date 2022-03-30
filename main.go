package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	svr := NewServer(10)
	client := NewClient()
	p1 := client.NewPlayer()
	g := NewGameEngine(p1.id)
	g.AddClient(client)
	client.Connect(svr)
	svr.Run()

	// 客户端预测
	g.OpenForecast()
	// 客户端、服务端对账
	g.OpenReconciliation()

	// 游戏引擎
	ebiten.SetWindowSize(1280, 720)        //窗口大小
	ebiten.SetWindowTitle("Hello, World!") //窗口标题
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
