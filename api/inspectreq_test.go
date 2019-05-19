package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeadersHander(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(headers [][2]string) args {
		r, err := http.NewRequest("GET", "/headers", nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, item := range headers {
			r.Header.Add(item[0], item[1])
		}
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result map[string]string
	}{
		{
			"TestHeadersHander1",
			createTestCase([][2]string{{"X-TEST", "Test"}}),
			map[string]string{"X-Test": "Test"},
		},
		{
			"TestHeadersHander2",
			createTestCase([][2]string{{"X-TEST", "Test1"}, {"X-tEST", "Test2"}}),
			map[string]string{"X-Test": "Test1,Test2"},
		},
		{
			"TestHeadersHander3",
			createTestCase([][2]string{{"X-TEST", "Test1"}, {"X-tEST", "Test2"}, {"X-tEST1", "Test3"}}),
			map[string]string{"X-Test": "Test1,Test2", "X-Test1": "Test3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HeadersHander(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			var body struct{ Headers map[string]string }
			json.Unmarshal(tt.args.w.Body.Bytes(), &body)
			headers := body.Headers
			for k, v := range tt.result {
				if headers[k] != v {
					t.Errorf("handler returned wrong response json body: got {\"%v\": \"%v\"} want {\"%v\": \"%v\"}",
						k, headers[k], k, v)
				}
			}
		})
	}
}

func TestIPHander(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(addr string) args {
		r, err := http.NewRequest("GET", "/ip", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.RemoteAddr = addr
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result string
	}{
		{"TestIPHander1", createTestCase("192.168.1.100:10086"), "192.168.1.100"},
		{"TestIPHander2", createTestCase("8.8.8.8:8888"), "8.8.8.8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IPHander(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			var body struct{ Origin string }
			json.Unmarshal(tt.args.w.Body.Bytes(), &body)
			origin := body.Origin
			if origin != tt.result {
				t.Errorf("handler returned wrong response json body: got {\"origin\": \"%v\"} want {\"origin\": \"%v\"}",
					origin, tt.result)
			}
		})
	}
}

func TestUserAgentHander(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(userAgent string) args {
		r, err := http.NewRequest("GET", "/ip", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("User-Agent", userAgent)
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result string
	}{
		{"TestUserAgentHander", createTestCase("curl/7.54.0"), "curl/7.54.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UserAgentHander(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			var body struct {
				UserAgent string `json:"user-agent"`
			}
			json.Unmarshal(tt.args.w.Body.Bytes(), &body)
			userAgent := body.UserAgent
			if userAgent != tt.result {
				t.Errorf("handler returned wrong response json body: got {\"user-agent\": \"%v\"} want {\"user-agent\": \"%v\"}",
					userAgent, tt.result)
			}
		})
	}
}
