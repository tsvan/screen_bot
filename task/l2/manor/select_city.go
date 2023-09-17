package manor

import (
	"clicker_bot/internal"
	"context"
	"fmt"
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

		//находим меню выбора города
		point, err := f.Screen.FindOnScreen(f.imgToFindPath)
		if err != nil {
			return fmt.Errorf("find img err:%w", err)
		}
		if point == nil {
			continue
		}

		// вводим количество семян для сдачи
		fmt.Println(f.SeedNumber)
		f.action.KeyPress(strconv.Itoa(f.SeedNumber))
		//Выбираем город
		f.action.Click(point.X+20, point.Y+5, false)
		time.Sleep(200 * time.Millisecond)
		f.action.Click(point.X+20, point.Y+15*f.CityNumber, false)

		time.Sleep(opts.DelayAfter * time.Millisecond)
		return nil
	}
}
