package huma_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	huma "github.com/valuetechdev/huma-go"
)

func Example() {
	client, err := huma.New(&huma.ClientCredentials{
		ClientId:     "your-client-id",
		ClientSecret: "your-client-secret",
	})
	if err != nil {
		panic(err)
	}

	res, err := client.ListUsersWithResponse(context.Background(), &huma.ListUsersParams{})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.StatusCode())
}

func Example_withCustomHTTPClient() {
	client, err := huma.New(
		&huma.ClientCredentials{
			ClientId:     "your-client-id",
			ClientSecret: "your-client-secret",
		},
		huma.WithHttpClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
	)
	if err != nil {
		panic(err)
	}

	_ = client
}

func Example_withRequestInterceptor() {
	client, err := huma.New(
		&huma.ClientCredentials{
			ClientId:     "your-client-id",
			ClientSecret: "your-client-secret",
		},
		huma.WithRequestInterceptor(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("X-Request-ID", "trace-123")
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	_ = client
}
