package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO "
	case LevelWarn:
		return "WARN "
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	}
	return "?    "
}

func (l Level) ColorString() string {
	switch l {
	case LevelDebug:
		return ColorGray + l.String() + ColorReset
	case LevelInfo:
		return ColorCyan + l.String() + ColorReset
	case LevelWarn:
		return ColorYellow + l.String() + ColorReset
	case LevelError:
		return ColorRed + l.String() + ColorReset
	case LevelFatal:
		return ColorBold + ColorRed + l.String() + ColorReset
	}
	return l.String()
}

type Format int

const (
	FormatText Format = iota
	FormatJSON
)

type Logger struct {
	mu            sync.Mutex
	level         atomic.Int32
	format        Format
	out           io.Writer
	fields        []interface{}
	async         bool
	logCh         chan logEntry
	wg            sync.WaitGroup
	includeCaller atomic.Bool
	enabled       atomic.Bool // قابلیت استپ و استارت لاگ‌ها
}

type logEntry struct {
	level  Level
	msg    string
	fields []interface{}
	ts     time.Time
	file   string
	line   int
}

func enableWindowsANSI() {
	var mode uint32
	stdoutHandle := syscall.Handle(os.Stdout.Fd())
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")

	r, _, _ := procGetConsoleMode.Call(uintptr(stdoutHandle), uintptr(unsafe.Pointer(&mode)))
	if r != 0 {
		mode |= 0x0004
		procSetConsoleMode.Call(uintptr(stdoutHandle), uintptr(mode))
	}
}

func init() {
	if runtime.GOOS == "windows" {
		enableWindowsANSI()
	}
}

func New(level Level) *Logger {
	l := &Logger{
		format: FormatText,
		out:    os.Stdout,
	}
	l.level.Store(int32(level))
	l.includeCaller.Store(true)
	l.enabled.Store(true) // پیش‌فرض لاگ‌ها روشن است
	return l
}

func NewJSON(level Level) *Logger {
	l := &Logger{
		format: FormatJSON,
		out:    os.Stdout,
	}
	l.level.Store(int32(level))
	l.includeCaller.Store(true)
	l.enabled.Store(true)
	return l
}

func NewAsync(level Level, bufferSize int) *Logger {
	l := &Logger{
		format: FormatText,
		out:    os.Stdout,
		async:  true,
		logCh:  make(chan logEntry, bufferSize),
	}
	l.level.Store(int32(level))
	l.includeCaller.Store(true)
	l.enabled.Store(true)
	l.wg.Add(1)
	go l.asyncWriter()
	return l
}

// SetIncludeCaller برای کنترل خاموش/روشن کردن گرفتن اطلاعات فایل
func (l *Logger) SetIncludeCaller(b bool) {
	l.includeCaller.Store(b)
}

// SetEnabled برای استپ و استارت کردن لاگ‌ها از بیرون
func (l *Logger) SetEnabled(b bool) {
	l.enabled.Store(b)
}

// IsEnabled بررسی وضعیت روشن یا خاموش بودن لاگ‌ها
func (l *Logger) IsEnabled() bool {
	return l.enabled.Load()
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) SetLevel(level Level) {
	l.level.Store(int32(level))
}

func (l *Logger) SetFormat(format Format) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.format = format
}

func (l *Logger) WithFields(fields ...interface{}) *Logger {
	if len(fields)%2 != 0 {
		fields = append(fields, "<MISSING_VALUE>")
	}
	newLogger := &Logger{
		format: l.format,
		out:    l.out,
		fields: append(append([]interface{}{}, l.fields...), fields...),
	}
	newLogger.level.Store(l.level.Load())
	newLogger.includeCaller.Store(l.includeCaller.Load())
	newLogger.enabled.Store(l.enabled.Load())
	if l.async {
		newLogger.async = true
		newLogger.logCh = l.logCh
	}
	return newLogger
}

func (l *Logger) log(level Level, msg string, fields ...interface{}) {
	// اگر لاگ‌ها خاموش بودند، هیچ چیزی پردازش نکن
	if !l.enabled.Load() {
		return
	}

	if int32(level) < l.level.Load() {
		return
	}

	var file string
	var line int

	if l.includeCaller.Load() {
		_, file, line, _ = runtime.Caller(2)
	}

	if l.async {
		f := make([]interface{}, len(fields))
		copy(f, fields)
		select {
		case l.logCh <- logEntry{
			level:  level,
			msg:    msg,
			fields: f,
			ts:     time.Now(),
			file:   file,
			line:   line,
		}:
		default:
			l.write(logEntry{level: level, msg: msg, fields: fields, ts: time.Now(), file: file, line: line})
		}
		return
	}

	l.write(logEntry{level: level, msg: msg, fields: fields, ts: time.Now(), file: file, line: line})
}

func (l *Logger) asyncWriter() {
	defer l.wg.Done()
	for entry := range l.logCh {
		l.write(entry)
	}
}

func (l *Logger) write(e logEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	allFields := append(append([]interface{}{}, l.fields...), e.fields...)

	if l.format == FormatJSON {
		l.writeJSON(e, allFields)
	} else {
		l.writeText(e, allFields)
	}
}

func (l *Logger) writeText(e logEntry, fields []interface{}) {
	msgColor := ColorReset
	switch e.level {
	case LevelInfo:
		msgColor = ColorGreen
	case LevelWarn:
		msgColor = ColorYellow
	case LevelError:
		msgColor = ColorRed
	case LevelFatal:
		msgColor = ColorBold + ColorRed
	case LevelDebug:
		msgColor = ColorGray
	}

	timeStr := fmt.Sprintf("%s[%s]%s", ColorGray, e.ts.Format("2006-01-02 15:04:05"), ColorReset)
	fmt.Fprintf(l.out, "%s %s %s%s%s", timeStr, e.level.ColorString(), msgColor, e.msg, ColorReset)

	if len(fields) > 0 {
		fmt.Fprint(l.out, " ")
		for i := 0; i+1 < len(fields); i += 2 {
			fmt.Fprintf(l.out, "%s%v%s=%s%v%s ", ColorBlue, fields[i], ColorReset, ColorPurple, fields[i+1], ColorReset)
		}
	}

	if e.file != "" {
		fmt.Fprintf(l.out, "%s(%s:%d)%s", ColorGray, shortFile(e.file), e.line, ColorReset)
	}

	fmt.Fprintln(l.out)
}

func (l *Logger) writeJSON(e logEntry, fields []interface{}) {
	m := map[string]interface{}{
		"ts":    e.ts.Format(time.RFC3339Nano),
		"level": e.level.String(),
		"msg":   e.msg,
		"file":  shortFile(e.file),
		"line":  e.line,
	}
	for i := 0; i+1 < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		m[key] = fields[i+1]
	}
	b, _ := json.Marshal(m)
	b = append(b, '\n')
	l.out.Write(b)
}

func shortFile(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' || s[i] == '\\' {
			return s[i+1:]
		}
	}
	return s
}

func (l *Logger) Debug(msg string, fields ...interface{}) { l.log(LevelDebug, msg, fields...) }
func (l *Logger) Info(msg string, fields ...interface{})  { l.log(LevelInfo, msg, fields...) }
func (l *Logger) Warn(msg string, fields ...interface{})  { l.log(LevelWarn, msg, fields...) }
func (l *Logger) Error(msg string, fields ...interface{}) { l.log(LevelError, msg, fields...) }
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.log(LevelFatal, msg, fields...)
	l.Close()
	os.Exit(1)
}

func (l *Logger) Close() {
	if l.async {
		l.mu.Lock()
		if l.logCh != nil {
			close(l.logCh)
			l.logCh = nil
		}
		l.mu.Unlock()
		l.wg.Wait()
	}
}

func (l *Logger) Sync() {
	if l.async {
		for len(l.logCh) > 0 {
			time.Sleep(time.Millisecond)
		}
	}
	if s, ok := l.out.(interface{ Sync() error }); ok {
		s.Sync()
	}
}

func (l *Logger) Output() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}
