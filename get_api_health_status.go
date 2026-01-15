package gosumsub

import (
	"context"
	"net/http"
)

func (c *Client) GetAPIHealthStatus(ctx context.Context) error {
	apiRequest := request{
		Method:   http.MethodGet,
		Endpoint: "/resources/status/api",
		Params:   nil,
		Query:    nil,
		Header:   nil,
		Body:     nil,
		FullURL:  "",
	}

	_, err := c.execute(ctx, &apiRequest)
	if err != nil {
		return err
	}

	return nil
}
