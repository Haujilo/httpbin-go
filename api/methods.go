package api

import (
	"encoding/json"
	"net/http"
)

func fmtQueryString(r *http.Request) map[string]interface{} {
	queryString := make(map[string]interface{})
	for k, v := range r.URL.Query() {
		if len(v) == 1 {
			queryString[k] = v[0]
		} else {
			queryString[k] = v
		}
	}
	return queryString
}

func getFullURL(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	return scheme + r.Host + r.RequestURI
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	type JSON struct {
		Args    map[string]interface{} `json:"args"`
		Headers map[string]string      `json:"headers"`
		Origin  string                 `json:"origin"`
		URL     string                 `json:"url"`
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSON{
		Args:    fmtQueryString(r),
		Headers: fmtHeaders(r),
		Origin:  getIP(r),
		URL:     getFullURL(r),
	})
}
