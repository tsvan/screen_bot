package overall

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"image"
	"time"
)

type FindAndClickTask struct {
	Screen      *internal.Screen
	action      *internal.Action
	Path        string
	PointOffset image.Point
	Attempts    int
}

func NewFindAndClickTask(action *internal.Action, screen *internal.Screen, path string, pointOffset image.Point, attempts int) *FindAndClickTask {
	return &FindAndClickTask{action: action, Screen: screen, Path: path, PointOffset: pointOffset, Attempts: attempts}
}

func (f *FindAndClickTask) Exec(ctx context.Context, opts internal.TaskOpts) error {
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
		point, err := f.Screen.FindOnScreen(f.Path)
		if err != nil {
			return fmt.Errorf("find img err:%w", err)
		}
		if point == nil {
			continue
		}

		f.action.Click(point.X+f.PointOffset.X, point.Y+f.PointOffset.Y, true)

		time.Sleep(opts.DelayAfter * time.Millisecond)
		return nil
	}
}
