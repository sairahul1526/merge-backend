package util

import (
	"fmt"
	MODEL "merge-backend/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUseCases(t *testing.T, testCases []MODEL.Test) {
	for _, testCase := range testCases {
		fmt.Println(testCase.Title, "-", testCase.Description)

		// setup before testing
		testCase.PreRequest()

		// request
		payload := strings.NewReader(testCase.Body)

		req, err := http.NewRequest(testCase.Method, testCase.URL, payload)
		if err != nil {
			t.Errorf("Error while testing %s, Error - %s", testCase.URL, err.Error())
			return
		}

		response := httptest.NewRecorder()

		testCase.Request(response, req)

		// validate and clear any data after testing
		err = testCase.PostRequest(response.Body.Bytes())
		if err != nil {
			t.Errorf("Error while testing %s, Response - %s, Error - %s", testCase.URL, response.Body.String(), err.Error())
			return
		}
	}
}
