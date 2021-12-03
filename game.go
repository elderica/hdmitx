package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
)

var (
	backgroundColor = color.White
	lineColor       = color.Black
)

type Game struct {
	frames     []Frame
	frameindex int
}

func NewGame(filename string) *Game {
	payload, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("payload size: %d\n", len(payload))

	frames := MakeFrames(payload)

	log.Printf("len(frames) = %d\n", len(frames))
	//log.Printf("frames: %v\n", frames)

	return &Game{
		frames:     frames,
		frameindex: 0,
	}
}

func (g *Game) Update() error {
	d := len(g.frames)
	g.frameindex = (g.frameindex + 1) % d
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	// for y := 0; y < ScreenHeight; y += 2 {
	// 	ebitenutil.DrawLine(screen, 0, float64(y), ScreenWidth, float64(y), lineColor)
	// }

	frame := g.frames[g.frameindex]

	for j, b := range frame {
		for i := 7; i >= 0; i-- {
			if (b & (1 << i)) != 0 {
				y := 8*j + i
				ebitenutil.DrawLine(screen, 0, float64(y), ScreenWidth, float64(y), lineColor)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
