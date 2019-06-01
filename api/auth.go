package api

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
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

const realm = "httpbin-go@project.haujilo.xyz"

type digestAuthJSONResponse struct {
	Authenticated bool   `json:"authenticated"`
	User          string `json:"user"`
}

func digest(data, algorithm string) string {
	b := []byte(data)
	var sum string
	switch algorithm {
	case "SHA-256":
		sum = fmt.Sprintf("%x", sha256.Sum256(b))
	case "SHA-512":
		sum = fmt.Sprintf("%x", sha512.Sum512(b))
	default:
		sum = fmt.Sprintf("%x", md5.Sum(b))
	}
	return sum
}

func generateDigestAuthHeader(qop, algorithm string, stale bool) string {
	nonce := digest(time.Now().Format(time.RFC3339), algorithm)
	opaque := digest(fmt.Sprint(rand.Int()), algorithm)
	return fmt.Sprintf(
		"Digest realm=\"%s\", nonce=\"%s\", qop=\"%s\", opaque=\"%s\", algorithm=%s, stale=%s",
		realm, nonce, qop, opaque, algorithm, strings.ToUpper(strconv.FormatBool(stale)),
	)
}

func parseDigestAuthHeader(r *http.Request) (map[string]string, error) {
	var info map[string]string

	credentials := r.Header.Get("Authorization")
	if credentials == "" || !strings.HasPrefix(credentials, "Digest ") {
		return info, errors.New("Can't parse digest auth header")
	}

	info = make(map[string]string)
	for _, item := range strings.Split(credentials[7:], ",") {
		item = strings.TrimSpace(item)
		index := strings.Index(item, "=")
		k, v := item[:index], item[index+1:]
		if strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") {
			info[k] = v[1 : len(v)-1]
		} else {
			info[k] = v
		}
	}

	return info, nil
}

func digestAuth(r *http.Request, username, password string) (bool, error) {
	info, err := parseDigestAuthHeader(r)
	if err != nil {
		return false, err
	}

	if info["realm"] != realm || info["username"] != username {
		return false, nil
	}

	ha1 := digest(info["username"]+":"+info["realm"]+":"+password, info["algorithm"])
	var ha2 string
	if info["qop"] == "auth-int" {
		body := ""
		if r.Body != nil {
			b, _ := ioutil.ReadAll(r.Body)
			body = string(b)
		}
		ha2 = digest(r.Method+":"+info["uri"]+":"+digest(body, info["algorithm"]), info["algorithm"])
	} else {
		ha2 = digest(r.Method+":"+info["uri"], info["algorithm"])
	}

	var response string
	if info["qop"] == "auth-int" || info["qop"] == "auth" {
		response = digest(strings.Join([]string{ha1, info["nonce"], info["nc"], info["cnonce"], info["qop"], ha2}, ":"), info["algorithm"])
	} else {
		response = digest(strings.Join([]string{ha1, info["nonce"], ha2}, ":"), info["algorithm"])
	}

	return info["response"] == response, nil
}

func DigestAuthHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	paths := strings.Split(r.URL.Path, "/")
	qop, username, password, algorithm := "auth", "", "", "MD5"
	switch len(paths) {
	case 5:
		qop, username, password = paths[2], paths[3], paths[4]
	case 6:
		qop, username, password, algorithm = paths[2], paths[3], paths[4], strings.ToUpper(paths[5])
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if algorithm != "MD5" && algorithm != "SHA-256" && algorithm != "SHA-512" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authenticated, err := digestAuth(r, username, password)
	if err == nil && authenticated {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(digestAuthJSONResponse{Authenticated: true, User: username})
		return
	}
	w.Header().Set("WWW-Authenticate", generateDigestAuthHeader(qop, algorithm, authenticated))
	w.WriteHeader(http.StatusUnauthorized)

}
