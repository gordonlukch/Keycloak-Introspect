package keycloak_introspect

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"os"
	"github.com/ory/fosite"
}

// IntrospectTokenResponse is the response from the introspection endpoint

type IntrospectTokenResponse struct {
	Active   bool   `json:"active"`
	Scope    string `json:"scope"`
	ClientID string `json:"client_id"`
	Username string `json:"username"`
	TokenType string `json:"token_type"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}


// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}


// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// ...
	return &Example{
		// ...
	}, nil
}

func (e *Example) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// ...
	e.next.ServeHTTP(rw, req)
}