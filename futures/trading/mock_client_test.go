package trading

import (
	"context"
	"net/url"

	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

// MockClient is a mock implementation of the ClientInterface for testing
type MockClient struct {
	mock.Mock
}

func (m *MockClient) CallAPI(ctx context.Context, method string, endpoint string, queryParams url.Values, body []byte, sign bool) (*ApiResponse, *fasthttp.ResponseHeader, error) {
	args := m.Called(ctx, method, endpoint, queryParams, body, sign)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*fasthttp.ResponseHeader), args.Error(2)
	}
	return args.Get(0).(*ApiResponse), args.Get(1).(*fasthttp.ResponseHeader), args.Error(2)
}

// Ensure MockClient implements ClientInterface
var _ ClientInterface = (*MockClient)(nil)
