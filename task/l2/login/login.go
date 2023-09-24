package login

import (
	"clicker_bot/internal"
	"clicker_bot/task/overall"
	"context"
	"fmt"
	"image"
)

type Task struct {
	action *internal.Action
}

func NewTask(action *internal.Action) *Task {
	return &Task{action: action}
}

func (d *Task) Exec(ctx context.Context, opts internal.TaskOpts) error {
	manager := internal.NewTaskManager(d.Init())
	err := manager.Run(ctx, internal.RunOptSequence, 0)
	if err != nil {
		return fmt.Errorf("manager run err: %w", err)
	}
	return nil
}

func (d *Task) Init() []internal.TaskWithOpts {
	l2Start := internal.TaskWithOpts{
		Task: NewL2StartTask(),
		Opts: internal.TaskOpts{Name: "l2 start", DelayAfter: 10000},
	}

	loginTaskOpts := overall.FindAndActionTaskOpts{
		Screen:      nil,
		PointOffset: image.Point{X: 10, Y: 10},
		Attempts:    15,
	}

	loginTaskOpts.ActionFunc = func(x, y int) {
		d.action.Click(x+loginTaskOpts.PointOffset.X, y+loginTaskOpts.PointOffset.Y, true)
		// что бы курсор не перкрывал изображение
		if x > 30 && y > 30 {
			d.action.Move(x-30, y-30)
		}
	}

	l2Login1 := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\login.PNG", loginTaskOpts),
		Opts: internal.TaskOpts{Name: "l2 login first screen", DelayBefore: 1000},
	}
	l2Login2 := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\accept.PNG", loginTaskOpts),
		Opts: internal.TaskOpts{Name: "l2 login second screen", DelayBefore: 1000},
	}

	l2Login3 := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\enter.PNG", loginTaskOpts),
		Opts: internal.TaskOpts{Name: "l2 login third screen", DelayBefore: 1000},
	}
	l2Login4 := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\start.PNG", loginTaskOpts),
		Opts: internal.TaskOpts{Name: "l2 login fourth screen", DelayBefore: 1000},
	}

	return []internal.TaskWithOpts{l2Start, l2Login1, l2Login2, l2Login3, l2Login4}
}
