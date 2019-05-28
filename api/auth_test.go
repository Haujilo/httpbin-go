package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestBasicAuthHander(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(path, username, password string) args {
		r, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		if username != "" {
			r.SetBasicAuth(username, password)
		}
		return args{httptest.NewRecorder(), r}
	}
	type result struct {
		code     int
		response *basicAuthJSONResponse
	}
	tests := []struct {
		name   string
		args   args
		result result
	}{
		{"TestBasicAuthHander1", createTestCase("/basic-auth/test/test", "test", "test"), result{200, &basicAuthJSONResponse{Authenticated: true, User: "test"}}},
		{"TestBasicAuthHander2", createTestCase("/basic-auth/test/test", "", ""), result{401, nil}},
		{"TestBasicAuthHander3", createTestCase("/basic-auth/test/test/", "", ""), result{400, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BasicAuthHander(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != tt.result.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.result.code)
			}
			if tt.result.response != nil {
				var body basicAuthJSONResponse
				json.Unmarshal(tt.args.w.Body.Bytes(), &body)
				if !reflect.DeepEqual(*tt.result.response, body) {
					t.Errorf("handler returned wrong response json body: got %v want %v",
						body, *tt.result.response)
				}
			}
		})
	}
}

func TestBearerAuthHander(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(token string) args {
		r, err := http.NewRequest("GET", "/bearer", nil)
		if err != nil {
			t.Fatal(err)
		}
		if token != "" {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		return args{httptest.NewRecorder(), r}
	}
	type result struct {
		code     int
		response *bearerAuthJSONResponse
	}
	tests := []struct {
		name   string
		args   args
		result result
	}{
		{"TestBearerAuthHander1", createTestCase("mF_9.B5f-4.1JqM"), result{200, &bearerAuthJSONResponse{true, "mF_9.B5f-4.1JqM"}}},
		{"TestBearerAuthHander2", createTestCase(""), result{401, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BearerAuthHander(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != tt.result.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.result.code)
			}
			if tt.result.response != nil {
				var body bearerAuthJSONResponse
				json.Unmarshal(tt.args.w.Body.Bytes(), &body)
				if !reflect.DeepEqual(*tt.result.response, body) {
					t.Errorf("handler returned wrong response json body: got %v want %v",
						body, *tt.result.response)
				}
			}
		})
	}
}
