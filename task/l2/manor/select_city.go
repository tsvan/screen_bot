package manor

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"image"
	"strconv"
	"time"
)

type SelectCityTask struct {
	imgToFindPath string
	Screen        *internal.Screen
	action        *internal.Action
	Attempts      int
	SeedNumber    int
	CityNumber    int
	savedPoint    *image.Point
}

func NewSelectCityTask(imgToFindPath string, screen *internal.Screen, action *internal.Action, seedNumber, cityNumber, attempts int) *SelectCityTask {
	return &SelectCityTask{
		imgToFindPath: imgToFindPath,
		Screen:        screen,
		action:        action,
		SeedNumber:    seedNumber,
		CityNumber:    cityNumber,
		Attempts:      attempts,
	}
}

func (f *SelectCityTask) Exec(ctx context.Context, opts internal.TaskOpts) error {
	i := 0
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return nil
		default:

		}

		// считаем количество попыток
		if i > f.Attempts {
			return &internal.TaskErr{StatusCode: internal.AttemptsEnd}
		}
		i++

		time.Sleep(opts.DelayBefore * time.Millisecond)

		var clickPoint image.Point
		if f.savedPoint == nil {
			point, err := f.Screen.FindOnScreen(f.imgToFindPath)
			if err != nil {
				return fmt.Errorf("find img err:%w", err)
			}
			if point == nil {
				continue
			}
			f.savedPoint = point
			clickPoint = *point
		} else {
			checkPoint, err := f.Screen.FindOnScreen(f.imgToFindPath, f.savedPoint.X-20, f.savedPoint.Y-10, 170, 100)
			if err != nil {
				return fmt.Errorf("find img err:%w", err)
			}
			if checkPoint == nil {
				continue
			}
			clickPoint.X = f.savedPoint.X - 20 + checkPoint.X
			clickPoint.Y = f.savedPoint.Y - 10 + checkPoint.Y
		}

		// вводим количество семян для сдачи
		f.action.KeyPress(strconv.Itoa(f.SeedNumber))
		//Выбираем город
		f.action.Click(clickPoint.X+20, clickPoint.Y+5, false)
		// ждём появления выпадающего меню
		time.Sleep(DefaultDelay * time.Millisecond)
		f.action.Click(clickPoint.X+20, clickPoint.Y+30+15*f.CityNumber, false)

		time.Sleep(opts.DelayAfter * time.Millisecond)
		return nil
	}
}
