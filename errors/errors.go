package errors

import (
	"fmt"
	"strconv"
	"strings"
)

// ErrorCode کد خطای API
type ErrorCode int

const (
	CodeBadRequest         ErrorCode = 400
	CodeUnauthorized       ErrorCode = 401
	CodeForbidden          ErrorCode = 403
	CodeNotFound           ErrorCode = 404
	CodeMethodNotAllowed   ErrorCode = 405
	CodeConflict           ErrorCode = 409
	CodeRequestTooLarge    ErrorCode = 413
	CodeTooManyRequests    ErrorCode = 429
	CodeInternal           ErrorCode = 500
	CodeNotImplemented     ErrorCode = 501
	CodeBadGateway         ErrorCode = 502
	CodeServiceUnavailable ErrorCode = 503
	CodeGatewayTimeout     ErrorCode = 504
)

// APIError خطای برگشتی از API
type APIError struct {
	Code        ErrorCode
	Description string
	RetryAfter  int
	Parameters  map[string]interface{}
}

func (e *APIError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("codemeet api error [%d]: %s (retry_after=%ds)", e.Code, e.Description, e.RetryAfter)
	}
	return fmt.Sprintf("codemeet api error [%d]: %s", e.Code, e.Description)
}

// IsRetryable آیا خطا قابل تلاش مجدد است؟
func (e *APIError) IsRetryable() bool {
	switch e.Code {
	case CodeTooManyRequests,
		CodeBadGateway,
		CodeInternal,
		CodeServiceUnavailable,
		CodeGatewayTimeout:
		return true
	}
	// در صورت وجود retry_after قطعا قابل retry است
	return e.RetryAfter > 0
}

// Is نوعی از خطا
func (e *APIError) Is(target error) bool {
	if t, ok := target.(*APIError); ok {
		return e.Code == t.Code
	}
	return false
}

// As interface برای خطاهای قابل retry
type retryableErr interface {
	IsRetryable() bool
}

// ValidationError خطای اعتبارسنجی
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' - %s", e.Field, e.Message)
}

// NewValidationError ساخت خطای اعتبارسنجی
func NewValidationError(field, msg string) *ValidationError {
	return &ValidationError{Field: field, Message: msg}
}

// NetworkError خطای شبکه
type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

func (e *NetworkError) Unwrap() error { return e.Err }

// IsRetryable شبکه قابل تلاش مجدد است
func (e *NetworkError) IsRetryable() bool { return true }

// NewNetworkError ساخت خطای شبکه
func NewNetworkError(err error) *NetworkError {
	return &NetworkError{Err: err}
}

// IsRetryable بررسی قابل تلاش مجدد بودن خطا
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	if r, ok := err.(retryableErr); ok {
		return r.IsRetryable()
	}
	// اگر unwrapped داشته باشد
	if u := unwrap(err); u != nil {
		if r, ok := u.(retryableErr); ok {
			return r.IsRetryable()
		}
	}
	return false
}

func unwrap(err error) error {
	type unwrapper interface{ Unwrap() error }
	if u, ok := err.(unwrapper); ok {
		return u.Unwrap()
	}
	return nil
}

// ParseError تبدیل پاسخ خطا به APIError
// علاوه بر parameters، در صورت وجود retry_after در description هم استخراج می‌کند
func ParseError(code int, desc string, params map[string]interface{}) *APIError {
	e := &APIError{
		Code:        ErrorCode(code),
		Description: desc,
		Parameters:  params,
	}
	// اول از params
	if v, ok := params["retry_after"]; ok {
		switch t := v.(type) {
		case float64:
			e.RetryAfter = int(t)
		case int:
			e.RetryAfter = t
		case string:
			if n, err := strconv.Atoi(t); err == nil {
				e.RetryAfter = n
			}
		}
	}
	// اگر در params نبود، در description بگرد: "Too Many Requests: retry after 5"
	if e.RetryAfter == 0 {
		if n := extractRetryAfter(desc); n > 0 {
			e.RetryAfter = n
		}
	}
	return e
}

func extractRetryAfter(desc string) int {
	lower := strings.ToLower(desc)
	idx := strings.Index(lower, "retry after")
	if idx == -1 {
		idx = strings.Index(lower, "retry_after")
	}
	if idx == -1 {
		return 0
	}
	rest := desc[idx:]
	// پیدا کردن اولین عدد
	start := -1
	for i := 0; i < len(rest); i++ {
		if rest[i] >= '0' && rest[i] <= '9' {
			if start == -1 {
				start = i
			}
		} else if start != -1 {
			n, err := strconv.Atoi(rest[start:i])
			if err == nil {
				return n
			}
			start = -1
		}
	}
	if start != -1 {
		n, err := strconv.Atoi(rest[start:])
		if err == nil {
			return n
		}
	}
	return 0
}

// MultiError مجموع خطاها
type MultiError struct {
	Errors []error
}

func (m *MultiError) Error() string {
	parts := make([]string, len(m.Errors))
	for i, e := range m.Errors {
		parts[i] = e.Error()
	}
	return strings.Join(parts, "; ")
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

func (m *MultiError) HasError() bool {
	return len(m.Errors) > 0
}

// AsAPIError اگر err از نوع APIError است برمی‌گرداند
func AsAPIError(err error) (*APIError, bool) {
	if e, ok := err.(*APIError); ok {
		return e, true
	}
	return nil, false
}

// AsNetworkError اگر err از نوع NetworkError است
func AsNetworkError(err error) (*NetworkError, bool) {
	if e, ok := err.(*NetworkError); ok {
		return e, true
	}
	return nil, false
}

// AsValidationError اگر err از نوع ValidationError است
func AsValidationError(err error) (*ValidationError, bool) {
	if e, ok := err.(*ValidationError); ok {
		return e, true
	}
	return nil, false
}
