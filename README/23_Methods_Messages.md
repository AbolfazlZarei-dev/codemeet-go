# 23 — Messages

## Sending

```text
Send
SendText
SendHTML
SendMarkdown
SendWithKeyboard
```

## Forward / Copy

```text
Forward
Copy
```

## Editing

```text
EditText
EditTextInline
EditCaption
EditReplyMarkup
```

## Delete

```text
Delete
DeleteMessages
```

## Interaction

```text
SendChatAction
AnswerCallback
AnswerCallbackSimple
```

## Internal retry

Messages شامل `doWithRetry` و `doWithRetryMultipart` برای اجرای درخواست‌ها با policy داخلی است.

## Parse errors

`ParseError` پاسخ API را به خطای typed تبدیل می‌کند.
