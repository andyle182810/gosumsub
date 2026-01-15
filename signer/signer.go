package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"strconv"
	"sync"
	"time"
)

const base10 = 10

var (
	ErrEmptySecret = errors.New("secret cannot be empty")
	ErrEmptyMethod = errors.New("HTTP method cannot be empty")
	ErrEmptyURI    = errors.New("URI cannot be empty")
)

type Signer struct {
	mu     sync.Mutex
	hash   hash.Hash
	secret []byte
}

func NewSigner(secret string) (*Signer, error) {
	if secret == "" {
		return nil, ErrEmptySecret
	}

	secretBytes := []byte(secret)

	return &Signer{
		mu:     sync.Mutex{},
		hash:   hmac.New(sha256.New, secretBytes),
		secret: secretBytes,
	}, nil
}

func (s *Signer) Sign(timestamp time.Time, method, uri string, payload *[]byte) (string, error) {
	if method == "" {
		return "", ErrEmptyMethod
	}

	if uri == "" {
		return "", ErrEmptyURI
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.hash.Reset()

	s.hash.Write(strconv.AppendInt(nil, timestamp.Unix(), base10))
	s.hash.Write([]byte(method))
	s.hash.Write([]byte(uri))

	if payload != nil && len(*payload) > 0 {
		s.hash.Write(*payload)
	}

	return hex.EncodeToString(s.hash.Sum(nil)), nil
}

func (s *Signer) Verify(signature string, timestamp time.Time, method, uri string, payload []byte) (bool, error) {
	expected, err := s.Sign(timestamp, method, uri, &payload)
	if err != nil {
		return false, err
	}

	return hmac.Equal([]byte(expected), []byte(signature)), nil
}
