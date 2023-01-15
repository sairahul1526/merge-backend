package user

import (
	"encoding/json"
	"errors"
	CONSTANT "merge-backend/constant"
	INIT "merge-backend/init"
	MODEL "merge-backend/model"
	UTIL "merge-backend/util"
	"testing"
)

func TestItemsGet(t *testing.T) {
	INIT.Init()
	testCases := []MODEL.Test{
		{
			Title:       "Get Items",
			Description: "Get all items",
			Method:      "GET",
			URL:         "/user/item",
			Headers: map[string]interface{}{
				"apikey": CONSTANT.UserAPIKey,
			},
			Body:       "",
			PreRequest: func() {},
			Request:    ItemsGet,
			PostRequest: func(resp []byte) error {
				// check if valid
				response := MODEL.Response{}
				err := json.Unmarshal(resp, &response)
				if err != nil {
					return err
				}

				if response.Meta.Status != CONSTANT.StatusCodeOk {
					return errors.New("Wrong status - " + response.Meta.Status)
				}

				return nil
			},
		},
	}

	UTIL.TestUseCases(t, testCases)
}
