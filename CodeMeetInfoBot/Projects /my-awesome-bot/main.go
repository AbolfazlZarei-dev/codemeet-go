// main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"my-awesome-bot/app"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("فایل .env یافت نشد، تلاش برای استفاده از متغیرهای محیطی سیستم...")
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("متغیر محیطی BOT_TOKEN تنظیم نشده است! لطفاً فایل .env را بسازید.")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, token); err != nil {
		log.Fatalf("ربات با خطا متوقف شد: %v", err)
	}

	log.Println("ربات با موفقیت خاموش شد.")
}
