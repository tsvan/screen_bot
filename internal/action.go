package internal

import (
	"github.com/alitto/pond"
	"github.com/go-vgo/robotgo"
)

type Action struct {
	actionPool *pond.WorkerPool
	opts       *ActionOpts
}

type ActionOpts struct {
}

func NewAction(actionPool *pond.WorkerPool, opts *ActionOpts) *Action {
	return &Action{actionPool: actionPool, opts: opts}
}

func (a *Action) Click(x, y int, double bool) {
	a.actionPool.Submit(func() {
		robotgo.MoveClick(x, y, "left", double)
	})
}

func (a *Action) Move(x, y int) {
	a.actionPool.Submit(func() {
		robotgo.Move(x, y)
	})
}
