package gosumsub

import (
	"context"
	"encoding/json"
	"net/http"
)

type Resolution struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

type FileMetadata struct {
	FileName   string      `json:"fileName,omitempty"`
	FileType   string      `json:"fileType,omitempty"`
	FileSize   int         `json:"fileSize,omitempty"`
	Resolution *Resolution `json:"resolution,omitempty"`
}

type IDDocDef struct {
	Country      string `json:"country,omitempty"`
	IDDocType    string `json:"idDocType,omitempty"`
	IDDocSubType string `json:"idDocSubType,omitempty"`
}

type DocumentImageItem struct {
	ID           string        `json:"id,omitempty"`
	PreviewID    string        `json:"previewId,omitempty"`
	AddedDate    string        `json:"addedDate,omitempty"`
	FileMetadata *FileMetadata `json:"fileMetadata,omitempty"`
	IDDocDef     *IDDocDef     `json:"idDocDef,omitempty"`
	ReviewResult *ReviewResult `json:"reviewResult,omitempty"`
	Deactivated  bool          `json:"deactivated,omitempty"`
	AttemptID    string        `json:"attemptId,omitempty"`
	Source       string        `json:"source,omitempty"`
}

type DocumentImagesResponse struct {
	Items      []DocumentImageItem `json:"items,omitempty"`
	TotalItems int                 `json:"totalItems,omitempty"`
}

func (c *Client) GetInformationDocumentImages(ctx context.Context, applicantID string) (*DocumentImagesResponse, error) {
	if applicantID == "" {
		return nil, ErrApplicantIDRequired
	}

	apiRequest := request{
		Method:   http.MethodGet,
		Endpoint: "/resources/applicants/" + applicantID + "/metadata/resources",
		Params:   nil,
		Query:    nil,
		Header:   nil,
		Body:     nil,
		FullURL:  "",
	}

	body, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	var resp DocumentImagesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
