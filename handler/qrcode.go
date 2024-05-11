package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"image/color"
	"net/http"
	"strconv"
	"strings"
)

func HandleQrcode(c *gin.Context) {
	text := c.DefaultQuery("text", "null")
	level := c.DefaultQuery("level", "H")
	sizeQuery := c.DefaultQuery("size", fmt.Sprintf("%d", 300))
	colorQuery := c.DefaultQuery("color", "000000")

	// 转换字符串为int，并增加错误处理
	size, err := strconv.ParseInt(sizeQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size format"})
		return
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid QR code level"})
		return
	}

	qrc, err := qrcode.New(text, qrLevel)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to create QR code")
		return
	}
	qrc.DisableBorder = true

	rgbaColor, err := hexToRGBA(colorQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid color format")
		return
	}
	qrc.ForegroundColor = rgbaColor

	png, errEncode := qrc.PNG(int(size))
	if errEncode != nil {
		c.JSON(http.StatusBadRequest, "Failed to encode QR code image")
		return
	}

	c.Data(http.StatusOK, "image/png", png)

}

// hexToRGBA 将16进制颜色转换为RGBA。
func hexToRGBA(hex string) (color.RGBA, error) {
	hex = strings.TrimPrefix(strings.ToLower(hex), "#")
	if len(hex) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex color format")
	}
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
