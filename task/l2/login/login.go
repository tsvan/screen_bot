package login

import (
	"clicker_bot/internal"
	"clicker_bot/task/overall"
	"context"
	"fmt"
	"image"
)

type LoginTask struct {
	action *internal.Action
}

func NewLoginTask(action *internal.Action) *LoginTask {
	return &LoginTask{action: action}
}

func (d *LoginTask) Exec(ctx context.Context, opts internal.TaskOpts) error {
	manager := internal.NewTaskManager(d.Init())
	err := manager.Run(ctx, internal.RunOptSequence)
	if err != nil {
		return fmt.Errorf("manager run err: %w", err)
	}
	return nil
}

func (d *LoginTask) Init() []internal.TaskWithOpts {
	l2Start := internal.TaskWithOpts{
		Task: NewL2StartTask(),
		Opts: internal.TaskOpts{Name: "l2 start", DelayAfter: 10000},
	}

	l2Login1 := internal.TaskWithOpts{
		Task: overall.NewFindAndClickTask(
			d.action, nil, "F:\\projects\\go\\screen_bot\\static\\l2\\login.PNG", image.Point{X: 10, Y: 10}, 10),
		Opts: internal.TaskOpts{Name: "l2 login first screen", DelayBefore: 1000},
	}
	l2Login2 := internal.TaskWithOpts{
		Task: overall.NewFindAndClickTask(
			d.action, nil, "F:\\projects\\go\\screen_bot\\static\\l2\\accept.PNG", image.Point{X: 10, Y: 10}, 15),
		Opts: internal.TaskOpts{Name: "l2 login second screen", DelayBefore: 1000},
	}

	l2Login3 := internal.TaskWithOpts{
		Task: overall.NewFindAndClickTask(
			d.action, nil, "F:\\projects\\go\\screen_bot\\static\\l2\\enter.PNG", image.Point{X: 10, Y: 10}, 15),
		Opts: internal.TaskOpts{Name: "l2 login third screen", DelayBefore: 1000},
	}
	l2Login4 := internal.TaskWithOpts{
		Task: overall.NewFindAndClickTask(
			d.action, nil, "F:\\projects\\go\\screen_bot\\static\\l2\\start.PNG", image.Point{X: 10, Y: 10}, 15),
		Opts: internal.TaskOpts{Name: "l2 login fourth screen", DelayBefore: 1000},
	}

	return []internal.TaskWithOpts{l2Start, l2Login1, l2Login2, l2Login3, l2Login4}
}
