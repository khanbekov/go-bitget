// Package common provides shared constants, utilities, and helper functions
// used across the Bitget Go SDK. This includes WebSocket operations,
// authentication methods, and request signing utilities.
package common

// SignType represents the cryptographic signature algorithm used for API authentication.
type SignType string

const (
	// WebSocket authentication and operation constants
	WsAuthMethod    = "GET"          // HTTP method used for WebSocket authentication
	WsAuthPath      = "/user/verify" // API path for WebSocket user verification
	WsOpLogin       = "login"        // WebSocket operation type for login
	WsOpUnsubscribe = "unsubscribe"  // WebSocket operation type for unsubscribing from channels
	WsOpSubscribe   = "subscribe"    // WebSocket operation type for subscribing to channels

	// Supported signature algorithms for API authentication
	RSA    SignType = "RSA"    // RSA signature algorithm
	SHA256 SignType = "SHA256" // HMAC-SHA256 signature algorithm (most common)
)
