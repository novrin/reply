# reply

[![GoDoc](https://godoc.org/github.com/novrin/reply?status.svg)](https://pkg.go.dev/github.com/novrin/reply) 
![tests](https://github.com/novrin/reply/workflows/tests/badge.svg) 
[![Go Report Card](https://goreportcard.com/badge/github.com/novrin/reply)](https://goreportcard.com/report/github.com/novrin/reply)

`reply` is a Go HTTP response engine. It provides convenience methods for common HTTP responses.

### Installation

```shell
go get github.com/novrin/reply
``` 

## Usage

The example below, we use a reply Engine to compose replies to server requests in an application's handler functions.

```go
package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/yourusername/yourpackage/internal/database" 
	"github.com/novrin/reply"
)

// Application orchestrates replies to server requests.
type Application struct {
	db    *database.Queries
	reply reply.Engine // Use a reply Engine to write responses 
}

// NotFound renders the 404 template.
func (app *Application) NotFound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.reply.MethodNotAllowed(w, http.MethodGet)
		return
	}
	app.reply.Write(w, http.StatusNotFound, reply.Options{Template: "404.html"})
}

// Home renders the home template.
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.reply.MethodNotAllowed(w, http.MethodGet)
		return
	}
	if r.URL.Path != "/" {
		app.reply.NotFound(w)
		return
	}
	users, err := app.db.Users(r.Context())
	if err != nil {
		app.reply.InternalServerError(w, fmt.Errorf("failed to retrieve users: %s", err.Error()))
		return
	}
	app.reply.OK(w, reply.Options{
		Template: "home.html",
		Invoke:   "base",
		Data: struct { Users []database.Users }{ Users: users },
	})
}
```

## License

Copyright (c) 2023-present [novrin](https://github.com/novrin)

Licensed under [MIT License](./LICENSE)