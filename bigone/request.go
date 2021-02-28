package bigone

import (
	"auto-trade/global"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetHostURL(global.BigOneSetting.BASEAPI).
	SetTimeout(2 * time.Second)

// HTTPError HTTPError
type HTTPError struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"message,omitempty"`
}

func (err *HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Msg)
}

// HTTPRequest HTTPRequest
func HTTPRequest(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func decodeResponse(resp *resty.Response) ([]byte, error) {
	var body struct {
		HTTPError
		Data      json.RawMessage `json:"data,omitempty"`
		Code      int             `json:"code,omitempty"`
		Message   string          `json:"message,omitempty"`
		PageToken string          `json:"page_token,omitempty"`
	}

	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		if resp.IsError() {
			return nil, &HTTPError{
				Code: resp.StatusCode(),
				Msg:  resp.Status(),
			}
		}

		return nil, err
	}

	if body.Data == nil {
		return nil, &HTTPError{
			Code: -1,
			Msg:  body.Message,
		}
	}

	if body.HTTPError.Code > 0 {
		return nil, &body.HTTPError
	}

	return body.Data, nil
}

// UnmarshalResponse Unmarshal Response
func UnmarshalResponse(resp *resty.Response, v interface{}) error {
	data, err := decodeResponse(resp)
	if err != nil {
		return err
	}
	if v != nil {
		return json.Unmarshal(data, v)
	}

	return nil
}
