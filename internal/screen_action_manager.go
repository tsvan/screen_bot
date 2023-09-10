package internal

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

type RunOpt string

const (
	RunOptSequence       RunOpt = "Sequence"
	RunOptParallel       RunOpt = "Parallel"
	RunOptLoop           RunOpt = "Loop"
	RunOptSequenceStatus RunOpt = "SequenceStatus"
)

type ScreenActionManager struct {
	screenActions []ScreenActionWithOpts
}

func NewScreenActionManager(screenActions []ScreenActionWithOpts) *ScreenActionManager {
	return &ScreenActionManager{screenActions: screenActions}
}

func (s *ScreenActionManager) Run(ctx context.Context, runOpt RunOpt) error {
	if len(s.screenActions) == 0 {
		return errors.New("no screen actions to run")
	}
	switch runOpt {
	case RunOptSequence:
		return s.runSequence(ctx)
	case RunOptParallel:
		return s.runParallel(ctx)
	case RunOptLoop:
		return s.runLoop(ctx)
	case RunOptSequenceStatus:
		return s.runStatus(ctx)

	}

	return nil
}

func (s *ScreenActionManager) runSequence(ctx context.Context) error {
	for _, v := range s.screenActions {
		time.Sleep(v.Opts.DelayBefore * time.Millisecond)
		_, err := v.ScreenAction.Handle(ctx, v.Opts)
		if err != nil {
			return fmt.Errorf("screen action run err: %w", err)
		}
		time.Sleep(v.Opts.DelayAfter * time.Millisecond)
	}

	return nil
}

func (s *ScreenActionManager) runParallel(ctx context.Context) error {
	errs, ctx := errgroup.WithContext(ctx)
	for _, v := range s.screenActions {
		errs.Go(func() error {
			_, err := v.ScreenAction.Handle(ctx, v.Opts)
			if err != nil {
				return fmt.Errorf("screen action run err: %w", err)
			}
			return nil
		})
	}

	return errs.Wait()
}

func (s *ScreenActionManager) runLoop(ctx context.Context) error {
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return nil
		default:

		}

		for _, v := range s.screenActions {
			time.Sleep(v.Opts.DelayBefore * time.Millisecond)
			_, err := v.ScreenAction.Handle(ctx, v.Opts)
			if err != nil {
				return fmt.Errorf("screen action run err: %w", err)
			}
			time.Sleep(v.Opts.DelayAfter * time.Millisecond)
		}

	}
}

func (s *ScreenActionManager) runStatus(ctx context.Context) error {
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return nil
		default:

		}
		for _, v := range s.screenActions {
			time.Sleep(v.Opts.DelayBefore * time.Millisecond)
			status, err := v.ScreenAction.Handle(ctx, v.Opts)
			if err != nil {
				return fmt.Errorf("screen action run err: %w", err)
			}
			time.Sleep(v.Opts.DelayAfter * time.Millisecond)
			// Выполняем задачи друг за другом, пока не выполнится первая, оставльные не запускаем
			if status == "" || status == InProgress {
				break
			} else if status == Failed {
				return fmt.Errorf("screen action (%s) status - %s", v.Opts.Name, status)
			}
		}
	}

}
