package internal

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
	"image"
	"time"
)

type Screen struct {
	bounds image.Rectangle
	opts   ScreenOpts
}

type ScreenOpts struct {
	CaptureDelay time.Duration
}

func NewScreen(bounds image.Rectangle, opts ScreenOpts) *Screen {
	return &Screen{bounds: bounds, opts: opts}
}

func (s *Screen) GetMousePos() (int, int) {
	x, y := robotgo.GetMousePos()
	fmt.Println("pos: ", x, y)
	return x, y
}

func (s *Screen) CaptureScreen() image.Image {
	time.Sleep(s.opts.CaptureDelay * time.Millisecond)
	bit := robotgo.CaptureScreen(s.bounds.Min.X, s.bounds.Min.Y, s.bounds.Max.X, s.bounds.Max.Y)
	defer robotgo.FreeBitmap(bit)
	img := robotgo.ToImage(bit)
	return img
}

func (s *Screen) ImageOnScreen(subImg image.Image) bool {
	time.Sleep(s.opts.CaptureDelay * time.Millisecond)
	bit := robotgo.CaptureScreen(s.bounds.Min.X, s.bounds.Min.Y, s.bounds.Max.X, s.bounds.Max.Y)
	defer robotgo.FreeBitmap(bit)
	img := robotgo.ToImage(bit)
	gcv.FindImg(subImg, img)
	return true
}
