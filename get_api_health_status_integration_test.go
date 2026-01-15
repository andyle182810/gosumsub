package gosumsub_test

import (
	"context"
	"testing"
	"time"
)

func TestIntegration_GetAPIHealthStatus(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	err := client.GetAPIHealthStatus(t.Context())
	if err != nil {
		t.Fatalf("GetAPIHealthStatus failed: %v", err)
	}

	t.Log("API is healthy")
}

func TestIntegration_GetAPIHealthStatus_WithCancelledContext(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err := client.GetAPIHealthStatus(ctx)
	if err == nil {
		t.Fatal("expected error with cancelled context, got nil")
	}

	t.Logf("Expected error received: %v", err)
}

func TestIntegration_GetAPIHealthStatus_WithTimeout(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	err := client.GetAPIHealthStatus(ctx)
	if err != nil {
		t.Fatalf("GetAPIHealthStatus failed: %v", err)
	}

	t.Log("API is healthy (with timeout)")
}
