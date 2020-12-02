package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/shilfol/generate-natural-harmony/pkg/nh"
)

func main() {
	var p float64
	var filePath string
	flag.Float64Var(&p, "p", 0.0, "natural harmony parameter : [0.0, 1.0]")
	flag.StringVar(&filePath, "f", "sample.png", "convert file path")
	flag.Parse()

	f, err := os.Open(filePath)
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

	ci := nh.ConvertNaturalHarmony(img, &nh.NaturalHarmonyParam{P: p})

	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		if err := os.Mkdir("./output", 0755); err != nil {
			fmt.Println(err)
			return
		}
	}

	basePath := filepath.Base(filePath)
	sp := strings.Split(basePath, ".")
	genPath := filepath.Join("./output", sp[len(sp)-2]+"_natural_harmony."+sp[len(sp)-1])
	fmt.Println("converted:", genPath)

	o, err := os.Create(genPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer o.Close()

	// TODO: 入力で受けた拡張子と合わせる
	png.Encode(o, ci)
}
