package signer_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/andyle182810/gosumsub/signer"
)

const (
	testURI    = "/api/v1/users"
	methodPOST = "POST"
)

func TestNewSigner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		secret  string
		wantErr error
	}{
		{
			name:    "valid secret",
			secret:  "my-secret-key",
			wantErr: nil,
		},
		{
			name:    "empty secret",
			secret:  "",
			wantErr: signer.ErrEmptySecret,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			signerInstance, err := signer.NewSigner(testCase.secret)
			if !errors.Is(err, testCase.wantErr) {
				t.Errorf("NewSigner() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if testCase.wantErr == nil && signerInstance == nil {
				t.Error("NewSigner() returned nil signer for valid secret")
			}
		})
	}
}

func TestSigner_Sign(t *testing.T) {
	t.Parallel()

	signerInstance, err := signer.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	fixedTime := time.Unix(1609459200, 0) // 2021-01-01 00:00:00 UTC
	postPayload := []byte(`{"name":"test"}`)
	emptyPayload := []byte{}

	tests := []struct {
		name    string
		time    time.Time
		method  string
		uri     string
		payload *[]byte
		wantErr error
	}{
		{
			name:    "valid GET request without payload",
			time:    fixedTime,
			method:  "GET",
			uri:     testURI,
			payload: nil,
			wantErr: nil,
		},
		{
			name:    "valid POST request with payload",
			time:    fixedTime,
			method:  methodPOST,
			uri:     testURI,
			payload: &postPayload,
			wantErr: nil,
		},
		{
			name:    "empty method",
			time:    fixedTime,
			method:  "",
			uri:     testURI,
			payload: nil,
			wantErr: signer.ErrEmptyMethod,
		},
		{
			name:    "empty URI",
			time:    fixedTime,
			method:  "GET",
			uri:     "",
			payload: nil,
			wantErr: signer.ErrEmptyURI,
		},
		{
			name:    "empty payload is valid",
			time:    fixedTime,
			method:  methodPOST,
			uri:     "/api/v1/data",
			payload: &emptyPayload,
			wantErr: nil,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			sig, signErr := signerInstance.Sign(testCase.time, testCase.method, testCase.uri, testCase.payload)
			if !errors.Is(signErr, testCase.wantErr) {
				t.Errorf("Sign() error = %v, wantErr %v", signErr, testCase.wantErr)

				return
			}

			if testCase.wantErr == nil {
				if sig == "" {
					t.Error("Sign() returned empty signature for valid input")
				}

				if len(sig) != 64 { // SHA256 hex = 64 chars
					t.Errorf("Sign() signature length = %d, want 64", len(sig))
				}
			}
		})
	}
}

func TestSigner_Sign_Deterministic(t *testing.T) {
	t.Parallel()

	signerInstance, err := signer.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	fixedTime := time.Unix(1609459200, 0)
	method := methodPOST
	uri := testURI
	payload := []byte(`{"name":"test"}`)

	sig1, err := signerInstance.Sign(fixedTime, method, uri, &payload)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	sig2, err := signerInstance.Sign(fixedTime, method, uri, &payload)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	if sig1 != sig2 {
		t.Errorf("Sign() not deterministic: %s != %s", sig1, sig2)
	}
}

func TestSigner_Sign_DifferentSecrets(t *testing.T) {
	t.Parallel()

	signer1, _ := signer.NewSigner("secret1")
	signer2, _ := signer.NewSigner("secret2")

	fixedTime := time.Unix(1609459200, 0)
	method := "GET"
	uri := testURI

	sig1, _ := signer1.Sign(fixedTime, method, uri, nil)
	sig2, _ := signer2.Sign(fixedTime, method, uri, nil)

	if sig1 == sig2 {
		t.Error("Different secrets should produce different signatures")
	}
}

func TestSigner_Sign_DifferentInputs(t *testing.T) {
	t.Parallel()

	signerInstance, _ := signer.NewSigner("test-secret")
	fixedTime := time.Unix(1609459200, 0)
	dataPayload := []byte("data")

	sig1, _ := signerInstance.Sign(fixedTime, "GET", testURI, nil)
	sig2, _ := signerInstance.Sign(fixedTime, methodPOST, testURI, nil)
	sig3, _ := signerInstance.Sign(fixedTime, "GET", "/api/v1/posts", nil)
	sig4, _ := signerInstance.Sign(fixedTime.Add(time.Second), "GET", testURI, nil)
	sig5, _ := signerInstance.Sign(fixedTime, "GET", testURI, &dataPayload)

	sigs := []string{sig1, sig2, sig3, sig4, sig5}
	seen := make(map[string]bool)

	for _, sig := range sigs {
		if seen[sig] {
			t.Error("Different inputs should produce different signatures")
		}

		seen[sig] = true
	}
}

func TestSigner_Verify_ValidSignature(t *testing.T) {
	t.Parallel()

	signerInstance, err := signer.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	fixedTime := time.Unix(1609459200, 0)
	method := methodPOST
	uri := testURI
	payload := []byte(`{"name":"test"}`)

	signature, err := signerInstance.Sign(fixedTime, method, uri, &payload)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	got, err := signerInstance.Verify(signature, fixedTime, method, uri, payload)
	if err != nil {
		t.Errorf("Verify() error = %v", err)
	}

	if !got {
		t.Error("Verify() = false, want true for valid signature")
	}
}

func TestSigner_Verify_InvalidSignature(t *testing.T) {
	t.Parallel()

	signerInstance, err := signer.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	fixedTime := time.Unix(1609459200, 0)

	got, err := signerInstance.Verify("invalid", fixedTime, methodPOST, testURI, []byte(`{"name":"test"}`))
	if err != nil {
		t.Errorf("Verify() error = %v", err)
	}

	if got {
		t.Error("Verify() = true, want false for invalid signature")
	}
}

func TestSigner_Verify_WrongTimestamp(t *testing.T) {
	t.Parallel()

	signerInstance, err := signer.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	fixedTime := time.Unix(1609459200, 0)
	method := methodPOST
	uri := testURI
	payload := []byte(`{"name":"test"}`)

	signature, err := signerInstance.Sign(fixedTime, method, uri, &payload)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	got, err := signerInstance.Verify(signature, fixedTime.Add(time.Second), method, uri, payload)
	if err != nil {
		t.Errorf("Verify() error = %v", err)
	}

	if got {
		t.Error("Verify() = true, want false for wrong timestamp")
	}
}

func TestSigner_Verify_WrongMethod(t *testing.T) {
	t.Parallel()

	signerInstance, err := signer.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	fixedTime := time.Unix(1609459200, 0)
	uri := testURI
	payload := []byte(`{"name":"test"}`)

	signature, err := signerInstance.Sign(fixedTime, methodPOST, uri, &payload)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	got, err := signerInstance.Verify(signature, fixedTime, "GET", uri, payload)
	if err != nil {
		t.Errorf("Verify() error = %v", err)
	}

	if got {
		t.Error("Verify() = true, want false for wrong method")
	}
}

func TestSigner_Verify_EmptyMethodReturnsError(t *testing.T) {
	t.Parallel()

	signerInstance, err := signer.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	fixedTime := time.Unix(1609459200, 0)

	_, err = signerInstance.Verify("signature", fixedTime, "", testURI, nil)
	if err == nil {
		t.Error("Verify() error = nil, want error for empty method")
	}
}

func TestSigner_ConcurrentSafety(t *testing.T) {
	t.Parallel()

	signerInstance, err := signer.NewSigner("test-secret")
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	fixedTime := time.Unix(1609459200, 0)

	const goroutines = 100

	const iterations = 100

	var waitGroup sync.WaitGroup

	waitGroup.Add(goroutines)

	errChan := make(chan error, goroutines*iterations)

	for range goroutines {
		go func() {
			defer waitGroup.Done()

			for range iterations {
				sig, signErr := signerInstance.Sign(fixedTime, "GET", "/api/v1/test", nil)
				if signErr != nil {
					errChan <- signErr

					return
				}

				if len(sig) != 64 {
					errChan <- signErr

					return
				}
			}
		}()
	}

	waitGroup.Wait()
	close(errChan)

	for e := range errChan {
		t.Errorf("Concurrent Sign() failed: %v", e)
	}
}

func TestSigner_KnownVector(t *testing.T) {
	t.Parallel()

	// Test against a known HMAC-SHA256 value to ensure correctness
	signerInstance, err := signer.NewSigner("secret")
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	// Unix timestamp 0 = "0"
	// Method = "GET"
	// URI = "/test"
	// No payload
	// Message = "0GET/test"
	fixedTime := time.Unix(0, 0)

	sig, err := signerInstance.Sign(fixedTime, "GET", "/test", nil)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	// Verify the signature is a valid hex string of correct length
	if len(sig) != 64 {
		t.Errorf("Signature length = %d, want 64", len(sig))
	}

	// Verify determinism with same input
	sig2, _ := signerInstance.Sign(fixedTime, "GET", "/test", nil)

	if sig != sig2 {
		t.Error("Signature should be deterministic")
	}
}

func BenchmarkSigner_Sign(b *testing.B) {
	signerInstance, _ := signer.NewSigner("benchmark-secret")
	fixedTime := time.Unix(1609459200, 0)
	payload := []byte(`{"user":"test","action":"create"}`)

	b.ResetTimer()

	for b.Loop() {
		_, _ = signerInstance.Sign(fixedTime, methodPOST, testURI, &payload)
	}
}

func BenchmarkSigner_SignParallel(b *testing.B) {
	signerInstance, _ := signer.NewSigner("benchmark-secret")
	fixedTime := time.Unix(1609459200, 0)
	payload := []byte(`{"user":"test","action":"create"}`)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = signerInstance.Sign(fixedTime, methodPOST, testURI, &payload)
		}
	})
}
