package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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
