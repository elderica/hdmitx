package main

import (
	"flag"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	flagsourcefile string
)

func init() {
	flag.StringVar(&flagsourcefile, "sourcefile", "", "送信したいファイルのパス")
}

func fileExists(filename string) bool {
	fileinfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if fileinfo.IsDir() {
		return false
	}

	return true
}

func parseFlags() {
	flag.Parse()

	if flagsourcefile == "" {
		log.Fatalln("送信したいファイルのパスが指定されていません")
	}

	if !fileExists(flagsourcefile) {
		log.Fatalln("送信したいファイルが存在しません")
	}

	log.Printf("sourcefile: %s\n", flagsourcefile)
}

func main() {
	parseFlags()

	ebiten.SetWindowTitle("hdmitx")
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetFullscreen(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)

	if err := ebiten.RunGame(NewGame(flagsourcefile)); err != nil {
		log.Fatal(err)
	}
}
