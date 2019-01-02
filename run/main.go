package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"

	"golang.org/x/image/font/gofont/gobold"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/kyeett/bongo-cat/resources"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenHeight = 600
	screenWidth  = 800
)

var (
	catImage   *ebiten.Image
	bongoImage *ebiten.Image

	audioContext *audio.Context

	leftSound  []byte
	rightSound []byte

	fnt     font.Face
	fntBold font.Face
)

func init() {
	// Load font

	var err error
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	fnt = truetype.NewFace(tt, &truetype.Options{
		Size: 24,
		DPI:  dpi,
	})

	tt, err = truetype.Parse(gobold.TTF)
	if err != nil {
		log.Fatal(err)
	}
	fntBold = truetype.NewFace(tt, &truetype.Options{
		Size: 24,
		DPI:  dpi,
	})
	// Load audio
	sampleRate := 44100
	audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}

	s, err := wav.Decode(audioContext, audio.BytesReadSeekCloser(resources.Bongo1_wav))
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(s)
	if err != nil {
		log.Fatal(err)
	}
	leftSound = b

	s2, err := wav.Decode(audioContext, audio.BytesReadSeekCloser(resources.Bongo4_wav))
	if err != nil {
		log.Fatal(err)
	}
	b2, err := ioutil.ReadAll(s2)
	if err != nil {
		log.Fatal(err)
	}
	rightSound = b2

	// Load images
	img, _, err := image.Decode(bytes.NewReader(resources.Cat_png))
	if err != nil {
		log.Fatal(err)
	}
	catImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, _, err = image.Decode(bytes.NewReader(resources.Bongo_png))
	if err != nil {
		log.Fatal(err)
	}
	bongoImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

}

func update(screen *ebiten.Image) error {

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	leftPlaying := inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft)
	rightPlaying := inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.JustPressedTouchIDs()) > 0 {
		var x int
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, _ = ebiten.CursorPosition()
		}
		if len(inpututil.JustPressedTouchIDs()) > 0 {
			x, _ = ebiten.TouchPosition(inpututil.JustPressedTouchIDs()[0])
		}
		switch {
		case x < screenWidth/2:
			leftPlaying = true
		default:
			rightPlaying = true
		}
	}

	drawTable(screen)
	if !leftPlaying {
		drawLeftCat(screen, leftPlaying)
	}
	if !rightPlaying {
		drawRightCat(screen, rightPlaying)
	}
	drawBongo(screen)

	if leftPlaying {
		drawLeftCat(screen, leftPlaying)
		playAudio(leftSound)
	}

	if rightPlaying {
		drawRightCat(screen, rightPlaying)
		playAudio(rightSound)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))

	// text.Draw(screen, "- Left bongo:   press left side, left arrow or A", fnt, 16, screenHeight-48, color.Black)
	text.Draw(screen, "Left bongo:", fntBold, 16, screenHeight-48, color.Black)
	text.Draw(screen, "press left side, left arrow or A", fnt, 180, screenHeight-48, color.Black)
	text.Draw(screen, "Right bongo:", fntBold, 16, screenHeight-16, color.Black)
	text.Draw(screen, "press right side, right arrow or D", fnt, 180, screenHeight-16, color.Black)
	return nil
}

const tileSize = 400

func drawTable(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, 0, 175, 800, 365, color.RGBA{0xff, 0x00, 0x00, 0xff})
	table := catImage.SubImage(image.Rect(0, 2*tileSize, 2*tileSize, 3*tileSize)).(*ebiten.Image)
	screen.DrawImage(table, nil)
}

func drawLeftCat(screen *ebiten.Image, playing bool) {
	var lx, ly int
	if playing {
		ly = tileSize
	}
	leftCat := catImage.SubImage(image.Rect(lx, ly, lx+tileSize, ly+tileSize)).(*ebiten.Image)
	screen.DrawImage(leftCat, nil)
}

func drawRightCat(screen *ebiten.Image, playing bool) {
	var rx int = tileSize
	var ry int
	if playing {
		ry = tileSize
	}
	rightCat := catImage.SubImage(image.Rect(rx, ry, rx+tileSize, ry+tileSize)).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(tileSize), 0)
	screen.DrawImage(rightCat, op)
}

func drawBongo(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Rotate(math.Pi / 20)
	op.GeoM.Translate(150, 200)
	screen.DrawImage(bongoImage, op)
}

func playAudio(b []byte) {
	p, _ := audio.NewPlayerFromBytes(audioContext, b)
	p.Play()
}

func main() {
	ebiten.SetMaxTPS(10)

	if err := ebiten.Run(update, screenWidth, screenHeight, 0.5, "Bongo cat"); err != nil {
		log.Fatal(err)
	}
}
