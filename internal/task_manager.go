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

func (s *TaskManager) Run(ctx context.Context, runOpt RunOpt, reRunAttempts int) error {
	if len(s.tasks) == 0 {
		return errors.New("no tasks to run")
	}
	switch runOpt {
	case RunOptSequence:
		return s.runSequence(ctx)
	case RunOptParallel:
		return s.runParallel(ctx)
	case RunOptLoop:
		return s.runLoop(ctx, reRunAttempts)

	}

	return nil
}

func (s *TaskManager) runSequence(ctx context.Context) error {
START:
	for _, v := range s.tasks {

		if v.Opts.RunType == Disabled {
			continue
		}

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
					fmt.Println(fmt.Sprintf("opt: parallel, %s from %s", err, v.Opts.Name))
				}
			}
			return nil
		})
	}

	return errs.Wait()
}

func (s *TaskManager) runLoop(ctx context.Context, reRunAttempts int) error {
	// Для запуска части задач один раз
	firstTime := true
	failRunsCount := 0

	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return nil
		default:

		}

	START:
		// для полного перезапуска, например если игра не отвечает
		if failRunsCount > reRunAttempts {
			firstTime = true
		}

		for _, v := range s.tasks {
			// пропускаем таски
			if v.Opts.RunType == Disabled {
				continue
			}

			// запускаем один раз, либо полный перезапуск по goto
			if v.Opts.RunType == RunOnce && !firstTime {
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
					fmt.Println(fmt.Sprintf("opt: loop, %s from %s", err, v.Opts.Name))
					failRunsCount++
					goto START
				}
			}
			time.Sleep(v.Opts.DelayAfter * time.Millisecond)
		}

		// сбрасываем счётчик неуспешных циклов при полном прохождении всех задач
		failRunsCount = 0

	}
}
