package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func fmtQueryString(r *http.Request) map[string]interface{} {
	queryString := make(map[string]interface{})
	for k, v := range r.URL.Query() {
		if len(v) == 1 {
			queryString[k] = v[0]
		} else {
			queryString[k] = v
		}
	}
	return queryString
}

func fmtForm(r *http.Request) map[string]interface{} {
	if err := r.ParseForm(); err != nil {
		return nil
	}
	form := make(map[string]interface{})
	for k, v := range r.PostForm {
		if len(v) == 1 {
			form[k] = v[0]
		} else {
			form[k] = v
		}
	}
	return form
}

const maxMultipartFormSize = 0

func fmtMultipartForm(r *http.Request) (map[string]interface{}, map[string]interface{}) {
	if err := r.ParseMultipartForm(maxMultipartFormSize); err != nil {
		return nil, nil
	}
	form := make(map[string]interface{})
	files := make(map[string]interface{})
	for k, v := range r.MultipartForm.Value {
		if len(v) == 1 {
			form[k] = v[0]
		} else {
			form[k] = v
		}
	}
	for k, v := range r.MultipartForm.File {
		if len(v) == 1 {
			f, _ := v[0].Open()
			data, _ := ioutil.ReadAll(f)
			files[k] = string(data)
		} else {
			var datas []string
			for _, vv := range v {
				f, _ := vv.Open()
				data, _ := ioutil.ReadAll(f)
				datas = append(datas, string(data))
			}
			files[k] = datas
		}
	}
	return form, files
}

func getFullURL(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	return scheme + r.Host + r.RequestURI
}

type methodsGETJSONResponse struct {
	Args    map[string]interface{} `json:"args"`
	Headers map[string]string      `json:"headers"`
	Origin  string                 `json:"origin"`
	URL     string                 `json:"url"`
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methodsGETJSONResponse{
		Args:    fmtQueryString(r),
		Headers: fmtHeaders(r),
		Origin:  getIP(r),
		URL:     getFullURL(r),
	})
}

type methodsJSONResponse struct {
	Args    map[string]interface{} `json:"args"`
	Data    string                 `json:"data"`
	Files   map[string]interface{} `json:"files"`
	Form    map[string]interface{} `json:"form"`
	Headers map[string]string      `json:"headers"`
	JSON    interface{}            `json:"json"`
	Origin  string                 `json:"origin"`
	URL     string                 `json:"url"`
}

func methodsHander(wp *http.ResponseWriter, r *http.Request) {
	w := *wp
	response := methodsJSONResponse{
		Args:    fmtQueryString(r),
		Headers: fmtHeaders(r),
		Origin:  getIP(r),
		URL:     getFullURL(r),
	}
	contentType := r.Header.Get("Content-Type")
	switch {
	case contentType == "application/x-www-form-urlencoded":
		response.Form = fmtForm(r)
	case strings.HasPrefix(contentType, "multipart/form-data"):
		response.Form, response.Files = fmtMultipartForm(r)
	case contentType == "application/json":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &response.JSON)
	default:
		body, _ := ioutil.ReadAll(r.Body)
		response.Data = string(body)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	methodsHander(&w, r)
}

func PUTHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	methodsHander(&w, r)
}

func PATCHHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	methodsHander(&w, r)
}

func DELETEHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	methodsHander(&w, r)
}
