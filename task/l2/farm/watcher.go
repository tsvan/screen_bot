package farm

import (
	"clicker_bot/internal"
	"clicker_bot/task/l2"
	"fmt"
	"image"
	"image/color"
)

type Watcher struct {
	Screen   *internal.Screen
	GameInfo l2.GameInfo

	savedHpBar *image.Rectangle
	maxHpCount int

	savedMobHpBar *image.Rectangle
	maxMobHpCount int
}

func NewWatcher() *Watcher {
	return &Watcher{}
}

func (w *Watcher) WatchScreen() {
	w.GameInfo.Character.Hp = w.CharacterBar(color.RGBA{R: 121, G: 28, B: 17, A: 255})
	w.GameInfo.Enemy.Hp = w.MobBar(color.RGBA{R: 111, G: 23, B: 20, A: 255})
	if w.savedMobHpBar != nil {
		defeatedMobHp := w.Screen.FindColorCount(color.RGBA{R: 47, G: 26, B: 23, A: 255},
			w.savedMobHpBar.Min.X-2,
			w.savedMobHpBar.Min.Y-1,
			w.savedMobHpBar.Min.X+4,
			w.savedMobHpBar.Min.Y+2,
		)
		if defeatedMobHp > 0 {
			w.GameInfo.Enemy.WasDefeated = true
		} else {
			w.GameInfo.Enemy.WasDefeated = false
		}
	}
}

func (w *Watcher) GetWindowInfo() l2.GameInfo {
	return w.GameInfo
}

func (w *Watcher) FindSize() error {
	beginPoint, endPoint, err := w.findWindowSize()
	if err != nil || beginPoint == nil || endPoint == nil {
		return fmt.Errorf("can't find window size:%w", err)
	}
	w.GameInfo.GameWindow = l2.WindowSize{X: beginPoint.X, Y: beginPoint.X, W: endPoint.X - beginPoint.X, H: endPoint.Y - beginPoint.Y}
	return nil
}

func (w *Watcher) findWindowSize() (*image.Point, *image.Point, error) {
	beginPoint, err := w.Screen.FindOnScreen("\\static\\l2\\farm\\begin_window.png")
	if err != nil {
		return nil, nil, fmt.Errorf("can't find start window size:%w", err)
	}

	endPoint, err := w.Screen.FindOnScreen("\\static\\l2\\farm\\end_window.png")
	if err != nil {
		return nil, nil, fmt.Errorf("can't find end window size:%w", err)
	}

	return beginPoint, endPoint, nil
}

func (w *Watcher) CharacterBar(color color.RGBA) float64 {
	if w.savedHpBar == nil {
		count, pointMin, pointMax := w.Screen.FindColorBounds(color)

		if count > 0 {
			w.savedHpBar = &image.Rectangle{
				Min: pointMin,
				Max: pointMax,
			}
			w.maxHpCount = count
			return 100
		}
	} else {
		hp := w.Screen.FindColorCount(color,
			w.savedHpBar.Min.X-2,
			w.savedHpBar.Min.Y-1,
			w.savedHpBar.Max.X-w.savedHpBar.Min.X+4,
			w.savedHpBar.Max.Y-w.savedHpBar.Min.Y+2,
		)
		return (float64(hp) / float64(w.maxHpCount)) * 100
	}
	return 0
}

func (w *Watcher) MobBar(color color.RGBA) float64 {
	if w.savedMobHpBar == nil {
		count, pointMin, pointMax := w.Screen.FindColorBounds(color)

		if count > 0 {
			w.savedMobHpBar = &image.Rectangle{
				Min: pointMin,
				Max: pointMax,
			}
			w.maxMobHpCount = count
			return 100
		}
	} else {
		hp := w.Screen.FindColorCount(color,
			w.savedMobHpBar.Min.X-2,
			w.savedMobHpBar.Min.Y-1,
			w.savedMobHpBar.Max.X-w.savedMobHpBar.Min.X+4,
			w.savedMobHpBar.Max.Y-w.savedMobHpBar.Min.Y+2,
		)
		return (float64(hp) / float64(w.maxMobHpCount)) * 100
	}
	return 0
}
