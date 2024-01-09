package reply

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	cases := map[string]struct {
		code     int
		opts     Options
		wantCode int
		wantBody string
	}{
		"err on marshal": {
			code:     http.StatusOK,
			opts:     Options{Data: map[string]interface{}{"foo": make(chan int)}},
			wantCode: http.StatusInternalServerError,
			wantBody: "failed to marshal",
		},
		"ok": {
			code:     http.StatusOK,
			opts:     Options{Data: map[string]string{"foo": "bar"}},
			wantCode: http.StatusOK,
			wantBody: `{"foo":"bar"}`,
		},
		"created; multi-key": {
			code:     http.StatusCreated,
			opts:     Options{Data: map[string]string{"baz": "qux", "foo": "bar"}},
			wantCode: http.StatusCreated,
			wantBody: `{"baz":"qux","foo":"bar"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			JSONWriter{}.Write(w, c.code, c.opts)
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := w.Body.String(); !strings.Contains(got, c.wantBody) {
				t.Fatalf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestErrorJSON(t *testing.T) {
	cases := map[string]struct {
		code int
		want string
	}{
		"bad Request": {
			code: http.StatusBadRequest,
			want: `{"error":"Bad Request"}`,
		},
		"not found": {
			code: http.StatusNotFound,
			want: `{"error":"Not Found"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			JSONWriter{}.Error(w, c.code)
			if got := w.Code; got != c.code {
				t.Fatalf(errorString, got, c.code)
			}
			if got := w.Body.String(); got != c.want {
				t.Fatalf(errorString, got, c.want)
			}
		})
	}
}
