// Package futures provides a Go SDK for the Bitget futures trading API.
// It offers comprehensive functionality for futures trading operations including
// order management, account operations, market data retrieval, and real-time WebSocket connections.
//
// The package follows a service-oriented architecture with a fluent API pattern,
// allowing for intuitive method chaining when building requests.
//
// Example usage:
//
//	client := NewClient("api_key", "secret_key", "passphrase")
//	candles, err := client.NewCandlestickService().
//		Symbol("BTCUSDT").
//		ProductType(ProductTypeUSDTFutures).
//		Granularity("1m").
//		Limit("100").
//		Do(context.Background())
package futures

import (
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/khanbekov/go-bitget/common"
	"github.com/khanbekov/go-bitget/common/client"
	"github.com/khanbekov/go-bitget/common/types"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"net"
	"net/url"

	//jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"

	// NOTE: Subdirectory package imports removed to avoid import cycles
	// Factory methods will be implemented using interface{} returns
)

// BaseApiMainUrl is the main API endpoint for Bitget futures trading.
var (
	BaseApiMainUrl = "https://api.bitget.com"
)

// getApiEndpoint returns the base API endpoint URL.
// Currently returns the main production API URL.
func getApiEndpoint() string {
	return BaseApiMainUrl
}

// Client represents a Bitget futures API client.
// It contains all necessary credentials and configuration for interacting with the Bitget API.
type Client struct {
	// API credentials
	apiKey     string
	secretKey  string
	keyType    string
	passphrase string

	// HTTP client configuration
	BaseURL    string
	UserAgent  string
	fastClient *fasthttp.Client

	// Debugging and logging
	Debug  bool
	Logger zerolog.Logger

	// Request signing
	signer *common.Signer
}

// NewClient initializes a new Bitget futures API client with the provided credentials.
// This function must be called before using any SDK functionality.
//
// Parameters:
//   - apiKey: Your Bitget API key
//   - secretKey: Your Bitget secret key for request signing
//   - passphrase: Your API passphrase
//
// Returns a configured Client instance ready for use.
// Services can be created using the client.NewXXXService() pattern.
//
// Example:
//
//	client := NewClient("your_api_key", "your_secret_key", "your_passphrase")
//	candles, err := client.NewCandlestickService().Symbol("BTCUSDT").Do(ctx)
func NewClient(apiKey, secretKey, passphrase string) *Client {
	return &Client{
		apiKey:     apiKey,
		secretKey:  secretKey,
		passphrase: passphrase,
		signer:     common.NewSigner(secretKey),
		BaseURL:    getApiEndpoint(),
		UserAgent:  "Bitget/golang",
		fastClient: &fasthttp.Client{},
		Logger:     zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}

// callAPI sends an HTTP request to the specified Bitget API endpoint with automatic retry logic.
// It handles request signing, authentication headers, and error retry for transient failures.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - method: HTTP method (GET, POST, etc.)
//   - endpoint: API endpoint path
//   - queryParams: URL query parameters
//   - body: Request body for POST requests
//   - sign: Whether the request requires authentication signing
//
// Returns the API response, response headers, and any error encountered.
// Implements exponential backoff retry for retryable errors (network timeouts, etc.).
func (c *Client) CallAPI(ctx context.Context, method string, endpoint string, queryParams url.Values, body []byte, sign bool) (*client.ApiResponse, *fasthttp.ResponseHeader, error) {
	const maxRetries = 3
	var backoff = 1 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()

		// Build URL
		requestURL := c.GetUrl(endpoint)
		if queryParams != nil && len(queryParams) > 0 {
			requestURL += "?" + queryParams.Encode()
		}
		req.SetRequestURI(requestURL)

		// Set method and body
		req.Header.SetMethod(method)
		if method == "POST" {
			req.SetBody(body)
			req.Header.Set("Content-Type", "application/json")
		}

		// Sign the request if needed
		if sign {
			ts := common.TimestampMs()
			req.Header.Set("ACCESS-TIMESTAMP", ts)
			req.Header.Set("ACCESS-KEY", c.apiKey)
			req.Header.Set("ACCESS-PASSPHRASE", c.passphrase)
			req.Header.Set("locale", "en-US")

			var reqParamStr string
			if method == "GET" {
				reqParamStr = "?" + queryParams.Encode()
			} else {
				reqParamStr = string(body)
			}
			sign := c.signer.Sign(method, endpoint, reqParamStr, ts)
			req.Header.Set("ACCESS-SIGN", sign)
		}

		// Execute request
		done := make(chan error, 1)
		go func() {
			err := c.fastClient.DoTimeout(req, resp, 5*time.Second)
			done <- err
		}()

		select {
		case <-ctx.Done():
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			return nil, nil, ctx.Err()
		case err := <-done:
			if err != nil {
				// Handle retryable errors
				if isRetryableError(err) {
					fasthttp.ReleaseRequest(req)
					fasthttp.ReleaseResponse(resp)
					if attempt == maxRetries-1 {
						return nil, nil, err
					}
					time.Sleep(backoff)
					backoff *= 2
					continue
				}

				// Return non-retryable errors immediately
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(resp)
				return nil, nil, err
			}

			// Process response
			if resp.StatusCode() >= http.StatusBadRequest {
				apiErr := &types.APIError{}
				if err := jsoniter.Unmarshal(resp.Body(), apiErr); err != nil {
					fasthttp.ReleaseRequest(req)
					fasthttp.ReleaseResponse(resp)
					return nil, nil, fmt.Errorf("error parsing API response: %w", err)
				}
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(resp)
				return nil, nil, apiErr
			}

			// Success case
			var apiResp client.ApiResponse
			if err := jsoniter.Unmarshal(resp.Body(), &apiResp); err != nil {
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(resp)
				return nil, nil, err
			}

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			return &apiResp, &resp.Header, nil
		}
	}

	return nil, nil, fmt.Errorf("max retries exceeded")
}

// isRetryableError determines if an error is transient and worth retrying.
// Returns true for network timeouts, connection errors, and other temporary failures.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// FastHTTP-specific errors
	if errors.Is(err, fasthttp.ErrTimeout) ||
		errors.Is(err, fasthttp.ErrNoFreeConns) ||
		errors.Is(err, fasthttp.ErrDialTimeout) {
		return true
	}

	// Network errors
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Temporary() {
		return true
	}

	return false
}

// SetApiEndpoint sets a custom API endpoint URL for the client.
// This can be used to switch between different environments or use a proxy.
func (c *Client) SetApiEndpoint(url string) *Client {
	c.BaseURL = url
	return c
}

// GetUrl constructs the full URL by combining the base URL with the given endpoint.
func (c *Client) GetUrl(endpoint string) string {
	return c.BaseURL + endpoint
}

// Factory methods are not provided in the main client to avoid import cycles.
// Use the package-specific constructors instead:
//
// Account Services Example:
//   import "github.com/khanbekov/go-bitget/futures/account"
//   client := futures.NewClient(apiKey, secretKey, passphrase)
//   accountInfo := account.NewAccountInfoService(client)
//   result, err := accountInfo.Symbol("BTCUSDT").ProductType("USDT-FUTURES").Do(ctx)
//
// Market Services Example:
//   import "github.com/khanbekov/go-bitget/futures/market"
//   candles := market.NewCandlestickService(client)
//   result, err := candles.Symbol("BTCUSDT").ProductType("USDT-FUTURES").Do(ctx)
//
// Trading Services Example:
//   import "github.com/khanbekov/go-bitget/futures/trading"
//   order := trading.NewCreateOrderService(client)
//   result, err := order.Symbol("BTCUSDT").Side("buy").Size("0.01").Do(ctx)
//
// Position Services Example:
//   import "github.com/khanbekov/go-bitget/futures/position"
//   positions := position.NewAllPositionsService(client)
//   result, err := positions.ProductType("USDT-FUTURES").Do(ctx)
//
// This approach provides:
// - Strong type safety 
// - No import cycles
// - Clear package organization
// - Better IDE support with auto-completion
