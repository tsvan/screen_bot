package base

import (
	"clicker_bot/internal"
	"clicker_bot/script/base/screenaction"
	"context"
	"fmt"
	"image"
)

type DummyScript struct {
}

func NewDummyScript() *DummyScript {
	return &DummyScript{}
}

func (d*DummyScript) Start(ctx context.Context) error {
	manager := internal.NewScreenActionManager(d.InitScrAct())
	err:=manager.Run(ctx, internal.RunOptSequenceStatus)
	if err!= nil {
		return fmt.Errorf("manager run err: %w", err)
	}
	fmt.Println("dummy script running")
	return nil
}

func (d *DummyScript) InitScrAct() []internal.ScreenActionWithOpts {
	screen:= internal.NewScreen(image.Rectangle{
		Min: image.Point{X: 1, Y: 1},
		Max: image.Point{X: 100, Y: 100},
	}, internal.ScreenOpts{CaptureDelay:100})

	dummy1:= internal.ScreenActionWithOpts{
		ScreenAction: screenaction.NewDummyScreenAction(screen, internal.InProgress),
		Opts: internal.ScreenActionOpts{
			Name: "dummy1",
		},
	}

	dummy2:= internal.ScreenActionWithOpts{
		ScreenAction: screenaction.NewDummyScreenAction(screen,internal.InProgress),
		Opts: internal.ScreenActionOpts{
			DelayBefore: 200,
			DelayAfter:  300,
			Name: "dummy2",
		},
	}
	return []internal.ScreenActionWithOpts{dummy1, dummy2}
}

func (d *DummyScript) Stop() {
	fmt.Println("stop")
}
