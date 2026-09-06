# Method Signatures

امضاهای exported در package `methods` مطابق سورس نسخه 1.1.0:

```go
func (b *BotMethods) GetMe(ctx context.Context) (*models.User, error) {
```
```go
func (b *BotMethods) SetName(ctx context.Context, name string) error {
```
```go
func (b *BotMethods) GetName(ctx context.Context) (string, error) {
```
```go
func (b *BotMethods) GetNameWithLang(ctx context.Context, lang string) (string, error) {
```
```go
func (b *BotMethods) SetDescription(ctx context.Context, desc string) error {
```
```go
func (b *BotMethods) SetDescriptionWithLang(ctx context.Context, desc, lang string) error {
```
```go
func (b *BotMethods) GetDescription(ctx context.Context) (string, error) {
```
```go
func (b *BotMethods) GetDescriptionWithLang(ctx context.Context, lang string) (string, error) {
```
```go
func (b *BotMethods) SetShortDescription(ctx context.Context, desc string) error {
```
```go
func (b *BotMethods) SetShortDescriptionWithLang(ctx context.Context, desc, lang string) error {
```
```go
func (b *BotMethods) GetShortDescription(ctx context.Context) (string, error) {
```
```go
func (b *BotMethods) GetShortDescriptionWithLang(ctx context.Context, lang string) (string, error) {
```
```go
func (b *BotMethods) SetCommands(ctx context.Context, cmds []models.BotCommand, lang string) error {
```
```go
func (b *BotMethods) GetCommands(ctx context.Context, lang string) ([]models.BotCommand, error) {
```
```go
func (b *BotMethods) DeleteCommands(ctx context.Context, lang string) error {
```
```go
func (b *BotMethods) LogOut(ctx context.Context) (bool, error) {
```
```go
func (b *BotMethods) Close(ctx context.Context) (bool, error) {
```
```go
func (c *ChatMethods) GetChat(ctx context.Context, chatID string) (*models.Chat, error) {
```
```go
func (c *ChatMethods) GetChatMember(ctx context.Context, chatID, userID string) (*models.ChatMember, error) {
```
```go
func (c *ChatMethods) GetChatAdministrators(ctx context.Context, chatID string) ([]models.ChatMember, error) {
```
```go
func (c *ChatMethods) GetChatMemberCount(ctx context.Context, chatID string) (int, error) {
```
```go
func (c *ChatMethods) PinMessage(ctx context.Context, chatID string, messageID int, disableNotif bool) error {
```
```go
func (c *ChatMethods) UnpinMessage(ctx context.Context, chatID string, messageID int) error {
```
```go
func (c *ChatMethods) UnpinAllMessages(ctx context.Context, chatID string) error {
```
```go
func (c *ChatMethods) BanChatMember(ctx context.Context, chatID, userID string, untilDate int64, revokeMessages bool) error {
```
```go
func (c *ChatMethods) UnbanChatMember(ctx context.Context, chatID, userID string, onlyIfBanned bool) error {
```
```go
func (c *ChatMethods) RestrictChatMember(ctx context.Context, chatID, userID string, permissions *models.ChatPermissions, untilDate int64) error {
```
```go
func (c *ChatMethods) PromoteChatMember(ctx context.Context, chatID, userID string, rights *models.ChatAdministratorRights) error {
```
```go
func (c *ChatMethods) SetChatAdministratorCustomTitle(ctx context.Context, chatID, userID, customTitle string) error {
```
```go
func (c *ChatMethods) SetChatTitle(ctx context.Context, chatID, title string) error {
```
```go
func (c *ChatMethods) SetChatDescription(ctx context.Context, chatID, description string) error {
```
```go
func (c *ChatMethods) SetChatPhoto(ctx context.Context, chatID, photoPath string) error {
```
```go
func (c *ChatMethods) DeleteChatPhoto(ctx context.Context, chatID string) error {
```
```go
func (c *ChatMethods) SetChatPermissions(ctx context.Context, chatID string, permissions *models.ChatPermissions) error {
```
```go
func (c *ChatMethods) LeaveChat(ctx context.Context, chatID string) error {
```
```go
func (c *ChatMethods) ExportChatInviteLink(ctx context.Context, chatID string) (string, error) {
```
```go
func (c *ChatMethods) CreateChatInviteLink(ctx context.Context, req *models.CreateChatInviteLinkRequest) (*models.ChatInviteLink, error) {
```
```go
func (c *ChatMethods) EditChatInviteLink(ctx context.Context, req *models.EditChatInviteLinkRequest) (*models.ChatInviteLink, error) {
```
```go
func (c *ChatMethods) RevokeChatInviteLink(ctx context.Context, chatID, inviteLink string) (*models.ChatInviteLink, error) {
```
```go
func (m *Media) SendPhoto(ctx context.Context, chatID, photo, caption string) (*models.Message, error) {
```
```go
func (m *Media) SendPhotoWithParams(ctx context.Context, req *SendPhotoRequest) (*models.Message, error) {
```
```go
func (m *Media) SendVideo(ctx context.Context, chatID, video, caption string) (*models.Message, error) {
```
```go
func (m *Media) SendDocument(ctx context.Context, chatID, doc, caption string) (*models.Message, error) {
```
```go
func (m *Media) SendVoice(ctx context.Context, chatID, voice, caption string) (*models.Message, error) {
```
```go
func (m *Media) SendAudio(ctx context.Context, chatID, audio, caption string) (*models.Message, error) {
```
```go
func (m *Media) SendAnimation(ctx context.Context, chatID, animation, caption string) (*models.Message, error) {
```
```go
func (m *Media) SendSticker(ctx context.Context, chatID, sticker string) (*models.Message, error) {
```
```go
func (m *Media) SendVideoNote(ctx context.Context, chatID, videoNote string) (*models.Message, error) {
```
```go
func (m *Media) SendMediaGroup(ctx context.Context, chatID string, media []models.InputMedia) ([]models.Message, error) {
```
```go
func (m *Media) SendLocation(ctx context.Context, chatID string, lat, lng float64) (*models.Message, error) {
```
```go
func (m *Media) SendVenue(ctx context.Context, chatID string, lat, lng float64, title, address string) (*models.Message, error) {
```
```go
func (m *Media) SendContact(ctx context.Context, chatID, phone, firstName, lastName string) (*models.Message, error) {
```
```go
func (m *Media) SendPoll(ctx context.Context, chatID, question string, options []string) (*models.Message, error) {
```
```go
func (m *Media) SendDice(ctx context.Context, chatID, emoji string) (*models.Message, error) {
```
```go
func (m *Media) GetStickerSet(ctx context.Context, name string) (*models.StickerSet, error) {
```
```go
func (m *Media) UploadStickerFile(ctx context.Context, userID, stickerPath string) (*models.File, error) {
```
```go
func (m *Media) GetFile(ctx context.Context, fileID string) (*models.File, error) {
```
```go
func (m *Media) DownloadFile(ctx context.Context, filePath string) ([]byte, error) {
```
```go
func (m *Media) SetMessageReaction(ctx context.Context, chatID string, messageID int, reaction []models.ReactionType) error {
```
```go
func (m *Messages) Send(ctx context.Context, req *SendMessageRequest) (*models.Message, error) {
```
```go
func (m *Messages) SendText(ctx context.Context, chatID, text string) (*models.Message, error) {
```
```go
func (m *Messages) SendHTML(ctx context.Context, chatID, text string) (*models.Message, error) {
```
```go
func (m *Messages) SendMarkdown(ctx context.Context, chatID, text string) (*models.Message, error) {
```
```go
func (m *Messages) SendWithKeyboard(ctx context.Context, chatID, text string, markup interface{}) (*models.Message, error) {
```
```go
func (m *Messages) Forward(ctx context.Context, chatID, fromChatID string, messageID int) (*models.Message, error) {
```
```go
func (m *Messages) Copy(ctx context.Context, chatID, fromChatID string, messageID int, caption string) (*CopyMessageResult, error) {
```
```go
func (m *Messages) EditText(ctx context.Context, chatID string, messageID int, text string, parseMode models.ParseMode, replyMarkup *models.InlineKeyboardMarkup) error {
```
```go
func (m *Messages) EditTextInline(ctx context.Context, inlineMessageID, text string, parseMode models.ParseMode, replyMarkup *models.InlineKeyboardMarkup) error {
```
```go
func (m *Messages) EditCaption(ctx context.Context, chatID string, messageID int, caption string, parseMode models.ParseMode, replyMarkup *models.InlineKeyboardMarkup) error {
```
```go
func (m *Messages) EditReplyMarkup(ctx context.Context, chatID string, messageID int, markup *models.InlineKeyboardMarkup) error {
```
```go
func (m *Messages) Delete(ctx context.Context, chatID string, messageID int) error {
```
```go
func (m *Messages) DeleteMessages(ctx context.Context, chatID string, messageIDs []int) error {
```
```go
func (m *Messages) SendChatAction(ctx context.Context, chatID string, action models.ChatAction) error {
```
```go
func (m *Messages) AnswerCallback(ctx context.Context, req *models.AnswerCallbackRequest) error {
```
```go
func (m *Messages) AnswerCallbackSimple(ctx context.Context, callbackID, text string, showAlert bool) error {
```
```go
func (m *Messages) SetMyDescription(ctx context.Context, desc, lang string) error {
```
```go
func (m *Messages) ParseError(err error) *errors.APIError {
```
```go
func New(c *api.Client, r *retry.Policy, l *ratelimit.Limiter) *Methods {
```
```go
func (m *Methods) Messages() *Messages      { return m.messages }
```
```go
func (m *Methods) Media() *Media            { return m.media }
```
```go
func (m *Methods) Bot() *BotMethods         { return m.bot }
```
```go
func (m *Methods) Chat() *ChatMethods       { return m.chat }
```
```go
func (m *Methods) Webhook() *WebhookMethods { return m.webhook }
```
```go
func (m *Methods) Updates() *UpdatesMethods { return m.updates }
```
```go
func (m *Methods) Client() *api.Client { return m.api }
```
```go
func (u *UpdatesMethods) Get(ctx context.Context, p *GetUpdatesParams) ([]models.Update, error) {
```
```go
func (w *WebhookMethods) Set(ctx context.Context, req *models.SetWebhookRequest) error {
```
```go
func (w *WebhookMethods) GetInfo(ctx context.Context) (*models.WebhookInfo, error) {
```
```go
func (w *WebhookMethods) Delete(ctx context.Context) error {
```
```go
func (w *WebhookMethods) DeleteWithDrop(ctx context.Context, drop bool) error {
```
