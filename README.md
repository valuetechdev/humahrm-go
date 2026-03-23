[![go reference badge](https://pkg.go.dev/badge/github.com/valuetechdev/huma-go.svg)](https://pkg.go.dev/github.com/valuetechdev/huma-go)

# huma-go

Go API client for [Huma HRM]. It's generated with [oapi-codegen]

## Prerequisites

1. `clientId`
2. `clientSecret`

## Usage

```bash
go get github.com/valuetechdev/huma-go
```

```go
import "github.com/valuetechdev/huma-go"

func yourFunc() error {
	client, err := huma.New(&huma.ClientCredentials{
		ClientId:     "your-id",
		ClientSecret: "your-secret",
	})
	if err != nil {
		return fmt.Errorf("failed to init client: %w", err)
	}

	res, err := client.ListUsersWithResponse(context.Background(), &huma.ListUsersParams{})
	if err != nil {
		return fmt.Errorf("failed to search for users: %w", err)
	}

	// Do something with res

	return nil
}
```

## Things to know

- We convert the original [Huma API] from OpenAPI 3.1 to OpenAPI 3.0 with
  OpenAPI Overlay.
- We alter a lot of the `operationId` in the original spec for readability in
  `overlay.yaml`

[Huma HRM]: https://humahr.com/
[huma api]: https://demo.openapi.humahr.com/swagger-ui/index.html
[oapi-codegen]: https://github.com/oapi-codegen/oapi-codegen
