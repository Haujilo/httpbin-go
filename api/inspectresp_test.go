package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestResponseHeadersHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(path string) args {
		r, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result map[string]interface{}
	}{
		{"TestResponseHeadersHandler1", createTestCase("/response-headers?a=1&a=2&b=3"), map[string]interface{}{"Content-Length": "80", "Content-Type": "application/json", "a": []interface{}{"1", "2"}, "b": "3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResponseHeadersHandler(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			var body map[string]interface{}
			json.Unmarshal(tt.args.w.Body.Bytes(), &body)
			if !reflect.DeepEqual(tt.result, body) {
				t.Errorf("handler returned wrong response json body: got %v want %v",
					body, tt.result)
			}
		})
	}
}

func TestCacheHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(since, match string) args {
		r, err := http.NewRequest("GET", "/cache", nil)
		if err != nil {
			t.Fatal(err)
		}
		if since != "" {
			r.Header.Set("If-Modified-Since", since)
		}
		if match != "" {
			r.Header.Set("If-None-Match", match)
		}
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result int
	}{
		{"TestCacheHandler1", createTestCase("", ""), 200},
		{"TestCacheHandler2", createTestCase("Feb 17 00:00:00 2019 GMT", ""), 304},
		{"TestCacheHandler3", createTestCase("", "718508b3fd063ea9a0c081cc4f059ea9"), 304},
		{"TestCacheHandler4", createTestCase("Feb 17 00:00:00 2019 GMT", "718508b3fd063ea9a0c081cc4f059ea9"), 304},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CacheHandler(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != tt.result {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.result)
			}
			if status := tt.args.w.Code; status == http.StatusOK {
				if tt.args.w.Header().Get("Last-Modified") == "" || tt.args.w.Header()["ETag"][0] == "" {
					t.Errorf("handler lost header")
				}
			}
		})
	}
}

func TestETagHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(etag, match, notMatch string) args {
		r, err := http.NewRequest("GET", "/etag/"+etag, nil)
		if err != nil {
			t.Fatal(err)
		}
		if notMatch != "" {
			r.Header.Set("If-None-Match", notMatch)
		}
		if match != "" {
			r.Header.Set("If-Match", match)
		}
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result int
	}{
		{"TestETagHandler1", createTestCase("718508b3fd063ea9a0c081cc4f059ea9", "", ""), 200},
		{"TestETagHandler2", createTestCase("718508b3fd063ea9a0c081cc4f059ea9", "", "*"), 304},
		{"TestETagHandler3", createTestCase("718508b3fd063ea9a0c081cc4f059ea9", "", "718508b3fd063ea9a0c081cc4f059ea9"), 304},
		{"TestETagHandler4", createTestCase("718508b3fd063ea9a0c081cc4f059ea9", "", "718508b3fd063ea9a0c081cc4f095ea9"), 200},
		{"TestETagHandler5", createTestCase("718508b3fd063ea9a0c081cc4f059ea9", "718508b3fd063ea9a0c081cc4f095ea9", ""), 412},
		{"TestETagHandler6", createTestCase("718508b3fd063ea9a0c081cc4f059ea9", "*", ""), 200},
		{"TestETagHandler7", createTestCase("718508b3fd063ea9a0c081cc4f059ea9", "718508b3fd063ea9a0c081cc4f095ea9", "718508b3fd063ea9a0c081cc4f059ea9"), 304},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ETagHandler(tt.args.w, tt.args.r)
		})
		if status := tt.args.w.Code; status != tt.result {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.result)
		}
		if status := tt.args.w.Code; status == http.StatusOK || status == http.StatusNotModified {
			if tt.args.w.Header()["ETag"][0] == "" {
				t.Errorf("handler lost header")
			}
		}
	}
}
