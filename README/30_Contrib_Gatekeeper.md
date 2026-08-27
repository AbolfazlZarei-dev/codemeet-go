# 30 — `contrib/gatekeeper`

Gatekeeper یک سیستم CAPTCHA برای کاربران جدید گروه است.

## CAPTCHA types

- `math`
- `numbers`

## CaptchaConfig

```go
type CaptchaConfig struct {
    Type string
    Options int
    MinNumber int
    MaxNumber int
    Title string
    CorrectText string
    WrongText string
}
```

## Config

```go
type Config struct {
    ChallengeTimeout time.Duration
    WrongAnswersLimit int
    VerifiedTTL time.Duration
    Captcha CaptchaConfig
    WorkerCount int
    QueueSize int
    SendCaptchaAction func(ctx context.Context, chatID, userID, text string, keyboard *models.InlineKeyboardMarkup) (int, error)
    AnswerCallbackAction func(ctx context.Context, callbackID, text string, showAlert bool) error
    VerifyAction func(ctx context.Context, chatID, userID string)
    KickAction func(ctx context.Context, chatID, userID string)
    DeleteMessageAction func(ctx context.Context, chatID string, messageID int)
}
```

## Workflow

1. کاربر جدید وارد گروه می‌شود.
2. Middleware برای او challenge ایجاد می‌کند.
3. CAPTCHA با Inline Keyboard ارسال می‌شود.
4. Callback بررسی می‌شود.
5. پاسخ صحیح کاربر را برای `VerifiedTTL` تأیید می‌کند.
6. پاسخ‌های بیش از `WrongAnswersLimit` باعث fail شدن challenge می‌شوند.
7. challenge منقضی‌شده cleanup می‌شود.

## Worker Pool

Challengeها در queue قرار می‌گیرند و Workerها ارسال CAPTCHA را انجام می‌دهند.

## API

```text
DefaultConfig
New
Middleware
SendCaptcha
SendCaptchaWithConfig
Stats
Close
```

## مثال

```go
gk := gatekeeper.New(gatekeeper.Config{
    WorkerCount: 4,
    QueueSize: 100,
    VerifyAction: func(ctx context.Context, chatID, userID string) {
        log.Printf("verified %s", userID)
    },
    KickAction: func(ctx context.Context, chatID, userID string) {
        log.Printf("kick %s", userID)
    },
})

bot.Use(gk.Middleware())
defer gk.Close()
```

## Stats

- challenges_sent
- challenges_passed
- challenges_failed
- users_kicked
