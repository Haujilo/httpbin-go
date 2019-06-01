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

func TestHiddenBasicAuthHander(t *testing.T) {
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
		{"TestHiddenBasicAuthHander1", createTestCase("/hidden-basic-auth/test/test", "test", "test"), result{200, &basicAuthJSONResponse{Authenticated: true, User: "test"}}},
		{"TestHiddenBasicAuthHander2", createTestCase("/hidden-basic-auth/test/test", "", ""), result{404, nil}},
		{"TestHiddenBasicAuthHander3", createTestCase("/hidden-basic-auth/test/test/", "", ""), result{400, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HiddenBasicAuthHander(tt.args.w, tt.args.r)
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

func TestDigestAuthHander(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(path, credentials string) args {
		r, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		if credentials != "" {
			r.Header.Set("Authorization", credentials)
		}
		return args{httptest.NewRecorder(), r}
	}
	type result struct {
		code     int
		response *digestAuthJSONResponse
	}
	tests := []struct {
		name   string
		args   args
		result result
	}{
		{"TestDigestAuthHander1", createTestCase("/digest-auth/auth/test/test", `Digest username="test", realm="httpbin-go@project.haujilo.xyz", nonce="75cbb4a2367b603fd3e7e4b67af388a0", uri="/digest-auth/auth/test/test", cnonce="M2Q2YWQxOWQ0YmZkZDJlMjY1MGYxMTMyMTlmNGNhYmY=", nc=00000001, qop=auth, response="d6cae8c7775d18b1159367cd6ac83c1f", opaque="f154b9bfc203c3adf417bb3f08191633", algorithm="MD5"`), result{200, &digestAuthJSONResponse{Authenticated: true, User: "test"}}},
		{"TestDigestAuthHander2", createTestCase("/digest-auth/auth-int/test/test", `Digest username="test", realm="httpbin-go@project.haujilo.xyz", nonce="6d2e1732e28e0b0e706504b9e1814aed", uri="/digest-auth/auth-int/test/test", cnonce="ODc4Y2QxNGQ0MDkyYjkyZTNiNzEzYTFjNGY1YjdiMDk=", nc=00000001, qop=auth-int, response="c9b232792d2bfb287bc0c2e7f3d63a9a", opaque="855914e07cac2d33fa54ed8e1f790440", algorithm="MD5"`), result{200, &digestAuthJSONResponse{Authenticated: true, User: "test"}}},
		{"TestDigestAuthHander3", createTestCase("/digest-auth/auth-int/test/test/MD5", `Digest username="test", realm="httpbin-go@project.haujilo.xyz", nonce="b90d5196857db630d611b2188597ec58", uri="/digest-auth/auth-int/test/test/MD5", cnonce="ODY4YWVhZTYyOGI3ODYxMTk0YzU1NTU4ODUwZjdhZTI=", nc=00000001, qop=auth-int, response="b60dd9f566449559706f3eb7edea752c", opaque="e51cecd28bfacf2aa61e2876920f083c", algorithm="MD5"`), result{200, &digestAuthJSONResponse{Authenticated: true, User: "test"}}},
		{"TestDigestAuthHander4", createTestCase("/digest-auth/auth-int/test/test/MD5", ""), result{401, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DigestAuthHander(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != tt.result.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.result.code)
			}
			if tt.result.response != nil {
				var body digestAuthJSONResponse
				json.Unmarshal(tt.args.w.Body.Bytes(), &body)
				if !reflect.DeepEqual(*tt.result.response, body) {
					t.Errorf("handler returned wrong response json body: got %v want %v",
						body, *tt.result.response)
				}
			}
		})
	}
}
