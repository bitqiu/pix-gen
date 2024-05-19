package handler

import (
	"fmt"
	"github.com/bitqiu/pix-gen/pkg/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HandleCaptcha 处理验证码生成请求的处理程序
func HandleCaptcha(c *gin.Context) {

	// 获取并解析 width 和 height 参数，如果不存在则使用默认值
	widthStr := c.DefaultQuery("width", fmt.Sprintf("%d", 120))
	heightStr := c.DefaultQuery("height", fmt.Sprintf("%d", 30))
	code := c.Query("code")

	// 调用 captcha 包生成验证码
	captchaImage, err := captcha.GenerateCaptcha(widthStr, heightStr, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回验证码图像
	c.Data(http.StatusOK, "image/png", captchaImage)

}
