package restuss

import (
	"fmt"
	"net/http"
)

// AuthProvider expose the methods necessary to perform authenticated calls
type AuthProvider interface {
	AddAuthHeaders(*http.Request)
	Prepare(url string, client *http.Client) error
}

// BasicAuthProvider represent the basic auth method
type BasicAuthProvider struct {
	username string
	password string
	token    string
}

// NewBasicAuthProvider returns a new BasicAuthProvider
func NewBasicAuthProvider(username string, password string) *BasicAuthProvider {
	return &BasicAuthProvider{username: username, password: password}
}

// AddAuthHeaders add auth headers
func (b *BasicAuthProvider) AddAuthHeaders(r *http.Request) {
	panic("not supported yet")
}

// Prepare performs tasks required pre-auth, it should be called before AddAuthHeaders can be used
func (b *BasicAuthProvider) Prepare(url string, c *http.Client) error {
	//todo perform api request and get the token -> b.token
	//todo then use the token in addAuthHeaders
	panic("not supported yet")
	return nil
}

// KeyAuthProvider represent the key based auth method
type KeyAuthProvider struct {
	accessKey string
	secretKey string
}

// NewKeyAuthProvider returns a new KeyAuthProvider
func NewKeyAuthProvider(accessKey string, secretKey string) *KeyAuthProvider {
	return &KeyAuthProvider{accessKey: accessKey, secretKey: secretKey}
}

// Prepare performs tasks required pre-auth, it should be called before AddAuthHeaders can be used
func (k *KeyAuthProvider) Prepare(_ string, _ *http.Client) error {
	return nil
}

// AddAuthHeaders add auth headers
func (k *KeyAuthProvider) AddAuthHeaders(r *http.Request) {
	r.Header.Add(
		"X-ApiKeys",
		fmt.Sprintf("accessKey=%s; secretKey=%s", k.accessKey, k.secretKey),
	)
}
