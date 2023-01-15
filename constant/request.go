package constant

// required fields for api endpoints
var (
	// admin
	ItemAddAdminRequiredFields    = []string{"title", "stock"}
	ItemUpdateRequiredFields      = []string{"title", "stock"}
	UserLoginAdminRequiredFields  = []string{"email", "password"}
	UserSignUpAdminRequiredFields = []string{"name", "email", "password"}

	// user
	CartUpdateUserRequiredFields = []string{"item_id", "type"}
	UserLoginUserRequiredFields  = []string{"email", "password"}
	UserSignUpUserRequiredFields = []string{"name", "email", "password"}
)
