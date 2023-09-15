package overall

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"image"
	"time"
)

type FindAndActionTask struct {
	imgToFindPath string
	Opts          FindAndActionTaskOpts
}

type FindAndActionTaskOpts struct {
	Screen      *internal.Screen
	ActionFunc  func(x, y int)
	PointOffset image.Point
	Attempts    int
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
			return &internal.TaskErr{StatusCode: internal.AttemptsEnd}
		}
		i++

		time.Sleep(opts.DelayBefore * time.Millisecond)
		point, err := f.Opts.Screen.FindOnScreen(f.imgToFindPath)
		if err != nil {
			return fmt.Errorf("find img err:%w", err)
		}
		if point == nil {
			continue
		}
		f.Opts.ActionFunc(point.X, point.Y)

		time.Sleep(opts.DelayAfter * time.Millisecond)
		return nil
	}
}
