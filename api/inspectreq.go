package api

import (
	"encoding/json"
	"net/http"
)

func HeadersHander(w http.ResponseWriter, r *http.Request) {
	type JSON struct {
		Headers http.Header `json:"headers"`
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSON{Headers: r.Header})
}
