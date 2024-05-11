package main

import (
	"github.com/bitqiu/pix-gen/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/captcha", handler.HandleCaptcha)
	r.GET("/qrcode", handler.HandleQrcode)
	r.Run(":8080")
}
