package internal

import (
	"context"
	"fmt"
	"time"
)

type Status string

const (
	Complete   Status = "Complete"
	InProgress Status = "InProgress"
	Failed     Status = "Failed"
)

type TaskErrCode int

const (
	AttemptsEnd    TaskErrCode = 1
	InvalidCaptcha TaskErrCode = 2
)

type TaskOpts struct {
	Name        string
	DelayBefore time.Duration
	DelayAfter  time.Duration
	RunOnce     bool
}

type TaskWithOpts struct {
	Task Task
	Opts TaskOpts
}

type TaskErr struct {
	Err        error
	StatusCode TaskErrCode
}

type Task interface {
	Exec(ctx context.Context, opts TaskOpts) error
}

func (r *TaskErr) Error() string {
	return fmt.Sprintf("task error code %d. %s", r.StatusCode, r.Err)
}
