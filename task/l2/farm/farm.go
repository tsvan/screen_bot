package farm

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"time"
)

type Task struct {
	action *internal.Action
	Screen *internal.Screen
}

func NewTask(action *internal.Action) *Task {
	return &Task{action: action}
}

func (d *Task) Exec(ctx context.Context, opts internal.TaskOpts) error {
	time.Sleep(2 * time.Second)
	watcher := NewWatcher()
	err := d.Start(ctx, watcher)
	if err != nil {
		return fmt.Errorf("farm err:%w", err)
	}
	return nil
}

func (d *Task) Start(ctx context.Context, watcher *Watcher) error {
	err := watcher.FindSize()
	if err != nil {
		return fmt.Errorf("l2 size :%w", err)
	}

	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return nil
		default:

		}
		watcher.WatchScreen()
		d.Action(watcher)
		//fmt.Println(robotgo.GetMouseColor())
		time.Sleep(100 * time.Millisecond)
	}
}

func (d *Task) Action(watcher *Watcher) {
	info := watcher.GetWindowInfo()

	if info.Enemy.Hp == 100 {
		fmt.Println("attack")
		d.action.KeyPress("f1")
	}

	if info.Character.Hp <= 60 {
		fmt.Println("heal")
	}

	if info.Enemy.Hp == 0 && info.Enemy.WasDefeated {
		fmt.Println("pickup")
		d.action.KeyPress("f4")
	}

	if info.Enemy.Hp == 0 && !info.Enemy.WasDefeated && info.Character.Hp > 60 {
		// next target
		fmt.Println("next target")
		d.action.KeyPress("f2")
	}
	fmt.Println(info)
}
