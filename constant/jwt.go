package constant

const (
	// admin
	AdminJWTRefreshExpiry = 1440 // jwt refresh token expiry in min // 1 day
	AdminJWTAccessExpiry  = 60   // jwt access token expiry in min // 1 hour

	// user
	UserJWTRefreshExpiry = 43200 // jwt refresh token expiry in min // 1 month
	UserJWTAccessExpiry  = 60    // jwt access token expiry in min // 1 hour
)
