package internal

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
	"image"
)

// Соответсвие при поиске изображения
const ImgSearchRate = 0.9

type Screen struct {
	bounds *image.Rectangle
}

func NewScreen(bounds image.Rectangle) *Screen {
	return &Screen{bounds: &bounds}
}

func (s *Screen) GetMousePos() (int, int) {
	x, y := robotgo.GetMousePos()
	fmt.Println("pos: ", x, y)
	return x, y
}

func (s *Screen) CaptureScreen() image.Image {
	// если не передали в NewScreen делаем скрин всего экрана
	if s == nil {
		sx, sy := robotgo.GetScreenSize()
		bit := robotgo.CaptureScreen(0, 0, sx, sy)
		defer robotgo.FreeBitmap(bit)
		img := robotgo.ToImage(bit)
		return img
	}
	bit := robotgo.CaptureScreen(s.bounds.Min.X, s.bounds.Min.Y, s.bounds.Max.X, s.bounds.Max.Y)
	defer robotgo.FreeBitmap(bit)
	img := robotgo.ToImage(bit)
	return img
}

func (s *Screen) FindOnScreen(path string) (*image.Point, error) {
	subImg, _, err := robotgo.DecodeImg(path)
	//imgSub, err:=imgo.ReadPNG("F:\\projects\\go\\screen_bot\\static\\l2\\login.PNG")
	if err != nil {
		return nil, fmt.Errorf("img decode err:%w", err)
	}
	screen := s.CaptureScreen()
	_, rate, _, point := gcv.FindImg(subImg, screen)
	if rate < ImgSearchRate {
		return nil, nil
	}
	return &point, nil
}
