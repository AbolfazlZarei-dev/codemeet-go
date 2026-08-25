package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// Level سطح لاگ
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	}
	return "?"
}

// ColorString نمایش رنگی سطح لاگ برای ترمینال
func (l Level) ColorString() string {
	switch l {
	case LevelDebug:
		return ColorGray + l.String() + ColorReset
	case LevelInfo:
		return ColorGreen + l.String() + ColorReset
	case LevelWarn:
		return ColorYellow + l.String() + ColorReset
	case LevelError:
		return ColorRed + l.String() + ColorReset
	case LevelFatal:
		return ColorPurple + l.String() + ColorReset
	}
	return l.String()
}

// Format فرمت لاگ
type Format int

const (
	FormatText Format = iota
	FormatJSON
)

// Logger لاگر ساختاریافته
type Logger struct {
	mu     sync.Mutex
	level  Level
	format Format
	out    io.Writer
	fields []interface{}
	async  bool
	logCh  chan logEntry
	wg     sync.WaitGroup
}

type logEntry struct {
	level  Level
	msg    string
	fields []interface{}
	ts     time.Time
	file   string
	line   int
}

// enableWindowsANSI فعال‌سازی رنگ‌ها در ترمینال ویندوز
func enableWindowsANSI() {
	var mode uint32
	stdoutHandle := syscall.Handle(os.Stdout.Fd())
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")

	r, _, _ := procGetConsoleMode.Call(uintptr(stdoutHandle), uintptr(unsafe.Pointer(&mode)))
	if r != 0 {
		mode |= 0x0004 // ENABLE_VIRTUAL_TERMINAL_PROCESSING
		procSetConsoleMode.Call(uintptr(stdoutHandle), uintptr(mode))
	}
}

func init() {
	enableWindowsANSI()
}

func New(level Level) *Logger {
	return &Logger{
		level:  level,
		format: FormatText,
		out:    os.Stdout,
	}
}

func NewJSON(level Level) *Logger {
	return &Logger{
		level:  level,
		format: FormatJSON,
		out:    os.Stdout,
	}
}

func NewAsync(level Level, bufferSize int) *Logger {
	l := &Logger{
		level:  level,
		format: FormatText,
		out:    os.Stdout,
		async:  true,
		logCh:  make(chan logEntry, bufferSize),
	}
	l.wg.Add(1)
	go l.asyncWriter()
	return l
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) SetFormat(format Format) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.format = format
}

func (l *Logger) WithFields(fields ...interface{}) *Logger {
	newLogger := &Logger{
		level:  l.level,
		format: l.format,
		out:    l.out,
		fields: append(append([]interface{}{}, l.fields...), fields...),
	}
	return newLogger
}

func (l *Logger) log(level Level, msg string, fields ...interface{}) {
	if level < l.level {
		return
	}

	_, file, line, _ := runtime.Caller(2)

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
	// استفاده از ColorString برای رنگ‌بندی سطح لاگ
	fmt.Fprintf(l.out, "[%s] %s %s", e.ts.Format(time.RFC3339), e.level.ColorString(), e.msg)
	if len(fields) > 0 {
		fmt.Fprint(l.out, " ")
		for i := 0; i+1 < len(fields); i += 2 {
			fmt.Fprintf(l.out, "%v=%v ", fields[i], fields[i+1])
		}
	}
	if e.file != "" {
		fmt.Fprintf(l.out, " (%s:%d)", shortFile(e.file), e.line)
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
		close(l.logCh)
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
