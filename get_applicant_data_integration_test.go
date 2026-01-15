package gosumsub_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestIntegration_GetApplicantData(t *testing.T) {
	t.Parallel()

	applicantID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_APPLICANT_ID"))
	if applicantID == "" {
		t.Skip("skipping test: SUMSUB_TEST_APPLICANT_ID not set")
	}

	client := newTestClient(t)

	resp, err := client.GetApplicantData(t.Context(), applicantID)
	if err != nil {
		t.Fatalf("GetApplicantData failed: %v", err)
	}

	if resp.ID == "" {
		t.Fatal("expected non-empty ID")
	}

	if resp.ID != applicantID {
		t.Errorf("expected ID %q, got %q", applicantID, resp.ID)
	}

	t.Logf("Applicant ID: %s", resp.ID)
	t.Logf("External User ID: %s", resp.ExternalUserID)
	t.Logf("Type: %s", resp.Type)
}

func TestIntegration_GetApplicantData_VerifyInfo(t *testing.T) {
	t.Parallel()

	applicantID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_APPLICANT_ID"))
	if applicantID == "" {
		t.Skip("skipping test: SUMSUB_TEST_APPLICANT_ID not set")
	}

	client := newTestClient(t)

	resp, err := client.GetApplicantData(t.Context(), applicantID)
	if err != nil {
		t.Fatalf("GetApplicantData failed: %v", err)
	}

	if resp.Info != nil {
		t.Logf("First Name: %s", resp.Info.FirstName)
		t.Logf("Last Name: %s", resp.Info.LastName)
		t.Logf("Country: %s", resp.Info.Country)
		t.Logf("DOB: %s", resp.Info.Dob)
		t.Logf("ID Docs count: %d", len(resp.Info.IDDocs))
	}

	if resp.FixedInfo != nil {
		t.Logf("Gender: %s", resp.FixedInfo.Gender)
		t.Logf("Nationality: %s", resp.FixedInfo.Nationality)
		t.Logf("Addresses count: %d", len(resp.FixedInfo.Addresses))
	}
}

func TestIntegration_GetApplicantData_VerifyReview(t *testing.T) {
	t.Parallel()

	applicantID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_APPLICANT_ID"))
	if applicantID == "" {
		t.Skip("skipping test: SUMSUB_TEST_APPLICANT_ID not set")
	}

	client := newTestClient(t)

	resp, err := client.GetApplicantData(t.Context(), applicantID)
	if err != nil {
		t.Fatalf("GetApplicantData failed: %v", err)
	}

	if resp.Review != nil {
		t.Logf("Review ID: %s", resp.Review.ReviewID)
		t.Logf("Review Status: %s", resp.Review.ReviewStatus)
		t.Logf("Level Name: %s", resp.Review.LevelName)

		if resp.Review.ReviewResult != nil {
			t.Logf("Review Answer: %s", resp.Review.ReviewResult.ReviewAnswer)
		}
	}
}

func TestIntegration_GetApplicantData_WithCancelledContext(t *testing.T) {
	t.Parallel()

	applicantID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_APPLICANT_ID"))
	if applicantID == "" {
		t.Skip("skipping test: SUMSUB_TEST_APPLICANT_ID not set")
	}

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := client.GetApplicantData(ctx, applicantID)
	if err == nil {
		t.Fatal("expected error with cancelled context, got nil")
	}

	t.Logf("Expected error received: %v", err)
}

func TestIntegration_GetApplicantData_WithTimeout(t *testing.T) {
	t.Parallel()

	applicantID := strings.TrimSpace(os.Getenv("SUMSUB_TEST_APPLICANT_ID"))
	if applicantID == "" {
		t.Skip("skipping test: SUMSUB_TEST_APPLICANT_ID not set")
	}

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	resp, err := client.GetApplicantData(ctx, applicantID)
	if err != nil {
		t.Fatalf("GetApplicantData failed: %v", err)
	}

	if resp.ID == "" {
		t.Fatal("expected non-empty ID")
	}

	t.Log("Successfully retrieved applicant data with timeout context")
}

func TestIntegration_GetApplicantData_NonExistent(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	_, err := client.GetApplicantData(t.Context(), "non-existent-applicant-id-12345")
	if err == nil {
		t.Fatal("expected error for non-existent applicant, got nil")
	}

	t.Logf("Expected error received: %v", err)
}
