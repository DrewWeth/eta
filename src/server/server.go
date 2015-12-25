package server

import (
	// "fmt"
	"github.com/gorilla/mux"

	"github.com/drewweth/eta/src/app/controllers/commentcontroller"
	"github.com/drewweth/eta/src/app/controllers/homecontroller"
	"github.com/drewweth/eta/src/app/controllers/postcontroller"
	"github.com/drewweth/eta/src/app/controllers/subcontroller"
	"net/http"
)

type Page struct {
	title string
}

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

	router.HandleFunc("/", homecontroller.IndexHandler)
	router.HandleFunc("/r", subcontroller.InsertSubHandler).Methods("POST") // Create sub

	router.HandleFunc("/r/{sub}", subcontroller.ShowHandler).Methods("GET")
	router.HandleFunc("/r/{sub}", postcontroller.InsertPostHandler).Methods("POST") // Create post

	router.HandleFunc("/r/{sub}/comments/{id}", postcontroller.ShowHandler)
	router.HandleFunc("/r/{sub}/comments", commentcontroller.PostHandler).Methods("POST") // Create comment



	// router.HandleFunc("/r/{sub}/comments/upvote", commentcontroller.PostHandler).Methods("POST") // Upvote a comment
	// router.HandleFunc("/r/{sub}/upvote", postcontroller.InsertPostHandler).Methods("POST") // Upvote a post

	router.HandleFunc("/createcache/{id}", commentcontroller.CreateCache)
	router.HandleFunc("/massinsert/{id}", commentcontroller.MassInsertHandler)

	return router
}
