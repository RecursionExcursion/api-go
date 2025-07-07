package api

import (
	"fmt"
	"net/http"
	"strings"
)

type PathBuilder struct {
	base string
}

func NewPathBuilder(base string) PathBuilder {
	if !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	cleanBase := strings.TrimSuffix(base, "/")
	return PathBuilder{base: cleanBase}
}

func (pb *PathBuilder) Append(parts ...string) PathBuilder {
	joined := strings.Join(parts, "/")
	joined = strings.Trim(joined, "/")
	newPath := pb.base + "/" + joined
	return PathBuilder{base: newPath}
}

func (pb *PathBuilder) Methods() HTTPMethods {

	assign := func(method string) string {
		return fmt.Sprintf("%v %v", method, pb.base)
	}

	return HTTPMethods{
		GET:    assign("GET"),
		POST:   assign("POST"),
		PUT:    assign("PUT"),
		PATCH:  assign("PATCH"),
		DELETE: assign("DELETE"),
	}
}

func (rh *RouteHandler) handleHttp() (string, func(http.ResponseWriter, *http.Request)) {
	if rh.Handler == nil {
		panic("handler is nil for route " + rh.MethodAndPath)
	}
	return rh.MethodAndPath, pipeMiddleware(rh.Middleware...)(rh.Handler)
}

// pipeMiddleware passes request through middleware fns from left to right.
// Returns a single middleware that wraps the final handler.
func pipeMiddleware(mws ...Middleware) Middleware {
	return func(final HandlerFn) HandlerFn {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}
