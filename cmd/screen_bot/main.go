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
	dummyScript := base.NewDummyScript()
	startBot(ctx, dummyScript)

	go closeHandler(dummyScript, exit, "q")
	<-exit
	go func() {
		<-ctx.Done()
	}()
}

func startBot(ctx context.Context, script internal.Script) {
	go func() {
		err := script.Start(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
}

func closeHandler(script internal.Script, exit chan bool, closeKey string) {
	fmt.Println(fmt.Sprintf("--- Please press %s to exit---", closeKey))

	hook.Register(hook.KeyDown, []string{closeKey}, func(e hook.Event) {
		script.Stop()
		exit <- true
	})

	s := hook.Start()
	<-hook.Process(s)
}
