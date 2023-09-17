package overall

import (
	"clicker_bot/internal"
	"context"
	"time"
)

type ActionTask struct {
	ActionFunc func()
}

func NewActionTask(ActionFunc func()) *ActionTask {
	return &ActionTask{ActionFunc: ActionFunc}
}

func (f *ActionTask) Exec(ctx context.Context, opts internal.TaskOpts) error {
	time.Sleep(opts.DelayBefore * time.Millisecond)
	f.ActionFunc()
	time.Sleep(opts.DelayAfter * time.Millisecond)
	return nil
}
