package internal

import (
	"bytes"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
)

// Соответсвие при поиске изображения
const ImgSearchRate = 0.9

const TesseractPath = "F:\\soft\\tesseract\\tesseract.exe"

type Screen struct {
	bounds *image.Rectangle
}

func NewScreen(bounds *image.Rectangle) *Screen {
	return &Screen{bounds: bounds}
}

func (s *Screen) GetMousePos() (int, int) {
	x, y := robotgo.GetMousePos()
	fmt.Println("pos: ", x, y)
	return x, y
}

// CaptureScreen (x, y, w, h int)
func (s *Screen) CaptureScreen(opt ...int) image.Image {
	if len(opt) == 0 {
		sx, sy := robotgo.GetScreenSize()
		bit := robotgo.CaptureScreen(0, 0, sx, sy)
		defer robotgo.FreeBitmap(bit)
		img := robotgo.ToImage(bit)
		return img
	}
	if len(opt) == 4 {
		bit := robotgo.CaptureScreen(opt[0], opt[1], opt[2], opt[3])
		defer robotgo.FreeBitmap(bit)
		img := robotgo.ToImage(bit)
		return img
	}
	return nil
}

func (s *Screen) FindOnScreen(path string, opt ...int) (*image.Point, error) {
	subImg, _, err := robotgo.DecodeImg(path)
	//imgSub, err:=imgo.ReadPNG("F:\\projects\\go\\screen_bot\\static\\l2\\login.PNG")
	if err != nil {
		return nil, fmt.Errorf("img decode err:%w", err)
	}
	screen := s.CaptureScreen(opt...)
	_, rate, _, point := gcv.FindImg(subImg, screen)
	if rate < ImgSearchRate {
		return nil, nil
	}
	return &point, nil
}

func (s *Screen) FindText(imgPath string) (string, error) {
	cmd := exec.Command(TesseractPath, imgPath, "-", "-l", "rus+eng")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run tesseract: %s", errb.String())
	}
	return outb.String(), nil
}

func (s *Screen) SaveScreen(path string, opt ...int) error {
	screen := s.CaptureScreen(opt...)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("can't create tmp image: %w", err)
	}
	defer f.Close()
	if err = jpeg.Encode(f, screen, nil); err != nil {
		return fmt.Errorf("failed to encode img: %w", err)
	}
	return nil
}
