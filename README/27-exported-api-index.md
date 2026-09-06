# Exported API Index

این فهرست از سورس نسخه 1.1.0 استخراج شده است.

## `codemeet`

```go
func (d *DashboardWriter) Write(p []byte) (int, error) {
```

```go
func (d *DashboardWriter) GetLogs() []string {
```

```go
func WithBaseURL(url string) Option {
```

```go
func WithHTTPClient(c *http.Client) Option {
```

```go
func WithTimeout(d time.Duration) Option {
```

```go
func WithRetry(p *retry.Policy) Option {
```

```go
func WithRateLimit(rps int) Option {
```

```go
func WithRateLimitBurst(rps, burst int) Option {
```

```go
func WithCache(ttl time.Duration) Option {
```

```go
func WithShardedCache(shards int, ttl time.Duration) Option {
```

```go
func WithLogger(l *logger.Logger) Option {
```

```go
func WithDashboardAuth(user, pass string) Option {
```

```go
func WithConsole(c *console.Console) Option {
```

```go
func WithBroadcaster(bc *broadcast.Broadcaster) Option {
```

```go
func WithMiddleware(mws ...dispatcher.MiddlewareFunc) Option {
```

```go
func New(token string, opts ...Option) (*Bot, error) {
```

```go
func (b *Bot) API() *methods.Methods              { return b.methods }
```

```go
func (b *Bot) Dispatcher() *dispatcher.Dispatcher { return b.dispatcher }
```

```go
func (b *Bot) Cache() Cache                       { return b.cache }
```

```go
func (b *Bot) Logger() *logger.Logger             { return b.logger }
```

```go
func (b *Bot) RateLimiter() *ratelimit.Limiter    { return b.rateLimit }
```

```go
func (b *Bot) RetryPolicy() *retry.Policy         { return b.retry }
```

```go
func (b *Bot) Token() string                      { return b.token }
```

```go
func (b *Bot) BaseURL() string                    { return b.baseURL }
```

```go
func (b *Bot) StartPolling(ctx context.Context, cfg polling.Config) error {
```

```go
func (b *Bot) StartWebhook(ctx context.Context, cfg webhook.Config) error {
```

```go
func (b *Bot) StartPollingDashboard(cfg polling.Config) {
```

```go
func (b *Bot) StopPollingDashboard() {
```

```go
func (b *Bot) StartDashboard(ctx context.Context, addr string) error {
```

```go
func (b *Bot) RunMode() string {
```

```go
func (b *Bot) SetWebhook(ctx context.Context, url, secretToken string) error {
```

```go
func (b *Bot) DeleteWebhook(ctx context.Context) error {
```

```go
func (b *Bot) GetMe(ctx context.Context) (*models.User, error) {
```

```go
func (b *Bot) ResetMe() {
```

```go
func (b *Bot) OnCommand(cmd string, h func(ctx context.Context, msg *models.Message)) {
```

```go
func (b *Bot) OnMessage(h func(ctx context.Context, msg *models.Message)) {
```

```go
func (b *Bot) OnCallback(h func(ctx context.Context, cq *models.CallbackQuery)) {
```

```go
func (b *Bot) OnText(text string, h func(ctx context.Context, msg *models.Message)) {
```

```go
func (b *Bot) OnRegex(pattern string, h func(ctx context.Context, msg *models.Message, matches []string)) {
```

```go
func (b *Bot) Fallback(h func(ctx context.Context, u *models.Update)) {
```

```go
func (b *Bot) Use(mw ...dispatcher.MiddlewareFunc) {
```

```go
func (b *Bot) Send(ctx context.Context, chatID, text string) (*models.Message, error) {
```

```go
func (b *Bot) SendHTML(ctx context.Context, chatID, text string) (*models.Message, error) {
```

```go
func (b *Bot) SendWithKeyboard(ctx context.Context, chatID, text string, markup interface{}) (*models.Message, error) {
```

```go
func (b *Bot) Reply(ctx context.Context, msg *models.Message, text string) (*models.Message, error) {
```

```go
func (b *Bot) AnswerCallback(ctx context.Context, callbackID, text string, showAlert bool) error {
```

```go
func (b *Bot) Close() error {
```

```go
func (b *Bot) HealthCheck(ctx context.Context) error {
```

```go
func (b *Bot) Stats() api.StatsSnapshot {
```

```go
func (b *Bot) Uptime() time.Duration {
```

```go
func WithoutLogger() Option {
```

## `api`

```go
func (r *countReader) Read(p []byte) (int, error) {
```

```go
func (s *Stats) Record(latency time.Duration, success bool, bytesIn, bytesOut int64) {
```

```go
func (s *Stats) AvgLatency() time.Duration {
```

```go
func (s *Stats) Snapshot() StatsSnapshot {
```

```go
func NewCircuitBreaker(threshold int, reset time.Duration) *CircuitBreaker {
```

```go
func (c *CircuitBreaker) Allow() bool {
```

```go
func (c *CircuitBreaker) RecordSuccess() {
```

```go
func (c *CircuitBreaker) RecordFailure() {
```

```go
func (c *CircuitBreaker) State() int {
```

```go
func NewClient(baseURL, token string, log *logger.Logger) *Client {
```

```go
func (c *Client) SetHTTPClient(hc *http.Client) {
```

```go
func (c *Client) SetTimeout(d time.Duration) {
```

```go
func (c *Client) Stats() *Stats                { return c.stats }
```

```go
func (c *Client) StatsSnapshot() StatsSnapshot { return c.stats.Snapshot() }
```

```go
func (c *Client) Breaker() *CircuitBreaker     { return c.breaker }
```

```go
func (c *Client) UserAgent() string {
```

```go
func (c *Client) DownloadFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
```

```go
func (c *Client) RequestWithParams(ctx context.Context, method string, params interface{}) (*Response, error) {
```

```go
func (c *Client) Request(ctx context.Context, method string) (*Response, error) {
```

```go
func (c *Client) RequestWithMultipart(ctx context.Context, method string, fields map[string]string, files map[string]string) (*Response, error) {
```

```go
func (c *Client) RequestWithForm(ctx context.Context, method string, form map[string]string, files map[string]string) (*Response, error) {
```

```go
func (c *Client) Close() error {
```

```go
func WithRequestID(ctx context.Context, id string) context.Context {
```

```go
func GetRequestID(ctx context.Context) string {
```

```go
func (r *Response) Decode(v interface{}) error {
```

```go
func (r *Response) AsBool() (bool, error) {
```

```go
func (r *Response) ParametersAsRetryAfter() int {
```

## `cache`

```go
func New(ttl time.Duration) *Cache {
```

```go
func (c *Cache) Get(key string) (interface{}, bool) {
```

```go
func (c *Cache) Set(key string, value interface{}) {
```

```go
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
```

```go
func (c *Cache) SetForever(key string, value interface{}) {
```

```go
func (c *Cache) Delete(key string) {
```

```go
func (c *Cache) Len() int {
```

```go
func (c *Cache) Keys() []string {
```

```go
func (c *Cache) Clear() {
```

```go
func (c *Cache) GetOrSet(key string, fn func() interface{}) interface{} {
```

```go
func (c *Cache) GetOrSetWithTTL(key string, ttl time.Duration, fn func() interface{}) interface{} {
```

```go
func (c *Cache) Close() {
```

```go
func NewSharded(shards int, ttl time.Duration) *ShardedCache {
```

```go
func (s *ShardedCache) Get(key string) (interface{}, bool) { return s.getShard(key).Get(key) }
```

```go
func (s *ShardedCache) Set(key string, value interface{})  { s.getShard(key).Set(key, value) }
```

```go
func (s *ShardedCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
```

```go
func (s *ShardedCache) SetForever(key string, value interface{}) {
```

```go
func (s *ShardedCache) Delete(key string) { s.getShard(key).Delete(key) }
```

```go
func (s *ShardedCache) GetOrSet(key string, fn func() interface{}) interface{} {
```

```go
func (s *ShardedCache) GetOrSetWithTTL(key string, ttl time.Duration, fn func() interface{}) interface{} {
```

```go
func (s *ShardedCache) Len() int {
```

```go
func (s *ShardedCache) Keys() []string {
```

```go
func (s *ShardedCache) Clear() {
```

```go
func (s *ShardedCache) Close() {
```

## `database`

```go
func New(dbPath string) (*Database, error) {
```

```go
func (d *Database) Close() error {
```

```go
func (d *Database) SQL() *sql.DB {
```

```go
func (d *Database) SaveGroup(ctx context.Context, chatID, title, groupType string) error {
```

```go
func (d *Database) SetKV(ctx context.Context, namespace, key, value string, ttl time.Duration) error {
```

```go
func (d *Database) GetKV(ctx context.Context, namespace, key string) (string, error) {
```

```go
func (d *Database) DeleteKV(ctx context.Context, namespace, key string) error {
```

```go
func (d *Database) CleanupExpired(ctx context.Context) error {
```

```go
func (d *Database) SaveUser(ctx context.Context, userID, chatID, firstName, username string, isBot bool) error {
```

```go
func (d *Database) GetUser(ctx context.Context, userID string) (*User, error) {
```

```go
func (d *Database) BanUser(ctx context.Context, userID string) error {
```

```go
func (d *Database) UnbanUser(ctx context.Context, userID string) error {
```

```go
func (d *Database) IsBanned(ctx context.Context, userID string) (bool, error) {
```

## `dispatcher`

```go
func New(maxWorkers int) *Dispatcher {
```

```go
func (d *Dispatcher) Use(mw ...MiddlewareFunc) {
```

```go
func (d *Dispatcher) Handle(updateType string, h HandlerFunc) {
```

```go
func (d *Dispatcher) OnMessage(h func(ctx context.Context, msg *models.Message)) {
```

```go
func (d *Dispatcher) OnCallback(h func(ctx context.Context, cq *models.CallbackQuery)) {
```

```go
func (d *Dispatcher) OnCommand(cmd string, h func(ctx context.Context, msg *models.Message)) {
```

```go
func (d *Dispatcher) OnText(text string, h func(ctx context.Context, msg *models.Message)) {
```

```go
func (d *Dispatcher) Fallback(h HandlerFunc) {
```

```go
func (d *Dispatcher) Dispatch(ctx context.Context, update *models.Update) {
```

```go
func (d *Dispatcher) Stop() {
```

```go
func (d *Dispatcher) OnRegex(pattern string, h func(ctx context.Context, msg *models.Message, matches []string)) error {
```

```go
func (d *Dispatcher) Stats() (dispatched, dropped, panics int64) {
```

## `errors`

```go
func (e *APIError) Error() string {
```

```go
func (e *APIError) IsRetryable() bool {
```

```go
func (e *ValidationError) Error() string {
```

```go
func NewValidationError(field, msg string) *ValidationError {
```

```go
func (e *NetworkError) Error() string {
```

```go
func (e *NetworkError) Unwrap() error     { return e.Err }
```

```go
func (e *NetworkError) IsRetryable() bool { return true }
```

```go
func NewNetworkError(err error) *NetworkError {
```

```go
func IsRetryable(err error) bool {
```

```go
func ParseError(code int, desc string, params map[string]interface{}) *APIError {
```

```go
func (m *MultiError) Error() string {
```

```go
func (m *MultiError) Unwrap() []error {
```

```go
func (m *MultiError) Add(err error) {
```

```go
func (m *MultiError) HasError() bool {
```

```go
func AsAPIError(err error) (*APIError, bool) {
```

```go
func AsNetworkError(err error) (*NetworkError, bool) {
```

```go
func AsValidationError(err error) (*ValidationError, bool) {
```

## `groupmanager`

```go
func ParseDuration(s string) (time.Duration, error) {
```

```go
func New(cfg Config) *GroupManager {
```

```go
func (gm *GroupManager) Middleware() dispatcher.MiddlewareFunc {
```

```go
func (gm *GroupManager) CheckLockViolation(msg *models.Message) string {
```

```go
func (s *SyncStore) Set(key string, val interface{})    { s.data.Store(key, val) }
```

```go
func (s *SyncStore) Get(key string) (interface{}, bool) { return s.data.Load(key) }
```

```go
func (s *SyncStore) Delete(key string)                  { s.data.Delete(key) }
```

```go
func NewLockManager() *LockManager { return &LockManager{locks: make(map[string]map[string]bool)} }
```

```go
func (lm *LockManager) SetLock(chatID, lockType string) {
```

```go
func (lm *LockManager) RemoveLock(chatID, lockType string) {
```

```go
func (lm *LockManager) IsLocked(chatID, lockType string) bool {
```

```go
func (lm *LockManager) GetLocksList(chatID string) []string {
```

```go
func NewWarnManager(maxWarns int) *WarnManager {
```

```go
func (wm *WarnManager) AddWarning(chatID, userID string) int {
```

```go
func (wm *WarnManager) GetWarnings(chatID, userID string) int {
```

```go
func (wm *WarnManager) ResetWarnings(chatID, userID string) {
```

```go
func NewNoteManager() *NoteManager { return &NoteManager{notes: make(map[string]string)} }
```

```go
func (nm *NoteManager) AddNote(chatID, keyword, text string) {
```

```go
func (nm *NoteManager) GetNote(chatID, keyword string) (string, bool) {
```

```go
func NewBlacklistManager() *BlacklistManager {
```

```go
func (bm *BlacklistManager) AddWord(chatID, word string) {
```

```go
func (bm *BlacklistManager) RemoveWord(chatID, word string) {
```

```go
func (bm *BlacklistManager) CheckViolation(chatID, text string) bool {
```

```go
func NewFloodManager(limit int, window time.Duration) *FloodManager {
```

```go
func (fm *FloodManager) IsFlooding(userID string) bool {
```

## `logger`

```go
func (l Level) String() string {
```

```go
func (l Level) ColorString() string {
```

```go
func New(level Level) *Logger {
```

```go
func NewJSON(level Level) *Logger {
```

```go
func NewAsync(level Level, bufferSize int) *Logger {
```

```go
func (l *Logger) SetIncludeCaller(b bool) {
```

```go
func (l *Logger) SetEnabled(b bool) {
```

```go
func (l *Logger) IsEnabled() bool {
```

```go
func (l *Logger) SetOutput(w io.Writer) {
```

```go
func (l *Logger) SetLevel(level Level) {
```

```go
func (l *Logger) SetFormat(format Format) {
```

```go
func (l *Logger) WithFields(fields ...interface{}) *Logger {
```

```go
func (l *Logger) Debug(msg string, fields ...interface{}) {
```

```go
func (l *Logger) Info(msg string, fields ...interface{}) {
```

```go
func (l *Logger) Warn(msg string, fields ...interface{}) {
```

```go
func (l *Logger) Error(msg string, fields ...interface{}) {
```

```go
func (l *Logger) Fatal(msg string, fields ...interface{}) {
```

```go
func (l *Logger) Close() {
```

```go
func (l *Logger) Sync() {
```

```go
func (l *Logger) Output() io.Writer {
```

## `methods`

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

## `middleware`

```go
func Recovery(log *logger.Logger) MiddlewareFunc {
```

```go
func Logging(log *logger.Logger) MiddlewareFunc {
```

```go
func RateLimit(perUser int, window time.Duration) MiddlewareFunc {
```

```go
func NewMetricsCounter() *MetricsCounter {
```

```go
func (m *MetricsCounter) Inc(updateType string) {
```

```go
func (m *MetricsCounter) Snapshot() map[string]int64 {
```

```go
func Metrics(counter *MetricsCounter) MiddlewareFunc {
```

```go
func Timeout(d time.Duration) MiddlewareFunc {
```

```go
func BotOnly() MiddlewareFunc {
```

```go
func UserOnly() MiddlewareFunc {
```

```go
func AdminOnly(isAdmin func(userID string) bool) MiddlewareFunc {
```

```go
func Blacklist(blacklist func(userID string) bool) MiddlewareFunc {
```

```go
func Whitelist(whitelist func(userID string) bool) MiddlewareFunc {
```

## `polling`

```go
func DefaultConfig() Config {
```

```go
func New(c *api.Client, d *dispatcher.Dispatcher, log *logger.Logger, cfg Config) *Poller {
```

```go
func NewWithLimiter(c *api.Client, d *dispatcher.Dispatcher, log *logger.Logger, cfg Config, lim *ratelimit.Limiter, rtry *retry.Policy) *Poller {
```

```go
func (p *Poller) Start(ctx context.Context) error {
```

```go
func (p *Poller) Offset() int  { return p.offset }
```

```go
func (p *Poller) ResetOffset() { p.offset = 0 }
```

## `ratelimit`

```go
func New(rate int) *Limiter {
```

```go
func NewWithBurst(rate, burst int) *Limiter {
```

```go
func (l *Limiter) Wait(ctx context.Context) error {
```

```go
func (l *Limiter) TryWait() bool {
```

```go
func (l *Limiter) WaitTimeout(timeout time.Duration) bool {
```

```go
func (l *Limiter) Available() int { return len(l.tokens) }
```

```go
func (l *Limiter) Rate() int      { return l.rate }
```

```go
func (l *Limiter) Total() int64   { return l.total.Load() }
```

```go
func (l *Limiter) Dropped() int64 { return l.dropped.Load() }
```

```go
func (l *Limiter) Close() {
```

```go
func NewConcurrencyLimiter(maxConcurrent int) *ConcurrencyLimiter {
```

```go
func (c *ConcurrencyLimiter) Acquire(ctx context.Context) error {
```

```go
func (c *ConcurrencyLimiter) Release() {
```

## `retry`

```go
func DefaultPolicy() *Policy {
```

```go
func AggressivePolicy() *Policy {
```

```go
func ConservativePolicy() *Policy {
```

```go
func (p *Policy) Do(ctx context.Context, fn func(ctx context.Context) error) error {
```

## `webhook`

```go
func DefaultConfig() Config {
```

```go
func New(c *api.Client, d *dispatcher.Dispatcher, log *logger.Logger, cfg Config) *Server {
```

```go
func (s *Server) Start(ctx context.Context) error {
```

```go
func (s *Server) Stats() (requests, errors int64) {
```

## `ai`

```go
func New(cfg Config) *Client {
```

```go
func (c *Client) StreamChat(ctx context.Context, req ChatRequest, onChunk func(content string)) error {
```

```go
func (c *Client) Ask(ctx context.Context, req ChatRequest) (string, error) {
```

## `antilink`

```go
func DefaultConfig() Config {
```

```go
func New(cfg Config) *AntiLink {
```

```go
func (al *AntiLink) IsBlocked(text string, entities []models.MessageEntity) (bool, string) {
```

```go
func (al *AntiLink) Middleware() dispatcher.MiddlewareFunc {
```

## `antispam`

```go
func DefaultConfig() Config {
```

```go
func New(cfg Config) *AntiSpam {
```

```go
func (as *AntiSpam) BanUser(userID string) {
```

```go
func (as *AntiSpam) UnbanUser(userID string) {
```

```go
func (as *AntiSpam) Stats() map[string]int64 {
```

```go
func (as *AntiSpam) Middleware() dispatcher.MiddlewareFunc {
```

## `forcejoin`

```go
func DefaultConfig() Config {
```

```go
func New(cfg Config) *ForceJoin {
```

```go
func (fj *ForceJoin) Middleware() dispatcher.MiddlewareFunc {
```

```go
func (fj *ForceJoin) ClearUserCache(userID string) {
```

```go
func (fj *ForceJoin) Stats() map[string]int64 {
```

## `gatekeeper`

```go
func DefaultConfig() Config {
```

```go
func New(cfg Config) *Gatekeeper {
```

```go
func (gk *Gatekeeper) Middleware() dispatcher.MiddlewareFunc {
```

```go
func (gk *Gatekeeper) Stats() map[string]int64 {
```

```go
func (gk *Gatekeeper) Close() {
```

## `maintenancemode`

```go
func DefaultConfig() Config {
```

```go
func New(cfg Config) *MaintenanceMode {
```

```go
func (mm *MaintenanceMode) SetEnabled(enabled bool) {
```

```go
func (mm *MaintenanceMode) IsEnabled() bool {
```

```go
func (mm *MaintenanceMode) Stop() {
```

```go
func (mm *MaintenanceMode) Middleware() dispatcher.MiddlewareFunc {
```

## `profanityfilter`

```go
func DefaultConfig() Config {
```

```go
func New(cfg Config) *ProfanityFilter {
```

```go
func (pf *ProfanityFilter) Middleware() dispatcher.MiddlewareFunc {
```

## `reminder`

```go
func New(cfg Config) *Reminder {
```

```go
func (r *Reminder) Middleware() dispatcher.MiddlewareFunc {
```

## `shopengine`

```go
func New(cfg Config) *ShopEngine {
```

```go
func (se *ShopEngine) Middleware() dispatcher.MiddlewareFunc {
```

```go
func (se *ShopEngine) ShowMainMenu(ctx context.Context, chatID string) {
```

## `vpndetector`

```go
func DefaultConfig() Config {
```

```go
func New(cfg Config) *VPNDetector {
```

```go
func (vd *VPNDetector) VPNDetectorMiddleware() dispatcher.MiddlewareFunc {
```

```go
func (vd *VPNDetector) Stats() map[string]int64 {
```

## `warnsystem`

```go
func DefaultConfig() Config {
```

```go
func New(cfg Config) *WarnSystem {
```

```go
func (ws *WarnSystem) AddWarning(ctx context.Context, chatID, userID string) int {
```

```go
func (ws *WarnSystem) ResetWarnings(chatID, userID string) bool {
```

```go
func (ws *WarnSystem) GetWarnings(chatID, userID string) int {
```

```go
func (ws *WarnSystem) Stats() map[string]int64 {
```

```go
func (ws *WarnSystem) Close() {
```
