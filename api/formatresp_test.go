package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func TestXMLHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func() args {
		r, err := http.NewRequest("GET", "/xml", nil)
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
		{"TestXMLHandler1", createTestCase(), xmlBody},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			XMLHandler(tt.args.w, tt.args.r)
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

func TestHTMLHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func() args {
		r, err := http.NewRequest("GET", "/html", nil)
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
		{"TestHTMLHandler1", createTestCase(), htmlBody},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HTMLHandler(tt.args.w, tt.args.r)
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

func TestGZipHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(name string) struct {
		name   string
		args   args
		result gzipJSONResponse
	} {
		r, err := http.NewRequest("GET", "/gzip", nil)
		if err != nil {
			t.Fatal(err)
		}
		return struct {
			name   string
			args   args
			result gzipJSONResponse
		}{name, args{httptest.NewRecorder(), r}, gzipJSONResponse{true, fmtHeaders(r), r.Method, getIP(r)}}
	}
	tests := []struct {
		name   string
		args   args
		result gzipJSONResponse
	}{
		createTestCase("TestHTMLHandler1"),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GZipHandler(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			reader, err := gzip.NewReader(bytes.NewReader(tt.args.w.Body.Bytes()))
			defer reader.Close()
			if err != nil {
				t.Error(err)
			}
			content, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Error(err)
			}
			var body gzipJSONResponse
			json.Unmarshal(content, &body)
			if !reflect.DeepEqual(tt.result, body) {
				t.Errorf("handler returned wrong response json body: got %v want %v",
					body, tt.result)
			}
		})
	}
}
