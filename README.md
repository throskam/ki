# ki

Thin wrapper around the Go standard Mux.

## Installation

```bash
go get -u github.com/throskam/ki
```

## Features

- Compatible with the Go standard Mux
- Method based routing
- Middlewares (see [middlewares](./middlewares))
- Sub-routers
- Groups
- Named routes

## Examples

Very simple example:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/throskam/ki"
	"github.com/throskam/ki/middlewares"
)

func main() {
	router := ki.NewRouter()

	router.Use(middlewares.Logger())

	router.Get("/{$}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello world"))
	})

	_ = http.ListenAndServe(":8080", router)
}

```

see [examples](./examples) for more complex examples.
