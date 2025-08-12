package futures

import (
	"github.com/khanbekov/go-bitget/common/client"
)

// ClientInterface re-exports the common client interface for backward compatibility.
type ClientInterface = client.ClientInterface

// Ensure that Client implements ClientInterface
var _ ClientInterface = (*Client)(nil)
