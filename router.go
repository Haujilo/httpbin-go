package main

import (
	"net/http"

	"github.com/Haujilo/httpbin-go/api"
)

func route(mux *http.ServeMux) {
	patterns := map[string]func(w http.ResponseWriter, r *http.Request){
		"/delete":      api.DELETEHandler,
		"/get":         api.GETHandler,
		"/patch":       api.PATCHHandler,
		"/post":        api.POSTHandler,
		"/put":         api.PUTHandler,
		"/basic-auth/": api.BasicAuthHander,
		"/status/":     api.StatusHander,
		"/headers":     api.HeadersHander,
		"/ip":          api.IPHander,
		"/user-agent":  api.UserAgentHander,
	}

	for endpoint, hander := range patterns {
		mux.HandleFunc(endpoint, hander)
	}
}
