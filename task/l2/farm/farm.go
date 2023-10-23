package farm

import (
	"clicker_bot/internal"
	"clicker_bot/task/l2"
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

		err = watcher.WatchScreen()
		if err != nil {
			fmt.Println(err)
		}
		d.Action(watcher)
		//fmt.Println(robotgo.GetMouseColor())
		fmt.Println(watcher.GameInfo.Enemy)
		time.Sleep(200 * time.Millisecond)
	}
}

func (d *Task) Action(watcher *Watcher) {
	info := watcher.GetWindowInfo()

	d.Attack(info)
	d.Heal(info, 50, "f5")
	d.AfterMobDefeat(info)
	d.NextTarget(info)
}

func (d *Task) Attack(info l2.GameInfo) {
	if info.Enemy.Hp == 100 {
		//fmt.Println("attack")
		d.action.KeyPress("f6")
		d.action.KeyPress("f1")
		return
	}
	if info.Enemy.Hp > 0 {
		d.action.KeyPress("f1")
	}
}

func (d *Task) Heal(info l2.GameInfo, hpVal float64, button string) {
	if info.Character.Hp <= hpVal {
		//fmt.Println("heal")
		d.action.KeyPress(button)
	}
}

func (d *Task) AfterMobDefeat(info l2.GameInfo) {
	if info.Character.HaveTarget && info.Enemy.Hp == 0 {
		//fmt.Println("pickup")
		d.action.KeyPress("f4")
		time.Sleep(200 * time.Millisecond)
		d.action.KeyPress("f4")
		time.Sleep(200 * time.Millisecond)
		d.action.KeyPress("f4")
		time.Sleep(200 * time.Millisecond)
		d.action.KeyPress("f4")
		time.Sleep(200 * time.Millisecond)
		d.action.KeyPress("f4")
		time.Sleep(200 * time.Millisecond)
		d.action.KeyPress("f4")
		time.Sleep(200 * time.Millisecond)
		d.action.KeyPress("f4")
		time.Sleep(200 * time.Millisecond)

		//собрать манор
		d.action.KeyPress("f7")
		time.Sleep(100 * time.Millisecond)

		d.action.KeyPress("esc")
		time.Sleep(100 * time.Millisecond)

	}
}

func (d *Task) NextTarget(info l2.GameInfo) {
	if info.Enemy.Hp == 0 && !info.Character.HaveTarget && info.Character.Hp > 50 {
		// next target
		//fmt.Println("next target")
		d.action.KeyPress("f2")
	}
}
