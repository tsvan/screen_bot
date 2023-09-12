package task

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"time"
)

type DummyTask struct {
	Screen *internal.Screen
	status internal.Status
}

func NewDummyTask(screen *internal.Screen, status internal.Status) *DummyTask {
	return &DummyTask{Screen: screen, status: status}
}

func (d *DummyTask) Exec(ctx context.Context, opts internal.TaskOpts) (internal.Status, error) {
	fmt.Println("dummy task exec -", opts.Name)
	d.Screen.CaptureScreen()
	time.Sleep(2 * time.Second)
	return d.status, nil
}
