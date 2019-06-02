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
