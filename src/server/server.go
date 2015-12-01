package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Server serves while allowing cross origin access.
type Server struct {
	r *mux.Router
}

// NewEaseServer creates a new handler for Ease.
func NewServer() *Server {
	return &Server{r: createRoutingMux()}
}

// ServeHTTP serves requests from the EaseServer's mux while allowing
// cross origin access.
func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}

	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

// createRoutingMux sets up the routing for the server.
func createRoutingMux() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the Eta!")
	})

	return router
}
