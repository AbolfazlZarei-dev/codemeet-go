# Sending Messages & Media

## Text

```go
bot.Send(ctx, chatID, "Hello")
bot.SendHTML(ctx, chatID, "<b>Hello</b>")
bot.SendWithKeyboard(ctx, chatID, "Choose", markup)
```

## Reply

```go
bot.Reply(ctx, msg, "پاسخ")
```

## Media

```go
bot.API().Media().SendPhoto(ctx, chatID, "./a.jpg", "caption")
bot.API().Media().SendVideo(ctx, chatID, "./a.mp4", "caption")
bot.API().Media().SendDocument(ctx, chatID, "./a.pdf", "caption")
bot.API().Media().SendVoice(ctx, chatID, "./a.ogg", "caption")
bot.API().Media().SendAudio(ctx, chatID, "./a.mp3", "caption")
```

## Advanced Photo

```go
bot.API().Media().SendPhotoWithParams(ctx, &methods.SendPhotoRequest{
    ChatID: chatID,
    Photo: "./a.jpg",
    Caption: "caption",
    ParseMode: models.ParseModeHTML,
})
```

## Download

```go
f, _ := bot.API().Media().GetFile(ctx, fileID)
data, _ := bot.API().Media().DownloadFile(ctx, f.FilePath)
```
