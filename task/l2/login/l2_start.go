package login

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"os/exec"
)

type L2StartTask struct {
}

func NewL2StartTask() *L2StartTask {
	return &L2StartTask{}
}

func (d *L2StartTask) Exec(_ context.Context, opts internal.TaskOpts) error {
	err := d.KillProcess("AsteriosGame.exe")
	if err != nil {
		return fmt.Errorf("can't kill fpid: %w", err)
	}

	cmd := exec.Command("F:\\games\\asterios\\Asterios.exe", "/autoplay")
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("l2 start failed: %w", err)
	}
	return nil
}

// KillProcess пришлось это использовать. robotgo.Kill не работает даже из под админа
func (d *L2StartTask) KillProcess(name string) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			return err
		}
		if n == name {
			return p.Kill()
		}
	}
	return nil
}
