package groupmanager

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/contrib/antilink"
	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

type CommandsConfig struct {
	Ban, Kick, Mute, Unmute, Warn, Unwarn string
	Lock, Unlock, LockAll, UnlockAll      string
	Pin, Unpin, Admins, Members           string
	Help, Info, Rules, SetRules           string
	Purge, Report, ID, Gban, Ungban       string
	LockWord, UnlockWord, Note            string
	Start, Stop, Settings, SlowMode       string
}

type MessagesConfig struct {
	BanMsg, MuteMsg, WarnMsg           string
	LockSuccess, UnlockSuccess         string
	FloodMsg, NoPermission, NotReplied string
	WelcomeMsg, GoodbyeMsg             string
	PrivateStartMsg, GroupStartMsg     string
	NotBotAdminMsg                     string
}

type Config struct {
	BotID    string
	AdminIDs []string
	Commands CommandsConfig
	Messages MessagesConfig

	// نمونه ضد لینک برای بررسی دقیق لینک‌ها
	AntiLink *antilink.AntiLink

	EnableLocks      bool
	EnableWarns      bool
	EnableNotes      bool
	EnableFlood      bool
	EnableBlacklist  bool
	EnableWelcome    bool
	EnableGBan       bool
	EnableRules      bool
	EnableAutoDelete bool
	EnableLockBots   bool

	MaxWarns          int
	FloodLimit        int
	FloodWindow       time.Duration
	FloodMuteDuration time.Duration

	DeleteMessageAction func(ctx context.Context, chatID string, messageID int) error
	RestrictUserAction  func(ctx context.Context, chatID, userID string, untilDate int64, perms *models.ChatPermissions) error
	BanUserAction       func(ctx context.Context, chatID, userID string) error
	UnbanUserAction     func(ctx context.Context, chatID, userID string) error
	SendMessageAction   func(ctx context.Context, chatID, text string, markup *models.InlineKeyboardMarkup) (int, error)
	GetChatMemberAction func(ctx context.Context, chatID, userID string) (*models.ChatMember, error)
	PinMessageAction    func(ctx context.Context, chatID string, messageID int, disableNotif bool) error
	UnpinMessageAction  func(ctx context.Context, chatID string, messageID int) error
	GetChatAdminsAction func(ctx context.Context, chatID string) ([]models.ChatMember, error)
	GetChatMemberCount  func(ctx context.Context, chatID string) (int, error)
}

type GroupManager struct {
	cfg        Config
	admins     map[string]struct{}
	adminCache sync.Map
	antiLink   *antilink.AntiLink

	lockManager      *LockManager
	warnManager      *WarnManager
	noteManager      *NoteManager
	blacklistManager *BlacklistManager
	floodManager     *FloodManager

	activeGroups  *SyncStore
	gbans         *SyncStore
	mutedUsers    *SyncStore
	rulesStore    *SyncStore
	slowModeStore *SyncStore
}

func New(cfg Config) *GroupManager {
	if cfg.Commands.Ban == "" {
		cfg.Commands = CommandsConfig{
			Ban: "بن", Kick: "اخراج", Mute: "سکوت", Unmute: "آزاد",
			Warn: "اخطار", Unwarn: "حذف اخطار", Lock: "قفل", Unlock: "باز کردن قفل",
			Pin: "سنجاق", Unpin: "حذف سنجاق", Admins: "ادمین ها", Members: "اعضا",
			Help: "راهنما", Info: "اطلاعات", Rules: "قوانین", SetRules: "تنظیم قوانین",
			LockAll: "قفل همه", UnlockAll: "باز کردن همه", Purge: "پاکسازی",
			Report: "گزارش", ID: "آیدی", Gban: "گبن", Ungban: "حذف گبن",
			LockWord: "قفل کلمه", UnlockWord: "باز کردن کلمه", Note: "یادداشت",
			Start: "start", Stop: "stop", Settings: "تنظیمات", SlowMode: "اسلوومود",
		}
	}
	if cfg.Messages.PrivateStartMsg == "" {
		cfg.Messages.PrivateStartMsg = "<b>🤖 سلام!</b>\nمن ربات مدیریت گروه هستم."
	}
	if cfg.Messages.GroupStartMsg == "" {
		cfg.Messages.GroupStartMsg = "<b>✅ سیستم فعال شد</b>\nگروه تحت حفاظت قرار گرفت."
	}
	if cfg.Messages.NotBotAdminMsg == "" {
		cfg.Messages.NotBotAdminMsg = "<b>🚫 خطای دسترسی</b>\nمن ادمین نیستم!"
	}
	if cfg.Messages.WelcomeMsg == "" {
		cfg.Messages.WelcomeMsg = "👋 خوش آمدی {user}!"
	}
	if cfg.Messages.GoodbyeMsg == "" {
		cfg.Messages.GoodbyeMsg = "👋 خداحافظ {user}!"
	}
	if cfg.Messages.BanMsg == "" {
		cfg.Messages.BanMsg = "<b>🚫 اخراج شد</b>\n👤 {user}\n🛡 {reason}"
	}
	if cfg.Messages.MuteMsg == "" {
		cfg.Messages.MuteMsg = "<b>🤐 سکوت شد</b>\n👤 {user}\n⏱ {duration}"
	}
	if cfg.Messages.WarnMsg == "" {
		cfg.Messages.WarnMsg = "<b>⚠️ اخطار {current}/{max}</b>\n👤 {user}\n🛡 {reason}"
	}
	if cfg.Messages.LockSuccess == "" {
		cfg.Messages.LockSuccess = "🔒 قفل {lockType} فعال شد."
	}
	if cfg.Messages.UnlockSuccess == "" {
		cfg.Messages.UnlockSuccess = "🔓 قفل {lockType} باز شد."
	}
	if cfg.Messages.NotReplied == "" {
		cfg.Messages.NotReplied = "❌ روی پیام شخص ریپلای بزنید."
	}
	if cfg.Messages.FloodMsg == "" {
		cfg.Messages.FloodMsg = "🌊 لطفاً اسپم نکنید! شما موقتاً محدود شدید."
	}
	if cfg.Messages.NoPermission == "" {
		cfg.Messages.NoPermission = "🚫 شما دسترسی لازم برای این کار را ندارید."
	}

	gm := &GroupManager{
		cfg: cfg, admins: make(map[string]struct{}),
		activeGroups: &SyncStore{}, gbans: &SyncStore{}, mutedUsers: &SyncStore{},
		rulesStore: &SyncStore{}, slowModeStore: &SyncStore{},
		antiLink: cfg.AntiLink,
	}
	for _, id := range cfg.AdminIDs {
		gm.admins[id] = struct{}{}
	}

	if cfg.EnableLocks {
		gm.lockManager = NewLockManager()
	}
	if cfg.EnableWarns {
		gm.warnManager = NewWarnManager(cfg.MaxWarns)
	}
	if cfg.EnableNotes {
		gm.noteManager = NewNoteManager()
	}
	if cfg.EnableBlacklist {
		gm.blacklistManager = NewBlacklistManager()
	}
	if cfg.EnableFlood {
		gm.floodManager = NewFloodManager(cfg.FloodLimit, cfg.FloodWindow)
	}

	return gm
}

func (gm *GroupManager) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u == nil || u.Message == nil {
				next(ctx, u)
				return
			}
			msg := u.Message
			chatID := msg.Chat.ID
			userID := ""
			if msg.From != nil {
				userID = msg.From.ID
			}

			if len(msg.NewChatMembers) > 0 {
				for _, member := range msg.NewChatMembers {
					if member.ID == gm.cfg.BotID {
						continue
					}
					if gm.cfg.EnableGBan {
						if _, banned := gm.gbans.Get(member.ID); banned {
							_ = gm.cfg.BanUserAction(ctx, chatID, member.ID)
							continue
						}
					}
					if gm.cfg.EnableLockBots && member.IsBot && !gm.isBotAdmin(userID) {
						_ = gm.cfg.BanUserAction(ctx, chatID, member.ID)
						continue
					}
					if gm.cfg.EnableWelcome {
						_, _ = gm.cfg.SendMessageAction(ctx, chatID, strings.ReplaceAll(gm.cfg.Messages.WelcomeMsg, "{user}", member.FullName()), nil)
					}
				}
				return
			}
			if msg.LeftChatMember != nil && gm.cfg.EnableWelcome {
				_, _ = gm.cfg.SendMessageAction(ctx, chatID, strings.ReplaceAll(gm.cfg.Messages.GoodbyeMsg, "{user}", msg.LeftChatMember.FullName()), nil)
				return
			}

			if msg.Text != "" {
				lowerText := strings.ToLower(strings.TrimLeft(msg.Text, "/#!."))
				var cmdStr, args string

				if strings.HasPrefix(lowerText, gm.cfg.Commands.Unlock) {
					cmdStr = gm.cfg.Commands.Unlock
					args = strings.TrimSpace(lowerText[len(gm.cfg.Commands.Unlock):])
				} else if strings.HasPrefix(lowerText, gm.cfg.Commands.SetRules) {
					cmdStr = gm.cfg.Commands.SetRules
					args = strings.TrimSpace(lowerText[len(gm.cfg.Commands.SetRules):])
				} else if strings.HasPrefix(lowerText, gm.cfg.Commands.LockWord) {
					cmdStr = gm.cfg.Commands.LockWord
					args = strings.TrimSpace(lowerText[len(gm.cfg.Commands.LockWord):])
				} else if strings.HasPrefix(lowerText, gm.cfg.Commands.UnlockWord) {
					cmdStr = gm.cfg.Commands.UnlockWord
					args = strings.TrimSpace(lowerText[len(gm.cfg.Commands.UnlockWord):])
				} else if strings.HasPrefix(lowerText, gm.cfg.Commands.Note+" ") {
					cmdStr = gm.cfg.Commands.Note
					args = strings.TrimSpace(lowerText[len(gm.cfg.Commands.Note):])
				} else if strings.HasPrefix(lowerText, gm.cfg.Commands.SlowMode) {
					cmdStr = gm.cfg.Commands.SlowMode
					args = strings.TrimSpace(lowerText[len(gm.cfg.Commands.SlowMode):])
				} else {
					parts := strings.Fields(lowerText)
					if len(parts) > 0 {
						cmdStr = parts[0]
						if len(parts) > 1 {
							args = strings.Join(parts[1:], " ")
						}
					}
				}

				if cmdStr != "" {
					if cmdStr == gm.cfg.Commands.Start || cmdStr == gm.cfg.Commands.Stop {
						if gm.isBotAdmin(userID) || gm.isChatAdmin(ctx, userID, chatID) {
							if cmdStr == gm.cfg.Commands.Start {
								if msg.Chat.IsPrivate() {
									_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.PrivateStartMsg, nil)
									return
								}
								if gm.cfg.BotID != "" && gm.cfg.GetChatMemberAction != nil {
									botMember, err := gm.cfg.GetChatMemberAction(ctx, chatID, gm.cfg.BotID)
									if err != nil || (botMember.Status != "creator" && botMember.Status != "administrator") {
										_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.NotBotAdminMsg, nil)
										return
									}
								}
								gm.activeGroups.Set(chatID, true)
								_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.GroupStartMsg, nil)
							} else {
								gm.activeGroups.Delete(chatID)
								_, _ = gm.cfg.SendMessageAction(ctx, chatID, "❌ ربات خاموش شد.", nil)
							}
							return
						}
					}

					if _, isActive := gm.activeGroups.Get(chatID); !isActive {
						next(ctx, u)
						return
					}

					if cmdStr == gm.cfg.Commands.Help || cmdStr == gm.cfg.Commands.ID || (gm.cfg.EnableRules && cmdStr == gm.cfg.Commands.Rules) || cmdStr == gm.cfg.Commands.Report {
						if gm.handlePublicCommand(ctx, msg, userID, cmdStr) {
							return
						}
					}

					if gm.cfg.EnableAutoDelete && (strings.HasPrefix(msg.Text, "/") || strings.HasPrefix(msg.Text, "#")) {
						if !gm.isBotAdmin(userID) && !gm.isChatAdmin(ctx, userID, chatID) {
							_ = gm.cfg.DeleteMessageAction(ctx, chatID, msg.MessageID)
							return
						}
					}

					if gm.isBotAdmin(userID) || gm.isChatAdmin(ctx, userID, chatID) {
						if gm.handleAdminCommand(ctx, msg, userID, cmdStr, args) {
							return
						}
					}
				} else {
					if _, isActive := gm.activeGroups.Get(chatID); !isActive {
						next(ctx, u)
						return
					}
					if gm.cfg.EnableNotes {
						if val, ok := gm.noteManager.GetNote(chatID, lowerText); ok {
							_, _ = gm.cfg.SendMessageAction(ctx, chatID, val, nil)
							return
						}
					}
				}
			}

			if _, isActive := gm.activeGroups.Get(chatID); !isActive {
				next(ctx, u)
				return
			}

			if !gm.isBotAdmin(userID) && !gm.isChatAdmin(ctx, userID, chatID) {
				if muteUntil, ok := gm.mutedUsers.Get(chatID + "_" + userID); ok {
					if time.Now().Unix() < muteUntil.(int64) {
						_ = gm.cfg.DeleteMessageAction(ctx, chatID, msg.MessageID)
						return
					} else {
						gm.mutedUsers.Delete(chatID + "_" + userID)
					}
				}

				if val, ok := gm.slowModeStore.Get(chatID + "_" + userID); ok {
					if time.Since(val.(time.Time)) < 10*time.Second {
						_ = gm.cfg.DeleteMessageAction(ctx, chatID, msg.MessageID)
						return
					}
				}

				if gm.cfg.EnableGBan {
					if _, banned := gm.gbans.Get(userID); banned {
						_ = gm.cfg.BanUserAction(ctx, chatID, userID)
						return
					}
				}
				if gm.cfg.EnableFlood && gm.floodManager.IsFlooding(userID) {
					_ = gm.cfg.DeleteMessageAction(ctx, chatID, msg.MessageID)
					muteUntil := time.Now().Add(gm.cfg.FloodMuteDuration).Unix()
					_ = gm.cfg.RestrictUserAction(ctx, chatID, userID, muteUntil, &models.ChatPermissions{})
					_, _ = gm.cfg.SendMessageAction(ctx, chatID, gm.cfg.Messages.FloodMsg, nil)
					return
				}
				if gm.cfg.EnableBlacklist && gm.blacklistManager.CheckViolation(chatID, msg.Text) {
					_ = gm.cfg.DeleteMessageAction(ctx, chatID, msg.MessageID)
					if gm.cfg.EnableWarns && gm.warnManager.AddWarning(chatID, userID) >= gm.cfg.MaxWarns {
						_ = gm.cfg.BanUserAction(ctx, chatID, userID)
					}
					return
				}
				if gm.cfg.EnableLocks {
					if violation := gm.CheckLockViolation(msg); violation != "" {
						_ = gm.cfg.DeleteMessageAction(ctx, chatID, msg.MessageID)
						if gm.cfg.EnableWarns {
							warns := gm.warnManager.AddWarning(chatID, userID)
							if warns >= gm.cfg.MaxWarns {
								userName := ""
								if msg.From != nil {
									userName = msg.From.FullName()
								}
								_, _ = gm.cfg.SendMessageAction(ctx, chatID, strings.ReplaceAll(gm.cfg.Messages.BanMsg, "{user}", userName), nil)
							} else {
								userName := ""
								if msg.From != nil {
									userName = msg.From.FullName()
								}
								warnMsg := strings.ReplaceAll(gm.cfg.Messages.WarnMsg, "{user}", userName)
								warnMsg = strings.ReplaceAll(warnMsg, "{current}", fmt.Sprintf("%d", warns))
								warnMsg = strings.ReplaceAll(warnMsg, "{max}", fmt.Sprintf("%d", gm.cfg.MaxWarns))
								warnMsg = strings.ReplaceAll(warnMsg, "{reason}", violation)
								_, _ = gm.cfg.SendMessageAction(ctx, chatID, warnMsg, nil)
							}
						}
						return
					}
				}
			}
			next(ctx, u)
		}
	}
}
