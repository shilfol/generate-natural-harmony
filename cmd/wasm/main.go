package main // mainでないとダメそう, TODO: 調査する

import (
	"strconv"
	"syscall/js"

	"github.com/shilfol/generate-natural-harmony/pkg/nh"
)

func main() {
	document := js.Global().Get("document")

	cb := func(this js.Value, arg []js.Value) interface{} {
		// 読み書きを行うcanvasとcontextの取得
		outputCanvas := document.Call("getElementById", "outputNaturalHarmony")
		ctx := outputCanvas.Call("getContext", "2d")

		cvWidth := outputCanvas.Get("width")
		cvHeight := outputCanvas.Get("height")

		// canvasのbyte列をuint8の形で読み込む
		data := ctx.Call("getImageData", 0, 0, cvWidth, cvHeight).Get("data")
		uarr := js.Global().Get("Uint8ClampedArray").New(data)
		garr := make([]byte, data.Get("length").Int())

		// Go空間に読み込む
		_ = js.CopyBytesToGo(garr, uarr)

		// rangeの値からパラメータ取得
		rv := document.Call("getElementById", "range")
		p, _ := strconv.ParseFloat(rv.Get("value").String(), 64)
		nhp := nh.NaturalHarmonyParam{
			P: p,
		}

		// ナチュラルハーモニーする
		nh.ConvertNaturalHarmonyFromBytes(garr, cvWidth.Int(), cvHeight.Int(), &nhp)

		// objectを再利用してjs空間に戻す
		_ = js.CopyBytesToJS(uarr, garr)

		// canvasに書き戻す
		imageData := js.Global().Get("ImageData").New(uarr, cvWidth, cvHeight)
		ctx.Call("putImageData", imageData, 0, 0)

		println("enjoy Natural Harmony!")

		return "wasm value"
	}

	// rangeの変更に追従するcbとして登録
	inputElement := document.Call("getElementById", "range")
	inputElement.Call("addEventListener", "change", js.FuncOf(cb))

	// listenerに登録した関数を終了させない
	select {}
}
