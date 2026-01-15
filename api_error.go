package gosumsub

import "fmt"

type APIError struct {
	Description   string `json:"description"`
	Code          int    `json:"code"`
	CorrelationID string `json:"correlationId"`
	ErrorCode     int    `json:"errorCode"`
	ErrorName     string `json:"errorName"`
}

func (e *APIError) Error() string {
	msg := e.Description
	if msg == "" {
		msg = "API error"
	}

	details := fmt.Sprintf("code: %d", e.Code)

	if e.ErrorCode != 0 {
		details += fmt.Sprintf(", errorCode: %d", e.ErrorCode)
	}

	if e.ErrorName != "" {
		details += ", errorName: " + e.ErrorName
	}

	if e.CorrelationID != "" {
		details += ", correlationId: " + e.CorrelationID
	}

	return fmt.Sprintf("%s (%s)", msg, details)
}
