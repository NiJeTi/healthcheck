# Healthcheck

Go package `healthcheck` provides a simple yet very convenient application health checking tools.

## Features

- Lightweight and easy to integrate
- Supports custom health check functions
- Ready to be integrated with any HTTP-server of choice

## Installation

```sh
go get github.com/nijeti/healthcheck
```

## Usage

Here is an example of how to use the healthcheck library:

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/nijeti/healthcheck"
)

func main() {
	// Create a new health checker
	hc := healthcheck.New(
		healthcheck.WithSimpleProbe(
			"database", func(ctx context.Context) error {
				// Perform your health check
				return nil
			},
		),
	)

	// Set up the HTTP handler for the health check endpoint
	http.HandleFunc(
		"/health", func(writer http.ResponseWriter, request *http.Request) {
			// Check all registered probes
			status := hc.Handle(request.Context())
			
			_, err := writer.Write([]byte(status.String()))
			if err != nil {
				log.Println("failed to write response:", err)
			}
		},
	)

	// Start the server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
```

## Examples

The following projects has successfully integrated Healthcheck:

- [cinema-keeper](https://github.com/NiJeTi/cinema-keeper)

## License

This project is licensed under the Apache-2.0 License.
