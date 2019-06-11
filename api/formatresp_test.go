package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRobotTxtHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func() args {
		r, err := http.NewRequest("GET", "/robots.txt", nil)
		if err != nil {
			t.Fatal(err)
		}
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result string
	}{
		{"TestRobotTxtHandler1", createTestCase(), robotTxt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RobotTxtHandler(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			body, err := ioutil.ReadAll(tt.args.w.Body)
			if err != nil {
				t.Error(err)
			}
			if string(body) != tt.result {
				t.Errorf("http body error, got %v want %v", tt.result, body)
			}
		})
	}
}

func TestJsonHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func() args {
		r, err := http.NewRequest("GET", "/json", nil)
		if err != nil {
			t.Fatal(err)
		}
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result string
	}{
		{"TestJsonHandler1", createTestCase(), jsonBody},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			JsonHandler(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			body, err := ioutil.ReadAll(tt.args.w.Body)
			if err != nil {
				t.Error(err)
			}
			if string(body) != tt.result {
				t.Errorf("http body error, got %v want %v", tt.result, body)
			}
		})
	}
}
