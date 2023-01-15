package model

type Response struct {
	Data struct{} `json:"data"`
	Meta struct {
		DevMessage  string `json:"dev_message"`
		Message     string `json:"message"`
		MessageType string `json:"message_type"`
		Status      string `json:"status"`
	} `json:"meta"`
}
