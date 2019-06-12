package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getRedirectURL(r *http.Request, n int) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	if n != 1 {
		url := scheme + r.Host + r.RequestURI
		return fmt.Sprintf("%s/%d", url[:strings.LastIndex(url, "/")], n-1)
	}
	return scheme + r.Host + "/get"
}

func AbsoluteRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	paths := strings.Split(r.URL.Path, "/")
	pathsLength := len(paths)
	if pathsLength != 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	n, err := strconv.Atoi(paths[pathsLength-1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if n < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, getRedirectURL(r, n), http.StatusFound)
}
