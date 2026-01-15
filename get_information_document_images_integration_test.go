package gosumsub_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestIntegration_GetInformationDocumentImages(t *testing.T) {
	t.Parallel()

	applicantID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_APPLICANT_ID"))
	if applicantID == "" {
		t.Skip("skipping test: SUMSUB_TEST_APPLICANT_ID not set")
	}

	client := newTestClient(t)

	resp, err := client.GetInformationDocumentImages(t.Context(), applicantID)
	if err != nil {
		t.Fatalf("GetInformationDocumentImages failed: %v", err)
	}

	t.Logf("Total Items: %d", resp.TotalItems)
	t.Logf("Items count: %d", len(resp.Items))
}

func TestIntegration_GetInformationDocumentImages_VerifyItems(t *testing.T) {
	t.Parallel()

	applicantID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_APPLICANT_ID"))
	if applicantID == "" {
		t.Skip("skipping test: SUMSUB_TEST_APPLICANT_ID not set")
	}

	client := newTestClient(t)

	resp, err := client.GetInformationDocumentImages(t.Context(), applicantID)
	if err != nil {
		t.Fatalf("GetInformationDocumentImages failed: %v", err)
	}

	for i, item := range resp.Items {
		t.Logf("Item %d:", i)
		t.Logf("  ID: %s", item.ID)
		t.Logf("  PreviewID: %s", item.PreviewID)
		t.Logf("  AddedDate: %s", item.AddedDate)
		t.Logf("  Source: %s", item.Source)
		t.Logf("  AttemptID: %s", item.AttemptID)
		t.Logf("  Deactivated: %v", item.Deactivated)

		if item.FileMetadata != nil {
			t.Logf("  FileMetadata:")
			t.Logf("    FileName: %s", item.FileMetadata.FileName)
			t.Logf("    FileType: %s", item.FileMetadata.FileType)
			t.Logf("    FileSize: %d", item.FileMetadata.FileSize)

			if item.FileMetadata.Resolution != nil {
				t.Logf("    Resolution: %dx%d", item.FileMetadata.Resolution.Width, item.FileMetadata.Resolution.Height)
			}
		}

		if item.IDDocDef != nil {
			t.Logf("  IDDocDef:")
			t.Logf("    Country: %s", item.IDDocDef.Country)
			t.Logf("    IDDocType: %s", item.IDDocDef.IDDocType)
			t.Logf("    IDDocSubType: %s", item.IDDocDef.IDDocSubType)
		}

		if item.ReviewResult != nil {
			t.Logf("  ReviewAnswer: %s", item.ReviewResult.ReviewAnswer)
		}
	}
}

func TestIntegration_GetInformationDocumentImages_WithCancelledContext(t *testing.T) {
	t.Parallel()

	applicantID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_APPLICANT_ID"))
	if applicantID == "" {
		t.Skip("skipping test: SUMSUB_TEST_APPLICANT_ID not set")
	}

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := client.GetInformationDocumentImages(ctx, applicantID)
	if err == nil {
		t.Fatal("expected error with cancelled context, got nil")
	}

	t.Logf("Expected error received: %v", err)
}

func TestIntegration_GetInformationDocumentImages_WithTimeout(t *testing.T) {
	t.Parallel()

	applicantID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_APPLICANT_ID"))
	if applicantID == "" {
		t.Skip("skipping test: SUMSUB_TEST_APPLICANT_ID not set")
	}

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	resp, err := client.GetInformationDocumentImages(ctx, applicantID)
	if err != nil {
		t.Fatalf("GetInformationDocumentImages failed: %v", err)
	}

	t.Logf("Successfully retrieved %d document images with timeout context", len(resp.Items))
}

func TestIntegration_GetInformationDocumentImages_NonExistent(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	_, err := client.GetInformationDocumentImages(t.Context(), "non-existent-applicant-id-12345")
	if err == nil {
		t.Fatal("expected error for non-existent applicant, got nil")
	}

	t.Logf("Expected error received: %v", err)
}
