package gosumsub

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrTooManyQueryParams = errors.New("allowedQueryParams exceeds maximum of 4")
	ErrLevelNameRequired  = errors.New("levelName is required")
)

type ApplicantIdentifiers struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type Redirect struct {
	AllowedQueryParams []string `json:"allowedQueryParams,omitempty"`
	SuccessURL         string   `json:"successUrl,omitempty"`
	RejectURL          string   `json:"rejectUrl,omitempty"`
	SignKey            string   `json:"signKey,omitempty"`
}

type GenerateExternalWebSDKLinkRequest struct {
	TTLInSecs            int64                 `json:"ttlInSecs,omitempty"`
	UserID               string                `json:"userId,omitempty"`
	LevelName            string                `json:"levelName"`
	ApplicantIdentifiers *ApplicantIdentifiers `json:"applicantIdentifiers,omitempty"`
	Redirect             *Redirect             `json:"redirect,omitempty"`
}

type GenerateExternalWebSDKLinkResponse struct {
	URL string `json:"url"`
}

func (c *Client) GenerateExternalWebSDKLink(
	ctx context.Context,
	req *GenerateExternalWebSDKLinkRequest,
) (*GenerateExternalWebSDKLinkResponse, error) {
	if req.LevelName == "" {
		return nil, ErrLevelNameRequired
	}

	if req.Redirect != nil && len(req.Redirect.AllowedQueryParams) > 4 {
		return nil, ErrTooManyQueryParams
	}

	if req.TTLInSecs == 0 {
		req.TTLInSecs = 1800
	}

	apiRequest := request{
		Method:   http.MethodPost,
		Endpoint: "/resources/sdkIntegrations/levels/-/websdkLink",
		Params:   req,
		Query:    nil,
		Header:   nil,
		Body:     nil,
		FullURL:  "",
	}

	body, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return nil, err
	}

	var resp GenerateExternalWebSDKLinkResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
