package screenaction

import (
	"clicker_bot/internal"
	"context"
	"fmt"
	"time"
)

type DummyScreenAction struct {
	Screen *internal.Screen
	status internal.Status
}

func NewDummyScreenAction(screen *internal.Screen, status internal.Status) *DummyScreenAction {
	return &DummyScreenAction{Screen: screen, status: status}
}

func (d *DummyScreenAction) Handle(ctx context.Context, opts internal.ScreenActionOpts) (internal.Status, error) {
	fmt.Println("dummy screen action handling -", opts.Name)
	d.Screen.CaptureScreen()
	time.Sleep(2*time.Second)
	return d.status, nil
}

func (d *DummyScreenAction) Close() {
	fmt.Println("close")
}
