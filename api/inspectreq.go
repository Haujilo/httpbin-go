package api

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

func fmtHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	for k, v := range r.Header {
		headers[k] = strings.Join(v, ",")
	}
	return headers
}

func getIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func HeadersHander(w http.ResponseWriter, r *http.Request) {
	type JSON struct {
		Headers map[string]string `json:"headers"`
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSON{Headers: fmtHeaders(r)})
}

func IPHander(w http.ResponseWriter, r *http.Request) {
	type JSON struct {
		Origin string `json:"origin"`
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSON{Origin: getIP(r)})
}

func UserAgentHander(w http.ResponseWriter, r *http.Request) {
	type JSON struct {
		UserAgent string `json:"user-agent"`
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSON{UserAgent: r.Header.Get("User-Agent")})
}
