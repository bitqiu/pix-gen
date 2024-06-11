package handler

import (
	"bytes"
	"fmt"
	"github.com/bitqiu/pix-gen/fonts"
	"github.com/bitqiu/pix-gen/pkg/captcha"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"image/png"
	"net/http"
)

// HandleCaptcha 处理验证码生成请求的处理程序
func HandleCaptcha(c *gin.Context) {

	// 获取并解析 width 和 height 参数，如果不存在则使用默认值
	width := cast.ToInt(c.DefaultQuery("width", "120"))
	height := cast.ToInt(c.DefaultQuery("height", "30"))
	code := c.Query("code")

	// 调用 captcha 包生成验证码
	captchaImage, err := generateCaptcha(width, height, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回验证码图像
	c.Data(http.StatusOK, "image/png", captchaImage)

}

// generateCaptcha 生成验证码图片
func generateCaptcha(width, height int, code string) ([]byte, error) {
	// 初始化验证码生成器
	cap := captcha.New()
	// 设置干扰模式
	cap.SetDisturbance(captcha.NORMAL)

	// 读取字体文件
	fontBytes, err := fonts.FontsFS.ReadFile("MiSans-Normal.ttf")
	if err != nil {
		return nil, fmt.Errorf("invalid font file")
	}

	// 添加字体到验证码生成器
	err = cap.AddFontFromBytes(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid font")
	}

	// 检查 width 和 height 的边界条件
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("width and height must be positive integers")
	}

	// 设置验证码图片的大小
	cap.SetSize(width, height)

	// 生成新的验证码
	img := cap.CreateCustom(code)

	// 将图像编码为 PNG 并输出二进制图像数据
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return nil, fmt.Errorf("failed to encode image")
	}

	return buffer.Bytes(), nil
}
