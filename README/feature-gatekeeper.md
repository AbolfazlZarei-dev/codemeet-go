# Gatekeeper

سیستم CAPTCHA/Challenge برای تأیید انسان.

## Types

```text
math
numbers
```

## Defaults

```text
ChallengeTimeout=120s
WrongAnswersLimit=2
VerifiedTTL=24h
WorkerCount=16
QueueSize=256
Captcha.Options=4
Captcha.MinNumber=2
Captcha.MaxNumber=9
```

## استفاده

```go
gk := gatekeeper.New(gatekeeper.DefaultConfig())
bot.Use(gk.Middleware())
defer gk.Close()
```

## CAPTCHA Config

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

## Actionها

Config شامل callbackهای:

```text
SendCaptchaAction
EditMessageAction
AnswerCallbackAction
VerifyAction
KickAction
DeleteMessageAction
```

است.

## Metrics

```go
gk.Stats()
```

Gatekeeper از worker pool و queue استفاده می‌کند.
