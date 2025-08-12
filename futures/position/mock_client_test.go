package position

import (
	"context"
	"net/url"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

// MockClient is a mock implementation of the ClientInterface for testing
type MockClient struct {
	mock.Mock
}

func (m *MockClient) CallAPI(ctx context.Context, method string, endpoint string, queryParams url.Values, body []byte, sign bool) (*futures.ApiResponse, *fasthttp.ResponseHeader, error) {
	args := m.Called(ctx, method, endpoint, queryParams, body, sign)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*fasthttp.ResponseHeader), args.Error(2)
	}
	return args.Get(0).(*futures.ApiResponse), args.Get(1).(*fasthttp.ResponseHeader), args.Error(2)
}

// Ensure MockClient implements ClientInterface
var _ futures.ClientInterface = (*MockClient)(nil)