package overall

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"github.com/go-vgo/robotgo"
	"image"
	"time"
)

type FindAndClickTask struct {
	Screen      *internal.Screen
	Path        string
	PointOffset image.Point
}

func NewFindAndClickTask(screen *internal.Screen, path string, pointOffset image.Point) *FindAndClickTask {
	return &FindAndClickTask{Screen: screen, Path: path, PointOffset: pointOffset}
}

func (f *FindAndClickTask) Exec(ctx context.Context, opts internal.TaskOpts) error {
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return nil
		default:

		}

		time.Sleep(opts.DelayBefore * time.Millisecond)
		point, err := f.Screen.FindOnScreen(f.Path)
		if err != nil {
			return fmt.Errorf("find img err:%w", err)
		}
		if point == nil {
			continue
		}

		robotgo.MoveClick(point.X+10, point.Y+10)

		time.Sleep(opts.DelayAfter * time.Millisecond)
		return nil
	}
}
