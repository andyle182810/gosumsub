package gosumsub

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/andyle182810/gosumsub/signer"
)

const UserAgent = "sumsub-go-sdk"

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Signer interface {
	Sign(timestamp time.Time, method, uri string, payload *[]byte) (string, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClockFunc func() time.Time

var (
	// Client initialization errors.
	ErrEmptyBaseURL = errors.New("base URL cannot be empty")
	ErrEmptyToken   = errors.New("token cannot be empty")
	ErrEmptySecret  = errors.New("secret cannot be empty")

	// Request lifecycle errors.
	ErrNilRequest    = errors.New("request is nil")
	ErrRequestEncode = errors.New("failed to encode request body")
	ErrRequestSign   = errors.New("failed to sign request")

	// HTTP / transport errors.
	ErrHTTPFailure      = errors.New("http request failed")
	ErrUnexpectedStatus = errors.New("unexpected http status code")
)

type Option func(*Client)

type Client struct {
	baseURL    string
	logger     Logger
	debug      bool
	token      string
	signer     Signer
	clock      ClockFunc
	httpClient HTTPClient
}

func NewClient(baseURL, token, secret string, opts ...Option) (*Client, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return nil, ErrEmptyBaseURL
	}

	baseURL = strings.TrimSuffix(baseURL, "/")

	if strings.TrimSpace(token) == "" {
		return nil, ErrEmptyToken
	}

	if strings.TrimSpace(secret) == "" {
		return nil, ErrEmptySecret
	}

	signer, err := signer.NewSigner(secret)
	if err != nil {
		return nil, err
	}

	client := &Client{
		baseURL:    baseURL,
		logger:     slog.Default(),
		debug:      false,
		token:      token,
		signer:     signer,
		clock:      time.Now,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

func WithLogger(logger Logger) Option {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

func WithHTTPClient(httpClient HTTPClient) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

func WithClock(clock ClockFunc) Option {
	return func(c *Client) {
		if clock != nil {
			c.clock = clock
		}
	}
}

func WithSigner(signer Signer) Option {
	return func(c *Client) {
		if signer != nil {
			c.signer = signer
		}
	}
}

func (c *Client) logDebug(msg string, attrs ...any) {
	if c.debug {
		c.logger.Debug(msg, attrs...)
	}
}

func (c *Client) buildRequest(req *request) error {
	if req == nil {
		return ErrNilRequest
	}

	fullURL := c.baseURL + req.Endpoint

	headers := http.Header{}
	if req.Header != nil {
		headers = req.Header.Clone()
	}

	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", UserAgent)
	headers.Set("X-App-Token", c.token)

	if len(req.Query) > 0 {
		fullURL += "?" + req.Query.Encode()
	}

	now := c.clock()
	headers.Set("X-App-Access-Ts", strconv.FormatInt(now.Unix(), 10))

	var bodyBytes *[]byte

	if req.Params != nil {
		encodedBody, err := json.Marshal(req.Params)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrRequestEncode, err)
		}

		bodyBytes = &encodedBody
		req.Body = bytes.NewReader(encodedBody)
	}

	sig, err := c.signer.Sign(now, req.Method, req.Endpoint, bodyBytes)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRequestSign, err)
	}

	headers.Set("X-App-Access-Sig", sig)

	c.logDebug("http request", "url", fullURL)

	if bodyBytes != nil {
		c.logDebug("http request body", "body", string(*bodyBytes))
	}

	req.FullURL = fullURL
	req.Header = headers

	return nil
}

func (c *Client) handleErrorResponse(statusCode int, body []byte) error {
	if len(body) > 0 && json.Valid(body) {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Code != 0 {
			return &apiErr
		}
	}

	return fmt.Errorf("%w: %d", ErrUnexpectedStatus, statusCode)
}

func (c *Client) execute(ctx context.Context, req *request) ([]byte, error) {
	body, _, err := c.executeWithContentType(ctx, req)

	return body, err
}

func (c *Client) executeWithContentType(ctx context.Context, req *request) ([]byte, string, error) {
	if err := c.buildRequest(req); err != nil {
		return nil, "", err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		req.Method,
		req.FullURL,
		req.Body,
	)
	if err != nil {
		return nil, "", err
	}

	httpReq.Header = req.Header

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrHTTPFailure, err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	c.logDebug("http response status", "status", resp.StatusCode)
	c.logDebug("http response body", "body", string(responseBody))

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, "", c.handleErrorResponse(resp.StatusCode, responseBody)
	}

	return responseBody, resp.Header.Get("Content-Type"), nil
}
