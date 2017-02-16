package restuss

import (
	"fmt"
	"net/http"
)

type AuthProvider interface {
	AddAuthHeaders(*http.Request)
	Prepare(url string, client *http.Client) error
}

type basicAuthProvider struct {
	username string
	password string
	token    string
}

func (b *basicAuthProvider) AddAuthHeaders(r *http.Request) {
	panic("not supported yet")
}
func (b *basicAuthProvider) Prepare(url string, c *http.Client) error {
	//todo perform api request and het the token -> b.token
	//todo then use the token in addAuthHeaders
	panic("not supported yet")
	return nil
}

func NewBasicAuthProvider(username string, password string) AuthProvider {
	return &basicAuthProvider{username: username, password: password}
}

type keyAuthProvider struct {
	accessKey string
	secretKey string
}

func (k *keyAuthProvider) Prepare(_ string, _ *http.Client) error {
	return nil
}

func (k *keyAuthProvider) AddAuthHeaders(r *http.Request) {
	r.Header.Add(
		"X-ApiKeys",
		fmt.Sprintf("accessKey=%s; secretKey=%s", k.accessKey, k.secretKey),
	)
}

func NewKeyAuthProvider(accessKey string, secretKey string) AuthProvider {
	return &keyAuthProvider{accessKey: accessKey, secretKey: secretKey}
}
