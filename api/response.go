package api

import (
	"encoding/json"
	"fmt"
)

// Response پاسخ استاندارد API کدمیت
type Response struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result,omitempty"`
	ErrorCode   int             `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
	Parameters  interface{}     `json:"parameters,omitempty"`
	HTTPStatus  int             `json:"-"`
}

// Decode تبدیل Result به ساختار دلخواه
func (r *Response) Decode(v interface{}) error {
	if r == nil {
		return fmt.Errorf("response is nil")
	}
	if len(r.Result) == 0 {
		return nil
	}
	return json.Unmarshal(r.Result, v)
}

// DecodeInto.ToBoolean — متد کمکی برای پاسخ‌های true/false
func (r *Response) AsBool() (bool, error) {
	if r == nil || len(r.Result) == 0 {
		return false, nil
	}
	var b bool
	if err := json.Unmarshal(r.Result, &b); err != nil {
		return false, err
	}
	return b, nil
}

// ParametersAsRetryAfter گرفتن retry_after از parameters
func (r *Response) ParametersAsRetryAfter() int {
	if r == nil || r.Parameters == nil {
		return 0
	}
	if m, ok := r.Parameters.(map[string]interface{}); ok {
		if v, ok := m["retry_after"]; ok {
			switch t := v.(type) {
			case float64:
				return int(t)
			case int:
				return t
			}
		}
	}
	return 0
}
