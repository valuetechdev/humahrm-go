package humahrm_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	humahrm "github.com/valuetechdev/humahrm-go"
)

func Example() {
	client, err := humahrm.New(&humahrm.ClientCredentials{
		ClientId:     "your-client-id",
		ClientSecret: "your-client-secret",
	})
	if err != nil {
		panic(err)
	}

	res, err := client.ListUsersWithResponse(context.Background(), &humahrm.ListUsersParams{})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.StatusCode())
}

func Example_withCustomHTTPClient() {
	client, err := humahrm.New(
		&humahrm.ClientCredentials{
			ClientId:     "your-client-id",
			ClientSecret: "your-client-secret",
		},
		humahrm.WithHttpClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
	)
	if err != nil {
		panic(err)
	}

	_ = client
}

func Example_withRequestInterceptor() {
	client, err := humahrm.New(
		&humahrm.ClientCredentials{
			ClientId:     "your-client-id",
			ClientSecret: "your-client-secret",
		},
		humahrm.WithRequestInterceptor(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("X-Request-ID", "trace-123")
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	_ = client
}
