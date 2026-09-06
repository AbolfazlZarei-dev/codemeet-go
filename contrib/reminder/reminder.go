package reminder

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

type Config struct {
	SendAction func(ctx context.Context, chatID, text string) error
}

type Reminder struct {
	cfg   Config
	tasks sync.Map // key: taskID, value: *time.Timer
}

func New(cfg Config) *Reminder {
	return &Reminder{cfg: cfg}
}

// Middleware برای پردازش دستورات یادآوری (مثال: /remind 10m خرید شیر)
func (r *Reminder) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u.Message != nil && (strings.HasPrefix(u.Message.Text, "/remind") || strings.HasPrefix(u.Message.Text, "یادآور")) {
				r.handleReminderCommand(ctx, u.Message)
				return
			}
			next(ctx, u)
		}
	}
}

func (r *Reminder) handleReminderCommand(ctx context.Context, msg *models.Message) {
	text := strings.TrimSpace(msg.Text)
	// حذف کلمه /remind یا یادآور از ابتدای متن
	if strings.HasPrefix(text, "/remind") {
		text = strings.TrimPrefix(text, "/remind")
	} else {
		text = strings.TrimPrefix(text, "یادآور")
	}
	text = strings.TrimSpace(text)

	parts := strings.Fields(text)
	if len(parts) < 2 {
		r.cfg.SendAction(ctx, msg.Chat.ID, "❌ فرمت اشتباه است.\nمثال: `/remind 10m خرید شیر`")
		return
	}

	timeStr := parts[0]
	taskText := strings.Join(parts[1:], " ")

	duration, err := parseDuration(timeStr)
	if err != nil {
		r.cfg.SendAction(ctx, msg.Chat.ID, "❌ زمان نامعتبر است. مثال معتبر: `10m` یا `2h` یا `1d`")
		return
	}

	taskID := fmt.Sprintf("%s_%d", msg.From.ID, time.Now().UnixNano())

	// ذخیره در مپ برای امکان کنسل کردن در آینده
	timer := time.AfterFunc(duration, func() {
		alertText := fmt.Sprintf("⏰ یادآوری برای %s:\n\n%s", msg.From.FullName(), taskText)
		r.cfg.SendAction(context.Background(), msg.Chat.ID, alertText)
		r.tasks.Delete(taskID)
	})

	r.tasks.Store(taskID, timer)

	r.cfg.SendAction(ctx, msg.Chat.ID, fmt.Sprintf("✅ یادآوری تنظیم شد.\n⏱ زمان: %s\n📝 متن: %s", duration.String(), taskText))
}

func parseDuration(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration")
	}

	unit := s[len(s)-1]
	numStr := s[:len(s)-1]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	switch unit {
	case 's':
		return time.Duration(num) * time.Second, nil
	case 'm':
		return time.Duration(num) * time.Minute, nil
	case 'h':
		return time.Duration(num) * time.Hour, nil
	case 'd':
		return time.Duration(num) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown time unit")
	}
}
