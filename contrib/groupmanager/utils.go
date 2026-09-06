package groupmanager

import (
	"strings"
	"sync"
	"time"
)

// SyncStore
type SyncStore struct{ data sync.Map }

func (s *SyncStore) Set(key string, val interface{})    { s.data.Store(key, val) }
func (s *SyncStore) Get(key string) (interface{}, bool) { return s.data.Load(key) }
func (s *SyncStore) Delete(key string)                  { s.data.Delete(key) }

// LockManager
type LockManager struct {
	mu    sync.RWMutex
	locks map[string]map[string]bool
}

func NewLockManager() *LockManager { return &LockManager{locks: make(map[string]map[string]bool)} }

func (lm *LockManager) getLockType(input string) string {
	switch strings.ToLower(input) {
	case "لینک", "link", "links":
		return "links"
	case "عکس", "photo":
		return "photo"
	case "ویدیو", "video":
		return "video"
	case "ویس", "voice":
		return "voice"
	case "فایل", "document":
		return "document"
	case "استیکر", "sticker":
		return "sticker"
	case "گیف", "gif", "animation":
		return "gif"
	case "فوروارد", "forward":
		return "forward"
	case "متن", "text":
		return "text"
	case "فحش", "fosh":
		return "fosh"
	case "عربی", "arabic":
		return "arabic"
	case "مدیا", "media":
		return "media"
	default:
		return strings.ToLower(input)
	}
}

func (lm *LockManager) SetLock(chatID, lockType string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	if lm.locks[chatID] == nil {
		lm.locks[chatID] = make(map[string]bool)
	}
	lm.locks[chatID][lm.getLockType(lockType)] = true
}

func (lm *LockManager) RemoveLock(chatID, lockType string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	if lm.locks[chatID] != nil {
		delete(lm.locks[chatID], lm.getLockType(lockType))
	}
}

func (lm *LockManager) IsLocked(chatID, lockType string) bool {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	if locks, ok := lm.locks[chatID]; ok {
		return locks[lockType]
	}
	return false
}

func (lm *LockManager) GetLocksList(chatID string) []string {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	var list []string
	if locks, ok := lm.locks[chatID]; ok {
		for k := range locks {
			list = append(list, k)
		}
	}
	return list
}

// WarnManager
type WarnManager struct {
	mu       sync.RWMutex
	warns    map[string]int
	MaxWarns int
}

func NewWarnManager(maxWarns int) *WarnManager {
	if maxWarns <= 0 {
		maxWarns = 3
	}
	return &WarnManager{warns: make(map[string]int), MaxWarns: maxWarns}
}

func (wm *WarnManager) AddWarning(chatID, userID string) int {
	key := chatID + "_" + userID
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.warns[key]++
	return wm.warns[key]
}

func (wm *WarnManager) GetWarnings(chatID, userID string) int {
	key := chatID + "_" + userID
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.warns[key]
}

func (wm *WarnManager) ResetWarnings(chatID, userID string) {
	key := chatID + "_" + userID
	wm.mu.Lock()
	defer wm.mu.Unlock()
	delete(wm.warns, key)
}

// NoteManager
type NoteManager struct {
	mu    sync.RWMutex
	notes map[string]string
}

func NewNoteManager() *NoteManager { return &NoteManager{notes: make(map[string]string)} }

func (nm *NoteManager) AddNote(chatID, keyword, text string) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	nm.notes[chatID+"_"+strings.ToLower(keyword)] = text
}

func (nm *NoteManager) GetNote(chatID, keyword string) (string, bool) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	val, ok := nm.notes[chatID+"_"+strings.ToLower(keyword)]
	return val, ok
}

// BlacklistManager
type BlacklistManager struct {
	mu    sync.RWMutex
	words map[string][]string
}

func NewBlacklistManager() *BlacklistManager {
	return &BlacklistManager{words: make(map[string][]string)}
}

func (bm *BlacklistManager) AddWord(chatID, word string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.words[chatID] = append(bm.words[chatID], strings.ToLower(word))
}

func (bm *BlacklistManager) RemoveWord(chatID, word string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	word = strings.ToLower(word)
	var newList []string
	for _, w := range bm.words[chatID] {
		if w != word {
			newList = append(newList, w)
		}
	}
	bm.words[chatID] = newList
}

func (bm *BlacklistManager) CheckViolation(chatID, text string) bool {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	lowerText := strings.ToLower(text)
	for _, word := range bm.words[chatID] {
		if strings.Contains(lowerText, word) {
			return true
		}
	}
	return false
}

// FloodManager
type FloodManager struct {
	mu     sync.Mutex
	users  map[string][]time.Time
	limit  int
	window time.Duration
}

func NewFloodManager(limit int, window time.Duration) *FloodManager {
	if limit <= 0 {
		limit = 5
	}
	if window <= 0 {
		window = 5 * time.Second
	}
	return &FloodManager{users: make(map[string][]time.Time), limit: limit, window: window}
}

func (fm *FloodManager) IsFlooding(userID string) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-fm.window)
	valid := fm.users[userID][:0]
	for _, t := range fm.users[userID] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	valid = append(valid, now)
	fm.users[userID] = valid
	return len(valid) > fm.limit
}
