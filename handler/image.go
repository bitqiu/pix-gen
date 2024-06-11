package handler

import (
	"bytes"
	"fmt"
	"github.com/bitqiu/pix-gen/fonts"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/spf13/cast"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"net/http"
)

// HandleImage 是处理生成二维码请求的处理程序
func HandleImage(c *gin.Context) {
	text := c.DefaultQuery("text", "null")                       // 获取二维码的内容，默认为 "null"
	width := cast.ToInt(c.DefaultQuery("width", "500"))          // 获取图片宽度，默认为 800
	height := cast.ToInt(c.DefaultQuery("height", "100"))        // 获取图片高度，默认为 600
	tipText := c.DefaultQuery("tipText", "请通过图片和复制的地址核对一样后进行转账") // 获取自定义文字，默认为 "Custom Text"

	// 调用 generateImage 函数生成图像
	imageData, err := generateImage(text, tipText, width, height)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回生成的图像
	c.Data(http.StatusOK, "image/png", imageData)
}

// generateImage 生成带有指定文字的图像
func generateImage(text, tipText string, width, height int) ([]byte, error) {
	// 读取字体数据
	fontBytes, err := fonts.FontsFS.ReadFile("MiSans-Normal.ttf")
	if err != nil {
		return nil, fmt.Errorf("读取字体文件出错: %v", err)
	}

	// 解析字体
	parsedFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("解析字体出错: %v", err)
	}

	// 根据图像尺寸动态计算字体大小
	area := float64(width * height)
	fontSize := math.Sqrt(area / float64(100))

	// 创建一个新的 RGBA 图像，背景为白色
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	bgColor := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 初始化 freetype 上下文
	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(parsedFont)
	ctx.SetFontSize(fontSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(image.Black)

	// 计算文字的尺寸
	opts := &truetype.Options{
		Size: fontSize,
		DPI:  72,
	}
	face := truetype.NewFace(parsedFont, opts)
	textWidth, textHeight := getTextDimensions(face, text)
	tipTextWidth, tipTextHeight := getTextDimensions(face, tipText)

	// 计算文本块的总高度
	totalTextHeight := textHeight + tipTextHeight + 10 // 两行文字间距 10 像素

	// 计算起始点，使文本水平和垂直居中
	x := (width - textWidth) / 2
	tipX := (width - tipTextWidth) / 2
	y := (height - totalTextHeight) / 2

	// 绘制主文字
	pt := freetype.Pt(x, y+int(ctx.PointToFixed(fontSize)>>6))
	if _, err := ctx.DrawString(text, pt); err != nil {
		return nil, fmt.Errorf("绘制文字出错: %v", err)
	}

	// 用红色绘制自定义提示文字
	ctx.SetSrc(image.NewUniform(color.RGBA{255, 0, 0, 255}))
	tipY := y + textHeight + 10
	tipPt := freetype.Pt(tipX, tipY+int(ctx.PointToFixed(fontSize)>>6))
	if _, err := ctx.DrawString(tipText, tipPt); err != nil {
		return nil, fmt.Errorf("绘制自定义文字出错: %v", err)
	}

	// 将图像编码为 PNG 格式
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("编码图像出错: %v", err)
	}

	return buf.Bytes(), nil
}

// getTextDimensions 计算给定文字的宽度和高度
func getTextDimensions(face font.Face, text string) (int, int) {
	bounds, _ := font.BoundString(face, text)
	width := (bounds.Max.X - bounds.Min.X).Ceil()
	height := (bounds.Max.Y - bounds.Min.Y).Ceil()
	return width, height
}
