package nh

import (
	"image"
)

func ConvertNaturalHarmony(img image.Image) image.Image {
	b := img.Bounds()
	ci := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			ci.Set(x, y, img.At(x, y))
		}
	}
	return ci 
}