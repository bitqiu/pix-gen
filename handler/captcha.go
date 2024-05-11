package handler

import (
	"bytes"
	"fmt"
	"github.com/bitqiu/pix-gen/captcha"
	"github.com/bitqiu/pix-gen/fonts"
	"github.com/gin-gonic/gin"
	"image/png"
	"net/http"
	"strconv"
)

var cap *captcha.Captcha

var defaultWidth = 120
var defaultHeight = 30

func HandleCaptcha(c *gin.Context) {
	cap = captcha.New()
	cap.SetDisturbance(captcha.NORMAL)
	//cap.SetFrontColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	fontBytes, err := fonts.FontsFS.ReadFile("MiSans-Normal.ttf")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid"})
		return
	}
	err = cap.AddFontFromBytes(fontBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid font"})
		return
	}

	// 解析width和height参数，如果不存在则使用默认值
	widthStr := c.DefaultQuery("width", fmt.Sprintf("%d", defaultWidth))
	heightStr := c.DefaultQuery("height", fmt.Sprintf("%d", defaultHeight))

	// 转换字符串为int，并增加错误处理
	width, err := strconv.ParseInt(widthStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid width format"})
		return
	}
	height, err := strconv.ParseInt(heightStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid height format"})
		return
	}

	// 检查code参数，如果存在则尝试解码验证码
	code := c.Query("code")

	// 检查width和height的边界条件
	if width <= 0 || height <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Width and height must be positive integers"})
		return
	}
	cap.SetSize(int(width), int(height))

	// 生成新的验证码
	img := cap.CreateCustom(code)

	// 直接输出二进制图像数据
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image"})
		return
	}
	c.Data(http.StatusOK, "image/png", buffer.Bytes())
}
