package main

import (
	"clicker_bot/internal"
	"clicker_bot/script/base"
	"context"
	"fmt"
	hook "github.com/robotn/gohook"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var exit = make(chan bool)
	dummy := base.NewDummyMainTask()
	startBot(ctx, dummy)

	go closeHandler(exit, "q")
	<-exit
	go func() {
		<-ctx.Done()
	}()
}

func startBot(ctx context.Context, script internal.Task) {
	go func() {
		_, err := script.Exec(ctx, internal.TaskOpts{})
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
