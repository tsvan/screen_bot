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
}

func NewScreen() *Screen {
	return &Screen{}
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
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't get root directory path")
	}

	subImg, _, err := robotgo.DecodeImg(rootPath + path)
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
	rootPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("can't get root directory path")
	}

	cmd := exec.Command(TesseractPath, rootPath+imgPath, "-", "-l", "rus+eng")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run tesseract: %s", errb.String())
	}
	return outb.String(), nil
}

func (s *Screen) SaveScreen(path string, opt ...int) error {
	screen := s.CaptureScreen(opt...)

	rootPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("can't get root directory path")
	}

	f, err := os.Create(rootPath + path)
	if err != nil {
		return fmt.Errorf("can't create tmp image: %w", err)
	}

	if err = jpeg.Encode(f, screen, nil); err != nil {
		return fmt.Errorf("failed to encode img: %w", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("can't close tmp image: %w", err)
	}

	return nil
}
