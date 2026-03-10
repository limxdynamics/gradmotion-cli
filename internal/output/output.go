package output

import (
	"encoding/json"
	"fmt"
	"os"
)

type Meta struct {
	TraceID    string `json:"trace_id,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
	DurationMS int64  `json:"duration_ms,omitempty"`
	Command    string `json:"command,omitempty"`
	Profile    string `json:"profile,omitempty"`
	Endpoint   string `json:"endpoint,omitempty"`
	Method     string `json:"method,omitempty"`
	Status     int    `json:"status,omitempty"`
	RetryCount int    `json:"retry_count,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Hint    string `json:"hint,omitempty"`
}

type Envelope struct {
	Success bool       `json:"success"`
	Data    any        `json:"data"`
	Meta    Meta       `json:"meta"`
	Error   *ErrorInfo `json:"error"`
}

func Print(env Envelope, human, quiet bool) error {
	switch {
	case human:
		return printHuman(env)
	case quiet:
		return printQuiet(env)
	default:
		return printJSON(env)
	}
}

func printJSON(env Envelope) error {
	b, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(os.Stdout, string(b))
	return err
}

func printQuiet(env Envelope) error {
	if env.Success {
		if env.Data == nil {
			_, err := fmt.Fprintln(os.Stdout, "{}")
			return err
		}
		b, err := json.MarshalIndent(env.Data, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(os.Stdout, string(b))
		return err
	}

	b, err := json.MarshalIndent(env.Error, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(os.Stdout, string(b))
	return err
}

func printHuman(env Envelope) error {
	if env.Success {
		_, _ = fmt.Fprintln(os.Stdout, "OK")
		if env.Data == nil {
			return nil
		}
		b, err := json.MarshalIndent(env.Data, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(os.Stdout, string(b))
		return err
	}

	_, _ = fmt.Fprintln(os.Stdout, "FAILED")
	if env.Error != nil {
		_, _ = fmt.Fprintf(os.Stdout, "code: %s\n", env.Error.Code)
		_, _ = fmt.Fprintf(os.Stdout, "message: %s\n", env.Error.Message)
		if env.Error.Hint != "" {
			_, _ = fmt.Fprintf(os.Stdout, "hint: %s\n", env.Error.Hint)
		}
	}
	return nil
}
