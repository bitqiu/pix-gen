package captcha

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// Image 图片类型
type Image struct {
	*image.RGBA
}

// NewImage 创建一个新的图片
func NewImage(w, h int) *Image {
	img := &Image{image.NewRGBA(image.Rect(0, 0, w, h))}
	return img
}

// sign 返回数字的符号
func sign(x int) int {
	if x > 0 {
		return 1
	}
	return -1
}

// DrawLine 画直线
// 使用Bresenham算法 (https://zh.wikipedia.org/zh-cn/布雷森漢姆直線演算法)
// x1,y1 起点坐标；x2,y2 终点坐标；c 颜色
func (img *Image) DrawLine(x1, y1, x2, y2 int, c color.Color) {
	dx, dy, flag := int(math.Abs(float64(x2-x1))),
		int(math.Abs(float64(y2-y1))),
		false
	if dy > dx {
		flag = true
		x1, y1 = y1, x1
		x2, y2 = y2, x2
		dx, dy = dy, dx
	}
	ix, iy := sign(x2-x1), sign(y2-y1)
	n2dy := dy * 2
	n2dydx := (dy - dx) * 2
	d := n2dy - dx
	for x1 != x2 {
		if d < 0 {
			d += n2dy
		} else {
			y1 += iy
			d += n2dydx
		}
		if flag {
			img.Set(y1, x1, c)
		} else {
			img.Set(x1, y1, c)
		}
		x1 += ix
	}
}

// drawCircle8 绘制圆的八个对称点
func (img *Image) drawCircle8(xc, yc, x, y int, c color.Color) {
	img.Set(xc+x, yc+y, c)
	img.Set(xc-x, yc+y, c)
	img.Set(xc+x, yc-y, c)
	img.Set(xc-x, yc-y, c)
	img.Set(xc+y, yc+x, c)
	img.Set(xc-y, yc+x, c)
	img.Set(xc+y, yc-x, c)
	img.Set(xc-y, yc-x, c)
}

// DrawCircle 画圆
// xc,yc 圆心坐标；r 半径；fill 是否填充；c 颜色
func (img *Image) DrawCircle(xc, yc, r int, fill bool, c color.Color) {
	size := img.Bounds().Size()
	// 如果圆在图片可见区域外，直接退出
	if xc+r < 0 || xc-r >= size.X || yc+r < 0 || yc-r >= size.Y {
		return
	}
	x, y, d := 0, r, 3-2*r
	for x <= y {
		if fill {
			for yi := x; yi <= y; yi++ {
				img.drawCircle8(xc, yc, x, yi, c)
			}
		} else {
			img.drawCircle8(xc, yc, x, y, c)
		}
		if d < 0 {
			d = d + 4*x + 6
		} else {
			d = d + 4*(x-y) + 10
			y--
		}
		x++
	}
}

// DrawString 写字
// font 字体；c 颜色；str 字符串；fontsize 字体大小
func (img *Image) DrawString(font *truetype.Font, c color.Color, str string, fontsize float64) {
	ctx := freetype.NewContext()
	// 默认72dpi
	ctx.SetDst(img)
	ctx.SetClip(img.Bounds())
	ctx.SetSrc(image.NewUniform(c))
	ctx.SetFontSize(fontsize)
	ctx.SetFont(font)
	// 写入文字的位置
	pt := freetype.Pt(0, int(-fontsize/6)+ctx.PointToFixed(fontsize).Ceil())
	ctx.DrawString(str, pt)
}

// Rotate 旋转图像
// angle 旋转角度
func (img *Image) Rotate(angle float64) image.Image {
	return new(rotate).Rotate(angle, img.RGBA).transformRGBA()
}

// FillBkg 填充背景
// c 背景颜色
func (img *Image) FillBkg(c image.Image) {
	draw.Draw(img, img.Bounds(), c, image.ZP, draw.Over)
}

// distortTo 添加水波纹效果
// amplude 振幅；period 周期
// 来源：https://github.com/dchest/captcha/blob/master/image.go
func (img *Image) distortTo(amplude float64, period float64) {
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	oldm := img.RGBA

	dx := 1.4 * math.Pi / period
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)
			rgba := oldm.RGBAAt(x+int(xo), y+int(yo))
			if rgba.A > 0 {
				oldm.SetRGBA(x, y, rgba)
			}
		}
	}
}

// inBounds 判断坐标是否在图像边界内
func inBounds(b image.Rectangle, x, y float64) bool {
	return x >= float64(b.Min.X) && x < float64(b.Max.X) && y >= float64(b.Min.Y) && y < float64(b.Max.Y)
}

// rotate 结构体用于图像旋转
type rotate struct {
	dx, dy, sin, cos, neww, newh float64
	src                          *image.RGBA
}

// radian 将角度转换为弧度
func radian(angle float64) float64 {
	return angle * math.Pi / 180.0
}

// Rotate 旋转图像
func (r *rotate) Rotate(angle float64, src *image.RGBA) *rotate {
	r.src = src
	srsize := src.Bounds().Size()
	width, height := srsize.X, srsize.Y

	// 源图四个角的坐标（以图像中心为坐标系原点）
	srcwp, srchp := float64(width)*0.5, float64(height)*0.5
	srcx1, srcy1 := -srcwp, srchp
	srcx2, srcy2 := srcwp, srchp
	srcx3, srcy3 := -srcwp, -srchp
	srcx4, srcy4 := srcwp, -srchp

	r.sin, r.cos = math.Sincos(radian(angle))
	// 旋转后的四角坐标
	desx1, desy1 := r.cos*srcx1+r.sin*srcy1, -r.sin*srcx1+r.cos*srcy1
	desx2, desy2 := r.cos*srcx2+r.sin*srcy2, -r.sin*srcx2+r.cos*srcy2
	desx3, desy3 := r.cos*srcx3+r.sin*srcy3, -r.sin*srcx3+r.cos*srcy3
	desx4, desy4 := r.cos*srcx4+r.sin*srcy4, -r.sin*srcx4+r.cos*srcy4

	// 新的高度和宽度
	r.neww = math.Max(math.Abs(desx4-desx1), math.Abs(desx3-desx2)) + 0.5
	r.newh = math.Max(math.Abs(desy4-desy1), math.Abs(desy3-desy2)) + 0.5
	r.dx = -0.5*r.neww*r.cos - 0.5*r.newh*r.sin + srcwp
	r.dy = 0.5*r.neww*r.sin - 0.5*r.newh*r.cos + srchp
	return r
}

// pt 计算旋转后的坐标
func (r *rotate) pt(x, y int) (float64, float64) {
	return float64(-y)*r.sin + float64(x)*r.cos + r.dy,
		float64(y)*r.cos + float64(x)*r.sin + r.dx
}

// transformRGBA 将旋转后的图像转换为 RGBA
func (r *rotate) transformRGBA() image.Image {
	srcb := r.src.Bounds()
	b := image.Rect(0, 0, int(r.neww), int(r.newh))
	dst := image.NewRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			sx, sy := r.pt(x, y)
			if inBounds(srcb, sx, sy) {
				// 消除锯齿填色
				c := bili.RGBA(r.src, sx, sy)
				off := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
				dst.Pix[off+0] = c.R
				dst.Pix[off+1] = c.G
				dst.Pix[off+2] = c.B
				dst.Pix[off+3] = c.A
			}
		}
	}
	return dst
}
