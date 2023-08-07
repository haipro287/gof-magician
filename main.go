package main

import (
	"bufio"
	"gof-magician/assets/collection/slice"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/ericpauley/go-quantize/quantize"
	"golang.org/x/image/draw"
)

func decodePng(name string) (image.Image, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func decodeGif(name string) (*gif.GIF, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	g, err := gif.DecodeAll(r)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func encodeGif(name string, g *gif.GIF) error {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	gif.EncodeAll(f, g)
	return nil
}

func addFrame(g *gif.GIF, img image.Paletted, delay int) {
	g.Image = append(g.Image, &img)
	g.Delay = append(g.Delay, delay)
}

func weigth2Img(weight string) (image.Image, error) {
	bg, err := decodePng("assets/NuggetBackground.png")
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(bg.Bounds())
	draw.Draw(img, bg.Bounds(), bg, image.Point{0, 0}, draw.Src)

	weight = weight + "g"
	text2Num := []image.Image{}
	for _, c := range weight {
		name := string(c)
		size := 120
		if c == '.' {
			name = "dot"
			size = 45
		}
		img, err := decodePng("assets/digit/" + name + ".png")
		if err != nil {
			return nil, err
		}

		resized := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx()*size/img.Bounds().Dy(), size))

		draw.CatmullRom.Scale(resized, resized.Rect, img, img.Bounds(), draw.Over, nil)
		text2Num = append(text2Num, resized)
	}

	marginY := -12

	numWidths := slice.Map(text2Num, func(img image.Image) int { return img.Bounds().Dx() })
	numWidth := slice.Reduce(numWidths, func(a, b int) int { return a + b }, 0)
	_ = numWidth

	left := 0

	for i := range text2Num {
		ml := 0
		if i != 0 {
			ml = text2Num[i-1].Bounds().Dx() + marginY*2
		}
		left += ml
		marginTop := 0
		if weight[i] == '.' {
			marginTop += 73
		}
		if weight[i] == 'g' {
			marginTop += 22
		}

		draw.Draw(img, text2Num[i].Bounds().Add(image.Point{150 - (numWidth+len(text2Num)*marginY*2)/2 + left + marginY, 90 + marginTop}), text2Num[i], image.Point{0, 0}, draw.Over)
	}

	return img, nil
}

func main() {
	start := time.Now()

	decodedGif, err := decodeGif("assets/gif/GoldNugget_var4_tier1.gif")
	if err != nil {
		panic(err)
	}

	img, err := weigth2Img("4.0")
	if err != nil {
		panic(err)
	}

	q := quantize.MedianCutQuantizer{
		Aggregation:    quantize.Mean,
		Weighting:      nil,
		AddTransparent: true,
	}

	p := q.Quantize(make([]color.Color, 0, 256), img)
	palettedImg := image.NewPaletted(img.Bounds(), p)
	draw.Draw(palettedImg, palettedImg.Rect, img, image.Point{0, 0}, draw.Over)

	newGif := &gif.GIF{
		Image: decodedGif.Image,
		Delay: decodedGif.Delay,
	}
	addFrame(newGif, *palettedImg, 200)

	err = encodeGif("test.gif", newGif)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	log.Printf("Processing gif took %s", elapsed)
}
