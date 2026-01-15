package gosumsub

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

var ErrApplicantIDRequired = errors.New("applicantId is required")

type IDDoc struct {
	IDDocType   string `json:"idDocType,omitempty"`
	Country     string `json:"country,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	FirstNameEn string `json:"firstNameEn,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	LastNameEn  string `json:"lastNameEn,omitempty"`
	IssuedDate  string `json:"issuedDate,omitempty"`
	ValidUntil  string `json:"validUntil,omitempty"`
	Number      string `json:"number,omitempty"`
	Dob         string `json:"dob,omitempty"`
	MrzLine1    string `json:"mrzLine1,omitempty"`
	MrzLine2    string `json:"mrzLine2,omitempty"`
	MrzLine3    string `json:"mrzLine3,omitempty"`
	Termless    bool   `json:"termless,omitempty"`
}

type ApplicantInfo struct {
	FirstName   string  `json:"firstName,omitempty"`
	FirstNameEn string  `json:"firstNameEn,omitempty"`
	LastName    string  `json:"lastName,omitempty"`
	LastNameEn  string  `json:"lastNameEn,omitempty"`
	Dob         string  `json:"dob,omitempty"`
	Country     string  `json:"country,omitempty"`
	IDDocs      []IDDoc `json:"idDocs,omitempty"`
}

type Address struct {
	SubStreet        string `json:"subStreet,omitempty"`
	SubStreetEn      string `json:"subStreetEn,omitempty"`
	Street           string `json:"street,omitempty"`
	StreetEn         string `json:"streetEn,omitempty"`
	State            string `json:"state,omitempty"`
	StateEn          string `json:"stateEn,omitempty"`
	Town             string `json:"town,omitempty"`
	TownEn           string `json:"townEn,omitempty"`
	PostCode         string `json:"postCode,omitempty"`
	FormattedAddress string `json:"formattedAddress,omitempty"`
}

type FixedInfo struct {
	Gender      string    `json:"gender,omitempty"`
	Nationality string    `json:"nationality,omitempty"`
	Addresses   []Address `json:"addresses,omitempty"`
}

type AgreementItem struct {
	ID         string   `json:"id,omitempty"`
	AcceptedAt string   `json:"acceptedAt,omitempty"`
	Source     string   `json:"source,omitempty"`
	Type       string   `json:"type,omitempty"`
	RecordIDs  []string `json:"recordIds,omitempty"`
}

type Agreement struct {
	Items      []AgreementItem `json:"items,omitempty"`
	CreatedAt  string          `json:"createdAt,omitempty"`
	AcceptedAt string          `json:"acceptedAt,omitempty"`
	Source     string          `json:"source,omitempty"`
	RecordIDs  []string        `json:"recordIds,omitempty"`
}

type NfcVerificationSettings struct {
	Mode string `json:"mode,omitempty"`
}

type DocSetField struct {
	Name               string `json:"name,omitempty"`
	Required           bool   `json:"required,omitempty"`
	Prefill            any    `json:"prefill,omitempty"`
	ImmutableIfPresent bool   `json:"immutableIfPresent,omitempty"`
}

type DocSet struct {
	IDDocSetType            string                   `json:"idDocSetType,omitempty"`
	Types                   []string                 `json:"types,omitempty"`
	VideoRequired           string                   `json:"videoRequired,omitempty"`
	CaptureMode             string                   `json:"captureMode,omitempty"`
	UploaderMode            string                   `json:"uploaderMode,omitempty"`
	NfcVerificationSettings *NfcVerificationSettings `json:"nfcVerificationSettings,omitempty"`
	Fields                  []DocSetField            `json:"fields,omitempty"`
}

type RequiredIDDocs struct {
	DocSets []DocSet `json:"docSets,omitempty"`
}

type ReviewResult struct {
	ReviewAnswer string `json:"reviewAnswer,omitempty"`
}

type Review struct {
	ReviewID              string        `json:"reviewId,omitempty"`
	AttemptID             string        `json:"attemptId,omitempty"`
	AttemptCnt            int           `json:"attemptCnt,omitempty"`
	ElapsedSincePendingMs int64         `json:"elapsedSincePendingMs,omitempty"`
	ElapsedSinceQueuedMs  int64         `json:"elapsedSinceQueuedMs,omitempty"`
	Reprocessing          bool          `json:"reprocessing,omitempty"`
	LevelName             string        `json:"levelName,omitempty"`
	LevelAutoCheckMode    any           `json:"levelAutoCheckMode,omitempty"`
	CreateDate            string        `json:"createDate,omitempty"`
	ReviewDate            string        `json:"reviewDate,omitempty"`
	ReviewResult          *ReviewResult `json:"reviewResult,omitempty"`
	ReviewStatus          string        `json:"reviewStatus,omitempty"`
	Priority              int           `json:"priority,omitempty"`
}

type ApplicantData struct {
	ID                string          `json:"id,omitempty"`
	CreatedAt         string          `json:"createdAt,omitempty"`
	Key               string          `json:"key,omitempty"`
	ClientID          string          `json:"clientId,omitempty"`
	InspectionID      string          `json:"inspectionId,omitempty"`
	ExternalUserID    string          `json:"externalUserId,omitempty"`
	Info              *ApplicantInfo  `json:"info,omitempty"`
	FixedInfo         *FixedInfo      `json:"fixedInfo,omitempty"`
	Email             string          `json:"email,omitempty"`
	Phone             string          `json:"phone,omitempty"`
	PhoneCountry      string          `json:"phoneCountry,omitempty"`
	ApplicantPlatform string          `json:"applicantPlatform,omitempty"`
	Agreement         *Agreement      `json:"agreement,omitempty"`
	RequiredIDDocs    *RequiredIDDocs `json:"requiredIdDocs,omitempty"`
	Review            *Review         `json:"review,omitempty"`
	Lang              string          `json:"lang,omitempty"`
	Type              string          `json:"type,omitempty"`
	Notes             []string        `json:"notes,omitempty"`
}

func (c *Client) GetApplicantData(ctx context.Context, applicantID string) (*ApplicantData, error) {
	if applicantID == "" {
		return nil, ErrApplicantIDRequired
	}

	apiRequest := request{
		Method:   http.MethodGet,
		Endpoint: "/resources/applicants/" + applicantID + "/one",
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

	var resp ApplicantData
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
