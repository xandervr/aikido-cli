package client

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	Status  int
	Code    string
	Message string
	Body    []byte
}

func (e *APIError) Error() string {
	if e.Code != "" && e.Message != "" {
		return fmt.Sprintf("aikido api: %d %s: %s", e.Status, e.Code, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("aikido api: %d: %s", e.Status, e.Message)
	}
	return fmt.Sprintf("aikido api: %d", e.Status)
}

func parseAPIError(status int, body []byte) *APIError {
	e := &APIError{Status: status, Body: body}
	var shape struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		Detail  string `json:"detail"`
	}
	if json.Unmarshal(body, &shape) == nil {
		e.Code = shape.Error
		e.Message = shape.Message
		if e.Message == "" {
			e.Message = shape.Detail
		}
	}
	if e.Message == "" {
		e.Message = string(body)
	}
	return e
}
