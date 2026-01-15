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

const testDocumentImagesResponseBody = `{
	"items": [
		{
			"id": "449741312",
			"previewId": "1865335706",
			"addedDate": "2025-08-22 02:25:27",
			"fileMetadata": {
				"fileName": "CCCD-Front.jpg",
				"fileType": "jpeg",
				"fileSize": 129069,
				"resolution": {
					"width": 1152,
					"height": 2048
				}
			},
			"idDocDef": {
				"country": "VNM",
				"idDocType": "ID_CARD",
				"idDocSubType": "FRONT_SIDE"
			},
			"reviewResult": {
				"reviewAnswer": "GREEN"
			},
			"deactivated": false,
			"attemptId": "sdqwI",
			"source": "fileupload"
		},
		{
			"id": "1856600123",
			"previewId": "1956012165",
			"addedDate": "2025-08-22 02:25:32",
			"fileMetadata": {
				"fileName": "CCCD-Back.jpg",
				"fileType": "jpeg",
				"fileSize": 138459,
				"resolution": {
					"width": 1152,
					"height": 2048
				}
			},
			"idDocDef": {
				"country": "VNM",
				"idDocType": "ID_CARD",
				"idDocSubType": "BACK_SIDE"
			},
			"reviewResult": {
				"reviewAnswer": "GREEN"
			},
			"deactivated": false,
			"attemptId": "sdqwI",
			"source": "fileupload"
		},
		{
			"id": "516556662",
			"previewId": "2078923157",
			"addedDate": "2025-08-22 02:32:54",
			"fileMetadata": {
				"fileName": "liveness_photo.jpg",
				"fileType": "jpeg",
				"fileSize": 131392,
				"resolution": {
					"width": 1280,
					"height": 720
				}
			},
			"idDocDef": {
				"country": "VNM",
				"idDocType": "SELFIE"
			},
			"reviewResult": {
				"reviewAnswer": "GREEN"
			},
			"deactivated": false,
			"attemptId": "sdqwI",
			"source": "liveness"
		}
	],
	"totalItems": 3
}`

func TestGetInformationDocumentImages_Success(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testDocumentImagesResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetInformationDocumentImages(t.Context(), "68a7d46b8a6f58bf219053c6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.TotalItems != 3 {
		t.Errorf("expected TotalItems 3, got %d", resp.TotalItems)
	}

	if len(resp.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(resp.Items))
	}
}

func TestGetInformationDocumentImages_VerifyFirstItem(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testDocumentImagesResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetInformationDocumentImages(t.Context(), "68a7d46b8a6f58bf219053c6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Items) < 1 {
		t.Fatal("expected at least 1 item")
	}

	item := resp.Items[0]

	if item.ID != "449741312" {
		t.Errorf("expected ID '449741312', got %q", item.ID)
	}

	if item.PreviewID != "1865335706" {
		t.Errorf("expected PreviewID '1865335706', got %q", item.PreviewID)
	}

	if item.Source != "fileupload" {
		t.Errorf("expected Source 'fileupload', got %q", item.Source)
	}

	if item.AttemptID != "sdqwI" {
		t.Errorf("expected AttemptID 'sdqwI', got %q", item.AttemptID)
	}

	if item.Deactivated != false {
		t.Errorf("expected Deactivated false, got %v", item.Deactivated)
	}
}

func TestGetInformationDocumentImages_VerifyFileMetadata(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testDocumentImagesResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetInformationDocumentImages(t.Context(), "68a7d46b8a6f58bf219053c6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Items) < 1 {
		t.Fatal("expected at least 1 item")
	}

	item := resp.Items[0]

	if item.FileMetadata == nil {
		t.Fatal("expected FileMetadata to be non-nil")
	}

	if item.FileMetadata.FileName != "CCCD-Front.jpg" {
		t.Errorf("expected FileName 'CCCD-Front.jpg', got %q", item.FileMetadata.FileName)
	}

	if item.FileMetadata.FileType != "jpeg" {
		t.Errorf("expected FileType 'jpeg', got %q", item.FileMetadata.FileType)
	}

	if item.FileMetadata.FileSize != 129069 {
		t.Errorf("expected FileSize 129069, got %d", item.FileMetadata.FileSize)
	}

	if item.FileMetadata.Resolution == nil {
		t.Fatal("expected Resolution to be non-nil")
	}

	if item.FileMetadata.Resolution.Width != 1152 {
		t.Errorf("expected Width 1152, got %d", item.FileMetadata.Resolution.Width)
	}

	if item.FileMetadata.Resolution.Height != 2048 {
		t.Errorf("expected Height 2048, got %d", item.FileMetadata.Resolution.Height)
	}
}

func TestGetInformationDocumentImages_VerifyIDDocDef(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testDocumentImagesResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetInformationDocumentImages(t.Context(), "68a7d46b8a6f58bf219053c6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Items) < 1 {
		t.Fatal("expected at least 1 item")
	}

	item := resp.Items[0]

	if item.IDDocDef == nil {
		t.Fatal("expected IDDocDef to be non-nil")
	}

	if item.IDDocDef.Country != "VNM" {
		t.Errorf("expected Country 'VNM', got %q", item.IDDocDef.Country)
	}

	if item.IDDocDef.IDDocType != "ID_CARD" {
		t.Errorf("expected IDDocType 'ID_CARD', got %q", item.IDDocDef.IDDocType)
	}

	if item.IDDocDef.IDDocSubType != "FRONT_SIDE" {
		t.Errorf("expected IDDocSubType 'FRONT_SIDE', got %q", item.IDDocDef.IDDocSubType)
	}
}

func TestGetInformationDocumentImages_VerifyReviewResult(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testDocumentImagesResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetInformationDocumentImages(t.Context(), "68a7d46b8a6f58bf219053c6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Items) < 1 {
		t.Fatal("expected at least 1 item")
	}

	item := resp.Items[0]

	if item.ReviewResult == nil {
		t.Fatal("expected ReviewResult to be non-nil")
	}

	if item.ReviewResult.ReviewAnswer != "GREEN" {
		t.Errorf("expected ReviewAnswer 'GREEN', got %q", item.ReviewResult.ReviewAnswer)
	}
}

func TestGetInformationDocumentImages_VerifySelfieItem(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(testDocumentImagesResponseBody)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetInformationDocumentImages(t.Context(), "68a7d46b8a6f58bf219053c6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Items) < 3 {
		t.Fatal("expected at least 3 items")
	}

	selfieItem := resp.Items[2]

	if selfieItem.IDDocDef == nil {
		t.Fatal("expected IDDocDef to be non-nil")
	}

	if selfieItem.IDDocDef.IDDocType != "SELFIE" {
		t.Errorf("expected IDDocType 'SELFIE', got %q", selfieItem.IDDocDef.IDDocType)
	}

	if selfieItem.Source != "liveness" {
		t.Errorf("expected Source 'liveness', got %q", selfieItem.Source)
	}
}

func TestGetInformationDocumentImages_ApplicantIDRequired(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{}")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetInformationDocumentImages(t.Context(), "")
	if err == nil {
		t.Fatal("expected error for missing applicant ID, got nil")
	}

	if !errors.Is(err, gosumsub.ErrApplicantIDRequired) {
		t.Errorf("expected ErrApplicantIDRequired, got %v", err)
	}
}

func TestGetInformationDocumentImages_HTTPError(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: nil,
		err:      context.DeadlineExceeded,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetInformationDocumentImages(t.Context(), "test-applicant-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetInformationDocumentImages_NotFound(t *testing.T) {
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

	_, err := client.GetInformationDocumentImages(t.Context(), "non-existent-id")
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestGetInformationDocumentImages_NonOKStatus(t *testing.T) {
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

	_, err := client.GetInformationDocumentImages(t.Context(), "test-applicant-id")
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestGetInformationDocumentImages_InvalidJSON(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("invalid json")),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	_, err := client.GetInformationDocumentImages(t.Context(), "test-applicant-id")
	if err == nil {
		t.Fatal("expected error for invalid JSON response, got nil")
	}
}

func TestGetInformationDocumentImages_EmptyResponse(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"items":[],"totalItems":0}`)),
		},
		err: nil,
	}

	client := newMockClient(t, httpClient)

	resp, err := client.GetInformationDocumentImages(t.Context(), "test-applicant-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.TotalItems != 0 {
		t.Errorf("expected TotalItems 0, got %d", resp.TotalItems)
	}

	if len(resp.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(resp.Items))
	}
}
