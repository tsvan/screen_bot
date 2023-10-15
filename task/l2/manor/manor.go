package manor

import (
	"clicker_bot/internal"
	"clicker_bot/task/l2/login"
	"clicker_bot/task/overall"
	"context"
	"fmt"
	"image"
	"os"
	"strconv"
)

const (
	DefaultDelay    = 100
	DefaultAttempts = 30
)

type Task struct {
	action *internal.Action
}

func NewTask(action *internal.Action) *Task {
	return &Task{action: action}
}

func (d *Task) Exec(ctx context.Context, opts internal.TaskOpts) error {
	manager := internal.NewTaskManager(d.Init())
	err := manager.Run(ctx, internal.RunOptLoop, DefaultAttempts)
	if err != nil {
		return fmt.Errorf("manager run err: %w", err)
	}
	return nil
}

func (d *Task) Init() []internal.TaskWithOpts {
	l2Login := internal.TaskWithOpts{
		Task: login.NewTask(d.action),
		Opts: internal.TaskOpts{Name: "l2 login", DelayAfter: 15000, RunType: internal.Disabled},
	}

	//Sit/Stand в игре нужно назначить на f10
	l2Sit := internal.TaskWithOpts{
		Task: overall.NewActionTask(func() {
			d.action.KeyPress("f10")
		}),
		Opts: internal.TaskOpts{Name: "after login", DelayBefore: 21000, RunType: internal.Disabled},
	}

	startManor := internal.TaskWithOpts{
		Task: overall.NewActionTask(func() {
			d.action.KeyPress("f1")
		}),
		Opts: internal.TaskOpts{Name: "manor start", DelayAfter: DefaultDelay, DelayBefore: DefaultDelay},
	}
	ManorTaskOpts := overall.FindAndActionTaskOpts{PointOffset: image.Point{X: 10, Y: 2}, Attempts: DefaultAttempts, WithSave: true, SavedH: 150, SavedW: 100}

	ManorTaskOpts.ActionFunc = func(x, y int) {
		d.action.Click(x+ManorTaskOpts.PointOffset.X, y+ManorTaskOpts.PointOffset.Y, false)
	}

	handOverManor := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\handover_manor.PNG", ManorTaskOpts),
		Opts: internal.TaskOpts{Name: "hand over manor", DelayBefore: DefaultDelay},
	}

	captcha := internal.TaskWithOpts{
		Task: NewCaptchaTaskTask("\\static\\l2\\check_manor.PNG", nil, d.action, DefaultAttempts),
		Opts: internal.TaskOpts{Name: "manor captcha", DelayBefore: DefaultDelay},
	}

	SelectSeedTaskOpts := overall.FindAndActionTaskOpts{PointOffset: image.Point{X: 20, Y: 22}, Attempts: DefaultAttempts, WithSave: true, SavedW: 200, SavedH: 100}
	SelectSeedTaskOpts.ActionFunc = func(x, y int) {
		d.action.Click(x+SelectSeedTaskOpts.PointOffset.X, y+SelectSeedTaskOpts.PointOffset.Y, true)
	}

	selectSeed := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\seed.PNG", SelectSeedTaskOpts),
		Opts: internal.TaskOpts{Name: "select seed type", DelayBefore: DefaultDelay},
	}

	citySelect := internal.TaskWithOpts{
		Task: NewSelectCityTask("\\static\\l2\\select_city.PNG", nil, d.action, d.GetSeedNumber(), d.GetCityNumber(), DefaultAttempts),
		Opts: internal.TaskOpts{Name: "select city and seed count", DelayBefore: DefaultDelay},
	}

	ManorTaskOpts.PointOffset = image.Point{X: 10, Y: 2}
	confirmSeeds := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\confirm_seed.PNG", ManorTaskOpts),
		Opts: internal.TaskOpts{Name: "confirm seed count and city", DelayBefore: DefaultDelay},
	}

	sellSeeds := internal.TaskWithOpts{
		Task: overall.NewFindAndActionTask("\\static\\l2\\sell_seed.PNG", ManorTaskOpts),
		Opts: internal.TaskOpts{Name: "sell seeds", DelayBefore: DefaultDelay},
	}

	return []internal.TaskWithOpts{l2Login, l2Sit, startManor, handOverManor, captcha, selectSeed, citySelect, confirmSeeds, sellSeeds}
}

func (d *Task) GetSeedNumber() int {
	seedNumber, err := strconv.Atoi(os.Getenv("MANOR_SEED_NUMBER"))
	if err != nil {
		panic("can't parse seed number")
	}
	return seedNumber
}

func (d *Task) GetCityNumber() int {
	cityNumber, err := strconv.Atoi(os.Getenv("MANOR_CITY_NUMBER"))
	if err != nil {
		panic("can't parse city number")
	}
	return cityNumber
}
