package manor

import (
	"clicker_bot/internal"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TmpImgPath = "\\static\\screens\\tmp.png"

type CaptchaTask struct {
	imgToFindPath string
	Screen        *internal.Screen
	action        *internal.Action
	Attempts      int
}

func NewCaptchaTaskTask(imgToFindPath string, screen *internal.Screen, action *internal.Action, attempts int) *CaptchaTask {
	return &CaptchaTask{imgToFindPath: imgToFindPath, Screen: screen, Attempts: attempts}
}

func (f *CaptchaTask) Exec(ctx context.Context, opts internal.TaskOpts) error {
	i := 0
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return nil
		default:

		}

		// считаем количество попыток
		if i > f.Attempts {
			return &internal.TaskErr{StatusCode: internal.AttemptsEnd}
		}
		i++

		time.Sleep(opts.DelayBefore * time.Millisecond)

		//находим область капчи
		point, err := f.Screen.FindOnScreen(f.imgToFindPath)
		if err != nil {
			return fmt.Errorf("find img err:%w", err)
		}
		if point == nil {
			continue
		}

		//сохраняем скрин с капчей
		err = f.Screen.SaveScreen(TmpImgPath, point.X, point.Y-50, 100, 50)
		if err != nil {
			return fmt.Errorf("can't save captcha:%w", err)
		}

		text, err := f.Screen.FindText(TmpImgPath)
		captchaRes, err := f.decodeCaptcha(text)
		if err != nil {
			return &internal.TaskErr{StatusCode: internal.InvalidCaptcha, Err: err}
		}
		strRes := strconv.Itoa(captchaRes)
		for _, v := range strRes {
			f.action.KeyPress(string(v))
		}

		//после ввода капчи жмём подтверждение
		time.Sleep(100 * time.Microsecond)
		f.action.Click(point.X+10, point.Y+10, false)

		time.Sleep(opts.DelayAfter * time.Millisecond)
		return nil
	}
}

func (f *CaptchaTask) decodeCaptcha(captcha string) (int, error) {
	captcha = strings.TrimRight(captcha, "\r\n")
	captcha = strings.ReplaceAll(captcha, " ", "")
	if last := len(captcha) - 1; last >= 0 {
		captcha = captcha[:last]
	}

	res := strings.Split(captcha, "+")
	if len(res) != 2 {
		return 0, errors.New("invalid captcha decode len")
	}
	firstNum, err := strconv.Atoi(res[0])
	if err != nil {
		return 0, fmt.Errorf("can't decode string. %s. %w", res, err)
	}
	secondNum, err := strconv.Atoi(res[1])
	if err != nil {
		return 0, fmt.Errorf("can't decode string. %s. %w", res, err)
	}
	return firstNum + secondNum, nil
}
