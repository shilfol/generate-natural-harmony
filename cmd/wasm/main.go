package main // mainでないとダメそう, TODO: 調査する

import (
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/shilfol/generate-natural-harmony/pkg/nh"
)

func main() {
	document := js.Global().Get("document")

	cb := func(this js.Value, arg []js.Value) interface{} {
		outputCanvas := document.Call("getElementById", "outputNaturalHarmony")
		ctx := outputCanvas.Call("getContext", "2d")

		cvWidth := outputCanvas.Get("width")
		cvHeight := outputCanvas.Get("height")

		data := ctx.Call("getImageData", 0, 0, cvWidth, cvHeight).Get("data")
		uarr := js.Global().Get("Uint8Array").New(data)
		garr := make([]byte, data.Get("length").Int())

		_ = js.CopyBytesToGo(garr, uarr)

		rv := document.Call("getElementById", "range")
		p, _ := strconv.ParseFloat(rv.Get("value").String(), 64)

		println("param: ", p)
		nhp := nh.NaturalHarmonyParam{
			P: p,
		}
		writeCanvas := func(r, g, b, a, x, y int) {
			str := fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
			ctx.Set("fillStyle", str)
			ctx.Call("fillRect", x, y, 1, 1)
		}
		nh.ConvertNaturalHarmonyFromBytes(garr, cvWidth.Int(), cvHeight.Int(), &nhp, writeCanvas)

		_ = js.CopyBytesToJS(uarr, garr)

		println(garr[0])
		println(garr[len(garr)-1])
		println(cvWidth.Int(), cvHeight.Int(), cvHeight.Int()*cvWidth.Int(), len(garr))

		return "wasm value"
	}

	inputElement := document.Call("getElementById", "range")
	inputElement.Call("addEventListener", "change", js.FuncOf(cb))

	// listenerに登録した関数を終了させない
	select {}
}
