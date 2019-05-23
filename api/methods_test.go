package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
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
		result methodsGETJSONResponse
	}{
		{
			"TestGETHandler",
			createTestCase("localhost:1121", "/get?a=1&a=2&b=3", "127.0.0.1:1121", [][2]string{{"X-TEST", "Test"}}),
			methodsGETJSONResponse{
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
			var body methodsGETJSONResponse
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

func TestPOSTHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(host, path, addr string, headers [][2]string, body io.Reader) args {
		r, err := http.NewRequest("POST", path, body)
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
		result methodsJSONResponse
	}{
		{
			"TestPOSTHandler1",
			createTestCase("localhost:1121", "/post?a=1&a=2&b=3", "127.0.0.1:1121", [][2]string{{"Content-Type", "application/x-www-form-urlencoded"}, {"X-TEST", "Test"}}, func() io.Reader {
				data := url.Values{}
				data.Add("test1", "test1")
				data.Add("test1", "test1-1")
				data.Add("test2", "test2")
				return strings.NewReader(data.Encode())
			}()),
			methodsJSONResponse{
				Args: map[string]interface{}{
					"a": []interface{}{"1", "2"},
					"b": "3",
				},
				Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded", "X-Test": "Test"},
				Origin:  "127.0.0.1",
				URL:     "http://localhost:1121/post?a=1\u0026a=2\u0026b=3",
				Form: map[string]interface{}{
					"test1": []interface{}{"test1", "test1-1"},
					"test2": "test2",
				},
			},
		},
		{
			"TestPOSTHandler2",
			func() args {
				var body bytes.Buffer
				w := multipart.NewWriter(&body)
				defer w.Close()
				w.SetBoundary("1eecdfc9a2d50fa834599aad3e2fb433a5a87887435d5e026e1d651da96a")
				w.WriteField("test1", "test1")
				w.WriteField("test1", "test1-1")
				w.WriteField("test2", "test2")
				ww, _ := w.CreateFormFile("testfile", "testfile")
				ww.Write([]byte("abc\n"))
				return createTestCase("localhost:1121", "/post?a=1&a=2&b=3", "127.0.0.1:1121", [][2]string{{"Content-Type", w.FormDataContentType()}, {"X-TEST", "Test"}}, &body)
			}(),
			methodsJSONResponse{
				Args: map[string]interface{}{
					"a": []interface{}{"1", "2"},
					"b": "3",
				},
				Headers: map[string]string{"Content-Type": "multipart/form-data; boundary=1eecdfc9a2d50fa834599aad3e2fb433a5a87887435d5e026e1d651da96a", "X-Test": "Test"},
				Origin:  "127.0.0.1",
				URL:     "http://localhost:1121/post?a=1\u0026a=2\u0026b=3",
				Form: map[string]interface{}{
					"test1": []interface{}{"test1", "test1-1"},
					"test2": "test2",
				},
				Files: map[string]interface{}{
					"testfile": "abc\n",
				},
			},
		},
		{
			"TestPOSTHandler3",
			createTestCase("localhost:1121", "/post?a=1&a=2&b=3", "127.0.0.1:1121", [][2]string{{"Content-Type", "application/json"}, {"X-TEST", "Test"}}, strings.NewReader("{\"a\": \"1\"}")),
			methodsJSONResponse{
				Args: map[string]interface{}{
					"a": []interface{}{"1", "2"},
					"b": "3",
				},
				Headers: map[string]string{"Content-Type": "application/json", "X-Test": "Test"},
				Origin:  "127.0.0.1",
				URL:     "http://localhost:1121/post?a=1\u0026a=2\u0026b=3",
				JSON:    map[string]interface{}{"a": "1"},
			},
		},
		{
			"TestPOSTHandler4",
			createTestCase("localhost:1121", "/post?a=1&a=2&b=3", "127.0.0.1:1121", [][2]string{{"X-TEST", "Test"}}, strings.NewReader("abcdefgh")),
			methodsJSONResponse{
				Args: map[string]interface{}{
					"a": []interface{}{"1", "2"},
					"b": "3",
				},
				Headers: map[string]string{"X-Test": "Test"},
				Origin:  "127.0.0.1",
				URL:     "http://localhost:1121/post?a=1\u0026a=2\u0026b=3",
				Data:    "abcdefgh",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			POSTHandler(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			var body methodsJSONResponse
			json.Unmarshal(tt.args.w.Body.Bytes(), &body)
			if !reflect.DeepEqual(tt.result, body) {
				t.Errorf("handler returned wrong response json body: got %v want %v",
					body, tt.result)
			}
		})
	}
}

func TestPUTHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(host, path, addr string, headers [][2]string, body io.Reader) args {
		r, err := http.NewRequest("PUT", path, body)
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
		result methodsJSONResponse
	}{
		{
			"TestPUTHandler1",
			createTestCase("localhost:1121", "/put?a=1&a=2&b=3", "127.0.0.1:1121", [][2]string{{"Content-Type", "application/x-www-form-urlencoded"}, {"X-TEST", "Test"}}, func() io.Reader {
				data := url.Values{}
				data.Add("test1", "test1")
				data.Add("test1", "test1-1")
				data.Add("test2", "test2")
				return strings.NewReader(data.Encode())
			}()),
			methodsJSONResponse{
				Args: map[string]interface{}{
					"a": []interface{}{"1", "2"},
					"b": "3",
				},
				Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded", "X-Test": "Test"},
				Origin:  "127.0.0.1",
				URL:     "http://localhost:1121/put?a=1\u0026a=2\u0026b=3",
				Form: map[string]interface{}{
					"test1": []interface{}{"test1", "test1-1"},
					"test2": "test2",
				},
			},
		},
		{
			"TestPUTHandler2",
			func() args {
				var body bytes.Buffer
				w := multipart.NewWriter(&body)
				defer w.Close()
				w.SetBoundary("1eecdfc9a2d50fa834599aad3e2fb433a5a87887435d5e026e1d651da96a")
				w.WriteField("test1", "test1")
				w.WriteField("test1", "test1-1")
				w.WriteField("test2", "test2")
				ww, _ := w.CreateFormFile("testfile", "testfile")
				ww.Write([]byte("abc\n"))
				return createTestCase("localhost:1121", "/put?a=1&a=2&b=3", "127.0.0.1:1121", [][2]string{{"Content-Type", w.FormDataContentType()}, {"X-TEST", "Test"}}, &body)
			}(),
			methodsJSONResponse{
				Args: map[string]interface{}{
					"a": []interface{}{"1", "2"},
					"b": "3",
				},
				Headers: map[string]string{"Content-Type": "multipart/form-data; boundary=1eecdfc9a2d50fa834599aad3e2fb433a5a87887435d5e026e1d651da96a", "X-Test": "Test"},
				Origin:  "127.0.0.1",
				URL:     "http://localhost:1121/put?a=1\u0026a=2\u0026b=3",
				Form: map[string]interface{}{
					"test1": []interface{}{"test1", "test1-1"},
					"test2": "test2",
				},
				Files: map[string]interface{}{
					"testfile": "abc\n",
				},
			},
		},
		{
			"TestPUTHandler3",
			createTestCase("localhost:1121", "/put?a=1&a=2&b=3", "127.0.0.1:1121", [][2]string{{"Content-Type", "application/json"}, {"X-TEST", "Test"}}, strings.NewReader("{\"a\": \"1\"}")),
			methodsJSONResponse{
				Args: map[string]interface{}{
					"a": []interface{}{"1", "2"},
					"b": "3",
				},
				Headers: map[string]string{"Content-Type": "application/json", "X-Test": "Test"},
				Origin:  "127.0.0.1",
				URL:     "http://localhost:1121/put?a=1\u0026a=2\u0026b=3",
				JSON:    map[string]interface{}{"a": "1"},
			},
		},
		{
			"TestPUTHandler4",
			createTestCase("localhost:1121", "/put?a=1&a=2&b=3", "127.0.0.1:1121", [][2]string{{"X-TEST", "Test"}}, strings.NewReader("abcdefgh")),
			methodsJSONResponse{
				Args: map[string]interface{}{
					"a": []interface{}{"1", "2"},
					"b": "3",
				},
				Headers: map[string]string{"X-Test": "Test"},
				Origin:  "127.0.0.1",
				URL:     "http://localhost:1121/put?a=1\u0026a=2\u0026b=3",
				Data:    "abcdefgh",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PUTHandler(tt.args.w, tt.args.r)
			if status := tt.args.w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			var body methodsJSONResponse
			json.Unmarshal(tt.args.w.Body.Bytes(), &body)
			if !reflect.DeepEqual(tt.result, body) {
				t.Errorf("handler returned wrong response json body: got %v want %v",
					body, tt.result)
			}
		})
	}
}
