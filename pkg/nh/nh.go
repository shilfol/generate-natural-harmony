package nh

import (
	"image"
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// NaturalHarmonyParam : Optional Parameter
// どれだけナチュラルハーモニーに近づけるか指定します
// 0に近づくほど元画像に近く, 1に近づくほど元画像から離れます
// 0.1~0.4くらいがおすすめです
type NaturalHarmonyParam struct {
	P float64
}

type processedMessage struct {
	x, y  int
	color color.Color
}

// ConvertNaturalHarmony convert Raw Image to Natual Harmony Image
func ConvertNaturalHarmony(img image.Image, nhp *NaturalHarmonyParam) image.Image {
	b := img.Bounds()
	ci := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			p := innerProcessNaturalHarmony(img.At(x, y), nhp)
			ci.Set(x, y, p)
		}
	}
	return ci
}

// ConvertNaturalHarmonyAsync convert Raw Image to Natual Harmony Image with goroutine
func ConvertNaturalHarmonyAsync(img image.Image, nhp *NaturalHarmonyParam) image.Image {
	b := img.Bounds()
	ci := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		go innerProcessNaturalHarmonyAsync(y, img, nhp, ci)
	}
	return ci
}

// ConvertNaturalHarmonyFromBytes convert raw bytes (destructive method)
func ConvertNaturalHarmonyFromBytes(raw []byte, w, h int, nhp *NaturalHarmonyParam) {
	nmax := 255.0
	// for y := 0; y < h; y++ {
	// 	id := 4 * y * w
	// 	for x := 0; x < w; x++ {
	// 		idx := id + x*4
	// 		p := raw[idx : idx+4] // r,g,b,a pixel value [0, 255]
	// 		co := colorful.Color{R: float64(p[0]) / nmax, G: float64(p[1]) / nmax, B: float64(p[2]) / nmax}
	// 		mc := innerColorProcessNaturalHarmony(co, nhp)
	// 		r, g, b, a := mc.RGBA()
	// 		ur, ug, ub, ua := convertUInt8(r, g, b, a)
	// 		raw[idx] = byte(ur)
	// 		raw[idx+1] = byte(ug)
	// 		raw[idx+2] = byte(ub)
	// 		raw[idx+3] = byte(ua)
	// 	}
	// }
	fc := make(chan int, h)
	for y := 0; y < h; y++ {
		id := 4 * y * w
		go func(id int) {
			for x := 0; x < w; x++ {
				idx := id + x*4
				p := raw[idx : idx+4] // r,g,b,a pixel value [0, 255]
				co := colorful.Color{R: float64(p[0]) / nmax, G: float64(p[1]) / nmax, B: float64(p[2]) / nmax}
				mc := innerColorProcessNaturalHarmony(co, nhp)
				r, g, b, a := mc.RGBA()
				ur, ug, ub, ua := convertUInt8(r, g, b, a)
				raw[idx] = byte(ur)
				raw[idx+1] = byte(ug)
				raw[idx+2] = byte(ub)
				raw[idx+3] = byte(ua)
			}
			fc <- id
		}(id)
	}
	for i := 0; i < h; i++ {
		<-fc
	}
}

func innerProcessNaturalHarmonyAsync(y int, img image.Image, nhp *NaturalHarmonyParam, ci *image.RGBA) {
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		color := innerProcessNaturalHarmony(img.At(x, y), nhp)
		ci.Set(x, y, color)
	}
}

func innerProcessNaturalHarmony(c color.Color, nhp *NaturalHarmonyParam) color.Color {
	nmax := 65535.0
	r, g, b, _ := c.RGBA()

	// convert colorful color
	co := colorful.Color{R: float64(r) / nmax, G: float64(g) / nmax, B: float64(b) / nmax}

	mc := innerColorProcessNaturalHarmony(co, nhp)
	mr, mg, mb, ma := mc.RGBA()
	ur, ug, ub, ua := convertUint16(mr, mg, mb, ma)
	nhc := color.RGBA64{ur, ug, ub, ua}
	return nhc
}

func innerColorProcessNaturalHarmony(c colorful.Color, nhp *NaturalHarmonyParam) colorful.Color {
	mc := processNaturalHarmonyHSV(c, nhp)
	return mc
}

func processNaturalHarmonyHSV(c colorful.Color, nhp *NaturalHarmonyParam) colorful.Color {
	h, s, v := c.Hsv()
	mh, ms, mv := mappingNaturalHarmonyHSV(h, s, v, nhp)
	mc := colorful.Hsv(mh, ms, mv)
	return mc
}

func processNaturalHarmonyHCL(co colorful.Color, nhp *NaturalHarmonyParam) colorful.Color {
	h, c, l := co.Hcl()
	mh, mc, ml := mappingNaturalHarmonyHCL(h, c, l, nhp)
	rc := colorful.Hcl(mh, mc, ml)
	return rc
}

func convertUint16(r, g, b, a uint32) (ur, ug, ub, ua uint16) {
	return uint16(r), uint16(g), uint16(b), uint16(a)
}

func convertUInt8(r, g, b, a uint32) (ur, ug, ub, ua int) {
	// intで返すがuint8の範囲([0,255])に収める
	div := 65535.0
	fr, fg, fb, fa := float64(r)/div, float64(g)/div, float64(b)/div, float64(a)/div
	return int(fr * 255.0), int(fg * 255.0), int(fb * 255.0), int(fa * 255.0)
}

func isNearlyZero(x float64) bool {
	delta := 1.0e-5
	return math.Abs(x-0.0000000000) < delta
}

func mappingNaturalHarmonyHSV(h, s, v float64, nhp *NaturalHarmonyParam) (hh, ss, vv float64) {
	// hue(色相)はそのまま
	hh = h
	// saturation(彩度)もいったんそのまま
	ss = s
	// value(明度)はhueに沿って変換する
	// 黄色の値をyh, 紫の値をphとするとcos(yh) = cos0 = 1, cos(ph) = cos pi = -1となるように定める
	// 0~360を0~2pi(rad)に変換するので, pi/180をかける
	radian := h * math.Pi / 180.0
	// 黄色のhは60(0~360)の値であるので, pi/3だけずらす
	// vは0~1, cosは-1~1なので2で割って1/2を足す
	cv := (math.Cos(radian-math.Pi/3.0))/2.0 + 0.5

	// もともとのvとdiffを取ってみる
	// dv := math.Abs(cv - v)
	// diffをいくらか小さくした値(p)をcvに近づける方向へ加算
	// vv = ((cv-v)/dv)*dv*p + v を整理したもの
	p := nhp.P
	if !isNearlyZero(h) || !isNearlyZero(s) {
		vv = p*(cv-v) + v
	} else {
		// 白 or 黒はそのまま出す
		vv = v
	}

	return
}

// 実験中 :　明度情報を使ったナチュラルハーモニー変換
// 暗い部分は青紫に近づけ, 明るい部分は黄色に近づけてみる
func mappingNaturalHarmonyHSVchangeHue(h, s, v float64, nhp *NaturalHarmonyParam) (hh, ss, vv float64) {
	// saturation(彩度)はそのまま
	ss = s
	// 後続処理は白or黒の場合除外する
	if isNearlyZero(h) && isNearlyZero(s) {
		vv = v
		hh = h
		return
	}

	// 元の色のラジアンを計算する
	rr := h * math.Pi / 180.0

	// valueは上述と同様の処理を加える
	p := nhp.P
	cv := (math.Cos(rr-math.Pi/3.0))/2.0 + 0.5
	vv = 0.5*p*(cv-v) + v

	// hue(色相)をvalueに沿って変換する
	// vが0に近い(黒寄り)なら青の方向に, 1に近い(白寄り)なら黄の方向に寄せる
	// 0~1を-pi~0(rad)に変換するので, -1した後にpiをかける
	vr := (v - 1.0) * math.Pi
	// 基準となる角度は黄色のpi/3
	radian := math.Pi / 3.0

	// 60(黄)~150(緑)~240(青)の範囲か240~330(赤)~60の範囲であるべき色を変える
	// vrが負の数を取ることに注意する
	if h >= 60.0 && h <= 240 {
		radian -= vr
	} else {
		radian += vr
	}

	// 差分を取りながら計算
	// パラメータもここで利用する
	// valueとは感度が違うので補正をかけて抑えておく
	rhh := 0.1*p*(radian-rr) + rr

	// [0, 2pi]の範囲に抑えておくために2piで剰余を取っておく
	// 剰余が取れないので力づくで
	tpi := 2.0 * math.Pi
	if rhh > tpi {
		rhh -= tpi
	} else if rhh < 0.0 {
		rhh += tpi
	}

	// ラジアンからhに戻す
	hh = rhh * 180.0 / math.Pi

	return
}

func mappingNaturalHarmonyHCL(h, c, l float64, nhp *NaturalHarmonyParam) (hh, cc, ll float64) {
	//WIP
	// hue(色相)はそのまま
	hh = h
	// chroma(彩度)もいったんそのまま
	cc = c
	// lightness(明度)はhueに沿って変換する
	// 以下はメソッドはhsvそのまま
	// 黄色の値をyh, 紫の値をphとするとcos(yh) = cos0 = 1, cos(ph) = cos pi = -1となるように定める
	// 0~360を0~2pi(rad)に変換するので, pi/180をかける
	radian := h * math.Pi / 180.0
	yr := 102.9 * math.Pi / 180.0
	// 黄色のhは90(0~360)の値であるので, pi/2だけずらす
	// vは0~1, cosは-1~1なので2で割って1/2を足す
	cl := (math.Cos(radian-yr))/2.0 + 0.5

	// もともとのvとdiffを取ってみる
	// dv := math.Abs(cv - v)
	// diffをいくらか小さくした値(p)をcvに近づける方向へ加算
	// vv = ((cv-v)/dv)*dv*p + v を整理したもの
	p := nhp.P
	if !isNearlyZero(c) {
		ll = p*(cl-l) + l
	} else {
		// TODO: 黒は弾いたが白を弾いてない
		ll = l
	}

	return
}
