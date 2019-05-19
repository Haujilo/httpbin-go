package api

import (
	"encoding/json"
	"net"
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

func IPHander(w http.ResponseWriter, r *http.Request) {
	type JSON struct {
		Origin string `json:"origin"`
	}
	w.Header().Set("Content-Type", "application/json")
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	json.NewEncoder(w).Encode(JSON{Origin: ip})
}

func UserAgentHander(w http.ResponseWriter, r *http.Request) {
	type JSON struct {
		UserAgent string `json:"user-agent"`
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSON{UserAgent: r.Header.Get("User-Agent")})
}
