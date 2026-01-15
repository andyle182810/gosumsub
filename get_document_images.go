package gosumsub

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var (
	ErrInspectionIDRequired = errors.New("inspectionId is required")
	ErrImageIDRequired      = errors.New("imageId is required")
)

type GetDocumentImageResponse struct {
	Data     []byte `json:"data"`
	MimeType string `json:"mimeType"`
}

func (r *GetDocumentImageResponse) GetBase64WithMime() string {
	if r.Data == nil {
		return ""
	}

	base64Data := base64.StdEncoding.EncodeToString(r.Data)

	return fmt.Sprintf("data:%s;base64,%s", r.MimeType, base64Data)
}

func (c *Client) GetDocumentImage(ctx context.Context, inspectionID, imageID string) (*GetDocumentImageResponse, error) {
	if inspectionID == "" {
		return nil, ErrInspectionIDRequired
	}

	if imageID == "" {
		return nil, ErrImageIDRequired
	}

	apiRequest := request{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/resources/inspections/%s/resources/%s", url.PathEscape(inspectionID), url.PathEscape(imageID)),
		Params:   nil,
		Query:    nil,
		Header:   nil,
		Body:     nil,
		FullURL:  "",
	}

	body, contentType, err := c.executeWithContentType(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	mimeType := contentType
	if mimeType == "" {
		mimeType = http.DetectContentType(body)
	}

	return &GetDocumentImageResponse{
		Data:     body,
		MimeType: mimeType,
	}, nil
}
