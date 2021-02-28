package bigone

import (
	"context"
)

// PingServer PingServer
func PingServer() (int64, error) {
	resp, err := HTTPRequest(context.Background()).Get("/ping")
	if err != nil {
		return 0, err
	}

	var data struct {
		Timestamp int64 `json:"timestamp,omitempty"`
	}

	if err := UnmarshalResponse(resp, &data); err != nil {
		return 0, err
	}

	return data.Timestamp, nil
}
