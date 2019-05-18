package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

func HeadersHander(w http.ResponseWriter, r *http.Request) {
	type JSON struct {
		Headers map[string]string `json:"headers"`
	}
	headers := make(map[string]string)
	for k, v := range r.Header {
		headers[k] = strings.Join(v, ",")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSON{Headers: headers})
}
