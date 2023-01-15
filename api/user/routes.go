package user

import "github.com/gorilla/mux"

// LoadUserRoutes - load all user routes with user prefix
func LoadUserRoutes(router *mux.Router) {
	userRoutes := router.PathPrefix("/user").Subrouter()

	// middlewares
	userRoutes.Use(APIKeyMiddleware)
	userRoutes.Use(CheckAuthToken)
	userRoutes.Use(CheckUserValid)

	// user
	userRoutes.HandleFunc("/cart", CartUpdate).Methods("PUT")

	// note
	userRoutes.HandleFunc("/item", ItemsGet).Methods("GET")

	// user
	userRoutes.HandleFunc("/login", UserLogin).Methods("POST")
	userRoutes.HandleFunc("/signup", UserSignUp).Methods("POST")
	userRoutes.HandleFunc("/refresh-token", UserRefreshToken).Methods("GET")

}
