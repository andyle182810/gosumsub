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

var testImageData = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} //nolint:gochecknoglobals

func TestGetDocumentImage_Success(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"image/png"}},
			Body:       io.NopCloser(strings.NewReader(string(testImageData))),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetDocumentImage(t.Context(), "test-inspection-id", "test-image-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.MimeType != "image/png" {
		t.Errorf("expected MimeType 'image/png', got %q", resp.MimeType)
	}

	if len(resp.Data) != len(testImageData) {
		t.Errorf("expected Data length %d, got %d", len(testImageData), len(resp.Data))
	}
}

func TestGetDocumentImage_SuccessWithJPEG(t *testing.T) {
	t.Parallel()

	jpegData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG magic bytes

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"image/jpeg"}},
			Body:       io.NopCloser(strings.NewReader(string(jpegData))),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetDocumentImage(t.Context(), "test-inspection-id", "test-image-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.MimeType != "image/jpeg" {
		t.Errorf("expected MimeType 'image/jpeg', got %q", resp.MimeType)
	}
}

func TestGetDocumentImage_DetectMimeTypeWhenHeaderMissing(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(string(testImageData))),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetDocumentImage(t.Context(), "test-inspection-id", "test-image-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.MimeType == "" {
		t.Error("expected MimeType to be detected, got empty string")
	}
}

func TestGetDocumentImage_InspectionIDRequired(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetDocumentImage(t.Context(), "", "test-image-id")
	if err == nil {
		t.Fatal("expected error for missing inspection ID, got nil")
	}

	if !errors.Is(err, gosumsub.ErrInspectionIDRequired) {
		t.Errorf("expected ErrInspectionIDRequired, got %v", err)
	}
}

func TestGetDocumentImage_ImageIDRequired(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetDocumentImage(t.Context(), "test-inspection-id", "")
	if err == nil {
		t.Fatal("expected error for missing image ID, got nil")
	}

	if !errors.Is(err, gosumsub.ErrImageIDRequired) {
		t.Errorf("expected ErrImageIDRequired, got %v", err)
	}
}

func TestGetDocumentImage_BothIDsRequired(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetDocumentImage(t.Context(), "", "")
	if err == nil {
		t.Fatal("expected error for missing IDs, got nil")
	}

	if !errors.Is(err, gosumsub.ErrInspectionIDRequired) {
		t.Errorf("expected ErrInspectionIDRequired (checked first), got %v", err)
	}
}

func TestGetDocumentImage_HTTPError(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: nil,
		err:      context.DeadlineExceeded,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetDocumentImage(t.Context(), "test-inspection-id", "test-image-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetDocumentImage_NotFound(t *testing.T) {
	t.Parallel()

	body := `{"description":"Resource not found","code":404,` +
		`"correlationId":"abc123","errorCode":1001,"errorName":"NOT_FOUND"}`
	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetDocumentImage(t.Context(), "non-existent-inspection", "non-existent-image")
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestGetDocumentImage_Unauthorized(t *testing.T) {
	t.Parallel()

	body := `{"description":"Unauthorized","code":401,` +
		`"correlationId":"abc123","errorCode":1002,"errorName":"UNAUTHORIZED"}`
	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetDocumentImage(t.Context(), "test-inspection-id", "test-image-id")
	if err == nil {
		t.Fatal("expected error for unauthorized, got nil")
	}
}

func TestGetDocumentImage_GetBase64WithMime(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"image/png"}},
			Body:       io.NopCloser(strings.NewReader(string(testImageData))),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetDocumentImage(t.Context(), "test-inspection-id", "test-image-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	base64Result := resp.GetBase64WithMime()
	if base64Result == "" {
		t.Error("expected non-empty base64 result")
	}

	expectedPrefix := "data:image/png;base64,"
	if !strings.HasPrefix(base64Result, expectedPrefix) {
		t.Errorf("expected base64 result to start with %q, got %q", expectedPrefix, base64Result)
	}
}

func TestGetDocumentImage_GetBase64WithMime_NilData(t *testing.T) {
	t.Parallel()

	resp := &gosumsub.GetDocumentImageResponse{
		Data:     nil,
		MimeType: "image/png",
	}

	result := resp.GetBase64WithMime()
	if result != "" {
		t.Errorf("expected empty string for nil data, got %q", result)
	}
}

func TestGetDocumentImage_GetBase64WithMime_EmptyData(t *testing.T) {
	t.Parallel()

	resp := &gosumsub.GetDocumentImageResponse{
		Data:     []byte{},
		MimeType: "image/png",
	}

	result := resp.GetBase64WithMime()
	expectedPrefix := "data:image/png;base64,"

	if !strings.HasPrefix(result, expectedPrefix) {
		t.Errorf("expected result to start with %q, got %q", expectedPrefix, result)
	}
}

func TestGetDocumentImage_SpecialCharactersInIDs(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"image/png"}},
			Body:       io.NopCloser(strings.NewReader(string(testImageData))),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetDocumentImage(t.Context(), "inspection/with/slashes", "image?with=query")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}
