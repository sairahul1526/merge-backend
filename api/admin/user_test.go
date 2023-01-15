package admin

import (
	"encoding/json"
	"errors"
	CONSTANT "merge-backend/constant"
	DB "merge-backend/database"
	INIT "merge-backend/init"
	MODEL "merge-backend/model"
	UTIL "merge-backend/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserSignUp(t *testing.T) {
	INIT.Init()
	testCases := []MODEL.Test{
		{
			Title:       "User Sign up",
			Description: "Signup with email, first time",
			Method:      "POST",
			URL:         "/admin/signup",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"name": "qwerty",
				"email": "qwerty@gmail.com",
				"password": "qwerty"
			}`,
			PreRequest: func() {},
			Request:    UserSignUp,
			PostRequest: func(resp []byte) error {
				// delete that user
				defer DB.MainDB.DeleteSQL(CONSTANT.UsersTable, map[string]string{
					"email": "qwerty@gmail.com",
				})

				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status != CONSTANT.StatusCodeOk {
					return errors.New("Sign up failed for the first time")
				}

				return nil
			},
		},
		{
			Title:       "User Sign up",
			Description: "Signup with same email, second time",
			Method:      "POST",
			URL:         "/admin/signup",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"name": "qwerty",
				"email": "qwerty@gmail.com",
				"password": "qwerty"
			}`,
			PreRequest: func() {
				// add user
				req, _ := http.NewRequest("POST", "/admin/signup", strings.NewReader(`{
					"name": "qwerty",
					"email": "qwerty@gmail.com",
					"password": "qwerty"
				}`))
				UserSignUp(httptest.NewRecorder(), req)
			},
			Request: UserSignUp,
			PostRequest: func(resp []byte) error {
				// delete that user
				defer DB.MainDB.DeleteSQL(CONSTANT.UsersTable, map[string]string{
					"email": "qwerty@gmail.com",
				})

				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Signup succeeded for duplicate email")
				}

				return nil
			},
		},
		{
			Title:       "User Sign up",
			Description: "Signup with wrong email",
			Method:      "POST",
			URL:         "/admin/signup",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"name": "qwerty",
				"email": "qwerty",
				"password": "qwerty"
			}`,
			PreRequest: func() {},
			Request:    UserSignUp,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Signup succeeded for wrong email")
				}

				return nil
			},
		},
		{
			Title:       "User Sign up",
			Description: "Signup with no email",
			Method:      "POST",
			URL:         "/admin/signup",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"name": "qwerty",
				"password": "qwerty"
			}`,
			PreRequest: func() {},
			Request:    UserSignUp,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Signup succeeded with no email")
				}

				return nil
			},
		},
	}

	UTIL.TestUseCases(t, testCases)
}

func TestUserLogin(t *testing.T) {
	INIT.Init()

	// add test user
	userID, _, _ := DB.MainDB.InsertWithUniqueID(CONSTANT.UsersTable, map[string]string{
		"name":     "qwerty",
		"email":    "qwerty@gmail.com",
		"password": "65e84be33532fb784c48129675f9eff3a682b27168c0ea744b2cf58ee02337c5",
		"role":     CONSTANT.UserAdmin,
	}, "id")
	// delete that user
	defer DB.MainDB.DeleteSQL(CONSTANT.UsersTable, map[string]string{
		"id": userID,
	})

	testCases := []MODEL.Test{
		{
			Title:       "User Login",
			Description: "Login with email, not signed up",
			Method:      "POST",
			URL:         "/admin/login",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"email": "qwerty2@gmail.com",
				"password": "qwerty"
			}`,
			PreRequest: func() {},
			Request:    UserLogin,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Login succeeded for email, which is not signed up")
				}

				return nil
			},
		},
		{
			Title:       "User Login",
			Description: "Login with email, already signed up",
			Method:      "POST",
			URL:         "/admin/login",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"email": "qwerty@gmail.com",
				"password": "qwerty"
			}`,
			PreRequest: func() {},
			Request:    UserLogin,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status != CONSTANT.StatusCodeOk {
					return errors.New("Login failed for email, which is already signed up")
				}

				return nil
			},
		},
		{
			Title:       "User Login",
			Description: "Login with normal role user",
			Method:      "POST",
			URL:         "/admin/login",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"email": "qwerty@gmail.com",
				"password": "qwerty"
			}`,
			PreRequest: func() {
				DB.MainDB.UpdateSQL(CONSTANT.UsersTable, map[string]string{
					"id": userID,
				}, map[string]string{
					"role": CONSTANT.UserNormal,
				})
			},
			Request: UserLogin,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Login succeeded for normal role user")
				}

				return nil
			},
		},
		{
			Title:       "User Login",
			Description: "Login with blocked user",
			Method:      "POST",
			URL:         "/admin/login",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"email": "qwerty@gmail.com",
				"password": "qwerty"
			}`,
			PreRequest: func() {
				DB.MainDB.UpdateSQL(CONSTANT.UsersTable, map[string]string{
					"id": userID,
				}, map[string]string{
					"status": CONSTANT.UserBlocked,
				})
			},
			Request: UserLogin,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Login succeeded for blocked user")
				}

				return nil
			},
		},
	}

	UTIL.TestUseCases(t, testCases)
}
