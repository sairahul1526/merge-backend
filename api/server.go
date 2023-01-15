package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// StartServer - start server using mux
func StartServer() {
	// ec2 router
	fmt.Println(http.ListenAndServe(":5000", &WithCORS{LoadRouter()}))
}

func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	// cors configuration
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	res.Header().Set("Access-Control-Allow-Headers", "*")

	if req.Method == "OPTIONS" {
		return
	}

	s.r.ServeHTTP(res, req)
}

// WithCORS .
type WithCORS struct {
	r *mux.Router
}
