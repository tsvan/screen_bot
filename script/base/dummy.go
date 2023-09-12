package base

import (
	"clicker_bot/internal"
	"clicker_bot/script/base/task"
	"context"
	"fmt"
	"image"
)

type DummyMainTask struct {
}

func NewDummyMainTask() *DummyMainTask {
	return &DummyMainTask{}
}

func (d *DummyMainTask) Exec(ctx context.Context, opts internal.TaskOpts) (internal.Status, error) {
	manager := internal.NewTaskManager(d.InitScrAct())
	err := manager.Run(ctx, internal.RunOptSequenceStatus)
	if err != nil {
		return internal.Failed, fmt.Errorf("manager run err: %w", err)
	}
	fmt.Println("dummy main task exec")
	return internal.Complete, nil
}

func (d *DummyMainTask) InitScrAct() []internal.TaskWithOpts {
	screen := internal.NewScreen(image.Rectangle{
		Min: image.Point{X: 1, Y: 1},
		Max: image.Point{X: 100, Y: 100},
	}, internal.ScreenOpts{CaptureDelay: 100})

	dummy1 := internal.TaskWithOpts{
		Task: task.NewDummyTask(screen, internal.InProgress),
		Opts: internal.TaskOpts{
			Name: "dummy1",
		},
	}

	dummy2 := internal.TaskWithOpts{
		Task: task.NewDummyTask(screen, internal.InProgress),
		Opts: internal.TaskOpts{
			DelayBefore: 200,
			DelayAfter:  300,
			Name:        "dummy2",
		},
	}
	return []internal.TaskWithOpts{dummy1, dummy2}
}
