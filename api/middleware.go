package api

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

/* Initalize logger to dictate logging action */
func LoggerMW(logger *log.Logger) Middleware {
	return func(next HandlerFn) HandlerFn {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next(w, r)
			logger.Printf("%v %v accessed %v in %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
		}
	}
}

func StaticBearerAuthMW(expectedKey string) Middleware {
	return func(next HandlerFn) HandlerFn {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := strings.TrimSpace(r.Header.Get(("Authorization")))
			fields := strings.Fields(auth)

			if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" || fields[1] != expectedKey {
				Response.Unauthorized(w, "Invalid token")
				return
			}

			next(w, r)
		}
	}
}

func HeaderAuthMW(validator func(token string) bool, header string, requireBearer bool) Middleware {
	return func(next HandlerFn) HandlerFn {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := strings.TrimSpace(r.Header.Get(header))
			fields := strings.Fields(auth)

			if requireBearer {
				if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" || !validator(fields[1]) {
					Response.Unauthorized(w, "Invalid token")
					return
				}
			} else {
				if len(fields) == 0 || !validator(fields[0]) {
					Response.Unauthorized(w, "Invalid token")
					return
				}
			}

			next(w, r)
		}
	}
}

// Place on the front of the mw pipe to handle all panics
func RecoveryMW(next HandlerFn) HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[panic] %v", err)
				Response.ServerError(w, "Internal server error")
			}
		}()
		next(w, r)
	}
}

func TimeoutMW(d time.Duration) Middleware {
	return func(next HandlerFn) HandlerFn {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			next(w, r.WithContext(ctx))
		}
	}
}

func CORSMW(origin string) Middleware {
	return func(next HandlerFn) HandlerFn {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next(w, r)
		}
	}
}
