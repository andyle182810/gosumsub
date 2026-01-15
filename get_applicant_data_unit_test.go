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

const testApplicantDataResponseBody = `{
	"id": "68a7d46b8a6f58bf219053c6",
	"createdAt": "2025-08-22 02:22:35",
	"key": "CDYEGPPSVCWKAT",
	"clientId": "abc.com",
	"inspectionId": "68a7d46b8a6f58bf219053c6",
	"externalUserId": "28",
	"info": {
		"firstName": "John",
		"firstNameEn": "John",
		"lastName": "Mock-Doe",
		"lastNameEn": "Mock-Doe",
		"dob": "2006-02-23",
		"country": "VNM",
		"idDocs": [
			{
				"idDocType": "ID_CARD",
				"country": "VNM",
				"firstName": "John",
				"lastName": "Mock-Doe",
				"number": "Mock-LMWH2VPLUB"
			}
		]
	},
	"fixedInfo": {
		"gender": "M",
		"nationality": "VNM",
		"addresses": [
			{
				"street": "123,Nguyen Van A",
				"town": "Ho Chi Minh",
				"postCode": "7000"
			}
		]
	},
	"email": "test@example.com",
	"phone": "+84939171027",
	"review": {
		"reviewId": "ujsDa",
		"reviewStatus": "completed",
		"reviewResult": {
			"reviewAnswer": "GREEN"
		}
	},
	"type": "individual"
}`

func TestGetApplicantData_Success(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testApplicantDataResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetApplicantData(t.Context(), "68a7d46b8a6f58bf219053c6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "68a7d46b8a6f58bf219053c6" {
		t.Errorf("expected ID '68a7d46b8a6f58bf219053c6', got %q", resp.ID)
	}

	if resp.ExternalUserID != "28" {
		t.Errorf("expected ExternalUserID '28', got %q", resp.ExternalUserID)
	}

	if resp.Type != "individual" {
		t.Errorf("expected Type 'individual', got %q", resp.Type)
	}
}

func TestGetApplicantData_VerifyInfo(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testApplicantDataResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetApplicantData(t.Context(), "68a7d46b8a6f58bf219053c6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Info == nil {
		t.Fatal("expected Info to be non-nil")
	}

	if resp.Info.FirstName != "John" {
		t.Errorf("expected FirstName 'John', got %q", resp.Info.FirstName)
	}

	if resp.Info.LastName != "Mock-Doe" {
		t.Errorf("expected LastName 'Mock-Doe', got %q", resp.Info.LastName)
	}

	if resp.Info.Country != "VNM" { //nolint:goconst
		t.Errorf("expected Country 'VNM', got %q", resp.Info.Country)
	}

	if len(resp.Info.IDDocs) != 1 {
		t.Fatalf("expected 1 IdDoc, got %d", len(resp.Info.IDDocs))
	}

	if resp.Info.IDDocs[0].IDDocType != "ID_CARD" {
		t.Errorf("expected IdDocType 'ID_CARD', got %q", resp.Info.IDDocs[0].IDDocType)
	}
}

func TestGetApplicantData_VerifyFixedInfo(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testApplicantDataResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetApplicantData(t.Context(), "68a7d46b8a6f58bf219053c6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.FixedInfo == nil {
		t.Fatal("expected FixedInfo to be non-nil")
	}

	if resp.FixedInfo.Gender != "M" {
		t.Errorf("expected Gender 'M', got %q", resp.FixedInfo.Gender)
	}

	if resp.FixedInfo.Nationality != "VNM" {
		t.Errorf("expected Nationality 'VNM', got %q", resp.FixedInfo.Nationality)
	}

	if len(resp.FixedInfo.Addresses) != 1 {
		t.Fatalf("expected 1 Address, got %d", len(resp.FixedInfo.Addresses))
	}

	if resp.FixedInfo.Addresses[0].Town != "Ho Chi Minh" {
		t.Errorf("expected Town 'Ho Chi Minh', got %q", resp.FixedInfo.Addresses[0].Town)
	}
}

func TestGetApplicantData_VerifyReview(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testApplicantDataResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetApplicantData(t.Context(), "68a7d46b8a6f58bf219053c6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Review == nil {
		t.Fatal("expected Review to be non-nil")
	}

	if resp.Review.ReviewID != "ujsDa" {
		t.Errorf("expected ReviewId 'ujsDa', got %q", resp.Review.ReviewID)
	}

	if resp.Review.ReviewStatus != "completed" {
		t.Errorf("expected ReviewStatus 'completed', got %q", resp.Review.ReviewStatus)
	}

	if resp.Review.ReviewResult == nil {
		t.Fatal("expected ReviewResult to be non-nil")
	}

	if resp.Review.ReviewResult.ReviewAnswer != "GREEN" {
		t.Errorf("expected ReviewAnswer 'GREEN', got %q", resp.Review.ReviewResult.ReviewAnswer)
	}
}

func TestGetApplicantData_ApplicantIDRequired(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{}")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetApplicantData(t.Context(), "")
	if err == nil {
		t.Fatal("expected error for missing applicant ID, got nil")
	}

	if !errors.Is(err, gosumsub.ErrApplicantIDRequired) {
		t.Errorf("expected ErrApplicantIDRequired, got %v", err)
	}
}

func TestGetApplicantData_HTTPError(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: nil,
		err:      context.DeadlineExceeded,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetApplicantData(t.Context(), "test-applicant-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetApplicantData_NotFound(t *testing.T) {
	t.Parallel()

	body := `{"description":"Applicant not found","code":404,` +
		`"correlationId":"abc123","errorCode":1001,"errorName":"NOT_FOUND"}`
	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(body)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetApplicantData(t.Context(), "non-existent-id")
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestGetApplicantData_NonOKStatus(t *testing.T) {
	t.Parallel()

	body := `{"description":"Unauthorized","code":401,` +
		`"correlationId":"abc123","errorCode":1002,"errorName":"UNAUTHORIZED"}`
	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader(body)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetApplicantData(t.Context(), "test-applicant-id")
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestGetApplicantData_InvalidJSON(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("invalid json")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetApplicantData(t.Context(), "test-applicant-id")
	if err == nil {
		t.Fatal("expected error for invalid JSON response, got nil")
	}
}

func TestGetApplicantData_EmptyResponse(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{}")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetApplicantData(t.Context(), "test-applicant-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "" {
		t.Errorf("expected empty ID, got %q", resp.ID)
	}
}
