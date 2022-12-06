package keycloakintrospect

import (
	"context"
	"errors"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
)

// IntrospectTokenResponse is the response from the introspection endpoint

// Config the plugin configuration.
type Config struct {
	Hostname     string `json:"hostname,omitempty"`
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	Realm        string `json:"realm,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// string returns a readable representation of the configuration.
type keycloak struct {
	name   string
	next   http.Handler
	config *Config
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Hostname) == 0 {
		return nil, errors.New("hostname is required")
	}
	if len(config.ClientID) == 0 {
		return nil, errors.New("client_id is required")
	}
	if len(config.ClientSecret) == 0 {
		return nil, errors.New("client_secret is required")
	}
	if len(config.Realm) == 0 {
		return nil, errors.New("realm is required")
	}
	return &keycloak{
		name:   name,
		next:   next,
		config: config,
	}, nil
}

// ServeHTTP implements the http.Handler interface.
// nolint: contextcheck
func (k *keycloak) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	client := gocloak.NewClient(k.config.Hostname)
	ctx := context.Background()

	const BearerSchema = "Bearer "
	authHeader := req.Header.Get("Authorization")
	token := authHeader[len(BearerSchema):]

	rptResult, err := client.RetrospectToken(ctx, token, k.config.ClientID, k.config.ClientSecret, k.config.Realm)
	if err != nil {
		http.Error(rw, "Token Inspection Error : "+err.Error(), http.StatusUnauthorized)
		return
	}
	if rptResult.Active == nil {
		http.Error(rw, "Token Inspection Failed : Bearer Token is not Active", http.StatusUnauthorized)
		return
	}

	k.next.ServeHTTP(rw, req)
}
