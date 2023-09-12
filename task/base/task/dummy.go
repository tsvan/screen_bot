package task

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"time"
)

type DummyTask struct {
	Screen *internal.Screen
}

func NewDummyTask(screen *internal.Screen) *DummyTask {
	return &DummyTask{Screen: screen}
}

func (d *DummyTask) Exec(ctx context.Context, opts internal.TaskOpts) error {
	fmt.Println("dummy task exec -", opts.Name)
	d.Screen.CaptureScreen()
	time.Sleep(2 * time.Second)
	return nil
}
