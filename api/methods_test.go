package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGETHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(host, path, addr string, headers [][2]string) args {
		r, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		r.RequestURI = path
		r.Host = host
		r.RemoteAddr = addr
		for _, item := range headers {
			r.Header.Add(item[0], item[1])
		}
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result responseGETHandler
	}{
		{
			"TestGETHandler",
			createTestCase("localhost:1121", "/get?a=1&a=2&b=3", "127.0.0.1:1121", [][2]string{{"X-TEST", "Test"}}),
			responseGETHandler{
				Args: map[string]interface{}{
					"a": []interface{}{"1", "2"},
					"b": "3",
				},
				Headers: map[string]string{"X-Test": "Test"},
				Origin:  "127.0.0.1",
				URL:     "http://localhost:1121/get?a=1\u0026a=2\u0026b=3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GETHandler(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			var body responseGETHandler
			json.Unmarshal(tt.args.w.Body.Bytes(), &body)
			for k, v := range tt.result.Args {
				if !reflect.DeepEqual(body.Args[k], v) {
					t.Errorf("handler returned wrong response json body: got %v want %v",
						body.Args, tt.result.Args)
				}
			}
			for k, v := range tt.result.Headers {
				if body.Headers[k] != v {
					t.Errorf("handler returned wrong response json body: got %v want %v",
						body.Headers, tt.result.Headers)
				}
			}
			if tt.result.Origin != body.Origin {
				t.Errorf("handler returned wrong response json body: got %v want %v",
					body.Origin, tt.result.Origin)
			}
			if tt.result.URL != body.URL {
				t.Errorf("handler returned wrong response json body: got %v want %v",
					body.URL, tt.result.URL)
			}
		})
	}
}
