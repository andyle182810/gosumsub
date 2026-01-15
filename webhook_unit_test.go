package gosumsub_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1" //nolint:gosec
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andyle182810/gosumsub"
	"github.com/labstack/echo/v4"
)

const (
	testSecretKey     = "secret"
	testEchoSecretKey = "echo_test_secret"
	testApplicantBody = `{"applicantId":"app456"}`
)

func computeHMACSHA256(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))

	return hex.EncodeToString(mac.Sum(nil))
}

func computeHMACSHA512(payload, secret string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(payload))

	return hex.EncodeToString(mac.Sum(nil))
}

func computeHMACSHA1(payload, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(payload))

	return hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyWebhookDigest_ValidHMACSHA1(t *testing.T) {
	t.Parallel()

	payload := []byte("someText")
	testKey := "SoMe_SeCrEt_KeY"
	digestHex := "f6e92ffe371718694d46e28436f76589312df8db"

	err := gosumsub.VerifyWebhookDigest(payload, testKey, gosumsub.AlgoHMACSHA1, digestHex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyWebhookDigest_ValidHMACSHA256(t *testing.T) {
	t.Parallel()

	payload := []byte("test payload")
	secretKey := "secret123"
	digestHex := computeHMACSHA256("test payload", secretKey)

	err := gosumsub.VerifyWebhookDigest(payload, secretKey, gosumsub.AlgoHMACSHA256, digestHex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyWebhookDigest_ValidHMACSHA512(t *testing.T) {
	t.Parallel()

	payload := []byte("another payload")
	secretKey := "another_secret"
	digestHex := computeHMACSHA512("another payload", secretKey)

	err := gosumsub.VerifyWebhookDigest(payload, secretKey, gosumsub.AlgoHMACSHA512, digestHex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyWebhookDigest_EmptyDigest(t *testing.T) {
	t.Parallel()

	err := gosumsub.VerifyWebhookDigest([]byte("test"), testSecretKey, gosumsub.AlgoHMACSHA256, "")
	if err == nil {
		t.Fatal("expected error for empty digest, got nil")
	}

	if !errors.Is(err, gosumsub.ErrEmptyDigest) {
		t.Errorf("expected ErrEmptyDigest, got %v", err)
	}
}

func TestVerifyWebhookDigest_EmptySecretKey(t *testing.T) {
	t.Parallel()

	err := gosumsub.VerifyWebhookDigest([]byte("test"), "", gosumsub.AlgoHMACSHA256, "abc123")
	if err == nil {
		t.Fatal("expected error for empty secret key, got nil")
	}

	if !errors.Is(err, gosumsub.ErrEmptySecretKey) {
		t.Errorf("expected ErrEmptySecretKey, got %v", err)
	}
}

func TestVerifyWebhookDigest_UnsupportedAlgorithm(t *testing.T) {
	t.Parallel()

	err := gosumsub.VerifyWebhookDigest([]byte("test"), testSecretKey, "HMAC_MD5_HEX", "abc123")
	if err == nil {
		t.Fatal("expected error for unsupported algorithm, got nil")
	}
}

func TestVerifyWebhookDigest_MalformedDigestHex(t *testing.T) {
	t.Parallel()

	err := gosumsub.VerifyWebhookDigest([]byte("test"), testSecretKey, gosumsub.AlgoHMACSHA256, "invalid_hex")
	if err == nil {
		t.Fatal("expected error for malformed digest hex, got nil")
	}

	if !errors.Is(err, gosumsub.ErrMalformedDigest) {
		t.Errorf("expected ErrMalformedDigest, got %v", err)
	}
}

func TestVerifyWebhookDigest_DigestMismatch(t *testing.T) {
	t.Parallel()

	wrongDigest := "0000000000000000000000000000000000000000000000000000000000000000"
	err := gosumsub.VerifyWebhookDigest([]byte("test"), testSecretKey, gosumsub.AlgoHMACSHA256, wrongDigest)

	if err == nil {
		t.Fatal("expected error for digest mismatch, got nil")
	}

	if !errors.Is(err, gosumsub.ErrDigestMismatch) {
		t.Errorf("expected ErrDigestMismatch, got %v", err)
	}
}

func TestVerifyWebhookRequest_ValidSHA256(t *testing.T) {
	t.Parallel()

	body := "test body content 1"
	secretKey := "secret_key_1"
	digestHex := computeHMACSHA256(body, secretKey)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/webhook", strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, digestHex)

	err = gosumsub.VerifyWebhookRequest(req, secretKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyWebhookRequest_ValidSHA1(t *testing.T) {
	t.Parallel()

	body := "test body content 2"
	secretKey := "secret_key_2"
	digestHex := computeHMACSHA1(body, secretKey)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/webhook", strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA1)
	req.Header.Set(gosumsub.HeaderDigest, digestHex)

	err = gosumsub.VerifyWebhookRequest(req, secretKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyWebhookRequest_ValidSHA512(t *testing.T) {
	t.Parallel()

	body := "test body content 3"
	secretKey := "secret_key_3"
	digestHex := computeHMACSHA512(body, secretKey)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/webhook", strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA512)
	req.Header.Set(gosumsub.HeaderDigest, digestHex)

	err = gosumsub.VerifyWebhookRequest(req, secretKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyWebhookRequest_EmptyBody(t *testing.T) {
	t.Parallel()

	body := ""
	secretKey := "secret_key_4"
	digestHex := computeHMACSHA256(body, secretKey)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/webhook", strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, digestHex)

	err = gosumsub.VerifyWebhookRequest(req, secretKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyWebhookRequest_MissingDigestHeader(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/webhook", strings.NewReader("test"))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)

	err = gosumsub.VerifyWebhookRequest(req, testSecretKey)
	if err == nil {
		t.Fatal("expected error for missing digest header, got nil")
	}
}

func TestVerifyWebhookRequest_MissingAlgorithmHeader(t *testing.T) {
	t.Parallel()

	body := "test"
	digestHex := computeHMACSHA256(body, testSecretKey)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/webhook", strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set(gosumsub.HeaderDigest, digestHex)

	err = gosumsub.VerifyWebhookRequest(req, testSecretKey)
	if err == nil {
		t.Fatal("expected error for missing algorithm header, got nil")
	}
}

func TestVerifyWebhookRequest_InvalidDigestFormat(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/webhook", strings.NewReader("test"))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, "invalid_hex")

	err = gosumsub.VerifyWebhookRequest(req, testSecretKey)
	if err == nil {
		t.Fatal("expected error for invalid digest format, got nil")
	}
}

func TestVerifyWebhookRequestWithBody_ValidRequest(t *testing.T) {
	t.Parallel()

	body := `{"applicantId":"test123"}`
	secretKey := "my_secret"
	digestHex := computeHMACSHA256(body, secretKey)

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, digestHex)

	returnedBody, err := gosumsub.VerifyWebhookRequestWithBody(req, secretKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(returnedBody) != body {
		t.Errorf("expected body %q, got %q", body, string(returnedBody))
	}
}

func TestVerifyWebhookRequestWithBody_BodyReadableAfterVerification(t *testing.T) {
	t.Parallel()

	body := `{"applicantId":"test123"}`
	secretKey := "my_secret"
	digestHex := computeHMACSHA256(body, secretKey)

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, digestHex)

	_, err := gosumsub.VerifyWebhookRequestWithBody(req, secretKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bodyAgain, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("failed to read body again: %v", err)
	}

	if string(bodyAgain) != body {
		t.Errorf("expected body %q after re-read, got %q", body, string(bodyAgain))
	}
}

func TestVerifyWebhookRequestWithBody_InvalidSignature(t *testing.T) {
	t.Parallel()

	body := `{"test":"data"}`
	wrongDigest := "0000000000000000000000000000000000000000000000000000000000000000"

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, wrongDigest)

	returnedBody, err := gosumsub.VerifyWebhookRequestWithBody(req, testSecretKey)
	if err == nil {
		t.Fatal("expected error for invalid signature, got nil")
	}

	if returnedBody != nil {
		t.Errorf("expected nil body on error, got %q", string(returnedBody))
	}
}

func TestVerifyWebhookRequestWithBody_MissingHeaders(t *testing.T) {
	t.Parallel()

	body := `{"test":"data"}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))

	returnedBody, err := gosumsub.VerifyWebhookRequestWithBody(req, testSecretKey)
	if err == nil {
		t.Fatal("expected error for missing headers, got nil")
	}

	if returnedBody != nil {
		t.Errorf("expected nil body on error, got %q", string(returnedBody))
	}
}

func TestWebhookMiddleware_ValidSignature(t *testing.T) {
	t.Parallel()

	body := `{"applicantId":"app123"}`
	secretKey := "test_secret"
	digest := computeHMACSHA256(body, secretKey)

	called := false
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		called = true

		readBody, _ := io.ReadAll(request.Body)
		if string(readBody) != body {
			t.Errorf("expected body %q, got %q", body, string(readBody))
		}

		writer.WriteHeader(http.StatusOK)
	})

	middleware := gosumsub.WebhookMiddleware(secretKey)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader([]byte(body)))
	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, digest)

	recorder := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(recorder, req)

	if !called {
		t.Error("expected handler to be called")
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestWebhookMiddleware_InvalidSignature(t *testing.T) {
	t.Parallel()

	body := `{"applicantId":"app123"}`
	secretKey := "test_secret"

	called := false
	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		called = true
	})

	middleware := gosumsub.WebhookMiddleware(secretKey)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader([]byte(body)))
	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, "invalid")

	recorder := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(recorder, req)

	if called {
		t.Error("expected handler not to be called")
	}

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestEchoWebhookMiddleware_ValidSignature(t *testing.T) {
	t.Parallel()

	digest := computeHMACSHA256(testApplicantBody, testEchoSecretKey)

	echoInstance := echo.New()
	called := false

	handler := func(ctx echo.Context) error {
		called = true

		readBody, _ := io.ReadAll(ctx.Request().Body)
		if string(readBody) != testApplicantBody {
			t.Errorf("expected body %q, got %q", testApplicantBody, string(readBody))
		}

		return ctx.String(http.StatusOK, "ok")
	}

	middleware := gosumsub.EchoWebhookMiddleware(testEchoSecretKey)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader([]byte(testApplicantBody)))
	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, digest)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoInstance.NewContext(req, rec)

	err := wrappedHandler(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Error("expected handler to be called")
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestEchoWebhookMiddleware_InvalidSignature(t *testing.T) {
	t.Parallel()

	echoInstance := echo.New()
	called := false

	handler := func(ctx echo.Context) error {
		called = true

		return ctx.String(http.StatusOK, "ok")
	}

	middleware := gosumsub.EchoWebhookMiddleware(testEchoSecretKey)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader([]byte(testApplicantBody)))
	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, "invalid_digest")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoInstance.NewContext(req, rec)

	err := wrappedHandler(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if called {
		t.Error("expected handler not to be called")
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestEchoWebhookMiddleware_MissingDigestHeader(t *testing.T) {
	t.Parallel()

	echoInstance := echo.New()
	called := false

	handler := func(ctx echo.Context) error {
		called = true

		return ctx.String(http.StatusOK, "ok")
	}

	middleware := gosumsub.EchoWebhookMiddleware(testEchoSecretKey)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader([]byte(testApplicantBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoInstance.NewContext(req, rec)

	err := wrappedHandler(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if called {
		t.Error("expected handler not to be called")
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestEchoWebhookMiddleware_WrongSecretKey(t *testing.T) {
	t.Parallel()

	digest := computeHMACSHA256(testApplicantBody, testEchoSecretKey)

	echoInstance := echo.New()
	called := false

	handler := func(ctx echo.Context) error {
		called = true

		return ctx.String(http.StatusOK, "ok")
	}

	middleware := gosumsub.EchoWebhookMiddleware("wrong_secret")
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader([]byte(testApplicantBody)))
	req.Header.Set(gosumsub.HeaderDigestAlg, gosumsub.AlgoHMACSHA256)
	req.Header.Set(gosumsub.HeaderDigest, digest)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoInstance.NewContext(req, rec)

	err := wrappedHandler(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if called {
		t.Error("expected handler not to be called")
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}
