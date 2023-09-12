package l2

import (
	"clicker_bot/internal"
	"clicker_bot/script/l2/task"
	"context"
	"fmt"
	"image"
)

type ManorTask struct {
}

func NewManorTask() *ManorTask {
	return &ManorTask{}
}

func (d *ManorTask) Exec(ctx context.Context) error {
	manager := internal.NewTaskManager(d.Init())
	err := manager.Run(ctx, internal.RunOptSequence)
	if err != nil {
		return fmt.Errorf("manager run err: %w", err)
	}
	return nil
}

func (d *ManorTask) Init() []internal.TaskWithOpts {
	screen := internal.NewScreen(image.Rectangle{}, internal.ScreenOpts{})

	login := internal.TaskWithOpts{
		Task: task.NewLoginScreenAction(screen, internal.InProgress),
		Opts: internal.TaskOpts{
			Name: "login",
		},
	}

	return []internal.TaskWithOpts{login}
}
