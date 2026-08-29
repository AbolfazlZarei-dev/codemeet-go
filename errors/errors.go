package errors

import (
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"
)

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

func (e *APIError) IsRetryable() bool {
	switch e.Code {
	case CodeTooManyRequests,
		CodeBadGateway,
		CodeInternal,
		CodeServiceUnavailable,
		CodeGatewayTimeout:
		return true
	}
	return e.RetryAfter > 0
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' - %s", e.Field, e.Message)
}

func NewValidationError(field, msg string) *ValidationError {
	return &ValidationError{Field: field, Message: msg}
}

type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

func (e *NetworkError) Unwrap() error     { return e.Err }
func (e *NetworkError) IsRetryable() bool { return true }

func NewNetworkError(err error) *NetworkError {
	return &NetworkError{Err: err}
}

type retryableErr interface {
	IsRetryable() bool
}

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	var r retryableErr
	if stderrors.As(err, &r) {
		return r.IsRetryable()
	}
	return false
}

func ParseError(code int, desc string, params map[string]interface{}) *APIError {
	e := &APIError{
		Code:        ErrorCode(code),
		Description: desc,
		Parameters:  params,
	}
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

// Unwrap برای پشتیبانی از errors.Is و errors.As در زنجیره خطاها
func (m *MultiError) Unwrap() []error {
	return m.Errors
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

func (m *MultiError) HasError() bool {
	return len(m.Errors) > 0
}

// استفاده از errors.As استاندارد برای تطبیق دقیق
func AsAPIError(err error) (*APIError, bool) {
	var e *APIError
	if stderrors.As(err, &e) {
		return e, true
	}
	return nil, false
}

func AsNetworkError(err error) (*NetworkError, bool) {
	var e *NetworkError
	if stderrors.As(err, &e) {
		return e, true
	}
	return nil, false
}

func AsValidationError(err error) (*ValidationError, bool) {
	var e *ValidationError
	if stderrors.As(err, &e) {
		return e, true
	}
	return nil, false
}
