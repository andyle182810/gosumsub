package gosumsub_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/andyle182810/gosumsub"
)

const testSDKLinkResponseBody = `{"url":"https://in.sumsub.com/websdk/abc123"}`

func TestGenerateExternalWebSDKLink_Success(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testSDKLinkResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "test-user-123",
		LevelName:            "basic-kyc-level",
		ApplicantIdentifiers: nil,
		Redirect:             nil,
	}

	resp, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.URL != "https://in.sumsub.com/websdk/abc123" {
		t.Errorf("expected URL 'https://in.sumsub.com/websdk/abc123', got %q", resp.URL)
	}
}

func TestGenerateExternalWebSDKLink_WithAllFields(t *testing.T) {
	t.Parallel()

	responseBody := `{"url":"https://in.sumsub.com/websdk/xyz789"}`
	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(responseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs: 3600,
		UserID:    "user-456",
		LevelName: "full-kyc-level",
		ApplicantIdentifiers: &gosumsub.ApplicantIdentifiers{
			Email: "test@example.com",
			Phone: "+1234567890",
		},
		Redirect: &gosumsub.Redirect{
			AllowedQueryParams: []string{"param1", "param2"},
			SuccessURL:         "https://example.com/success",
			RejectURL:          "https://example.com/reject",
			SignKey:            "sign-key-123",
		},
	}

	resp, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.URL != "https://in.sumsub.com/websdk/xyz789" {
		t.Errorf("expected URL 'https://in.sumsub.com/websdk/xyz789', got %q", resp.URL)
	}
}

func TestGenerateExternalWebSDKLink_LevelNameRequired(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{}")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "test-user-123",
		LevelName:            "",
		ApplicantIdentifiers: nil,
		Redirect:             nil,
	}

	_, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err == nil {
		t.Fatal("expected error for missing level name, got nil")
	}

	if !errors.Is(err, gosumsub.ErrLevelNameRequired) {
		t.Errorf("expected ErrLevelNameRequired, got %v", err)
	}
}

func TestGenerateExternalWebSDKLink_TooManyQueryParams(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{}")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "",
		LevelName:            "basic-kyc-level",
		ApplicantIdentifiers: nil,
		Redirect: &gosumsub.Redirect{
			AllowedQueryParams: []string{"p1", "p2", "p3", "p4", "p5"},
			SuccessURL:         "",
			RejectURL:          "",
			SignKey:            "",
		},
	}

	_, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err == nil {
		t.Fatal("expected error for too many query params, got nil")
	}

	if !errors.Is(err, gosumsub.ErrTooManyQueryParams) {
		t.Errorf("expected ErrTooManyQueryParams, got %v", err)
	}
}

func TestGenerateExternalWebSDKLink_DefaultTTL(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testSDKLinkResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "",
		LevelName:            "basic-kyc-level",
		ApplicantIdentifiers: nil,
		Redirect:             nil,
	}

	_, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// TTL should be set to default 1800 internally
	if req.TTLInSecs != 1800 {
		t.Errorf("expected default TTL of 1800, got %d", req.TTLInSecs)
	}
}

func TestGenerateExternalWebSDKLink_HTTPError(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: nil,
		err:      context.DeadlineExceeded,
	}

	client := newMockClient(t, httpClient)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "",
		LevelName:            "basic-kyc-level",
		ApplicantIdentifiers: nil,
		Redirect:             nil,
	}

	_, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGenerateExternalWebSDKLink_NonOKStatus(t *testing.T) {
	t.Parallel()

	body := `{"description":"Bad Request","code":400,` +
		`"correlationId":"abc123","errorCode":1001,"errorName":"BAD_REQUEST"}`
	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(body)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "",
		LevelName:            "basic-kyc-level",
		ApplicantIdentifiers: nil,
		Redirect:             nil,
	}

	_, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestGenerateExternalWebSDKLink_InvalidJSON(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("invalid json")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "",
		LevelName:            "basic-kyc-level",
		ApplicantIdentifiers: nil,
		Redirect:             nil,
	}

	_, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err == nil {
		t.Fatal("expected error for invalid JSON response, got nil")
	}
}

func TestGenerateExternalWebSDKLink_MaxAllowedQueryParams(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testSDKLinkResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "",
		LevelName:            "basic-kyc-level",
		ApplicantIdentifiers: nil,
		Redirect: &gosumsub.Redirect{
			AllowedQueryParams: []string{"p1", "p2", "p3", "p4"},
			SuccessURL:         "",
			RejectURL:          "",
			SignKey:            "",
		},
	}

	_, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error with 4 query params: %v", err)
	}
}
