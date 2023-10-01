package manor

import (
	"clicker_bot/internal"
	"clicker_bot/task/l2/login"
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
	err := manager.Run(ctx, internal.RunOptLoop, 10)
	if err != nil {
		return fmt.Errorf("manager run err: %w", err)
	}
	return nil
}

func (d *Task) Init() []internal.TaskWithOpts {
	_ = internal.TaskWithOpts{
		Task: login.NewTask(d.action),
		Opts: internal.TaskOpts{Name: "l2 login", DelayAfter: 15000, RunOnce: true},
	}

	//Sit/Stand в игре нужно назначить на f10
	_ = internal.TaskWithOpts{
		Task: overall.NewActionTask(func() {
			d.action.KeyPress("f10")
		}),
		Opts: internal.TaskOpts{Name: "after login", DelayBefore: 21000, RunOnce: true},
	}

	startManor := internal.TaskWithOpts{
		Task: overall.NewActionTask(func() {
			d.action.KeyPress("f1")
		}),
		Opts: internal.TaskOpts{Name: "manor start", DelayAfter: 100, DelayBefore: 100},
	}
	ManorTaskOpts := overall.FindAndActionTaskOpts{
		PointOffset: image.Point{X: 10, Y: 2},
		Attempts:    15,
		WithSave:    true,
		SavedH:      150,
		SavedW:      100,
	}

	ManorTaskOpts.ActionFunc = func(x, y int) {
		d.action.Click(x+ManorTaskOpts.PointOffset.X, y+ManorTaskOpts.PointOffset.Y, false)
	}

	handOverManor := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\handover_manor.PNG", ManorTaskOpts),
		Opts: internal.TaskOpts{Name: "hand over manor", DelayBefore: 100},
	}

	captcha := internal.TaskWithOpts{
		Task: NewCaptchaTaskTask("\\static\\l2\\check_manor.PNG", nil, d.action, 10),
		Opts: internal.TaskOpts{Name: "manor captcha", DelayBefore: 100},
	}

	SelectSeedTaskOpts := overall.FindAndActionTaskOpts{
		PointOffset: image.Point{X: 20, Y: 22},
		Attempts:    10,
		WithSave:    true,
		SavedW:      200,
		SavedH:      100,
	}
	SelectSeedTaskOpts.ActionFunc = func(x, y int) {
		d.action.Click(x+SelectSeedTaskOpts.PointOffset.X, y+SelectSeedTaskOpts.PointOffset.Y, true)
	}

	selectSeed := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\seed.PNG", SelectSeedTaskOpts),
		Opts: internal.TaskOpts{Name: "select seed type", DelayBefore: 200},
	}

	citySelect := internal.TaskWithOpts{
		Task: NewSelectCityTask("\\static\\l2\\select_city.PNG", nil, d.action, 1, 3, 10),
		Opts: internal.TaskOpts{Name: "select city and seed count", DelayBefore: 100},
	}

	ManorTaskOpts.PointOffset = image.Point{X: 10, Y: 2}
	confirmSeeds := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\confirm_seed.PNG", ManorTaskOpts),
		Opts: internal.TaskOpts{Name: "confirm seed count and city", DelayBefore: 100},
	}

	sellSeeds := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\sell_seed.PNG", ManorTaskOpts),
		Opts: internal.TaskOpts{Name: "sell seeds", DelayBefore: 100},
	}

	return []internal.TaskWithOpts{startManor, handOverManor, captcha, selectSeed, citySelect, confirmSeeds, sellSeeds}
}
