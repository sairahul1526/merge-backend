package api

import (
	"encoding/json"
	"net/http"

	AdminAPI "merge-backend/api/admin"
	UserAPI "merge-backend/api/user"

	"github.com/gorilla/mux"
)

// HealthCheck .
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// for load balancer/beanstalk to know whether server/ec2 is healthy
	json.NewEncoder(w).Encode("ok")
}

// LoadRouter - get mux router with all the routes
func LoadRouter() *mux.Router {
	router := mux.NewRouter()

	AdminAPI.LoadAdminRoutes(router)
	UserAPI.LoadUserRoutes(router)

	router.Path("/").HandlerFunc(HealthCheck).Methods("GET")

	return router
}
