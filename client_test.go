package humahrm

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientInitialization(t *testing.T) {
	require := require.New(t)

	c, err := New(&ClientCredentials{
		ClientId:     os.Getenv("HUMA_CLIENT_ID"),
		ClientSecret: os.Getenv("HUMA_CLIENT_SECRET"),
	})
	require.NoError(err)

	// Authentication happens automatically on first request
	res, err := c.ListUsersWithResponse(context.Background(), &ListUsersParams{})
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode(), "status should be OK 200")
	require.NotNil(res.JSON200, "JSON200 should not be nil")
	require.NotEmpty(res.JSON200.Items, "JSON200.Items should not be empty")
}
