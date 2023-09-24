package overall

import (
	"clicker_bot/internal"
	"context"
	"errors"
	"fmt"
	"image"
	"time"
)

type FindAndActionTask struct {
	imgToFindPath string
	Opts          FindAndActionTaskOpts
	savedPoint    *image.Point
}

type FindAndActionTaskOpts struct {
	Screen         *internal.Screen
	ActionFunc     func(x, y int)
	PointOffset    image.Point
	Attempts       int
	WithSave       bool
	SavedW, SavedH int
}

func NewFindAndActionTask(imgToFindPath string, opts FindAndActionTaskOpts) *FindAndActionTask {
	return &FindAndActionTask{imgToFindPath: imgToFindPath, Opts: opts}
}

func (f *FindAndActionTask) Exec(ctx context.Context, opts internal.TaskOpts) error {
	i := 0
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return nil
		default:

		}

		// считаем количество попыток
		if i > f.Opts.Attempts {
			f.savedPoint = nil
			return &internal.TaskErr{Err: errors.New("attempts end"), StatusCode: internal.AttemptsEnd}
		}
		i++

		time.Sleep(opts.DelayBefore * time.Millisecond)

		if f.savedPoint == nil || !f.Opts.WithSave {
			point, err := f.Opts.Screen.FindOnScreen(f.imgToFindPath)
			if err != nil {
				return fmt.Errorf("find img err:%w", err)
			}
			if point == nil {
				continue
			}
			f.savedPoint = point
		} else {
			checkPoint, err := f.Opts.Screen.FindOnScreen(f.imgToFindPath, f.savedPoint.X-10, f.savedPoint.Y-10, f.Opts.SavedW, f.Opts.SavedH)
			if err != nil {
				return fmt.Errorf("find img err:%w", err)
			}
			if checkPoint == nil {
				continue
			}
		}

		f.Opts.ActionFunc(f.savedPoint.X, f.savedPoint.Y)

		time.Sleep(opts.DelayAfter * time.Millisecond)
		return nil
	}
}
