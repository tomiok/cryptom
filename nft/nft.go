package nft

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func CreateNFT() *os.File {
	blocks := 10
	blockw := 40
	img := image.NewRGBA(image.Rect(0, 0, blocks*blockw, 200))
	// Use these colors to get invalid RGB in the gradient.
	//c1, _ := colorful.Hex("#EEEF61")
	//c2, _ := colorful.Hex("#1E3140")

	for i := 0; i < blocks; i++ {
		draw.Draw(img, image.Rect(i*blockw, 0, (i+1)*blockw, 40), &image.Uniform{color.White}, image.Point{}, draw.Src)
		draw.Draw(img, image.Rect(i*blockw, 40, (i+1)*blockw, 80), &image.Uniform{color.Black}, image.Point{}, draw.Src)
		draw.Draw(img, image.Rect(i*blockw, 80, (i+1)*blockw, 120), &image.Uniform{color.White}, image.Point{}, draw.Src)
		draw.Draw(img, image.Rect(i*blockw, 120, (i+1)*blockw, 160), &image.Uniform{color.Black}, image.Point{}, draw.Src)
		draw.Draw(img, image.Rect(i*blockw, 160, (i+1)*blockw, 200), &image.Uniform{color.RGBA{
			R: 123,
			G: 250,
			B: 1,
			A: 59,
		}}, image.Point{}, draw.Src)
	}

	toimg, err := os.Create("colorblend.png")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil
	}
	defer toimg.Close()

	png.Encode(toimg, img)

	return toimg
}
