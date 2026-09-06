# Group Manager

ماژول جامع مدیریت گروه.

## قابلیت‌ها

```text
ban / kick / mute / unmute
warn / unwarn
lock / unlock / lock-all / unlock-all
pin / unpin
admins / members
help / info / rules / set-rules
purge / report / id
gban / ungban
lock-word / unlock-word / note
start / stop / settings / slow-mode
welcome / goodbye
flood protection
blacklist
```

## Config

گزینه‌های اصلی:

```text
BotID
AdminIDs
Commands
Messages
AntiLink
EnableLocks
EnableWarns
EnableNotes
EnableFlood
EnableBlacklist
EnableWelcome
EnableGBan
EnableRules
EnableAutoDelete
EnableLockBots
MaxWarns
FloodLimit
FloodWindow
FloodMuteDuration
```

همچنین Actionهای مربوط به Delete/Restrict/Ban/Unban/Send/Member/Pin/Admin/Count را دریافت می‌کند.

## ساخت

```go
gm := groupmanager.New(groupmanager.Config{
    BotID: bot.Token(),
    AdminIDs: []string{"admin-id"},
    EnableLocks: true,
    EnableWarns: true,
    EnableFlood: true,
})
bot.Use(gm.Middleware())
```

## Managers

```text
LockManager
WarnManager
NoteManager
BlacklistManager
FloodManager
SyncStore
```

## Lock

```go
gm.SetLock(chatID, "links")
gm.IsLocked(chatID, "links")
gm.RemoveLock(chatID, "links")
gm.GetLocksList(chatID)
```

## Warning

```go
gm.AddWarning(chatID, userID)
gm.GetWarnings(chatID, userID)
gm.ResetWarnings(chatID, userID)
```

## Note

```go
gm.AddNote(chatID, "rules", "قوانین...")
text, ok := gm.GetNote(chatID, "rules")
```

## Blacklist

```go
gm.AddWord(chatID, "spam")
gm.RemoveWord(chatID, "spam")
blocked := gm.CheckViolation(chatID, text)
```

## Flood

```go
if gm.IsFlooding(userID) { ... }
```

## Duration

```go
d, err := groupmanager.ParseDuration("30m")
```
