package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

type RunOpt string

const (
	RunOptSequence       RunOpt = "Sequence"
	RunOptParallel       RunOpt = "Parallel"
	RunOptLoop           RunOpt = "Loop"
	RunOptSequenceStatus RunOpt = "SequenceStatus"
)

type TaskManager struct {
	tasks []TaskWithOpts
}

func NewTaskManager(tasks []TaskWithOpts) *TaskManager {
	return &TaskManager{tasks: tasks}
}

func (s *TaskManager) Run(ctx context.Context, runOpt RunOpt) error {
	if len(s.tasks) == 0 {
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

func (s *TaskManager) runSequence(ctx context.Context) error {
	for _, v := range s.tasks {
		time.Sleep(v.Opts.DelayBefore * time.Millisecond)
		_, err := v.Task.Exec(ctx, v.Opts)
		if err != nil {
			return fmt.Errorf("screen action run err: %w", err)
		}
		time.Sleep(v.Opts.DelayAfter * time.Millisecond)
	}

	return nil
}

func (s *TaskManager) runParallel(ctx context.Context) error {
	errs, ctx := errgroup.WithContext(ctx)
	for _, v := range s.tasks {
		errs.Go(func() error {
			_, err := v.Task.Exec(ctx, v.Opts)
			if err != nil {
				return fmt.Errorf("screen action run err: %w", err)
			}
			return nil
		})
	}

	return errs.Wait()
}

func (s *TaskManager) runLoop(ctx context.Context) error {
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return nil
		default:

		}

		for _, v := range s.tasks {
			time.Sleep(v.Opts.DelayBefore * time.Millisecond)
			_, err := v.Task.Exec(ctx, v.Opts)
			if err != nil {
				return fmt.Errorf("screen action run err: %w", err)
			}
			time.Sleep(v.Opts.DelayAfter * time.Millisecond)
		}

	}
}

func (s *TaskManager) runStatus(ctx context.Context) error {
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return nil
		default:

		}
		for _, v := range s.tasks {
			time.Sleep(v.Opts.DelayBefore * time.Millisecond)
			status, err := v.Task.Exec(ctx, v.Opts)
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
