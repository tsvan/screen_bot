package task

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"time"
)

type DummyTask struct {
}

func NewDummyTask() *DummyTask {
	return &DummyTask{}
}

func (d *DummyTask) Exec(ctx context.Context, opts internal.TaskOpts) error {
	time.Sleep(opts.DelayBefore * time.Millisecond)
	fmt.Print("Dummy task exec")
	time.Sleep(opts.DelayAfter * time.Millisecond)
	return nil
}
