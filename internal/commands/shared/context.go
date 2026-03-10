package shared

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"gradmotion-cli/internal/auth"
	"gradmotion-cli/internal/client"
	"gradmotion-cli/internal/config"
	clilog "gradmotion-cli/internal/log"
	"gradmotion-cli/internal/output"
)

type Runtime struct {
	ConfigManager *config.Manager
	ProfileName   string
	Profile       config.Profile

	Timeout     time.Duration
	Retry       int
	Concurrency int

	Human    bool
	Quiet    bool
	Debug    bool
	ForceYes bool

	logger    *clilog.Logger
	apiClient *client.Client
}

var current *Runtime

func SetRuntime(rt *Runtime) {
	current = rt
}

func GetRuntime() (*Runtime, error) {
	if current == nil {
		return nil, errors.New("runtime not initialized")
	}
	return current, nil
}

func (rt *Runtime) SetLogger(logger *clilog.Logger) {
	rt.logger = logger
}

func (rt *Runtime) Logger() *clilog.Logger {
	return rt.logger
}

func (rt *Runtime) EnsureAPIClient() (*client.Client, error) {
	return rt.ensureAPIClientWithLogger(rt.logger)
}

// ensureAPIClientWithLogger 返回使用指定 logger 的 client；当 logger 为 nil 时不缓存（用于 --no-request-log）。
func (rt *Runtime) ensureAPIClientWithLogger(logger *clilog.Logger) (*client.Client, error) {
	if strings.TrimSpace(rt.Profile.BaseURL) == "" {
		return nil, errors.New("base_url is empty, run config set base_url <url> first")
	}
	key := strings.TrimSpace(rt.Profile.APIKey)
	if key == "" {
		s := auth.NewStore()
		v, found, err := s.Get(rt.ProfileName)
		if err == nil && found && strings.TrimSpace(v) != "" {
			key = strings.TrimSpace(v)
		}
	}
	if key == "" {
		return nil, errors.New("api key is empty, run auth login --api-key <key> first")
	}

	if logger == rt.logger && rt.apiClient != nil {
		return rt.apiClient, nil
	}
	c := client.New(
		rt.Profile.BaseURL,
		key,
		"gradmotion-cli/dev",
		rt.Timeout,
		rt.Retry,
		logger,
	)
	if logger == rt.logger {
		rt.apiClient = c
	}
	return c, nil
}

func EmitLocalSuccess(command string, data any) error {
	rt, err := GetRuntime()
	if err != nil {
		return err
	}
	env := output.Envelope{
		Success: true,
		Data:    data,
		Meta: output.Meta{
			Command: command,
			Profile: rt.ProfileName,
		},
		Error: nil,
	}
	return output.Print(env, rt.Human, rt.Quiet)
}

func EmitLocalError(command, code, msg, hint string) error {
	rt, err := GetRuntime()
	if err != nil {
		return err
	}
	env := output.Envelope{
		Success: false,
		Data:    nil,
		Meta: output.Meta{
			Command: command,
			Profile: rt.ProfileName,
		},
		Error: &output.ErrorInfo{
			Code:    code,
			Message: msg,
			Hint:    hint,
		},
	}
	return output.Print(env, rt.Human, rt.Quiet)
}

func CallAPI(command, method, endpoint string, body any, query map[string]string) error {
	return CallAPIWithOptions(command, method, endpoint, body, query, false, false)
}

func CallAPIAbsolute(command, method, endpoint string, body any, query map[string]string) error {
	return CallAPIWithOptions(command, method, endpoint, body, query, true, false)
}

// CallAPIWithNoRequestLog 与 CallAPI 相同，但可通过 noRequestLog 控制本次是否不向 stderr 输出请求元数据。
func CallAPIWithNoRequestLog(command, method, endpoint string, body any, query map[string]string, noRequestLog bool) error {
	return CallAPIWithOptions(command, method, endpoint, body, query, false, noRequestLog)
}

func CallAPIWithOptions(command, method, endpoint string, body any, query map[string]string, absolutePath, noRequestLog bool) error {
	rt, err := GetRuntime()
	if err != nil {
		return err
	}
	var c *client.Client
	if noRequestLog {
		c, err = rt.ensureAPIClientWithLogger(nil)
	} else {
		c, err = rt.EnsureAPIClient()
	}
	if err != nil {
		return EmitLocalError(command, "CLIENT_CONFIG_ERROR", err.Error(), "请先配置 base_url 和 api key")
	}

	res, err := c.Do(context.Background(), method, endpoint, body, query, absolutePath)
	if err != nil {
		return EmitLocalError(command, "NETWORK_ERROR", err.Error(), "请检查网络连通性与重试配置")
	}

	success := true
	if v, ok := res.Payload["success"].(bool); ok {
		success = v
	} else if res.Meta.Status >= 400 {
		success = false
	}

	data := any(res.Payload)
	if v, ok := res.Payload["data"]; ok {
		data = v
	}

	var errInfo *output.ErrorInfo
	if !success {
		code := fmt.Sprintf("%d", res.Meta.Status)
		if v, ok := res.Payload["code"]; ok {
			code = fmt.Sprintf("%v", v)
		}
		msg := "request failed"
		if v, ok := res.Payload["msg"].(string); ok && strings.TrimSpace(v) != "" {
			msg = v
		}
		errInfo = &output.ErrorInfo{
			Code:    code,
			Message: msg,
		}
	}

	env := output.Envelope{
		Success: success,
		Data:    data,
		Meta: output.Meta{
			TraceID:    res.Meta.TraceID,
			RequestID:  res.Meta.RequestID,
			DurationMS: res.Meta.DurationMS,
			Command:    command,
			Profile:    rt.ProfileName,
			Endpoint:   endpoint,
			Method:     strings.ToUpper(method),
			Status:     res.Meta.Status,
			RetryCount: res.Meta.RetryCount,
		},
		Error: errInfo,
	}
	return output.Print(env, rt.Human, rt.Quiet)
}

// CallAPIRawOutput 发起请求并将 data 中的日志内容以纯文本写入 w，不包在 JSON 信封中；失败时仍输出标准错误信封。
// noRequestLog 为 true 时本次请求不向 stderr 输出请求元数据。
func CallAPIRawOutput(command, method, endpoint string, body any, query map[string]string, absolutePath bool, w io.Writer, noRequestLog bool) error {
	rt, err := GetRuntime()
	if err != nil {
		return err
	}
	var c *client.Client
	if noRequestLog {
		c, err = rt.ensureAPIClientWithLogger(nil)
	} else {
		c, err = rt.EnsureAPIClient()
	}
	if err != nil {
		return EmitLocalError(command, "CLIENT_CONFIG_ERROR", err.Error(), "请先配置 base_url 和 api key")
	}

	res, err := c.Do(context.Background(), method, endpoint, body, query, absolutePath)
	if err != nil {
		return EmitLocalError(command, "NETWORK_ERROR", err.Error(), "请检查网络连通性与重试配置")
	}

	success := true
	if v, ok := res.Payload["success"].(bool); ok {
		success = v
	} else if res.Meta.Status >= 400 {
		success = false
	}

	if !success {
		code := fmt.Sprintf("%d", res.Meta.Status)
		if v, ok := res.Payload["code"]; ok {
			code = fmt.Sprintf("%v", v)
		}
		msg := "request failed"
		if v, ok := res.Payload["msg"].(string); ok && strings.TrimSpace(v) != "" {
			msg = v
		}
		return EmitLocalError(command, code, msg, "")
	}

	data := res.Payload["data"]
	raw := extractRawLogContent(data)
	if len(raw) == 0 {
		return nil
	}
	_, _ = w.Write(raw)
	if raw[len(raw)-1] != '\n' {
		_, _ = w.Write([]byte{'\n'})
	}
	return nil
}

// extractRawLogContent 从 data 中提取可原样输出的日志文本：data 为 string 则返回其字节；为 map 时尝试 content/log/text 等字段；否则返回 data 的 JSON 单行。
func extractRawLogContent(data any) []byte {
	if data == nil {
		return nil
	}
	if s, ok := data.(string); ok {
		return []byte(s)
	}
	if m, ok := data.(map[string]any); ok {
		for _, key := range []string{"content", "log", "text", "logs"} {
			if v, ok := m[key]; ok && v != nil {
				if s, ok := v.(string); ok {
					return []byte(s)
				}
			}
		}
	}
	b, _ := json.Marshal(data)
	return append(b, '\n')
}

func ParseTimeout(raw string) (time.Duration, error) {
	d, err := time.ParseDuration(strings.TrimSpace(raw))
	if err != nil {
		return 0, fmt.Errorf("invalid timeout %q", raw)
	}
	return d, nil
}

func ParseInt(raw string, fallback int) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}
