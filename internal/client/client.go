package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	clilog "gradmotion-cli/internal/log"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type RequestMeta struct {
	TraceID    string
	RequestID  string
	DurationMS int64
	Endpoint   string
	Method     string
	Status     int
	RetryCount int
}

type Result struct {
	Payload map[string]any
	RawBody []byte
	Meta    RequestMeta
}

type Client struct {
	http      *resty.Client
	baseURL   string
	baseAPI   string
	apiKey    string
	userAgent string
	logger    *clilog.Logger
}

func New(baseURL, apiKey, userAgent string, timeout time.Duration, retry int, logger *clilog.Logger) *Client {
	c := resty.New()
	c.SetTimeout(timeout)
	c.SetRetryCount(retry)
	c.SetRetryWaitTime(1 * time.Second)
	c.SetRetryMaxWaitTime(4 * time.Second)
	c.AddRetryCondition(func(r *resty.Response, err error) bool {
		if err != nil {
			return true
		}
		return r.StatusCode() >= 500
	})

	return &Client{
		http:      c,
		baseURL:   strings.TrimRight(baseURL, "/"),
		baseAPI:   strings.TrimRight(baseURL, "/") + "/api",
		apiKey:    apiKey,
		userAgent: userAgent,
		logger:    logger,
	}
}

func (c *Client) Do(ctx context.Context, method, endpoint string, body any, query map[string]string, absolutePath ...bool) (*Result, error) {
	traceID := uuid.NewString()
	start := time.Now()
	meta := RequestMeta{
		TraceID:  traceID,
		Endpoint: endpoint,
		Method:   strings.ToUpper(method),
	}

	req := c.http.R().
		SetContext(ctx).
		SetHeader("X-Api-Key", c.apiKey).
		SetHeader("X-Trace-Id", traceID).
		SetHeader("User-Agent", c.userAgent).
		SetHeader("Accept", "application/json")
	if body != nil {
		req.SetBody(body)
	}
	if len(query) > 0 {
		req.SetQueryParams(query)
	}

	targetURL := c.baseAPI + endpoint
	if len(absolutePath) > 0 && absolutePath[0] {
		if strings.HasPrefix(endpoint, "/") {
			targetURL = c.baseURL + endpoint
		} else {
			targetURL = c.baseURL + "/" + endpoint
		}
	}

	resp, err := req.Execute(method, targetURL)
	meta.DurationMS = time.Since(start).Milliseconds()
	if resp != nil {
		meta.Status = resp.StatusCode()
		meta.RequestID = firstNonEmpty(resp.Header().Get("X-Request-Id"), resp.Header().Get("Request-Id"))
		meta.RetryCount = resp.Request.Attempt
	}

	if err != nil {
		if c.logger != nil {
			c.logger.Error(map[string]any{
				"trace_id":    meta.TraceID,
				"request_id":  meta.RequestID,
				"endpoint":    endpoint,
				"method":      method,
				"status":      meta.Status,
				"duration_ms": meta.DurationMS,
				"error":       err.Error(),
			})
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}

	raw := resp.Body()
	payload := map[string]any{}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &payload); err != nil {
			return nil, errors.New("response is not valid JSON")
		}
	}

	if c.logger != nil {
		c.logger.Info(map[string]any{
			"trace_id":    meta.TraceID,
			"request_id":  meta.RequestID,
			"endpoint":    endpoint,
			"method":      method,
			"status":      meta.Status,
			"duration_ms": meta.DurationMS,
			"retry_count": meta.RetryCount,
		})
	}

	return &Result{
		Payload: payload,
		RawBody: raw,
		Meta:    meta,
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
