package admin

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

func TestItemAdd(t *testing.T) {
	INIT.Init()
	testCases := []MODEL.Test{
		{
			Title:       "Add item",
			Description: "Add new item",
			Method:      "POST",
			URL:         "/admin/item",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"title": "qwerty",
				"stock": "30"
			}`,
			PreRequest: func() {},
			Request:    ItemAdd,
			PostRequest: func(resp []byte) error {
				// delete that item
				defer DB.MainDB.DeleteSQL(CONSTANT.ItemsTable, map[string]string{
					"title": "qwerty",
				})

				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status != CONSTANT.StatusCodeOk {
					return errors.New("Add item failed")
				}

				return nil
			},
		},
		{
			Title:       "Add item",
			Description: "Add item with empty title",
			Method:      "POST",
			URL:         "/admin/item",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"title": "",
				"stock": "30"
			}`,
			PreRequest: func() {},
			Request:    ItemAdd,
			PostRequest: func(resp []byte) error {
				// delete that item
				defer DB.MainDB.DeleteSQL(CONSTANT.ItemsTable, map[string]string{
					"title": "qwerty",
				})

				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Add item succeeded even with empty title")
				}

				return nil
			},
		},
		{
			Title:       "Add item",
			Description: "Add item with negative stock",
			Method:      "POST",
			URL:         "/admin/item",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"title": "qwerty",
				"stock": "-1"
			}`,
			PreRequest: func() {},
			Request:    ItemAdd,
			PostRequest: func(resp []byte) error {
				// delete that item
				defer DB.MainDB.DeleteSQL(CONSTANT.ItemsTable, map[string]string{
					"title": "qwerty",
				})

				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Add item succeeded even with negative stock")
				}

				return nil
			},
		},
	}

	UTIL.TestUseCases(t, testCases)
}

func TestItemUpdate(t *testing.T) {
	INIT.Init()

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
			Title:       "Update item",
			Description: "Update item",
			Method:      "PUT",
			URL:         "/admin/item?item_id=" + itemID,
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"title": "qwerty 2",
				"stock": "13"
			}`,
			PreRequest: func() {},
			Request:    ItemUpdate,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status != CONSTANT.StatusCodeOk {
					return errors.New("Update item failed")
				}

				return nil
			},
		},
		{
			Title:       "Update item",
			Description: "Update item with empty title",
			Method:      "PUT",
			URL:         "/admin/item?item_id=" + itemID,
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"title": "",
				"stock": "30"
			}`,
			PreRequest: func() {},
			Request:    ItemUpdate,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Update item succeeded even with empty title")
				}

				return nil
			},
		},
		{
			Title:       "Update item",
			Description: "Update item with negative stock",
			Method:      "PUT",
			URL:         "/admin/item?item_id=" + itemID,
			Headers: map[string]interface{}{
				"apikey": CONSTANT.AdminAPIKey,
			},
			Body: `{
				"title": "qwerty",
				"stock": "-1"
			}`,
			PreRequest: func() {},
			Request:    ItemUpdate,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status == CONSTANT.StatusCodeOk {
					return errors.New("Update item succeeded even with negative stock")
				}

				return nil
			},
		},
	}

	UTIL.TestUseCases(t, testCases)
}
