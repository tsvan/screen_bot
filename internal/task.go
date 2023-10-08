package internal

import (
	"context"
	"fmt"
	"time"
)

type RunType int

const (
	Default  RunType = 0
	RunOnce  RunType = 1
	Disabled RunType = 2
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
	RunType     RunType
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
