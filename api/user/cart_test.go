package user

import (
	"encoding/json"
	"errors"
	CONSTANT "merge-backend/constant"
	DB "merge-backend/database"
	INIT "merge-backend/init"
	MODEL "merge-backend/model"
	UTIL "merge-backend/util"
	"testing"
)

func TestCartUpdate(t *testing.T) {
	INIT.Init()

	// add test user
	userID, _, _ := DB.MainDB.InsertWithUniqueID(CONSTANT.UsersTable, map[string]string{
		"name":     "qwerty",
		"email":    "qwerty@gmail.com",
		"password": "65e84be33532fb784c48129675f9eff3a682b27168c0ea744b2cf58ee02337c5",
	}, "id")
	// delete that user
	defer DB.MainDB.DeleteSQL(CONSTANT.UsersTable, map[string]string{
		"id": userID,
	})

	// add test item
	itemID, _, _ := DB.MainDB.InsertWithUniqueID(CONSTANT.ItemsTable, map[string]string{
		"title": "qwerty",
		"stock": "13",
	}, "id")
	// delete that item
	defer DB.MainDB.DeleteSQL(CONSTANT.ItemsTable, map[string]string{
		"id": itemID,
	})

	testCases := []MODEL.Test{
		{
			Title:       "Update cart",
			Description: "Add item to cart",
			Method:      "PUT",
			URL:         "/admin/cart",
			Headers: map[string]interface{}{
				"apikey":  CONSTANT.AdminAPIKey,
				"user_id": userID,
			},
			Body: `{
				"item_id": "` + itemID + `",
				"type": "add"
			}`,
			PreRequest: func() {},
			Request:    CartUpdate,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status != CONSTANT.StatusCodeOk {
					return errors.New("Add to cart failed")
				}

				return nil
			},
		},
		{
			Title:       "Update cart",
			Description: "Remove item from cart",
			Method:      "PUT",
			URL:         "/admin/cart",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"item_id": "` + itemID + `",
				"type": "remove"
			}`,
			PreRequest: func() {},
			Request:    CartUpdate,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status != CONSTANT.StatusCodeOk {
					return errors.New("Remove from cart failed")
				}

				return nil
			},
		},
		{
			Title:       "Update cart",
			Description: "Add item with no stock to cart",
			Method:      "PUT",
			URL:         "/admin/cart",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"item_id": "` + itemID + `",
				"type": "add"
			}`,
			PreRequest: func() {
				DB.MainDB.UpdateSQL(CONSTANT.ItemsTable, map[string]string{
					"id": itemID,
				}, map[string]string{
					"stock": "0",
				})
			},
			Request: CartUpdate,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Add item with no stock to cart succeeded")
				}

				return nil
			},
		},
	}

	UTIL.TestUseCases(t, testCases)
}
