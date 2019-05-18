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
	createTestCase := func(headers map[string]string) args {
		r, err := http.NewRequest("GET", "/headers", nil)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range headers {
			r.Header.Set(k, v)
		}
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result map[string]string
	}{
		{
			"TestHeadersHander",
			createTestCase(map[string]string{"X-TEST": "Test"}),
			map[string]string{"X-TEST": "Test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HeadersHander(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			var body struct{ Headers http.Header }
			json.Unmarshal(tt.args.w.Body.Bytes(), &body)
			headers := body.Headers
			for k, v := range tt.result {
				if headers.Get(k) != v {
					t.Errorf("handler returned wrong response json body: got {\"%v\": \"%v\"} want {\"%v\": \"%v\"}",
						k, headers.Get(k), k, v)
				}
			}
		})
	}
}
