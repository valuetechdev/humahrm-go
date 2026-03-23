//go:generate go tool -modfile=./go.tool.mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml ./api/openapi.yaml

// Package humahrm provides a Go client for the Huma HR API.
//
// The client handles OAuth2 authentication using client credentials and
// provides type-safe access to the Huma HR API endpoints.
//
// # Usage
//
// Create a new client with your credentials:
//
//	client, err := huma.New(&huma.ClientCredentials{
//		ClientId:     "your-client-id",
//		ClientSecret: "your-client-secret",
//	})
//	if err != nil {
//		return err
//	}
//
//	// Make API calls - authentication is handled automatically
//	resp, err := client.ListUsersWithResponse(ctx, &huma.ListUsersParams{})
//
// # Custom HTTP Client
//
// You can provide a custom http.Client for advanced use cases like proxies
// or custom timeouts:
//
//	client, err := huma.New(creds, huma.WithHttpClient(&http.Client{
//		Timeout: 30 * time.Second,
//	}))
package humahrm

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Client is the main entry point for interacting with the Huma HR API.
// It wraps the generated OpenAPI client with OAuth2 authentication.
type Client struct {
	token        *oauth2.Token
	tokenSource  oauth2.TokenSource
	interceptors []RequestEditorFn
	httpClient   *http.Client
	*ClientWithResponses
}

// ClientCredentials holds the OAuth2 client credentials used for authentication.
type ClientCredentials struct {
	ClientId     string
	ClientSecret string
}

// Option configures a [Client]. Use the With* functions to create Options.
type Option func(*Client)

// WithHttpClient sets a custom http.Client. Defaults to [http.DefaultClient].
func WithHttpClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithRequestInterceptor appends fn to the [Client] interceptors
func WithRequestInterceptor(fn RequestEditorFn) Option {
	return func(c *Client) {
		c.interceptors = append(c.interceptors, fn)
	}
}

// WithToken sets an initial OAuth2 token. The token will be used until it
// expires, then a new token will be fetched automatically using the client
// credentials. Use this to restore a previously saved token from [Client.Token].
func WithToken(token *oauth2.Token) Option {
	return func(c *Client) {
		c.token = token
	}
}

// New creates a new [Client] with the given credentials.
//
// The client automatically handles OAuth2 authentication using the client
// credentials flow. Tokens are obtained on first request and refreshed
// automatically when expired.
//
// Options can be provided to customize client behavior:
//   - [WithHttpClient]: Use a custom http.Client as the base transport
//   - [WithRequestInterceptor]: Add request interceptors
//   - [WithToken]: Restore a previously saved token
func New(creds *ClientCredentials, options ...Option) (*Client, error) {
	conf := &clientcredentials.Config{
		ClientID:     creds.ClientId,
		ClientSecret: creds.ClientSecret,
		TokenURL:     "https://demo.openapi.humahr.com/auth/oauth/token",
	}

	baseUrl := "https://demo.openapi.humahr.com/"
	client := &Client{httpClient: http.DefaultClient}

	for _, option := range options {
		option(client)
	}

	// Inject the custom http.Client into the context so oauth2 uses it as the base transport
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client.httpClient)
	tokenSource := conf.TokenSource(ctx)
	// If an initial token was provided, reuse it until it expires
	client.tokenSource = oauth2.ReuseTokenSource(client.token, tokenSource)
	oauthClient := oauth2.NewClient(ctx, client.tokenSource)

	clientOptions := []ClientOption{
		WithHTTPClient(oauthClient),
	}

	for _, interceptor := range client.interceptors {
		clientOptions = append(clientOptions, WithRequestEditorFn(interceptor))
	}

	c, err := NewClientWithResponses(
		baseUrl,
		clientOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init client: %w", err)
	}
	client.ClientWithResponses = c
	return client, nil
}

// Token returns the current OAuth2 token, fetching a new one if necessary.
// This can be used to inspect token details or for external integrations.
func (c *Client) Token() (*oauth2.Token, error) {
	return c.tokenSource.Token()
}
