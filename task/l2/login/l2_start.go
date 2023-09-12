package login

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"os/exec"
)

type L2StartTask struct {
}

func NewL2StartTask() *L2StartTask {
	return &L2StartTask{}
}

func (d *L2StartTask) Exec(_ context.Context, opts internal.TaskOpts) error {
	cmd := exec.Command("F:\\games\\asterios\\Asterios.exe", "/autoplay")
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("l2 start failed: %w", err)
	}
	return nil
}
