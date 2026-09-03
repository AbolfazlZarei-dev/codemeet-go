// app/bot.go
package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/codemeet"
	"github.com/AbolfazlZarei-dev/codemeet-go/middleware"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
	"github.com/AbolfazlZarei-dev/codemeet-go/polling"
)

func Run(ctx context.Context, token string) error {
	// 2. کانفیگ کاستوم برای HTTP/2 وmax Performance
	transport := &http.Transport{
		MaxIdleConns:          2000,
		MaxIdleConnsPerHost:   1000,
		MaxConnsPerHost:       500,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		DisableCompression:    false,
		WriteBufferSize:       128 * 1024,
		ReadBufferSize:        128 * 1024,
	}
	customHTTPClient := &http.Client{
		Transport: transport,
		Timeout:   0,
	}

	// 2. ساخت ربات با امکانات پیشرفته
	bot, err := codemeet.New(token,
		codemeet.WithHTTPClient(customHTTPClient),
		codemeet.WithShardedCache(64, 15*time.Minute),
		codemeet.WithRateLimitBurst(30, 100),
		codemeet.WithMiddleware(middleware.Recovery(nil)),
	)
	if err != nil {
		return fmt.Errorf("خطا در مقداردهی ربات: %w", err)
	}
	defer bot.Close()

	// 3. ثبت کامندها
	cmds := []models.BotCommand{
		{Command: "start", Description: "شروع و نمایش اطلاعات کاربر"},
		{Command: "info", Description: "نمایش مجدد اطلاعات شما"},
		{Command: "help", Description: "مرکز راهنما و پشتیبانی"},
		{Command: "stats", Description: "آمار و وضعیت ربات"},
	}
	if err := bot.API().Bot().SetCommands(ctx, cmds, ""); err != nil {
		log.Printf("⚠️ خطا در تنظیم کامندها: %v", err)
	} else {
		log.Println("✅ Bot commands registered successfully.")
	}

	// 4. رجیستر کردن هندلرها
	RegisterHandlers(bot)

	// 5. تنظیمات Long Polling
	cfg := polling.Config{
		Timeout:            30,
		PollInterval:       1 * time.Second,
		Limit:              100,
		BufferSize:         5000,
		DeleteWebhookFirst: true,
	}

	fmt.Println("🚀 Bot is starting with HTTP/2 and Sharded Cache...")
	return bot.StartPolling(ctx, cfg)
}
