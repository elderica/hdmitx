package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"

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

	flagsourcefile string
	flagframefile  string
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
	for y := 0; y < ScreenHeight; y += 2 {
		ebitenutil.DrawLine(screen, 0, float64(y), ScreenWidth, float64(y), lineColor)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func init() {
	flag.StringVar(&flagsourcefile, "sourcefile", "", "送信したいファイルのパス")
	flag.StringVar(&flagframefile, "framefile", "", "画像単位列を保存するファイルのパス")
}

func parseFlag() {
	flag.Parse()

	if flagsourcefile == "" {
		log.Fatalln("送信したいファイルのパスが指定されていません")
	}
	if abspath, err := filepath.Abs(flagsourcefile); err == nil {
		flagsourcefile = abspath
	} else {
		log.Fatalln(err)
	}
	if physpath, err := filepath.EvalSymlinks(flagsourcefile); err == nil {
		flagsourcefile = physpath
	} else {
		log.Fatalln(err)
	}
	log.Printf("sourcefile: %s\n", flagsourcefile)

	if flagframefile == "" {
		log.Fatalln("画像単位列を保存するファイルのパスが指定されていません")
	}
	if abspath, err := filepath.Abs(flagframefile); err == nil {
		flagframefile = abspath
	} else {
		log.Fatalln(err)
	}
	log.Printf("framefile: %s\n", flagframefile)
}

func main() {
	//parseFlag()

	flagsourcefile = "main.go"
	flagframefile = "/tmp/frames"

	frames, err := constructFramesFromFile(flagsourcefile)
	if err != nil {
		log.Fatalln(err)
	}

	for numbering, frame := range frames {
		framefilepath := fmt.Sprintf("%s.%d", flagframefile, numbering)
		os.WriteFile(framefilepath, frame, 0666)
	}

	ebiten.SetWindowTitle("hdmitx")
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetFullscreen(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
