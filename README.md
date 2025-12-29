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

## Usage

```go
package main

import (
	"net/http"

	"github.com/throskam/ki"
	"github.com/throskam/ki/middlewares"
)

func main() {
	router := ki.NewRouter()

	router.Use(middlewares.RequestID())
	router.Use(middlewares.RealIP())
	router.Use(middlewares.RequestLogger())
	router.Use(middlewares.Recoverer())

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
	})

	_ = http.ListenAndServe(":8080", router)
}


```

see [examples](./examples) for examples.

## License

MIT
