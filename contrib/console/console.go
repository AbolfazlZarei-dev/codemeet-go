package console

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// CommandFunc تابعی که هنگام اجرای یک دستور سفارشی صدا زده می‌شود
type CommandFunc func(ctx context.Context, args []string) string

type Console struct {
	commands map[string]CommandFunc
	mu       sync.RWMutex
}

func New() *Console {
	c := &Console{
		commands: make(map[string]CommandFunc),
	}
	c.Register("neofetch", func(ctx context.Context, args []string) string {
		return `
███╗   ██╗███████╗██╗  ██╗ ██████╗ ██████╗ ████████╗███████╗██████╗ 
████╗  ██║██╔════╝╚██╗██╔╝██╔═══██╗██╔══██╗╚══██╔══╝██╔════╝██╔══██╗
██╔██╗ ██║█████╗   ╚███╔╝ ██║   ██║██████╔╝   ██║   █████╗  ██████╔╝
██║╚██╗██║██╔══╝   ██╔██╗ ██║   ██║██╔══██╗   ██║   ██╔══╝  ██╔══██╗
██║ ╚████║███████╗██╔╝ ██╗╚██████╔╝██║  ██║   ██║   ███████╗██║  ██║
╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝ ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝
                                                                   
[+] CodeMeet Go Bot Framework
[+] Status: Fully Operational
[+] Terminal: Interactive Mode Enabled
`
	})
	return c
}

// Register ثبت یک دستور جدید برای ترمینال
func (c *Console) Register(cmd string, handler CommandFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.commands[strings.ToLower(cmd)] = handler
}

// ExecuteCommand برای اجرای دستورات از طریق داشبورد وب
func (c *Console) ExecuteCommand(ctx context.Context, cmd string) string {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return ""
	}

	parts := strings.Fields(cmd)
	command := strings.ToLower(parts[0])
	args := parts[1:]

	// 1. ابتدا بررسی کن آیا دستور سفارشی است
	c.mu.RLock()
	handler, exists := c.commands[command]
	c.mu.RUnlock()

	if exists {
		return handler(ctx, args)
	}

	// 2. اگر دستور سفارشی نبود، آن را به عنوان دستور سیستمی (OS) اجرا کن
	// (امکان اجرای هر دستوری که کاربر خواست)
	execCmd := exec.CommandContext(ctx, command, args...)
	execOut, err := execCmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error executing command '%s': %v\n%s", command, err, string(execOut))
	}

	return string(execOut)
}

// Start شروع خواندن از ترمینال (بلاکینگ - در یک گوروتین جداگانه اجرا شود)
func (c *Console) Start(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("💻 Console started. Type 'help' to see commands.")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}

				result := c.ExecuteCommand(ctx, line)
				if result != "" {
					fmt.Println("OUTPUT:", result)
				}
			}
		}
	}
}
