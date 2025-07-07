package api

import "net/http"

// Core server & routing types
type APIServer struct {
	Addr   string
	Server http.Server
	Router *http.ServeMux
}

// Handler and Middleware abstractions
type HandlerFn = func(http.ResponseWriter, *http.Request)
type Middleware = func(HandlerFn) HandlerFn

type RouteHandler struct {
	MethodAndPath string
	Handler       HandlerFn
	Middleware    []Middleware
}

// HTTP method helper struct
type HTTPMethods struct {
	GET    string
	POST   string
	PUT    string
	PATCH  string
	DELETE string
}

// Response interface types
type response = func(http.ResponseWriter, ...any)
type customResponse = func(http.ResponseWriter, int, ...any)

type ApiResponses struct {
	Ok              response
	Created         response
	NoContent       response
	ServerError     response
	NotFound        response
	Unauthorized    response
	Forbidden       response
	TooManyRequests response
	BadRequest      response
	Send            customResponse
	Gzip            func(w http.ResponseWriter, status int, data ...any)
	StreamBytes     func(w http.ResponseWriter, status int, bytes []byte, name string)
	StreamFile      func(w http.ResponseWriter, status int, binPath string, name string)
}
