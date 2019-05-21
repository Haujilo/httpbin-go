package main

import (
	"net/http"

	"github.com/Haujilo/httpbin-go/api"
)

func route(mux *http.ServeMux) {
	patterns := map[string]func(w http.ResponseWriter, r *http.Request){
		"/get":        api.GETHandler,
		"/post":       api.POSTHandler,
		"/headers":    api.HeadersHander,
		"/ip":         api.IPHander,
		"/user-agent": api.UserAgentHander,
	}

	for endpoint, hander := range patterns {
		mux.HandleFunc(endpoint, hander)
	}
}
