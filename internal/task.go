package internal

import (
	"context"
	"time"
)

type Status string

const (
	Complete   Status = "Complete"
	InProgress Status = "InProgress"
	Failed     Status = "Failed"
)

type TaskOpts struct {
	Name        string
	DelayBefore time.Duration
	DelayAfter  time.Duration
}

type TaskWithOpts struct {
	Task Task
	Opts TaskOpts
}

type Task interface {
	Exec(ctx context.Context, opts TaskOpts) (Status, error)
}
