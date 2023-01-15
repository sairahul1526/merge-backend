package util

import (
	"encoding/json"
	CONSTANT "merge-backend/constant"
	"net/http"
)

// SetReponse - set request response with status, message etc
func SetReponse(w http.ResponseWriter, status, msg, msgType, devMessage string, resp map[string]interface{}) {
	w.Header().Set("Status", "200")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{}
	response["data"] = resp
	response["meta"] = setMeta(status, msg, msgType, devMessage)
	json.NewEncoder(w).Encode(response)
}

// you will always have meta in any response (GET/POST/PUT/PATCH/DELETE)
// status - HTTP status codes like 200,201,400,500,503
// message - Any message which would be used by app to display or take action
// message_type - 1 : show dialog, 2 : show toast, else nothing
// dev_message - Useful for developers to debug
func setMeta(status, message, msgType, devMessage string) map[string]string {
	if len(message) == 0 {
		if status == CONSTANT.StatusCodeBadRequest {
			message = "Bad Request"
		} else if status == CONSTANT.StatusCodeServerError {
			message = "Server Error"
		}
	}
	return map[string]string{
		"status":       status,
		"message":      message,
		"message_type": msgType,
		"dev_message":  devMessage,
	}
}
