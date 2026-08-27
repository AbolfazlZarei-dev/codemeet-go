# CodeMeet Go

کتابخانه‌ی رسمی Go برای توسعه‌ی ربات‌های CodeMeet است. این پروژه یک لایه‌ی Go-native روی CodeMeet Bot API فراهم می‌کند تا توسعه‌دهنده بتواند بدون ساخت مستقیم درخواست‌های HTTP، ربات را بسازد، رویدادها را دریافت و dispatch کند، پیام و رسانه بفرستد، پروفایل ربات را مدیریت کند و قابلیت‌های عملیاتی مانند Retry، Rate Limit، Cache، Logging و Webhook/Long Polling را به کار بگیرد.

> Version: `1.0.0`  
> Author: **Abolfazl Zarei**  
> Repository: `github.com/AbolfazlZarei-dev/codemeet-go`

## نصب

```bash
go get github.com/AbolfazlZarei-dev/codemeet-go
```

## یک ربات ساده

```go
package main

import (
	"context"
	"log"

	codemeet "github.com/AbolfazlZarei-dev/codemeet-go"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
	"github.com/AbolfazlZarei-dev/codemeet-go/polling"
)

func main() {
	bot, err := codemeet.New("YOUR_BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()

	bot.OnCommand("start", func(ctx context.Context, msg *models.Message) {
		_, err := bot.Reply(ctx, msg, "سلام! ربات فعال است.")
		if err != nil {
			log.Println(err)
		}
	})

	cfg := polling.DefaultConfig()
	if err := bot.StartPolling(context.Background(), cfg); err != nil {
		log.Println(err)
	}
}
```

## راهنمای مستندات

- [01 - معرفی، نصب و شروع کار](01_Getting_Started.md)
- [02 - ارسال پیام و رسانه](02_Messages_and_Media.md)
- [03 - دریافت رویدادها: Polling و Webhook](03_Updates_Polling_Webhook.md)
- [04 - مدیریت Chat، گروه و کانال](04_Chat_Management.md)
- [05 - پروفایل و دستورات Bot](05_Bot_Profile_and_Commands.md)
- [06 - مدل‌ها و ساختار Update](06_Models.md)
- [07 - کیبوردها و تعامل با کاربر](07_Keyboards_and_Buttons.md)
- [08 - قابلیت‌های داخلی کتابخانه و معماری](08_Library_Architecture_and_Features.md)
- [09 - Retry، Rate Limit، Cache و Performance](09_Reliability_Performance.md)
- [10 - خطاها، عیب‌یابی و Shutdown](10_Errors_and_Troubleshooting.md)

## API خام

Endpoint اصلی Bot API:

```text
https://botapi.codemeet.chat/bot<TOKEN>/<METHOD>
```

کتابخانه در لایه‌ی `api` ساخت درخواست، JSON decoding، multipart upload، آمار درخواست‌ها و Circuit Breaker را مدیریت می‌کند.
