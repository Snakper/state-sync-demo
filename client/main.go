package main

import (
	"log"

	"client/src"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := src.NewGameEngine()
	src.ConnectToServer()
	// 游戏引擎
	ebiten.SetWindowSize(1280, 720)        //窗口大小
	ebiten.SetWindowTitle("Hello, World!") //窗口标题
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
