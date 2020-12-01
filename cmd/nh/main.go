package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/shilfol/generate-natural-harmony/pkg/nh"
)

func main() {
	path := "./sample.png"

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	ci := nh.ConvertNaturalHarmony(img)

	sp := strings.Split(path, ".")
	genPath := "." + sp[len(sp) - 2] + "_out." + sp[len(sp) - 1]
	fmt.Println("converted:", genPath)

	o, err := os.Create(genPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer o.Close()

	png.Encode(o, ci)
}