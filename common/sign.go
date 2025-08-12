package common

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

// Signer handles cryptographic signing of API requests for Bitget authentication.
// It supports both HMAC-SHA256 and RSA signature algorithms.
type Signer struct {
	secretKey []byte // The secret key used for signing requests
}

// NewSigner creates a new Signer instance with the provided secret key.
// The key will be used for HMAC-SHA256 signing or RSA signing depending on the method called.
func NewSigner(key string) *Signer {
	return &Signer{[]byte(key)}
}

// Sign creates an HMAC-SHA256 signature for API authentication.
// This is the most commonly used signing method for Bitget API requests.
//
// Parameters:
//   - method: HTTP method (GET, POST, etc.)
//   - requestPath: The API endpoint path
//   - body: Request body (empty string for GET requests)
//   - timesStamp: Current timestamp as string
//
// Returns a base64-encoded signature string.
func (p *Signer) Sign(method string, requestPath string, body string, timesStamp string) string {
	var payload strings.Builder
	payload.WriteString(timesStamp)
	payload.WriteString(method)
	payload.WriteString(requestPath)
	if body != "" && body != "?" {
		payload.WriteString(body)
	}
	hash := hmac.New(sha256.New, p.secretKey)
	hash.Write([]byte(payload.String()))

	result := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return result
}

// SignByRSA creates an RSA signature for API authentication.
// This is an alternative signing method for enhanced security.
//
// Parameters:
//   - method: HTTP method (GET, POST, etc.)
//   - requestPath: The API endpoint path
//   - body: Request body (empty string for GET requests)
//   - timesStamp: Current timestamp as string
//
// Returns a base64-encoded RSA signature string.
func (p *Signer) SignByRSA(method string, requestPath string, body string, timesStamp string) string {
	var payload strings.Builder
	payload.WriteString(timesStamp)
	payload.WriteString(method)
	payload.WriteString(requestPath)
	if body != "" && body != "?" {
		payload.WriteString(body)
	}

	sign, _ := RSASign([]byte(payload.String()), p.secretKey, crypto.SHA256)
	result := base64.StdEncoding.EncodeToString(sign)
	return result
}

// RSASign performs RSA-PKCS1v15 signing on the provided data.
// Supports both PKCS#1 and PKCS#8 private key formats.
//
// Parameters:
//   - src: The data to be signed
//   - priKey: RSA private key in PEM format
//   - hash: Hash algorithm to use (typically crypto.SHA256)
//
// Returns the signature bytes or an error if signing fails.
func RSASign(src []byte, priKey []byte, hash crypto.Hash) ([]byte, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("key is invalid format")
	}

	var pkixPrivateKey interface{}
	var err error
	if block.Type == "RSA PRIVATE KEY" {
		pkixPrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else if block.Type == "PRIVATE KEY" {
		pkixPrivateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	}

	h := hash.New()
	_, err = h.Write(src)
	if err != nil {
		return nil, err
	}

	bytes := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, pkixPrivateKey.(*rsa.PrivateKey), hash, bytes)
	if err != nil {
		return nil, err
	}

	return sign, nil
}
