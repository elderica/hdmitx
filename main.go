package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	SCREEN_WIDTH  = 1920
	SCREEN_HEIGHT = 1080
)

var (
	backgroundColor = color.White
	lineColor       = color.Black
)

type Game struct{}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	for y := 0; y < SCREEN_HEIGHT; y += 2 {
		ebitenutil.DrawLine(screen, 0, float64(y), SCREEN_WIDTH, float64(y), lineColor)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	ebiten.SetWindowTitle("hdmitx")
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetFullscreen(true)
	ebiten.SetVsyncEnabled(true)

	log.Printf("VSync: %v", ebiten.IsVsyncEnabled())

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
