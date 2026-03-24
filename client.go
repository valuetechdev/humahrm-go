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
	"errors"
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

	baseURL string

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

// WithCustomBaseURL returns an Option that overrides the default API base URL.
func WithCustomBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
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
	if creds == nil {
		return nil, errors.New("credentials must not be nil")
	}
	if creds.ClientId == "" {
		return nil, errors.New("client id must not be empty")
	}
	if creds.ClientSecret == "" {
		return nil, errors.New("client secret must not be empty")
	}

	client := &Client{
		httpClient: http.DefaultClient,
		baseURL:    "https://openapi.humahr.com",
	}

	for _, option := range options {
		option(client)
	}
	conf := &clientcredentials.Config{
		ClientID:     creds.ClientId,
		ClientSecret: creds.ClientSecret,
		TokenURL:     client.baseURL + "/auth/oauth/token",
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
		client.baseURL,
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
