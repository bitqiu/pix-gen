package handler

import (
	qc "github.com/bitqiu/pix-gen/pkg/qrcode"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HandleQrcode 是处理生成二维码请求的处理程序
func HandleQrcode(c *gin.Context) {
	text := c.DefaultQuery("text", "null")          // 获取二维码的内容，默认为 "null"
	level := c.DefaultQuery("level", "H")           // 获取错误校验级别，默认为 "H"
	sizeQuery := c.DefaultQuery("size", "300")      // 获取二维码大小，默认为 300
	colorQuery := c.DefaultQuery("color", "000000") // 获取前景颜色，默认为黑色
	marginQuery := c.DefaultQuery("margin", "0")    // 获取边距大小，默认为 0

	// 调用 qc 包生成二维码
	qrCode, err := qc.GenerateQRCode(text, level, sizeQuery, colorQuery, marginQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回二维码图像
	c.Data(http.StatusOK, "image/png", qrCode)
}
