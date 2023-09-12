package login

import (
	"clicker_bot/internal"
	"clicker_bot/task/overall"
	"context"
	"fmt"
	"image"
)

type LoginTask struct {
}

func NewLoginTask() *LoginTask {
	return &LoginTask{}
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
	//screen := internal.NewScreen(image.Rectangle{}, internal.ScreenOpts{})

	//l2Start := internal.TaskWithOpts{
	//	Task: NewL2StartTask(),
	//	Opts: internal.TaskOpts{Name: "l2 start"},
	//}

	l2Login1 := internal.TaskWithOpts{
		Task: overall.NewFindAndClickTask(nil, "F:\\projects\\go\\screen_bot\\static\\l2\\login.PNG", image.Point{X: 10, Y: 10}),
		Opts: internal.TaskOpts{Name: "l2 login first screen"},
	}
	l2Login2 := internal.TaskWithOpts{
		Task: overall.NewFindAndClickTask(nil, "F:\\projects\\go\\screen_bot\\static\\l2\\accept.PNG", image.Point{X: 10, Y: 10}),
		Opts: internal.TaskOpts{Name: "l2 login second screen"},
	}

	l2Login3 := internal.TaskWithOpts{
		Task: overall.NewFindAndClickTask(nil, "F:\\projects\\go\\screen_bot\\static\\l2\\enter.PNG", image.Point{X: 10, Y: 10}),
		Opts: internal.TaskOpts{Name: "l2 login third screen"},
	}
	l2Login4 := internal.TaskWithOpts{
		Task: overall.NewFindAndClickTask(nil, "F:\\projects\\go\\screen_bot\\static\\l2\\start.PNG", image.Point{X: 10, Y: 10}),
		Opts: internal.TaskOpts{Name: "l2 login fourth screen"},
	}

	return []internal.TaskWithOpts{l2Login1, l2Login2, l2Login3, l2Login4}
}
