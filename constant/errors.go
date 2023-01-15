package constant

import "errors"

var (
	JWTUnexpectedSigningMethodError = errors.New("Unexpected signing method")
	JWTInvalidTokenError            = errors.New("Invalid Token")
	JWTInvalidTokenExpiryError      = errors.New("Invalid Token Expiry")
	JWTRefreshTokenError            = errors.New("Refresh Token Used")
	JWTTokenExpiredError            = errors.New("Token Expired")
	JWTNotAcsessTokenError          = errors.New("Not an access token")
	JWTAcsessTokenError             = errors.New("Access Token Used")

	SQLCheckIfExistsEmptyError  = errors.New("No data available from specified filters")
	SQLInsertBodyEmptyError     = errors.New("Insert body is empty")
	SQLUpdateBodyEmptyError     = errors.New("Update body is empty")
	SQLDeleteAllNotAllowedError = errors.New("Delete all not allowed")
)
