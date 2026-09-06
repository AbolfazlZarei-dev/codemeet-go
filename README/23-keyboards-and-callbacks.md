# Keyboards & Callback

## Inline

```go
markup := models.NewInlineKeyboard(
    models.InlineRow(
        models.Btn("تأیید", "confirm"),
        models.Btn("لغو", "cancel"),
    ),
)
```

## URL

```go
models.URLBtn("Website", "https://example.com")
```

## WebApp

```go
models.WebAppBtn("Open App", "https://example.com/app")
```

## Reply

```go
models.NewReplyKeyboard(
    models.ReplyRow(
        models.KBtn("پروفایل"),
        models.KBtn("تنظیمات"),
    ),
)
```

## Callback lifecycle

```text
Button click
 → callback_query
 → OnCallback
 → AnswerCallback
 → optional Edit/Send
```

همیشه `callback_data` را به‌عنوان input کاربر validate کنید.
