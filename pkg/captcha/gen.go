package captcha

import (
	"bytes"
	"fmt"
	"github.com/bitqiu/pix-gen/fonts"
	"image/png"
	"strconv"
)

// GenerateCaptcha 生成验证码图片
func GenerateCaptcha(widthStr, heightStr, code string) ([]byte, error) {
	// 初始化验证码生成器
	cap := New()
	// 设置干扰模式
	cap.SetDisturbance(NORMAL)

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

	// 获取并解析 width 和 height 参数，如果不存在则使用默认值
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		return nil, fmt.Errorf("invalid width format")
	}
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		return nil, fmt.Errorf("invalid height format")
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
