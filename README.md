# <img src="https://golang.org/favicon.ico" width="24" alt="Go logo"> api-go/api (1.0.0)


A lightweight, dependency-free HTTP helper library for building Go APIs quickly and cleanly — with routing, structured responses, gzip streaming, and extensible middleware support.
## Features

- http.ServeMux-based routing
- Declarative route + method bindings (e.g., "GET /users")
- Middleware chaining
- Structured JSON responses
- GZIP + binary streaming support
- Predefined helpers for common status codes
- Auto-cleanup for streamed temp files
- Simple, idiomatic Go — no external dependencies

## Package Overview
### APIServer

The core HTTP server abstraction:

```go
type APIServer struct {
	Addr   string
	Server http.Server
	Router *http.ServeMux
}
```

#### Usage Example

```go
import "github.com/RecursionExcursion/api-go/api"

func main() {
	routes := []api.RouteHandler{
		{
			MethodAndPath: "GET /hello",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				api.Response.Ok(w, map[string]string{"message": "Hello, world!"})
			},
		},
	}

	server := api.NewApiServer(":8080", routes)
	server.ListenAndServe()
}
```

### RouteHandler

```go
type RouteHandler struct {
	MethodAndPath string
	Handler       HandlerFn
	Middleware    []Middleware
}
```

- MethodAndPath: Supports "METHOD /path" (e.g., "POST /users")
- Middleware: Executed left-to-right (wraps the handler)

Middleware signature:

```go
type Middleware = func(HandlerFn) HandlerFn
```

### Method Generator

Generate REST-style route strings dynamically:

```go
gen := api.HttpMethodGenerator("/users")
paths := gen("profile")

fmt.Println(paths.GET) // "GET /users/profile"
```

Helpful for grouping method routes under a shared base.

### Structured Responses

The api.Response object provides clean, consistent helpers for writing JSON responses:
#### Status Shortcuts

```go
api.Response.Ok(w, data)             // 200 OK
api.Response.Created(w, data)        // 201 Created
api.Response.NoContent(w)            // 204 No Content
api.Response.BadRequest(w, data)     // 400 Bad Request
api.Response.Unauthorized(w, data)   // 401 Unauthorized
api.Response.Forbidden(w, data)      // 403 Forbidden
api.Response.NotFound(w, data)       // 404 Not Found
api.Response.ServerError(w, data)    // 500 Internal Server Error
```

#### Custom Status

```go
api.Response.Send(w, 418, "I'm a teapot ☕")
```

#### GZIP Support

```go
api.Response.Gzip(w, 200, largePayload)
```

#### Binary File Streaming

```go
api.Response.StreamFile(w, 200, "/tmp/archive.zip", "download.zip")
```

Automatically sets headers and removes temp directory.

### Middleware Examples

Your own middleware functions can be chained per route:

```go
api.Middleware = func(next api.HandlerFn) api.HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request started")
		next(w, r)
		log.Println("Request finished")
	}
}
```

Or use built-in patterns:

- LoggerMW: Simple request logging
- KeyAuthMW: Header-based key authentication
- RecoveryMW: Catch and log panics
- TimeoutMW: Add per-request timeouts
- CustomKeyAuthMW(...): Supply custom validator, header, and bearer toggle

### Utilities
Parse Body to Type

```go
user, err := api.DecodeJSON[User](r)
```

Reads and unmarshals JSON into your struct.
 
📄 License

MIT © RecursionExcursion
