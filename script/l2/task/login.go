package task

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type LoginScreenAction struct {
	Screen *internal.Screen
	status internal.Status
}

func NewLoginScreenAction(screen *internal.Screen, status internal.Status) *LoginScreenAction {
	return &LoginScreenAction{Screen: screen, status: status}
}

func (d *LoginScreenAction) Exec(ctx context.Context, opts internal.TaskOpts) (internal.Status, error) {
	fmt.Println("dummy screen action handling -", opts.Name)
	cmd := exec.Command("F:\\games\\asterios\\Asterios.exe", "/autoplay")
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	//d.Screen.CaptureScreen()
	time.Sleep(2 * time.Second)
	return internal.Complete, nil
}

func (d *LoginScreenAction) Close() {
	fmt.Println("close")
}
