package captcha

import (
	"github.com/golang/freetype"
	"image/color"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

// TestNewImage 测试 NewImage 方法
func TestNewImage(t *testing.T) {
	img := NewImage(100, 50)
	if img.Bounds().Dx() != 100 || img.Bounds().Dy() != 50 {
		t.Errorf("NewImage: expected (100, 50), got (%d, %d)", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

// TestDrawLine 测试 DrawLine 方法
func TestDrawLine(t *testing.T) {
	img := NewImage(10, 10)
	red := color.RGBA{255, 0, 0, 255}
	img.DrawLine(1, 1, 8, 8, red)

	// 检查部分像素
	expectedPoints := []struct {
		x, y int
	}{
		{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8},
	}

	for _, p := range expectedPoints {
		if img.RGBAAt(p.x, p.y) != red {
			t.Errorf("DrawLine: expected pixel at (%d, %d) to be red", p.x, p.y)
		}
	}
}

// TestDrawCircle 测试 DrawCircle 方法
func TestDrawCircle(t *testing.T) {
	img := NewImage(20, 20)
	blue := color.RGBA{0, 0, 255, 255}
	img.DrawCircle(10, 10, 5, false, blue)

	// 检查部分像素
	expectedPoints := []struct {
		x, y int
	}{
		{10, 5}, {10, 15}, {5, 10}, {15, 10},
	}

	for _, p := range expectedPoints {
		if img.RGBAAt(p.x, p.y) != blue {
			t.Errorf("DrawCircle: expected pixel at (%d, %d) to be blue", p.x, p.y)
		}
	}
}

// TestDrawString 测试 DrawString 方法
func TestDrawString(t *testing.T) {
	img := NewImage(100, 50)
	black := color.RGBA{0, 0, 0, 255}

	// 使用 gofont 中的 goregular 字体
	font, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		t.Fatalf("failed to parse font: %v", err)
	}

	img.DrawString(font, black, "Test", 20)

	// 简单检查几个像素
	expectedPoints := []struct {
		x, y int
	}{
		{10, 10}, {30, 10}, {50, 10}, {70, 10},
	}

	for _, p := range expectedPoints {
		if img.RGBAAt(p.x, p.y) == black {
			t.Logf("DrawString: found black pixel at (%d, %d)", p.x, p.y)
		}
	}
}

// TestRotate 测试 Rotate 方法
func TestRotate(t *testing.T) {
	img := NewImage(10, 10)
	green := color.RGBA{0, 255, 0, 255}
	img.DrawLine(1, 1, 8, 8, green)
	rotated := img.Rotate(45)

	// 检查旋转后的部分像素
	expectedPoints := []struct {
		x, y int
	}{
		{5, 1}, {5, 8},
	}

	for _, p := range expectedPoints {
		if rotated.At(p.x, p.y) != green {
			t.Errorf("Rotate: expected pixel at (%d, %d) to be green", p.x, p.y)
		}
	}
}
