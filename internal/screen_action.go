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

type ScreenActionOpts struct {
	Name        string
	DelayBefore time.Duration
	DelayAfter  time.Duration
}

type ScreenActionWithOpts struct {
	ScreenAction ScreenAction
	Opts         ScreenActionOpts
}

type ScreenAction interface {
	Handle(ctx context.Context, opts ScreenActionOpts) (Status, error)
	Close()
}
