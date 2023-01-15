package constant

// server status codes
const (
	StatusCodeOk             = "200"
	StatusCodeCreated        = "201"
	StatusCodeBadRequest     = "400"
	StatusCodeForbidden      = "403"
	StatusCodeSessionExpired = "440"
	StatusCodeServerError    = "500"
	StatusCodeDuplicateEntry = "1000"
)

// type of alerts for frontend to show
const (
	NoDialog   = "0"
	ShowDialog = "1"
	ShowToast  = "2"
)

// user status
const (
	UserActive  = "1"
	UserBlocked = "2"
)

// user role
const (
	UserNormal = "1"
	UserAdmin  = "2"
)

// item status
const (
	ItemActive  = "1"
	ItemDeleted = "2"
)
