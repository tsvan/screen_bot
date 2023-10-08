package manor

import (
	"clicker_bot/internal"
	"context"
	"errors"
	"fmt"
	"image"
	"strconv"
	"strings"
	"time"
)

const TmpImgPath = "\\static\\screens\\tmp.jpeg"

type CaptchaTask struct {
	imgToFindPath string
	Screen        *internal.Screen
	action        *internal.Action
	Attempts      int
	savedPoint    *image.Point
}

func NewCaptchaTaskTask(imgToFindPath string, screen *internal.Screen, action *internal.Action, attempts int) *CaptchaTask {
	return &CaptchaTask{imgToFindPath: imgToFindPath, Screen: screen, action: action, Attempts: attempts}
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
			f.savedPoint = nil
			return &internal.TaskErr{StatusCode: internal.AttemptsEnd}
		}
		i++

		time.Sleep(opts.DelayBefore * time.Millisecond)

		var clickPoint image.Point
		if f.savedPoint == nil {
			point, err := f.Screen.FindOnScreen(f.imgToFindPath)
			if err != nil {
				return fmt.Errorf("find img err:%w", err)
			}
			if point == nil {
				continue
			}
			f.savedPoint = point
			clickPoint = *point
		} else {
			checkPoint, err := f.Screen.FindOnScreen(f.imgToFindPath, f.savedPoint.X-150, f.savedPoint.Y-200, 320, 400)
			if err != nil {
				return fmt.Errorf("find img err:%w", err)
			}
			if checkPoint == nil {
				continue
			}
			clickPoint.X = f.savedPoint.X - 150 + checkPoint.X
			clickPoint.Y = f.savedPoint.Y - 200 + checkPoint.Y
		}

		//сохраняем скрин с капчей
		err := f.Screen.SaveScreen(TmpImgPath, clickPoint.X, clickPoint.Y-50, 70, 50)
		if err != nil {
			return fmt.Errorf("can't save captcha:%w", err)
		}

		text, err := f.Screen.FindText(TmpImgPath)
		captchaRes, err := f.decodeCaptcha(text)
		if err != nil {
			//err = f.Screen.SaveScreen(fmt.Sprintf("\\static\\screens\\test\\%d.jpeg",time.Now().Unix()), clickPoint.X, clickPoint.Y-50, 70, 50)
			//if err != nil {
			//	fmt.Println(err)
			//}
			return &internal.TaskErr{StatusCode: internal.InvalidCaptcha, Err: err}
		}
		strRes := strconv.Itoa(captchaRes)
		for _, v := range strRes {
			f.action.KeyPress(string(v))
		}

		//после ввода капчи жмём подтверждение
		time.Sleep(100 * time.Microsecond)
		f.action.Click(clickPoint.X+10, clickPoint.Y+10, false)

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
