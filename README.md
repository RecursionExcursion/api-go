# 🌐 go-toolkit/api

A lightweight HTTP server and response helper package for quickly bootstrapping Go APIs with clean routing, structured responses, gzip support, and middleware chaining.

---

## 🚀 Features

- Custom router based on `http.ServeMux`
- Middleware chaining
- Route + method binding
- Structured JSON responses
- GZIP responses and binary streaming support
- Named helper functions for common HTTP statuses

---

## 📦 Package Overview

### 🔧 Server: `APIServer`

```go
type APIServer struct {
	Addr        string
	Server      http.Server
	Router      *http.ServeMux
	initalized  bool
}
```

✅ Usage

``` go
import "github.com/RecursionExcursion/go-toolkit/api"

func main() {
	server := api.NewApiServer(":8080")

	routes := []api.RouteHandler{
		{
			MethodAndPath: "GET /hello",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				api.Response.Ok(w, map[string]string{"message": "Hello, world!"})
			},
		},
	}

	server.Init(routes)
	server.ListenAndServe()
}
```

🧱 RouteHandler

``` go
type RouteHandler struct {
	MethodAndPath string
	Handler       HandlerFn
	Middleware    []Middleware
}
```

Supports middleware chaining via:

``` go
type Middleware = func(HandlerFn) HandlerFn
```
Middleware are executed left-to-right (top-down).
🧰 HTTP Method Generator

``` go
api.HttpMethodGenerator("/base")("users") // returns HTTPMethods{GET: "GET /base/users", ...}
``` 
Use it to cleanly assign route strings.
📤 Structured Responses: api.Response

Built-in helpers for consistent API responses:
✅ Status Shortcuts

``` go
api.Response.Ok(w, data)
api.Response.Created(w, data)
api.Response.BadRequest(w, data)
api.Response.Unauthorized(w, data)
api.Response.ServerError(w, data)
...
```

🔄 Custom Send

``` go
api.Response.Send(w, 418, "I'm a teapot ☕")
```

📦 GZIP Support

``` go
api.Response.Gzip(w, 200, data)
```

📁 File Streaming

``` go
api.Response.StreamFile(w, 200, "/tmp/file.zip", "download.zip")
```

🧼 Cleanup and Logging

    Temp files streamed with StreamFile are deleted automatically

    Stream errors are logged

    JSON errors result in 500 Internal Server Error

📄 License

MIT © RecursionExcursion
