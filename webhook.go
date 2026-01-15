package gosumsub

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1" //nolint:gosec
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	HeaderDigestAlg = "X-Payload-Digest-Alg"
	HeaderDigest    = "X-Payload-Digest"
)

const (
	AlgoHMACSHA1   = "HMAC_SHA1_HEX"
	AlgoHMACSHA256 = "HMAC_SHA256_HEX"
	AlgoHMACSHA512 = "HMAC_SHA512_HEX"
)

var (
	ErrEmptyDigest       = errors.New("empty digest")
	ErrEmptySecretKey    = errors.New("empty secret key")
	ErrMalformedDigest   = errors.New("malformed digest")
	ErrDigestMismatch    = errors.New("digest mismatch")
	ErrUnsupportedAlgo   = errors.New("unsupported algorithm")
	ErrEmptyPayload      = errors.New("empty payload")
	ErrMissingDigestAlgo = errors.New("missing digest algorithm header")
	ErrMissingDigest     = errors.New("missing digest header")
	ErrHMACWrite         = errors.New("failed to write payload to hmac")
)

const (
	WebhookTypeApplicantReviewed            = "applicantReviewed"
	WebhookTypeApplicantPending             = "applicantPending"
	WebhookTypeApplicantCreated             = "applicantCreated"
	WebhookTypeApplicantOnHold              = "applicantOnHold"
	WebhookTypeApplicantPersonalInfoChanged = "applicantPersonalInfoChanged"
	WebhookTypeApplicantPrechecked          = "applicantPrechecked"
	WebhookTypeApplicantDeleted             = "applicantDeleted"
	WebhookTypeApplicantLevelChanged        = "applicantLevelChanged"
	WebhookTypeApplicantReset               = "applicantReset"
	WebhookTypeApplicantActionPending       = "applicantActionPending"
	WebhookTypeApplicantActionReviewed      = "applicantActionReviewed"
	WebhookTypeApplicantActionOnHold        = "applicantActionOnHold"
	WebhookTypeApplicantWorkflowCompleted   = "applicantWorkflowCompleted"
	WebhookTypeVideoIdentStatusChanged      = "videoIdentStatusChanged"
)

const (
	ReviewAnswerGreen = "GREEN"
	ReviewAnswerRed   = "RED"
)

func VerifyWebhookDigest(payload []byte, secretKey, algo, digestHex string) error {
	if digestHex == "" {
		return ErrEmptyDigest
	}

	if secretKey == "" {
		return ErrEmptySecretKey
	}

	var hashFunc func() hash.Hash

	switch algo {
	case AlgoHMACSHA256:
		hashFunc = sha256.New
	case AlgoHMACSHA512:
		hashFunc = sha512.New
	case AlgoHMACSHA1:
		hashFunc = sha1.New
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedAlgo, algo)
	}

	mac := hmac.New(hashFunc, []byte(secretKey))
	if _, err := mac.Write(payload); err != nil {
		return fmt.Errorf("%w: %w", ErrHMACWrite, err)
	}

	expected := mac.Sum(nil)

	got, err := hex.DecodeString(digestHex)
	if err != nil {
		return ErrMalformedDigest
	}

	if !hmac.Equal(expected, got) {
		return ErrDigestMismatch
	}

	return nil
}

func VerifyWebhookRequest(request *http.Request, secretKey string) error {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	return VerifyWebhookDigest(
		body,
		secretKey,
		request.Header.Get(HeaderDigestAlg),
		request.Header.Get(HeaderDigest),
	)
}

func VerifyWebhookRequestWithBody(request *http.Request, secretKey string) ([]byte, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	request.Body = io.NopCloser(bytes.NewReader(body))

	err = VerifyWebhookDigest(
		body,
		secretKey,
		request.Header.Get(HeaderDigestAlg),
		request.Header.Get(HeaderDigest),
	)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func WebhookMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, err := VerifyWebhookRequestWithBody(request, secretKey)
			if err != nil {
				http.Error(writer, "unauthorized", http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func EchoWebhookMiddleware(secretKey string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			_, err := VerifyWebhookRequestWithBody(ctx.Request(), secretKey)
			if err != nil {
				return ctx.String(http.StatusUnauthorized, "unauthorized")
			}

			return next(ctx)
		}
	}
}
