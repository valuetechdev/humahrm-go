package humahrm

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

var validCreds = &ClientCredentials{
	ClientId:     "test-id",
	ClientSecret: "test-secret",
}

func TestNew_NilCredentials(t *testing.T) {
	c, err := New(nil)
	assert.Nil(t, c)
	assert.EqualError(t, err, "credentials must not be nil")
}

func TestNew_EmptyClientId(t *testing.T) {
	c, err := New(&ClientCredentials{
		ClientId:     "",
		ClientSecret: "secret",
	})
	assert.Nil(t, c)
	assert.EqualError(t, err, "client id must not be empty")
}

func TestNew_EmptyClientSecret(t *testing.T) {
	c, err := New(&ClientCredentials{
		ClientId:     "id",
		ClientSecret: "",
	})
	assert.Nil(t, c)
	assert.EqualError(t, err, "client secret must not be empty")
}

func TestNew_ValidCredentials(t *testing.T) {
	c, err := New(validCreds)
	require.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "https://openapi.humahr.com", c.baseURL)
}

func TestNew_WithCustomBaseURL(t *testing.T) {
	c, err := New(validCreds, WithCustomBaseURL("https://custom.example.com"))
	require.NoError(t, err)
	assert.Equal(t, "https://custom.example.com", c.baseURL)
}

func TestNew_WithHttpClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	c, err := New(validCreds, WithHttpClient(custom))
	require.NoError(t, err)
	assert.Equal(t, custom, c.httpClient)
}

func TestNew_WithRequestInterceptor(t *testing.T) {
	interceptor := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Test", "value")
		return nil
	}
	c, err := New(validCreds, WithRequestInterceptor(interceptor))
	require.NoError(t, err)
	assert.Len(t, c.interceptors, 1)
}

func TestNew_WithMultipleInterceptors(t *testing.T) {
	noop := func(ctx context.Context, req *http.Request) error { return nil }
	c, err := New(validCreds,
		WithRequestInterceptor(noop),
		WithRequestInterceptor(noop),
		WithRequestInterceptor(noop),
	)
	require.NoError(t, err)
	assert.Len(t, c.interceptors, 3)
}

func TestNew_WithToken(t *testing.T) {
	tok := &oauth2.Token{AccessToken: "cached-token"}
	c, err := New(validCreds, WithToken(tok))
	require.NoError(t, err)
	assert.Equal(t, tok, c.token)
}

func TestNew_WithAllOptions(t *testing.T) {
	custom := &http.Client{Timeout: 10 * time.Second}
	tok := &oauth2.Token{AccessToken: "tok"}
	noop := func(ctx context.Context, req *http.Request) error { return nil }

	c, err := New(validCreds,
		WithHttpClient(custom),
		WithCustomBaseURL("https://custom.example.com"),
		WithToken(tok),
		WithRequestInterceptor(noop),
	)
	require.NoError(t, err)
	assert.Equal(t, custom, c.httpClient)
	assert.Equal(t, "https://custom.example.com", c.baseURL)
	assert.Equal(t, tok, c.token)
	assert.Len(t, c.interceptors, 1)
	assert.NotNil(t, c.ClientWithResponses)
	assert.NotNil(t, c.tokenSource)
}

func TestClientInitialization(t *testing.T) {
	require := require.New(t)

	c, err := New(&ClientCredentials{
		ClientId:     os.Getenv("HUMA_CLIENT_ID"),
		ClientSecret: os.Getenv("HUMA_CLIENT_SECRET"),
	}, WithCustomBaseURL("https://demo.openapi.humahr.com"))
	require.NoError(err)

	// Authentication happens automatically on first request
	res, err := c.ListUsersWithResponse(context.Background(), &ListUsersParams{})
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode(), "status should be OK 200")
	require.NotNil(res.JSON200, "JSON200 should not be nil")
	require.NotEmpty(res.JSON200.Items, "JSON200.Items should not be empty")
}
