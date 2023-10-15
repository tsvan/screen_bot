package main

import (
	"clicker_bot/internal"
	base_task "clicker_bot/task"
	"clicker_bot/task/l2/farm"
	"clicker_bot/task/l2/manor"
	_ "clicker_bot/task/l2/manor"
	"context"
	"fmt"
	"github.com/alitto/pond"
	"github.com/joho/godotenv"
	hook "github.com/robotn/gohook"
	"log"
	"os"
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
	task := selectTask(actionPool)

	startBot(ctx, task)

	go closeHandler(exit, "q")
	<-exit
	go func() {
		<-ctx.Done()
	}()
}

func selectTask(actionPool *internal.Action) internal.Task {
	task := os.Getenv("TASK_TYPE")
	switch task {
	case "farm":
		return farm.NewTask(actionPool)
	case "manor":
		return manor.NewTask(actionPool)
	default:
		return base_task.NewDummyTask()

	}
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
