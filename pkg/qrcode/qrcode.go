package qrcode

import (
	"bytes"
	"fmt"
	"github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"
	"strings"
)

// colorNames 是颜色名称到16进制颜色值的映射表
var colorNames = map[string]string{
	"black":   "000000",
	"white":   "ffffff",
	"red":     "ff0000",
	"green":   "00ff00",
	"blue":    "0000ff",
	"yellow":  "ffff00",
	"cyan":    "00ffff",
	"magenta": "ff00ff",
	"gray":    "808080",
	"purple":  "800080",
	"orange":  "ffa500",
}

// GenerateQRCode 生成二维码图像
func GenerateQRCode(text, level, sizeQuery, colorQuery, marginQuery string) ([]byte, error) {
	// 转换字符串为int，并增加错误处理
	size, err := strconv.ParseInt(sizeQuery, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid size format")
	}

	// 转换边距字符串为int，并增加错误处理
	margin, err := strconv.ParseInt(marginQuery, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid margin format")
	}

	// 检查边距不能大于 size 的四分之一
	if margin > size/4 {
		return nil, fmt.Errorf("margin cannot be greater than one quarter of the size")
	}

	// 设置错误校验级别
	qrLevel := qrcode.Medium
	switch level {
	case "L":
		qrLevel = qrcode.Low
	case "M":
		qrLevel = qrcode.Medium
	case "Q":
		qrLevel = qrcode.High
	case "H":
		qrLevel = qrcode.Highest
	default:
		return nil, fmt.Errorf("invalid QR code level")
	}

	// 创建二维码对象
	qrc, err := qrcode.New(text, qrLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code")
	}

	// 禁用默认的边距
	qrc.DisableBorder = true

	// 获取前景颜色
	rgbaColor, err := getColor(colorQuery)
	if err != nil {
		return nil, fmt.Errorf("invalid color format")
	}
	qrc.ForegroundColor = rgbaColor

	// 计算二维码图片的实际大小
	qrCodeSize := int(size - 2*margin)
	qrImage, err := qrc.PNG(qrCodeSize)
	if err != nil {
		return nil, fmt.Errorf("failed to encode QR code image")
	}

	// 创建带有边距的新图像
	qrWithMargin := addMarginToQRCode(qrImage, int(size), int(margin))

	// 编码带有边距的二维码图像
	var pngBuffer bytes.Buffer
	err = png.Encode(&pngBuffer, qrWithMargin)
	if err != nil {
		return nil, fmt.Errorf("failed to encode QR code image with margin")
	}

	return pngBuffer.Bytes(), nil
}

// getColor 根据颜色名字或16进制颜色值返回RGBA颜色
func getColor(input string) (color.RGBA, error) {
	hex := strings.TrimPrefix(strings.ToLower(input), "#")

	// 如果是颜色名字，转换为16进制颜色值
	if hexValue, ok := colorNames[hex]; ok {
		hex = hexValue
	}

	if len(hex) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid color format")
	}

	return hexToRGBA(hex)
}

// hexToRGBA 将16进制颜色转换为RGBA
func hexToRGBA(hex string) (color.RGBA, error) {
	r, err := strconv.ParseUint(hex[:2], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	b, err := strconv.ParseUint(hex[4:], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}

// addMarginToQRCode 添加边距到二维码图像
func addMarginToQRCode(qrImage []byte, size int, margin int) image.Image {
	img, _ := png.Decode(bytes.NewReader(qrImage))

	// 创建带边距的新图像
	newImg := image.NewRGBA(image.Rect(0, 0, size, size))
	bgColor := color.White
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 将二维码图像绘制到新图像中，应用边距
	offset := margin
	draw.Draw(newImg, image.Rect(offset, offset, size-offset, size-offset), img, image.Point{}, draw.Over)

	return newImg
}
