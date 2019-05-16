package main

import (
	"net/http"

	"github.com/Haujilo/httpbin-go/api"
)

func route(mux *http.ServeMux) {
	patterns := map[string]func(w http.ResponseWriter, r *http.Request){
		"/headers": api.HeadersHander,
	}

	for endpoint, hander := range patterns {
		mux.HandleFunc(endpoint, hander)
	}
}
