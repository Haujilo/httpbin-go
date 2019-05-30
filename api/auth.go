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
	w.Header().Set("Content-Type", "text/plain")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	paths := strings.Split(r.URL.Path, "/")
	if len(paths) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, password, _ := r.BasicAuth()

	if username == paths[2] && password == paths[3] {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(basicAuthJSONResponse{Authenticated: true, User: username})
		return
	}

	w.WriteHeader(http.StatusUnauthorized)

}

func HiddenBasicAuthHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	paths := strings.Split(r.URL.Path, "/")
	if len(paths) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, password, _ := r.BasicAuth()

	if username == paths[2] && password == paths[3] {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(basicAuthJSONResponse{Authenticated: true, User: username})
		return
	}

	w.WriteHeader(http.StatusNotFound)

}

type bearerAuthJSONResponse struct {
	Authenticated bool   `json:"authenticated"`
	Token         string `json:"token"`
}

func BearerAuthHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	authorization := r.Header.Get("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bearerAuthJSONResponse{Authenticated: true, Token: authorization[7:]})
		return
	}

	w.WriteHeader(http.StatusUnauthorized)

}
