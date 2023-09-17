package internal

import (
	"fmt"
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
	robotgo.MoveClick(x, y, "left", double)
}

func (a *Action) Move(x, y int) {
	a.actionPool.Submit(func() {
		robotgo.Move(x, y)
	})
}

func (a *Action) KeyPress(key string) {
	//c pond ошибка, пока так оставил
	err := robotgo.KeyPress(key)
	if err != nil {
		fmt.Print("Key press err")
	}
}

func (a *Action) TypeStr(text string) {
	robotgo.TypeStr(text, 0, 300)
}
