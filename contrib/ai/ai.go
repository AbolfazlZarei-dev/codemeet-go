package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Models

// برای دریافت توکن API لطفا به سایت اصلی مراجعه کنید: https://openrouter.ai/
// برای مشاهده مدل‌های بیشتر به سایت https://openrouter.ai مراجعه کنید.
const (
	ModelDots3NotePreview = "dots-studio/dots-3-note-preview:free"
	ModelLingSante        = "inclusionai/ling-3.0-flash-sante:free"
	ModelLingFin          = "inclusionai/ling-3.0-flash-fin:free"
	ModelLFM2Embedding    = "liquid/lfm-2.5-embedding-350m:free"
	ModelFluxTTS          = "deepgram/flux-tts:free"
	ModelLFM2_2_6B        = "liquid/lfm-2.5-2.6b:free"
	ModelMiniMaxM3        = "minimax/minimax-m3:free"
	ModelMiniMaxM27       = "minimax/minimax-m2.7:free"
)

// Types

type Config struct {
	APIKey  string // توکن کاربر از openrouter.ai
	BaseURL string // پیش‌فرض: https://openrouter.ai/api/v1
}

type Client struct {
	cfg        Config
	httpClient *http.Client
}

type Message struct {
	Role    string `json:"role"` // "user", "assistant", "system"
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Initialization

func New(cfg Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://openrouter.ai/api/v1"
	}

	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // زمان بالا برای استریم شدن پاسخ‌های طولانی
		},
	}
}

// Methods

// StreamChat یک درخواست چت می‌فرستد و پاسخ را تکه به تکه (استریم) برمی‌گرداند.
// تابع onChunk به ازای هر کلمه یا بخشی که هوش مصنوعی تولید می‌کند صدا زده می‌شود.
func (c *Client) StreamChat(ctx context.Context, req ChatRequest, onChunk func(content string)) error {
	req.Stream = true

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.cfg.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	// هدرهای OpenRouter برای نمایش در لیست رنکینگ (اختیاری)
	httpReq.Header.Set("HTTP-Referer", "https://codemeet.chat")
	httpReq.Header.Set("X-Title", "CodeMeet Bot")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(errBody))
	}

	scanner := bufio.NewScanner(resp.Body)
	// افزایش بافر برای جلوگیری از خطا در پاسخ‌های طولانی
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &chunk); err == nil {
				if len(chunk.Choices) > 0 {
					content := chunk.Choices[0].Delta.Content
					if content != "" {
						onChunk(content) // ارسال تکه متن به ربات
					}
				}
			}
		}
	}

	return scanner.Err()
}

// Ask یک درخواست چت ساده (بدون استریم) می‌فرستد و کل پاسخ را برمی‌گرداند.
func (c *Client) Ask(ctx context.Context, req ChatRequest) (string, error) {
	req.Stream = false

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.cfg.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	httpReq.Header.Set("HTTP-Referer", "https://codemeet.chat")
	httpReq.Header.Set("X-Title", "CodeMeet Bot")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(errBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	return result.Choices[0].Message.Content, nil
}
