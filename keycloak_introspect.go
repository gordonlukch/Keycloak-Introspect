package keycloak_introspect

import (
	"context"
	"errors"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
)

// IntrospectTokenResponse is the response from the introspection endpoint

type Config struct {
	Hostname     string `json:"hostname,omitempty"`
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Realm        string `json:"realm,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type keycloak struct {
	name   string
	next   http.Handler
	config *Config
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Hostname) == 0 {
		return nil, errors.New("hostname is required")
	}
	if len(config.ClientId) == 0 {
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

func (k *keycloak) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	client := gocloak.NewClient(k.config.Hostname)
	ctx := context.Background()

	const BEARER_SCHEMA = "Bearer "
	authHeader := req.Header.Get("Authorization")
	token := authHeader[len(BEARER_SCHEMA):]

	rptResult, err := client.RetrospectToken(ctx, token, k.config.ClientId, k.config.ClientSecret, k.config.Realm)
	if err != nil {
		http.Error(rw, "Token Inspection: Failed", http.StatusUnauthorized)
		return
	}

	if !rptResult.Active {
		http.Error(rw, "Token is not active", http.StatusUnauthorized)
		return
	}

	k.next.ServeHTTP(rw, req)
}
