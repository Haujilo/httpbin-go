package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CacheHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("If-Modified-Since") == "" && r.Header.Get("If-None-Match") == "" {
		now := time.Now().UTC()
		w.Header().Set("Last-Modified", now.Format(http.TimeFormat))
		w.Header()["ETag"] = []string{digest(fmt.Sprint(time.Now().UnixNano()), "MD5")}
		GETHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusNotModified)
}

func CacheControlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := strings.Split(r.URL.Path, "/")
	if len(params) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	maxAge, err := strconv.Atoi(params[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("cache-control", fmt.Sprintf("public, max-age=%d", maxAge))
	GETHandler(w, r)
}

func ETagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := strings.Split(r.URL.Path, "/")
	if len(params) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	etag := params[2]
	ifNoneMatch := r.Header.Get("If-None-Match")
	ifMatch := r.Header.Get("If-Match")

	if ifNoneMatch != "" {
		if etag == ifNoneMatch || ifNoneMatch == "*" {
			w.WriteHeader(http.StatusNotModified)
			w.Header()["ETag"] = []string{etag}
			return
		}
	} else {
		if ifMatch != "" {
			if etag != ifMatch && ifMatch != "*" {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
		}
	}

	w.Header()["ETag"] = []string{etag}
	GETHandler(w, r)
}

func ResponseHeadersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	payload := fmtQueryString(r)
	payload["Content-Length"] = ""
	payload["Content-Type"] = "application/json"

	tmp, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentLength := len(tmp)
	size := contentLength + len(strconv.Itoa(contentLength)) + 1
	if len(strconv.Itoa(size)) > len(strconv.Itoa(contentLength)) {
		size++
	}
	payload["Content-Length"] = strconv.Itoa(size)

	header := w.Header()
	for k, v := range payload {
		switch v.(type) {
		case string:
			header.Add(k, v.(string))
		default:
			for _, vv := range v.([]string) {
				header.Add(k, vv)
			}
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(body)
	w.Write([]byte{'\n'})
}
