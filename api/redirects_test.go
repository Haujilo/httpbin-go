package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAbsoluteRedirectHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(path string) args {
		r, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		r.RequestURI = path
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result string
	}{
		{"TestAbsoluteRedirectHandler1", createTestCase("/absolute-redirect/2"), "/absolute-redirect/1"},
		{"TestAbsoluteRedirectHandler2", createTestCase("/absolute-redirect/1"), "/get"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AbsoluteRedirectHandler(tt.args.w, tt.args.r)
		})
		if status := tt.args.w.Code; status != http.StatusFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusFound)
		}
		if !strings.HasSuffix(tt.args.w.Header().Get("Location"), tt.result) {
			t.Errorf("handler returned wrong location header: got %v want %v",
				tt.args.w.Header().Get("Location"), tt.result)
		}
	}
}
