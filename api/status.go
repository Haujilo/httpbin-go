package api

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func randomStatusSelect(choices map[int]int) int {
	var totalWeight int
	for _, weight := range choices {
		totalWeight += weight
	}
	r := rand.Intn(totalWeight)
	var rStatus int
	for status, weight := range choices {
		r -= weight
		if r <= 0 {
			rStatus = status
			break
		}
	}
	return rStatus
}

func StatusHander(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	pathsLength := len(paths)
	if pathsLength < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if pathsLength > 3 {
		if pathsLength > 4 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if paths[pathsLength-1] != "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	statusCodesParam := paths[2]
	statusCodeChoices := make(map[int]int)
	for index, str := range strings.Split(statusCodesParam, ",") {
		if i := strings.Index(str, ":"); i > -1 {
			status, err := strconv.Atoi(str[:i])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			weight, err := strconv.Atoi(str[i+1:])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			statusCodeChoices[status] = weight
		} else {
			status, err := strconv.Atoi(str)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			statusCodeChoices[status] = index + 1
		}
	}
	status := randomStatusSelect(statusCodeChoices)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
}
