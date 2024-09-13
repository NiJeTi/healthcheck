# Healthcheck

Go package `healthcheck` provides a simple yet very convenient application health checking tools.

## Features

- Lightweight and easy to integrate
- Supports custom health check functions
- Has ready-to-run support for the 2 most popular Go HTTP servers:
  - `net/http` via `github.com/nijeti/healthcheck/servers/http`
  - `fasthttp` via `github.com/nijeti/healthcheck/servers/fasthttp`
- Ready to be integrated with any other server of choice

## Installation

### Library
```shell
go get github.com/nijeti/healthcheck
```

### `net/http` server
```shell
go get github.com/nijeti/healthcheck/servers/http
```

### `fasthttp` server
```shell
go get github.com/nijeti/healthcheck/servers/fasthttp
```

## Usage

Here is an example of how to use the healthcheck library with the standard HTTP server:

```go
package main

import (
	"context"
	"time"

	"github.com/nijeti/healthcheck"
	"github.com/nijeti/healthcheck/servers/http"
)

func main() {
	// Create a new health checker
	healthchecker := healthcheck.New(
		healthcheck.WithSimpleProbe(
			"database", func(ctx context.Context) error {
				// Perform your health check
				return nil
			},
		),
		healthcheck.WithTimeoutDegraded(1*time.Second),
		healthcheck.WithTimeoutUnhealthy(10*time.Second),
	)

	// Set up the HTTP server for health check endpoint 
	healthcheckServer := http.New(
		healthchecker,
		http.WithAddress(":8080"),
		http.WithRoute("/health"),
	)
	
	// Start the server
	healthcheckServer.Start()
	defer healthcheckServer.Stop()
}

```

## Examples

The following projects has successfully integrated Healthcheck:

- [cinema-keeper](https://github.com/NiJeTi/cinema-keeper)

## License

This project is licensed under the Apache-2.0 License.
