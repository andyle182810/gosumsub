package gosumsub_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestIntegration_GetDocumentImage(t *testing.T) {
	t.Parallel()

	inspectionID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_INSPECTION_ID"))
	if inspectionID == "" {
		t.Skip("skipping test: SUMSUB_TEST_INSPECTION_ID not set")
	}

	imageID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_IMAGE_ID"))
	if imageID == "" {
		t.Skip("skipping test: SUMSUB_TEST_IMAGE_ID not set")
	}

	client := newTestClient(t)

	resp, err := client.GetDocumentImage(t.Context(), inspectionID, imageID)
	if err != nil {
		t.Fatalf("GetDocumentImage failed: %v", err)
	}

	t.Logf("MimeType: %s", resp.MimeType)
	t.Logf("Data length: %d bytes", len(resp.Data))

	if resp.MimeType == "" {
		t.Error("expected non-empty MimeType")
	}

	if len(resp.Data) == 0 {
		t.Error("expected non-empty Data")
	}
}

func TestIntegration_GetDocumentImage_VerifyBase64(t *testing.T) {
	t.Parallel()

	inspectionID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_INSPECTION_ID"))
	if inspectionID == "" {
		t.Skip("skipping test: SUMSUB_TEST_INSPECTION_ID not set")
	}

	imageID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_IMAGE_ID"))
	if imageID == "" {
		t.Skip("skipping test: SUMSUB_TEST_IMAGE_ID not set")
	}

	client := newTestClient(t)

	resp, err := client.GetDocumentImage(t.Context(), inspectionID, imageID)
	if err != nil {
		t.Fatalf("GetDocumentImage failed: %v", err)
	}

	base64Result := resp.GetBase64WithMime()
	if base64Result == "" {
		t.Error("expected non-empty base64 result")
	}

	if !strings.HasPrefix(base64Result, "data:") {
		t.Errorf("expected base64 result to start with 'data:', got %q", base64Result[:min(50, len(base64Result))])
	}

	if !strings.Contains(base64Result, ";base64,") {
		t.Error("expected base64 result to contain ';base64,'")
	}

	t.Logf("Base64 result length: %d characters", len(base64Result))
	t.Logf("Base64 prefix: %s...", base64Result[:min(50, len(base64Result))])
}

func TestIntegration_GetDocumentImage_WithCancelledContext(t *testing.T) {
	t.Parallel()

	inspectionID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_INSPECTION_ID"))
	if inspectionID == "" {
		t.Skip("skipping test: SUMSUB_TEST_INSPECTION_ID not set")
	}

	imageID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_IMAGE_ID"))
	if imageID == "" {
		t.Skip("skipping test: SUMSUB_TEST_IMAGE_ID not set")
	}

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := client.GetDocumentImage(ctx, inspectionID, imageID)
	if err == nil {
		t.Fatal("expected error with cancelled context, got nil")
	}

	t.Logf("Expected error received: %v", err)
}

func TestIntegration_GetDocumentImage_WithTimeout(t *testing.T) {
	t.Parallel()

	inspectionID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_INSPECTION_ID"))
	if inspectionID == "" {
		t.Skip("skipping test: SUMSUB_TEST_INSPECTION_ID not set")
	}

	imageID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_IMAGE_ID"))
	if imageID == "" {
		t.Skip("skipping test: SUMSUB_TEST_IMAGE_ID not set")
	}

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	resp, err := client.GetDocumentImage(ctx, inspectionID, imageID)
	if err != nil {
		t.Fatalf("GetDocumentImage failed: %v", err)
	}

	t.Logf("Successfully retrieved document image with timeout context: %d bytes", len(resp.Data))
}

func TestIntegration_GetDocumentImage_NonExistent(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	_, err := client.GetDocumentImage(t.Context(), "non-existent-inspection-id", "non-existent-image-id")
	if err == nil {
		t.Fatal("expected error for non-existent IDs, got nil")
	}

	t.Logf("Expected error received: %v", err)
}

func TestIntegration_GetDocumentImage_InvalidInspectionID(t *testing.T) {
	t.Parallel()

	imageID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_IMAGE_ID"))
	if imageID == "" {
		t.Skip("skipping test: SUMSUB_TEST_IMAGE_ID not set")
	}

	client := newTestClient(t)

	_, err := client.GetDocumentImage(t.Context(), "invalid-inspection-id", imageID)
	if err == nil {
		t.Fatal("expected error for invalid inspection ID, got nil")
	}

	t.Logf("Expected error received: %v", err)
}

func TestIntegration_GetDocumentImage_InvalidImageID(t *testing.T) {
	t.Parallel()

	inspectionID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_INSPECTION_ID"))
	if inspectionID == "" {
		t.Skip("skipping test: SUMSUB_TEST_INSPECTION_ID not set")
	}

	client := newTestClient(t)

	_, err := client.GetDocumentImage(t.Context(), inspectionID, "invalid-image-id")
	if err == nil {
		t.Fatal("expected error for invalid image ID, got nil")
	}

	t.Logf("Expected error received: %v", err)
}
