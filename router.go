package main

import (
	"net/http"

	"github.com/Haujilo/httpbin-go/api"
)

func route(mux *http.ServeMux) {
	patterns := map[string]func(w http.ResponseWriter, r *http.Request){
		"/delete":             api.DELETEHandler,
		"/get":                api.GETHandler,
		"/patch":              api.PATCHHandler,
		"/post":               api.POSTHandler,
		"/put":                api.PUTHandler,
		"/basic-auth/":        api.BasicAuthHander,
		"/bearer":             api.BearerAuthHander,
		"/digest-auth/":       api.DigestAuthHander,
		"/hidden-basic-auth/": api.HiddenBasicAuthHander,
		"/status/":            api.StatusHander,
		"/headers":            api.HeadersHander,
		"/ip":                 api.IPHander,
		"/user-agent":         api.UserAgentHander,
		"/cache":              api.CacheHandler,
		"/cache/":             api.CacheControlHandler,
		"/etag/":              api.ETagHandler,
		"/response-headers":   api.ResponseHeadersHandler,
		"/deflate":            api.DeflateHandler,
		"/deny":               api.DenyHandler,
		"/encoding/utf8":      api.UTF8Handler,
		"/gzip":               api.GZipHandler,
		"/html":               api.HTMLHandler,
		"/json":               api.JsonHandler,
		"/robots.txt":         api.RobotTxtHandler,
		"/xml":                api.XMLHandler,
		"/absolute-redirect/": api.AbsoluteRedirectHandler,
	}

	for endpoint, hander := range patterns {
		mux.HandleFunc(endpoint, hander)
	}
}
