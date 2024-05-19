package captcha

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

// TestBilinearRGBA 测试 Bilinear 的 RGBA 方法
func TestBilinearRGBA(t *testing.T) {
	// 创建一个 2x2 的图像，分别填充不同的颜色
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{}, draw.Src) // 红色
	img.Set(1, 0, color.RGBA{0, 255, 0, 255})                                                         // 绿色
	img.Set(0, 1, color.RGBA{0, 0, 255, 255})                                                         // 蓝色
	img.Set(1, 1, color.RGBA{255, 255, 0, 255})                                                       // 黄色

	bilinear := Bilinear{}

	tests := []struct {
		x, y    float64
		want    color.RGBA
		message string
	}{
		{0.5, 0.5, color.RGBA{127, 127, 63, 255}, "中心点"}, // 中心点应为混合后的颜色
		{0, 0, color.RGBA{255, 0, 0, 255}, "左上角"},
		{1, 0, color.RGBA{0, 255, 0, 255}, "右上角"},
		{0, 1, color.RGBA{0, 0, 255, 255}, "左下角"},
		{1, 1, color.RGBA{255, 255, 0, 255}, "右下角"},
	}

	for _, tt := range tests {
		got := bilinear.RGBA(img, tt.x, tt.y)
		if got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.message, got, tt.want)
		}
	}
}

// TestFindLinearSrc 测试 findLinearSrc 函数
func TestFindLinearSrc(t *testing.T) {
	rect := image.Rect(0, 0, 2, 2)

	tests := []struct {
		sx, sy  float64
		want    BilinearSrc
		message string
	}{
		{0.5, 0.5, BilinearSrc{image.Pt(0, 0), image.Pt(1, 1), 0.25, 0.25, 0.25, 0.25}, "中心点"},
		{0, 0, BilinearSrc{image.Pt(0, 0), image.Pt(0, 0), 1.0, 0.0, 0.0, 0.0}, "左上角"},
		{1, 0, BilinearSrc{image.Pt(0, 0), image.Pt(1, 0), 0.5, 0.5, 0.0, 0.0}, "右上角"},
		{0, 1, BilinearSrc{image.Pt(0, 0), image.Pt(0, 1), 0.5, 0.0, 0.5, 0.0}, "左下角"},
		{1, 1, BilinearSrc{image.Pt(0, 0), image.Pt(1, 1), 0.25, 0.25, 0.25, 0.25}, "右下角"},
	}

	for _, tt := range tests {
		got := findLinearSrc(rect, tt.sx, tt.sy)
		if got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.message, got, tt.want)
		}
	}
}

// TestOffRGBA 测试 offRGBA 函数
func TestOffRGBA(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))

	tests := []struct {
		x, y    int
		want    int
		message string
	}{
		{0, 0, 0, "左上角"},
		{1, 0, 4, "右上角"},
		{0, 1, 8, "左下角"},
		{1, 1, 12, "右下角"},
	}

	for _, tt := range tests {
		got := offRGBA(img, tt.x, tt.y)
		if got != tt.want {
			t.Errorf("%s: got %d, want %d", tt.message, got, tt.want)
		}
	}
}
