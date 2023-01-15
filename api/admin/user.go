package admin

import (
	CONSTANT "merge-backend/constant"
	DB "merge-backend/database"
	"net/http"
	"strings"

	UTIL "merge-backend/util"
)

// UserMaintain - block/unblock user
func UserMaintain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	_, err = DB.MainDB.UpdateSQL(CONSTANT.UsersTable, map[string]string{
		"id": r.FormValue("user_id"),
	}, map[string]string{
		"status":     body["status"],
		"updated_by": r.Header.Get("user_id"),
	})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.NoDialog, "", response)
}

// UserLogin - login with email, password
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.UserLoginAdminRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, "", response)
		return
	}

	// check if user is valid and get user details
	user, err := DB.MainDB.SelectSQL(CONSTANT.UsersTable, []string{"id", "name", "email", "role", "status"}, map[string]string{"email": body["email"], "password": UTIL.GetMD5HashString(body["password"])})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}
	if len(user) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectCredentialsExistMessage, CONSTANT.ShowDialog, "", response)
		return
	}
	if !strings.EqualFold(user[0]["role"], CONSTANT.UserAdmin) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotAdminMessage, CONSTANT.ShowDialog, "", response)
		return
	}
	if !strings.EqualFold(user[0]["status"], CONSTANT.UserActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotAllowedMessage, CONSTANT.ShowDialog, "", response)
		return
	}

	// generate access and refresh token
	// access token - jwt token with short expiry added in header for authorization
	// refresh token - jwt token with long expiry to get new access token if expired
	// if refresh token expired, need to login
	accessToken, err := UTIL.CreateJWTToken(map[string]interface{}{"user_id": user[0]["id"], "role": CONSTANT.UserAdmin}, CONSTANT.AdminJWTAccessExpiry, false)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}
	refreshToken, err := UTIL.CreateJWTToken(map[string]interface{}{"user_id": user[0]["id"], "role": CONSTANT.UserAdmin}, CONSTANT.AdminJWTRefreshExpiry, true)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	response["user"] = user[0]
	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.NoDialog, "", response)
}

// UserSignUp - signup using name, email, password
func UserSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.UserSignUpAdminRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, "", response)
		return
	}

	// check if email is valid, based on regex
	if !UTIL.IsEmailValid(body["email"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UseValidEmailMessage, CONSTANT.ShowDialog, "", response)
		return
	}

	// check if email already exists
	if DB.MainDB.CheckIfExists(CONSTANT.UsersTable, map[string]string{"email": body["email"]}) == nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EmailExistMessage, CONSTANT.ShowDialog, "", response)
		return
	}

	// add user data and get user id
	userID, _, err := DB.MainDB.InsertWithUniqueID(CONSTANT.UsersTable, map[string]string{
		"password": UTIL.GetMD5HashString(body["password"]),
		"status":   CONSTANT.UserActive,
		"name":     body["name"],
		"email":    body["email"],
		"role":     CONSTANT.UserAdmin,
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	// generate access and refresh token
	// access token - jwt token with short expiry added in header for authorization
	// refresh token - jwt token with long expiry to get new access token if expired
	// if refresh token expired, need to login
	accessToken, err := UTIL.CreateJWTToken(map[string]interface{}{"user_id": userID, "role": CONSTANT.UserAdmin}, CONSTANT.AdminJWTAccessExpiry, false)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}
	refreshToken, err := UTIL.CreateJWTToken(map[string]interface{}{"user_id": userID, "role": CONSTANT.UserAdmin}, CONSTANT.AdminJWTRefreshExpiry, true)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	response["user_id"] = userID
	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.NoDialog, "", response)
}

// UserRefreshToken - refresh access token
func UserRefreshToken(w http.ResponseWriter, r *http.Request) {

	var response = make(map[string]interface{})

	// refresh token is already checked in middleware
	// generate new access token
	accessToken, err := UTIL.CreateJWTToken(map[string]interface{}{"user_id": r.Header.Get("user_id"), "role": CONSTANT.UserAdmin}, CONSTANT.AdminJWTAccessExpiry, false)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, err.Error(), response)
		return
	}

	response["access_token"] = accessToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.NoDialog, "", response)
}
