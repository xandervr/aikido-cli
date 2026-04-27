package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/xandervr/aikido-cli/internal/client"
)

const (
	ExitOK    = 0
	ExitAPI   = 1
	ExitAuth  = 2
	ExitUsage = 3
)

type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string { return e.Err.Error() }

func Exit(err error) {
	if err == nil {
		return
	}
	var ee *ExitError
	if errors.As(err, &ee) {
		fmt.Fprintln(os.Stderr, "error:", ee.Err)
		os.Exit(ee.Code)
	}
	var apiErr *client.APIError
	if errors.As(err, &apiErr) {
		code := ExitAPI
		if apiErr.Status == 401 || apiErr.Status == 403 {
			code = ExitAuth
		}
		fmt.Fprintln(os.Stderr, "error:", apiErr.Error())
		os.Exit(code)
	}
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(ExitAPI)
}
