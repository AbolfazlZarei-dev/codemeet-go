# Media Methods

دسترسی:

```go
media := bot.API().Media()
```

> وجود Method الزاماً فعال بودن Endpoint را تضمین نمی‌کند.

| Method | Endpoint |
|---|---|
| SendPhoto | sendPhoto |
| SendPhotoWithParams | sendPhoto |
| SendVideo | sendVideo |
| SendDocument | sendDocument |
| SendVoice | sendVoice |
| SendAudio | sendAudio |
| SendAnimation | sendAnimation |
| SendSticker | sendSticker |
| SendVideoNote | sendVideoNote |
| SendMediaGroup | sendMediaGroup |
| SendLocation | sendLocation |
| SendVenue | sendVenue |
| SendContact | sendContact |
| SendPoll | sendPoll |
| SendDice | sendDice |
| GetStickerSet | getStickerSet |
| UploadStickerFile | uploadStickerFile |
| GetFile | getFile |
| SetMessageReaction | setMessageReaction |

## فایل

```go
msg, err := media.SendPhoto(ctx, chatID, "./a.jpg", "عکس")
msg, err := media.SendDocument(ctx, chatID, "./a.pdf", "سند")
msg, err := media.SendVideo(ctx, chatID, "./a.mp4", "ویدیو")
```

## Request کامل عکس

```go
req := &methods.SendPhotoRequest{
    ChatID: chatID,
    Photo: "./a.jpg",
    Caption: "عکس",
    ParseMode: models.ParseModeHTML,
    ReplyToMessageID: messageID,
}
msg, err := media.SendPhotoWithParams(ctx, req)
```

## Media Group

```go
items := []models.InputMedia{
    {Type: "photo", Media: "./a.jpg"},
    {Type: "photo", Media: "./b.jpg"},
}
msgs, err := media.SendMediaGroup(ctx, chatID, items)
```

## Location / Venue / Contact

```go
media.SendLocation(ctx, chatID, 35.68, 51.38)
media.SendVenue(ctx, chatID, 35.68, 51.38, "Office", "Tehran")
media.SendContact(ctx, chatID, "+989000000000", "Abolfazl", "")
```

## Poll / Dice

```go
media.SendPoll(ctx, chatID, "Go?", []string{"Yes", "No"})
media.SendDice(ctx, chatID, "🎲")
```

## File

```go
f, err := media.GetFile(ctx, fileID)
data, err := media.DownloadFile(ctx, f.FilePath)
```

Helper دانلود فعلی حداکثر 20MB را در حافظه می‌خواند.

## Reaction

```go
err := media.SetMessageReaction(ctx, chatID, messageID, []models.ReactionType{
    {Type: "emoji", Emoji: "👍"},
})
```
