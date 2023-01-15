package model

import "net/http"

// Test - model for testing
type Test struct {
	Title       string
	Description string
	Method      string
	URL         string
	Headers     map[string]interface{}
	Body        string
	PreRequest  func()
	Request     func(w http.ResponseWriter, r *http.Request)
	PostRequest func(resp []byte) error
}
