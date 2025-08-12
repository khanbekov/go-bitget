package client

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/valyala/fasthttp"
)

// ClientInterface defines the contract for API clients.
// This allows for easy mocking and testing of services.
type ClientInterface interface {
	CallAPI(ctx context.Context, method string, endpoint string, queryParams url.Values, body []byte, sign bool) (*ApiResponse, *fasthttp.ResponseHeader, error)
}

// ApiResponse represents the standard response structure from the Bitget API.
type ApiResponse struct {
	Code        string          `json:"code"`
	Msg         string          `json:"msg"`
	RequestTime int64           `json:"requestTime"`
	Data        json.RawMessage `json:"data"`
}