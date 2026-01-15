package gosumsub_test

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/andyle182810/gosumsub"
)

func newTestClient(t *testing.T) *gosumsub.Client {
	t.Helper()

	token := strings.TrimSpace(os.Getenv("SUMSUB_APP_TOKEN"))
	if token == "" {
		t.Skip("skipping test: SUMSUB_APP_TOKEN not set")
	}

	secret := strings.TrimSpace(os.Getenv("SUMSUB_API_SECRET"))
	if secret == "" {
		t.Skip("skipping test: SUMSUB_API_SECRET not set")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))

	client, err := gosumsub.NewClient(
		"https://api.sumsub.com",
		token,
		secret,
		gosumsub.WithDebug(true),
		gosumsub.WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return client
}

type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return m.response, m.err
}

type mockSigner struct {
	signature string
	err       error
}

func (m *mockSigner) Sign(_ time.Time, _, _ string, _ *[]byte) (string, error) {
	return m.signature, m.err
}

func newMockClient(t *testing.T, httpClient gosumsub.HTTPClient) *gosumsub.Client {
	t.Helper()

	client, err := gosumsub.NewClient(
		"https://api.example.com",
		"test-token",
		"test-secret",
		gosumsub.WithHTTPClient(httpClient),
		gosumsub.WithSigner(&mockSigner{signature: "test-signature"}),
		gosumsub.WithClock(func() time.Time { return time.Unix(1234567890, 0) }),
	)
	if err != nil {
		t.Fatalf("failed to create mock client: %v", err)
	}

	return client
}
