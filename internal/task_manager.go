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
	RunOptSequence RunOpt = "Sequence"
	RunOptParallel RunOpt = "Parallel"
	RunOptLoop     RunOpt = "Loop"
)

type TaskManager struct {
	tasks []TaskWithOpts
}

func NewTaskManager(tasks []TaskWithOpts) *TaskManager {
	return &TaskManager{tasks: tasks}
}

func (s *TaskManager) Run(ctx context.Context, runOpt RunOpt) error {
	if len(s.tasks) == 0 {
		return errors.New("no tasks to run")
	}
	switch runOpt {
	case RunOptSequence:
		return s.runSequence(ctx)
	case RunOptParallel:
		return s.runParallel(ctx)
	case RunOptLoop:
		return s.runLoop(ctx)

	}

	return nil
}

func (s *TaskManager) runSequence(ctx context.Context) error {
START:
	for _, v := range s.tasks {
		time.Sleep(v.Opts.DelayBefore * time.Millisecond)
		err := v.Task.Exec(ctx, v.Opts)
		if err != nil {
			switch err.(type) {
			default:
				return fmt.Errorf("task run err: %w", err)
			case *TaskErr:
				fmt.Println(fmt.Sprintf("%s from %s", err, v.Opts.Name))
				goto START
			}
		}
		time.Sleep(v.Opts.DelayAfter * time.Millisecond)
	}

	return nil
}

func (s *TaskManager) runParallel(ctx context.Context) error {
	errs, ctx := errgroup.WithContext(ctx)
	for _, v := range s.tasks {
		errs.Go(func() error {
			err := v.Task.Exec(ctx, v.Opts)
			if err != nil {
				switch err.(type) {
				default:
					return fmt.Errorf("task run err: %w", err)
				case *TaskErr:
					fmt.Println(fmt.Sprintf("%s from %s", err, v.Opts.Name))
				}
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

	START:
		// Для запуска части задач один раз
		firstTime := true
		for _, v := range s.tasks {
			if v.Opts.RunOnce && !firstTime {
				continue
			}
			firstTime = false

			time.Sleep(v.Opts.DelayBefore * time.Millisecond)
			err := v.Task.Exec(ctx, v.Opts)
			if err != nil {
				switch err.(type) {
				default:
					return fmt.Errorf("task run err: %w", err)
				case *TaskErr:
					fmt.Println(fmt.Sprintf("%s from %s", err, v.Opts.Name))
					goto START
				}
			}
			time.Sleep(v.Opts.DelayAfter * time.Millisecond)
		}

	}
}
