package client

import (
	"strings"
	"testing"
)

func TestAPIError_Error_PrefersStructuredFields(t *testing.T) {
	e := &APIError{Status: 404, Code: "not_found", Message: "missing repo"}
	if !strings.Contains(e.Error(), "not_found") || !strings.Contains(e.Error(), "missing repo") {
		t.Fatalf("bad error string: %q", e.Error())
	}
}

func TestParseAPIError_DecodesAikidoShape(t *testing.T) {
	body := []byte(`{"error":"validation_failed","message":"name is required"}`)
	e := parseAPIError(400, body)
	if e.Code != "validation_failed" || e.Message != "name is required" {
		t.Fatalf("bad parse: %+v", e)
	}
}

func TestParseAPIError_FallsBackToBody(t *testing.T) {
	body := []byte("Internal Server Error")
	e := parseAPIError(500, body)
	if e.Status != 500 || !strings.Contains(e.Message, "Internal") {
		t.Fatalf("bad fallback: %+v", e)
	}
}
