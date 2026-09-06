# Message & Media Models

Package: `models`

## `Message`

```go
type Message struct {
    MessageID int `json:"message_id"`
    Date int64 `json:"date"`
    Chat *Chat `json:"chat"`
    From *User `json:"from,omitempty"`
    SenderChat *Chat `json:"sender_chat,omitempty"`
    Text string `json:"text,omitempty"`
    Caption string `json:"caption,omitempty"`
    Entities []MessageEntity `json:"entities,omitempty"`
    CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
    ReplyToMessage *Message `json:"reply_to_message,omitempty"`
    ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
    MediaGroupID string `json:"media_group_id,omitempty"`
    AuthorSignature string `json:"author_signature,omitempty"`
    EditDate int64 `json:"edit_date,omitempty"`
    HasProtectedContent bool `json:"has_protected_content,omitempty"`
    ForwardDate int64 `json:"forward_date,omitempty"`
    ForwardFrom *User `json:"forward_from,omitempty"`
    ForwardFromChat *Chat `json:"forward_from_chat,omitempty"`
    ForwardFromMessageID string `json:"forward_from_message_id,omitempty"`
    ForwardSignature string `json:"forward_signature,omitempty"`
    ForwardSenderName string `json:"forward_sender_name,omitempty"`
    ForwardOrigin *MessageOrigin `json:"forward_origin,omitempty"`
    NewChatMembers []User `json:"new_chat_members,omitempty"`
    LeftChatMember *User `json:"left_chat_member,omitempty"`
    Photo []PhotoSize `json:"photo,omitempty"`
    Video []Video `json:"video,omitempty"`
    Audio []Audio `json:"audio,omitempty"`
    Document []Document `json:"document,omitempty"`
    Animation []Animation `json:"animation,omitempty"`
    Voice []Voice `json:"voice,omitempty"`
    VideoNote []VideoNote `json:"video_note,omitempty"`
    Sticker []Sticker `json:"sticker,omitempty"`
    Contact *Contact `json:"contact,omitempty"`
    Location *Location `json:"location,omitempty"`
    Venue *Venue `json:"venue,omitempty"`
    Poll *Poll `json:"poll,omitempty"`
    Dice *Dice `json:"dice,omitempty"`
}
```

## `MessageEntity`

```go
type MessageEntity struct {
    Type string `json:"type"`
    Offset int `json:"offset"`
    Length int `json:"length"`
    URL string `json:"url,omitempty"`
    User *User `json:"user,omitempty"`
    Language string `json:"language,omitempty"`
    CustomEmoji string `json:"custom_emoji_id,omitempty"`
}
```

## `PhotoSize`

```go
type PhotoSize struct {
    FileID string `json:"file_id"`
    FileUniqueID string `json:"file_unique_id"`
    Width int `json:"width"`
    Height int `json:"height"`
    FileSize int64 `json:"file_size,omitempty"`
}
```

## `Video`

```go
type Video struct {
    FileID string `json:"file_id"`
    FileUniqueID string `json:"file_unique_id"`
    Width int `json:"width"`
    Height int `json:"height"`
    Duration int `json:"duration"`
    Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
    FileName string `json:"file_name,omitempty"`
    MimeType string `json:"mime_type,omitempty"`
    FileSize int64 `json:"file_size,omitempty"`
}
```

## `Audio`

```go
type Audio struct {
    FileID string `json:"file_id"`
    FileUniqueID string `json:"file_unique_id"`
    Duration int `json:"duration"`
    Performer string `json:"performer,omitempty"`
    Title string `json:"title,omitempty"`
    FileName string `json:"file_name,omitempty"`
    MimeType string `json:"mime_type,omitempty"`
    FileSize int64 `json:"file_size,omitempty"`
    Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
}
```

## `Document`

```go
type Document struct {
    FileID string `json:"file_id"`
    FileUniqueID string `json:"file_unique_id"`
    Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
    FileName string `json:"file_name,omitempty"`
    MimeType string `json:"mime_type,omitempty"`
    FileSize int64 `json:"file_size,omitempty"`
}
```

## `Animation`

```go
type Animation struct {
    FileID string `json:"file_id"`
    FileUniqueID string `json:"file_unique_id"`
    Width int `json:"width"`
    Height int `json:"height"`
    Duration int `json:"duration"`
    Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
    FileName string `json:"file_name,omitempty"`
    MimeType string `json:"mime_type,omitempty"`
    FileSize int64 `json:"file_size,omitempty"`
}
```

## `Voice`

```go
type Voice struct {
    FileID string `json:"file_id"`
    FileUniqueID string `json:"file_unique_id"`
    Duration int `json:"duration"`
    MimeType string `json:"mime_type,omitempty"`
    FileSize int64 `json:"file_size,omitempty"`
}
```

## `VideoNote`

```go
type VideoNote struct {
    FileID string `json:"file_id"`
    FileUniqueID string `json:"file_unique_id"`
    Length int `json:"length"`
    Duration int `json:"duration"`
    Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
    FileSize int64 `json:"file_size,omitempty"`
}
```

## `Contact`

```go
type Contact struct {
    PhoneNumber string `json:"phone_number"`
    FirstName string `json:"first_name"`
    LastName string `json:"last_name,omitempty"`
    UserID string `json:"user_id,omitempty"`
    VCard string `json:"vcard,omitempty"`
}
```

## `Location`

```go
type Location struct {
    Longitude float64 `json:"longitude"`
    Latitude float64 `json:"latitude"`
    HorizontalAccuracy float64 `json:"horizontal_accuracy,omitempty"`
    LivePeriod int `json:"live_period,omitempty"`
    Heading int `json:"heading,omitempty"`
    ProximityAlertRadius int `json:"proximity_alert_radius,omitempty"`
}
```

## `Venue`

```go
type Venue struct {
    Location *Location `json:"location"`
    Title string `json:"title"`
    Address string `json:"address"`
    FoursquareID string `json:"foursquare_id,omitempty"`
    FoursquareType string `json:"foursquare_type,omitempty"`
    GooglePlaceID string `json:"google_place_id,omitempty"`
    GooglePlaceType string `json:"google_place_type,omitempty"`
}
```

## `Poll`

```go
type Poll struct {
    ID string `json:"id"`
    Question string `json:"question"`
    Options []PollOption `json:"options"`
    TotalVoterCount int `json:"total_voter_count"`
    IsClosed bool `json:"is_closed"`
    IsAnonymous bool `json:"is_anonymous"`
    Type string `json:"type"`
    AllowsMultipleAnswers bool `json:"allows_multiple_answers"`
    CorrectOptionID int `json:"correct_option_id,omitempty"`
    Explanation string `json:"explanation,omitempty"`
}
```

## `PollOption`

```go
type PollOption struct {
    Text string `json:"text"`
    VoterCount int `json:"voter_count"`
}
```

## `Dice`

```go
type Dice struct {
    Emoji string `json:"emoji"`
    Value int `json:"value"`
}
```

## `Sticker`

```go
type Sticker struct {
    FileID string `json:"file_id"`
    FileUniqueID string `json:"file_unique_id"`
    Type string `json:"type"`
    Width int `json:"width"`
    Height int `json:"height"`
    IsAnimated bool `json:"is_animated"`
    IsVideo bool `json:"is_video"`
    Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
    Emoji string `json:"emoji,omitempty"`
    SetName string `json:"set_name,omitempty"`
    FileSize int64 `json:"file_size,omitempty"`
}
```

## `StickerSet`

```go
type StickerSet struct {
    Name string `json:"name"`
    Title string `json:"title"`
    StickerType string `json:"sticker_type"`
    IsAnimated bool `json:"is_animated"`
    IsVideo bool `json:"is_video"`
    Stickers []Sticker `json:"stickers"`
    Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
}
```

## `File`

```go
type File struct {
    FileID string `json:"file_id"`
    FileUniqueID string `json:"file_unique_id"`
    FileSize int64 `json:"file_size,omitempty"`
    FilePath string `json:"file_path,omitempty"`
}
```

## `InputMedia`

```go
type InputMedia struct {
    Type string `json:"type"`
    Media string `json:"media"`
    Caption string `json:"caption,omitempty"`
    ParseMode string `json:"parse_mode,omitempty"`
    CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
    Thumbnail string `json:"thumbnail,omitempty"`
    Width int `json:"width,omitempty"`
    Height int `json:"height,omitempty"`
    Duration int `json:"duration,omitempty"`
    Streaming bool `json:"supports_streaming,omitempty"`
    Title string `json:"title,omitempty"`
    Performer string `json:"performer,omitempty"`
    HasSpoiler bool `json:"has_spoiler,omitempty"`
}
```

## `ReactionType`

```go
type ReactionType struct {
    Type string `json:"type"`
    Emoji string `json:"emoji,omitempty"`
    CustomEmoji string `json:"custom_emoji_id,omitempty"`
}
```

## Message helpers

```go
msg.HasMedia()
msg.IsCommand()
msg.CommandName()
msg.CommandArgs()
```

## ParseMode

```go
models.ParseModeHTML
models.ParseModeMarkdown
models.ParseModeMarkdownV2
```

## ChatAction

```text
typing
upload_photo
record_video
upload_video
record_voice
upload_voice
upload_document
choose_sticker
find_location
record_video_note
upload_video_note
```

## Entity types

```text
bold, italic, underline, strikethrough, spoiler,
code, pre, text_link, text_mention, mention, hashtag,
cashtag, bot_command, url, email, phone_number,
blockquote, custom_emoji
```
