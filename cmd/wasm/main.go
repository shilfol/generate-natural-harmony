package main // mainでないとダメそう, TODO: 調査する

import "syscall/js"

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

		println(garr[0])
		println(garr[len(garr)-1])
		return "wasm value"
	}

	inputElement := document.Call("getElementById", "range")
	inputElement.Call("addEventListener", "change", js.FuncOf(cb))

	// listenerに登録した関数を終了させない
	select {}
}
