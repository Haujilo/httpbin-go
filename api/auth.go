package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type basicAuthJSONResponse struct {
	Authenticated bool   `json:"authenticated"`
	User          string `json:"user"`
}

func BasicAuthHander(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	paths := strings.Split(r.URL.Path, "/")
	if len(paths) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, password, _ := r.BasicAuth()

	if username == paths[2] && password == paths[3] {
		json.NewEncoder(w).Encode(basicAuthJSONResponse{Authenticated: true, User: username})
		return
	}

	w.WriteHeader(http.StatusUnauthorized)

}

type bearerAuthJSONResponse struct {
	Authenticated bool   `json:"authenticated"`
	Token         string `json:"token"`
}

func BearerAuthHander(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	authorization := r.Header.Get("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		json.NewEncoder(w).Encode(bearerAuthJSONResponse{Authenticated: true, Token: authorization[7:]})
		return
	}

	w.WriteHeader(http.StatusUnauthorized)

}
