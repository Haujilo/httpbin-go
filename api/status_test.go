package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusHander(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	createTestCase := func(statusCodesParams string) args {
		r, err := http.NewRequest("GET", statusCodesParams, nil)
		if err != nil {
			t.Fatal(err)
		}
		return args{httptest.NewRecorder(), r}
	}
	tests := []struct {
		name   string
		args   args
		result []int
	}{
		{"TestStatusHander1", createTestCase("/status/"), []int{400}},
		{"TestStatusHander2", createTestCase("/status/201"), []int{201}},
		{"TestStatusHander3", createTestCase("/status/502/"), []int{502}},
		{"TestStatusHander4", createTestCase("/status/abc/"), []int{400}},
		{"TestStatusHander5", createTestCase("/status/201/201"), []int{400}},
		{"TestStatusHander6", createTestCase("/status/201,401,500"), []int{201, 401, 500}},
		{"TestStatusHander7", createTestCase("/status/201:6,401,500:3"), []int{201, 401, 500}},
		{"TestStatusHander8", createTestCase("/status/201:6,401,500:3/"), []int{201, 401, 500}},
		{"TestStatusHander9", createTestCase("/status/201:6,401,abc:3/"), []int{400}},
		{"TestStatusHander10", createTestCase("/status/201:6,401,500:3/abc"), []int{400}},
		{"TestStatusHander11", createTestCase("/status/201:6,401,500:3/404"), []int{400}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StatusHander(tt.args.w, tt.args.r)
		})
		var flag bool
		status := tt.args.w.Code
		for _, rs := range tt.result {
			if status == rs {
				flag = true
				break
			}
		}
		if !flag {
			t.Errorf("handler returned wrong status code: got %v want status in %v",
				status, tt.result)
		}
	}
}
