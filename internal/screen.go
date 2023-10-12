package internal

import (
	"bytes"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
	"gocv.io/x/gocv"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"os"
	"os/exec"
)

// Соответсвие при поиске изображения
const ImgSearchRate = 0.9

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

func (s *Screen) FindColorCount(color color.RGBA, opt ...int) int {
	img := s.CaptureScreen(opt...)
	count := 0
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			imgColor := img.At(x, y)
			if imgColor == color {
				count++
			}
			//fmt.Printf("%v", color)
		}
	}
	return count
}

func (s *Screen) FindColorBounds(color color.RGBA) (int, image.Point, image.Point) {
	img := s.CaptureScreen()
	count := 0
	min := image.Point{}
	max := image.Point{}
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			imgColor := img.At(x, y)
			if imgColor == color {
				if count == 0 {
					min = image.Point{X: x, Y: y}
				}
				count++
				max = image.Point{X: x, Y: y}
			}
		}
	}
	return count, min, max
}

func image_2_array_pix(src image.Image) [][][3]float32 {
	bounds := src.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	iaa := make([][][3]float32, height)
	src_rgba := image.NewRGBA(src.Bounds())
	draw.Copy(src_rgba, image.Point{}, src, src.Bounds(), draw.Src, nil)

	for y := 0; y < height; y++ {
		row := make([][3]float32, width)
		for x := 0; x < width; x++ {
			idx_s := (y*width + x) * 4
			pix := src_rgba.Pix[idx_s : idx_s+4]
			fmt.Println(pix)
			row[x] = [3]float32{float32(pix[0]), float32(pix[1]), float32(pix[2])}
		}
		iaa[y] = row
	}
	fmt.Println("ok")
	return iaa
}

func (s *Screen) FindText(imgPath string) (string, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("can't get root directory path")
	}

	tesseratPath := os.Getenv("TESSERACT_PATH")
	cmd := exec.Command(tesseratPath, rootPath+imgPath, "-", "-l", "rus+eng", "-c", "tessedit_char_whitelist=+?0123456789")
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

	matGray := gocv.NewMat()
	matImg, err := gocv.ImageToMatRGB(screen)
	if err != nil {
		return fmt.Errorf("can't convert screen to gocv Mat: %w", err)
	}
	//gocv.AdaptiveThreshold( matImg,&matGray, 255.0, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, 5, 4.0)

	gocv.CvtColor(matImg, &matGray, gocv.ColorRGBToGray)
	gocv.IMWrite(rootPath+path, matGray)

	//grayScreen := image.NewGray(screen.Bounds())
	//draw.Draw(grayScreen, grayScreen.Bounds(), screen, screen.Bounds().Min, draw.Src)
	//f, err := os.Create(rootPath + path)
	//if err != nil {
	//	return fmt.Errorf("can't create tmp image: %w", err)
	//}
	//
	//if err = jpeg.Encode(f, grayScreen, &jpeg.Options{Quality: 100 }); err != nil {
	//	return fmt.Errorf("failed to encode img: %w", err)
	//}
	//
	//if err = f.Close(); err != nil {
	//	return fmt.Errorf("can't close tmp image: %w", err)
	//}

	return nil
}

//matGray := gocv.NewMat()
//matImg, err := gocv.ImageGrayToMatGray(grayScreen)
//if err != nil {
//return fmt.Errorf("can't convert screen to gocv Mat: %w", err)
//}
//gocv.AdaptiveThreshold( matImg,&matGray, 255.0, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, 5, 4.0)
//
////gocv.CvtColor(matImg, &matGray, gocv.ColorRGBToGray)
//gocv.IMWrite(rootPath + path, matGray)
