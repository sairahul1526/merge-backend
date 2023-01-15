package admin

import (
	CONSTANT "merge-backend/constant"
	DB "merge-backend/database"
	UTIL "merge-backend/util"
	"net/http"
	"strings"
)

// APIKeyMiddleware - check if api key is admin's
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check for admin apikey
		if !strings.EqualFold(r.Header.Get("apikey"), CONSTANT.AdminAPIKey) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeForbidden, CONSTANT.AdminAPIKeyInvalidMessage, CONSTANT.ShowDialog, "", map[string]interface{}{})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CheckAuthToken - verify access, refresh token and expiry
func CheckAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/login") || strings.Contains(r.URL.Path, "/signup") {
			// pass through
		} else if strings.Contains(r.URL.Path, "/refresh-token") {
			if UTIL.IsAccessToken(r.Header.Get("Authorization")) == nil {
				// if token is access type
				UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredRefreshMessage, CONSTANT.NoDialog, "", map[string]interface{}{})
				return
			} else {
				// if token is refresh type
				data, err := UTIL.ParseJWTToken(r.Header.Get("Authorization"))
				if err != nil {
					UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredMessage, CONSTANT.ShowDialog, "", map[string]interface{}{})
					return
				}

				// check if token has admin access
				if !strings.EqualFold(data["role"].(string), CONSTANT.UserAdmin) {
					UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.UserNotAdminMessage, CONSTANT.ShowDialog, "", map[string]interface{}{})
					return
				}

				// set user_id, company_id to header for further access
				r.Header.Set("user_id", data["user_id"].(string))
			}
		} else {
			// for all the other endpoints, other than login, signup, refresh
			// check if jwt token is access type and is valid, not expired
			if UTIL.IsAccessToken(r.Header.Get("Authorization")) == nil {
				// if token is access type
				data, err := UTIL.ParseJWTToken(r.Header.Get("Authorization"))
				if err != nil {
					UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredRefreshMessage, CONSTANT.NoDialog, "", map[string]interface{}{})
					return
				}

				// check if token has admin access
				if !strings.EqualFold(data["role"].(string), CONSTANT.UserAdmin) {
					UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.UserNotAdminMessage, CONSTANT.ShowDialog, "", map[string]interface{}{})
					return
				}

				// set user_id, company_id to header for further access
				r.Header.Set("user_id", data["user_id"].(string))
			} else {
				UTIL.SetReponse(w, CONSTANT.StatusCodeSessionExpired, CONSTANT.SessionExpiredRefreshMessage, CONSTANT.NoDialog, "", map[string]interface{}{})
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// CheckUserValid - check if user is valid
func CheckUserValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if len(r.Header.Get("user_id")) > 0 {
			// check if user id is valid
			user, err := DB.MainDB.SelectSQL(CONSTANT.UsersTable, []string{"status", "role"}, map[string]string{"id": r.Header.Get("user_id")})
			if err != nil {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), map[string]interface{}{})
				return
			}
			if len(user) == 0 {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotExistMessage, CONSTANT.ShowDialog, "", map[string]interface{}{})
				return
			}
			if !strings.EqualFold(user[0]["role"], CONSTANT.UserAdmin) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotAdminMessage, CONSTANT.ShowDialog, "", map[string]interface{}{})
				return
			}
			if !strings.EqualFold(user[0]["status"], CONSTANT.UserActive) {
				UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotAllowedMessage, CONSTANT.ShowDialog, "", map[string]interface{}{})
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
