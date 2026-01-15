package gosumsub_test

import (
	"context"
	"testing"
	"time"

	"github.com/andyle182810/gosumsub"
)

func TestIntegration_GenerateExternalWebSDKLink(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "test-user-integration",
		LevelName:            "idv-and-phone-verification",
		ApplicantIdentifiers: nil,
		Redirect:             nil,
	}

	resp, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err != nil {
		t.Fatalf("GenerateExternalWebSDKLink failed: %v", err)
	}

	if resp.URL == "" {
		t.Fatal("expected non-empty URL")
	}

	t.Logf("Generated SDK link: %s", resp.URL)
}

func TestIntegration_GenerateExternalWebSDKLink_WithAllFields(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs: 3600,
		UserID:    "test-user-full",
		LevelName: "idv-and-phone-verification",
		ApplicantIdentifiers: &gosumsub.ApplicantIdentifiers{
			Email: "integration-test@example.com",
			Phone: "+1234567890",
		},
		Redirect: &gosumsub.Redirect{
			AllowedQueryParams: []string{"utm_source", "utm_campaign"},
			SuccessURL:         "https://example.com/success",
			RejectURL:          "https://example.com/reject",
			SignKey:            "",
		},
	}

	resp, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err != nil {
		t.Fatalf("GenerateExternalWebSDKLink failed: %v", err)
	}

	if resp.URL == "" {
		t.Fatal("expected non-empty URL")
	}

	t.Logf("Generated SDK link with all fields: %s", resp.URL)
}

func TestIntegration_GenerateExternalWebSDKLink_WithCancelledContext(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "test-user-cancelled",
		LevelName:            "idv-and-phone-verification",
		ApplicantIdentifiers: nil,
		Redirect:             nil,
	}

	_, err := client.GenerateExternalWebSDKLink(ctx, req)
	if err == nil {
		t.Fatal("expected error with cancelled context, got nil")
	}

	t.Logf("Expected error received: %v", err)
}

func TestIntegration_GenerateExternalWebSDKLink_WithTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            0,
		UserID:               "test-user-timeout",
		LevelName:            "idv-and-phone-verification",
		ApplicantIdentifiers: nil,
		Redirect:             nil,
	}

	resp, err := client.GenerateExternalWebSDKLink(ctx, req)
	if err != nil {
		t.Fatalf("GenerateExternalWebSDKLink failed: %v", err)
	}

	if resp.URL == "" {
		t.Fatal("expected non-empty URL")
	}

	t.Log("Successfully generated SDK link with timeout context")
}

func TestIntegration_GenerateExternalWebSDKLink_CustomTTL(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	req := &gosumsub.GenerateExternalWebSDKLinkRequest{
		TTLInSecs:            600, // 10 minutes
		UserID:               "test-user-custom-ttl",
		LevelName:            "idv-and-phone-verification",
		ApplicantIdentifiers: nil,
		Redirect:             nil,
	}

	resp, err := client.GenerateExternalWebSDKLink(t.Context(), req)
	if err != nil {
		t.Fatalf("GenerateExternalWebSDKLink failed: %v", err)
	}

	if resp.URL == "" {
		t.Fatal("expected non-empty URL")
	}

	t.Logf("Generated SDK link with custom TTL: %s", resp.URL)
}
