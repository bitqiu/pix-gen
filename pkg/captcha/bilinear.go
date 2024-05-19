package captcha

import (
	"image"
	"image/color"
	"math"
)

// 声明双线性插值的全局变量
var bili = Bilinear{}

// 定义 Bilinear 结构体
type Bilinear struct{}

// 双线性插值方法，根据给定的坐标获取插值后的颜色值
func (Bilinear) RGBA(src *image.RGBA, x, y float64) color.RGBA {
	p := findLinearSrc(src.Bounds(), x, y)

	// 获取周围像素的偏移量
	off00 := offRGBA(src, p.low.X, p.low.Y)
	off01 := offRGBA(src, p.high.X, p.low.Y)
	off10 := offRGBA(src, p.low.X, p.high.Y)
	off11 := offRGBA(src, p.high.X, p.high.Y)

	var fr, fg, fb, fa float64

	fr += float64(src.Pix[off00+0]) * p.frac00
	fg += float64(src.Pix[off00+1]) * p.frac00
	fb += float64(src.Pix[off00+2]) * p.frac00
	fa += float64(src.Pix[off00+3]) * p.frac00

	fr += float64(src.Pix[off01+0]) * p.frac01
	fg += float64(src.Pix[off01+1]) * p.frac01
	fb += float64(src.Pix[off01+2]) * p.frac01
	fa += float64(src.Pix[off01+3]) * p.frac01

	fr += float64(src.Pix[off10+0]) * p.frac10
	fg += float64(src.Pix[off10+1]) * p.frac10
	fb += float64(src.Pix[off10+2]) * p.frac10
	fa += float64(src.Pix[off10+3]) * p.frac10

	fr += float64(src.Pix[off11+0]) * p.frac11
	fg += float64(src.Pix[off11+1]) * p.frac11
	fb += float64(src.Pix[off11+2]) * p.frac11
	fa += float64(src.Pix[off11+3]) * p.frac11

	var c color.RGBA
	c.R = uint8(fr + 0.5)
	c.G = uint8(fg + 0.5)
	c.B = uint8(fb + 0.5)
	c.A = uint8(fa + 0.5)
	return c
}

// BilinearSrc 结构体，用于存储插值所需的参数
type BilinearSrc struct {
	// 左上角和右下角的插值源
	low, high image.Point
	// 每个像素的权重，0 后缀表示上/左，1 后缀表示下/右
	frac00, frac01, frac10, frac11 float64
}

// findLinearSrc 查找双线性插值的源像素
func findLinearSrc(b image.Rectangle, sx, sy float64) BilinearSrc {
	maxX := float64(b.Max.X)
	maxY := float64(b.Max.Y)
	minX := float64(b.Min.X)
	minY := float64(b.Min.Y)
	lowX := math.Floor(sx - 0.5)
	lowY := math.Floor(sy - 0.5)
	if lowX < minX {
		lowX = minX
	}
	if lowY < minY {
		lowY = minY
	}

	highX := math.Ceil(sx - 0.5)
	highY := math.Ceil(sy - 0.5)
	if highX >= maxX {
		highX = maxX - 1
	}
	if highY >= maxY {
		highY = maxY - 1
	}

	// 下述变量中的 0 后缀表示上/左，1 后缀表示下/右

	// 每个周围像素的中心
	x00 := lowX + 0.5
	y00 := lowY + 0.5
	x01 := highX + 0.5
	y01 := lowY + 0.5
	x10 := lowX + 0.5
	y10 := highY + 0.5
	x11 := highX + 0.5
	y11 := highY + 0.5

	p := BilinearSrc{
		low:  image.Pt(int(lowX), int(lowY)),
		high: image.Pt(int(highX), int(highY)),
	}

	// 边缘情况处理。如果我们足够接近图像的边缘，限制插值源
	if lowX == highX && lowY == highY {
		p.frac00 = 1.0
	} else if sy-minY <= 0.5 && sx-minX <= 0.5 {
		p.frac00 = 1.0
	} else if maxY-sy <= 0.5 && maxX-sx <= 0.5 {
		p.frac11 = 1.0
	} else if sy-minY <= 0.5 || lowY == highY {
		p.frac00 = x01 - sx
		p.frac01 = sx - x00
	} else if sx-minX <= 0.5 || lowX == highX {
		p.frac00 = y10 - sy
		p.frac10 = sy - y00
	} else if maxY-sy <= 0.5 {
		p.frac10 = x11 - sx
		p.frac11 = sx - x10
	} else if maxX-sx <= 0.5 {
		p.frac01 = y11 - sy
		p.frac11 = sy - y01
	} else {
		p.frac00 = (x01 - sx) * (y10 - sy)
		p.frac01 = (sx - x00) * (y11 - sy)
		p.frac10 = (x11 - sx) * (sy - y00)
		p.frac11 = (sx - x10) * (sy - y01)
	}

	return p
}

// offRGBA 获取给定坐标的像素偏移量
func offRGBA(src *image.RGBA, x, y int) int {
	return (y-src.Rect.Min.Y)*src.Stride + (x-src.Rect.Min.X)*4
}
