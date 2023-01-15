package admin

import "github.com/gorilla/mux"

// LoadAdminRoutes - load all admin routes with admin prefix
func LoadAdminRoutes(router *mux.Router) {
	adminRoutes := router.PathPrefix("/admin").Subrouter()

	// middlewares
	adminRoutes.Use(APIKeyMiddleware)
	adminRoutes.Use(CheckAuthToken)
	adminRoutes.Use(CheckUserValid)

	// note
	adminRoutes.HandleFunc("/item", ItemAdd).Methods("POST")
	adminRoutes.HandleFunc("/item", ItemUpdate).Queries(
		"item_id", "{item_id}",
	).Methods("PUT")

	// user
	adminRoutes.HandleFunc("/user-maintain", UserMaintain).Queries(
		"user_id", "{user_id}",
	).Methods("PUT")
	adminRoutes.HandleFunc("/login", UserLogin).Methods("POST")
	adminRoutes.HandleFunc("/signup", UserSignUp).Methods("POST")
	adminRoutes.HandleFunc("/refresh-token", UserRefreshToken).Methods("GET")

}
