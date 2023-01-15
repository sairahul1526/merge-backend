package util

import (
	"encoding/json"
	"io/ioutil"
	"math"
	LOGGER "merge-backend/logger"
	"net/http"
	"strconv"
)

func GetNumberOfPages(count string, countPerPage int) int {
	ctn, _ := strconv.Atoi(count)
	return int(math.Ceil(float64(ctn) / float64(countPerPage)))
}

// GetPageNumber - get page number for pagination
func GetPageNumber(pageStr string) int {
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}
	return page
}

// RequiredFiledsCheck - check if all required fields are present
func RequiredFiledsCheck(body map[string]string, required []string) string {
	for _, field := range required {
		if len(body[field]) == 0 {
			return field
		}
	}
	return ""
}

// ReadRequestBodyToMap - read raw body from request
func ReadRequestBodyToMap(r *http.Request) (map[string]string, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LOGGER.Warn("ReadRequestBodyToMap", err)
		return map[string]string{}, err
	}
	defer r.Body.Close()

	body := map[string]string{}

	err = json.Unmarshal(b, &body)
	if err != nil {
		LOGGER.Warn("ReadRequestBodyToMap", err)
		return map[string]string{}, err
	}

	return body, nil
}
