package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/inpututil"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenHeight = 400
	screenWidth  = 800
)

var (
	catImage *ebiten.Image
)

func init() {

	f, err := ebitenutil.OpenFile("resources/cat.png")
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	catImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	tileSize := 400
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(0.25, 0.25)

	var leftCat, rightCat image.Rectangle

	// Left cat
	var lx, ly int
	if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		ly = tileSize
	}
	leftCat = image.Rect(lx, ly, lx+tileSize, ly+tileSize)
	screen.DrawImage(catImage.SubImage(leftCat).(*ebiten.Image), op)

	// Right cat
	var rx, ry int
	rx = tileSize
	if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		ry = tileSize
	}
	rightCat = image.Rect(rx, ry, rx+tileSize, ry+tileSize)
	op.GeoM.Translate(float64(tileSize), 0)
	screen.DrawImage(catImage.SubImage(rightCat).(*ebiten.Image), op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))

	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Bongo cat"); err != nil {
		log.Fatal(err)
	}
}
