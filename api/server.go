package api

import (
	"log"
	"net/http"
)

/* addr = :PORT */
func NewApiServer(addr string, routes []RouteHandler) *APIServer {
	router := http.NewServeMux()

	for _, r := range routes {
		router.HandleFunc(r.handleHttp())
	}

	return &APIServer{
		Addr:   addr,
		Router: router,
		Server: http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (s *APIServer) ListenAndServe() error {

	log.Printf("Server is listening on %s", s.Addr)
	return s.Server.ListenAndServe()
}
