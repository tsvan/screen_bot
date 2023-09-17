package main

import (
	"clicker_bot/internal"
	"clicker_bot/task/l2/manor"
	"context"
	"fmt"
	"github.com/alitto/pond"
	hook "github.com/robotn/gohook"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var exit = make(chan bool)
	pool := pond.New(1, 10)
	actionPool := internal.NewAction(pool, nil)
	task := manor.NewManorTask(actionPool)
	startBot(ctx, task)

	go closeHandler(exit, "q")
	<-exit
	go func() {
		<-ctx.Done()
	}()
}

func startBot(ctx context.Context, script internal.Task) {
	go func() {
		err := script.Exec(ctx, internal.TaskOpts{DelayBefore: 5000})
		if err != nil {
			fmt.Println(err)
		}
	}()
}

func closeHandler(exit chan bool, closeKey string) {
	fmt.Println(fmt.Sprintf("--- Please press %s to exit---", closeKey))

	hook.Register(hook.KeyDown, []string{closeKey}, func(e hook.Event) {
		exit <- true
	})

	s := hook.Start()
	<-hook.Process(s)
}
