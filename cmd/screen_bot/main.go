package main

import (
	"clicker_bot/internal"
	"clicker_bot/task/l2/farm"
	_ "clicker_bot/task/l2/manor"
	"context"
	"fmt"
	"github.com/alitto/pond"
	"github.com/joho/godotenv"
	hook "github.com/robotn/gohook"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var exit = make(chan bool)
	pool := pond.New(1, 10)
	actionPool := internal.NewAction(pool)
	task := farm.NewTask(actionPool)
	startBot(ctx, task)

	go closeHandler(exit, "q")
	<-exit
	go func() {
		<-ctx.Done()
	}()
}

func startBot(ctx context.Context, script internal.Task) {
	go func() {
		err := script.Exec(ctx, internal.TaskOpts{})
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
