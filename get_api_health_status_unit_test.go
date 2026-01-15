package gosumsub_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetAPIHealthStatus_Success(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{}")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	err := client.GetAPIHealthStatus(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAPIHealthStatus_HTTPError(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: nil,
		err:      context.DeadlineExceeded,
	}

	client := newMockClient(t, httpClient)

	err := client.GetAPIHealthStatus(t.Context())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAPIHealthStatus_NonOKStatus(t *testing.T) {
	t.Parallel()

	body := `{"description":"Internal Server Error","code":500,` +
		`"correlationId":"04fbcea5218053b7e06ca1236491a44f","errorCode":1000,"errorName":"SERVER_ERROR"}`
	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(body)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	err := client.GetAPIHealthStatus(t.Context())
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}
