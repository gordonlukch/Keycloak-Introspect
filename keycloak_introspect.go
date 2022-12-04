package keycloak_introspect

import (
	"context"
	"errors"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
)

// IntrospectTokenResponse is the response from the introspection endpoint

type IntrospectTokenResponse struct {
	Hostname     string `json:"hostname,omitempty"`
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Realm        string `json:"realm,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Hostname) == 0 {
		return nil, errors.New("The hostname is required")
	}
	if len(config.ClientId) == 0 {
		return nil, errors.New("The client_id is required")
	}
	if len(config.ClientSecret) == 0 {
		return nil, errors.New("The client_secret is required")
	}
	if len(config.Realm) == 0 {
		return nil, errors.New("The realm is required")
	}
	return &keycloak_introspect{
		name:   name,
		next:   next,
		config: config,
	}, nil
}

func (k *keycloak_introspect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	k.config.Hostname
	client := gocloak.NewClient(k.config.Hostname)
	ctx := context.Background()

	const BEARER_SCHEMA = "Bearer "
	authHeader := req.Header.Get("Authorization")
	token := authHeader[len(BEARER_SCHEMA):]

	rptResult, err := client.RetrospectToken(ctx, token, k.config.clientID, k.config.clientSecret, k.config.realm)
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
