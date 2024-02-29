package reply

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSONReply(t *testing.T) {
	cases := map[string]struct {
		code     int
		opts     Options
		wantErr  bool
		wantCode int
		wantBody string
	}{
		"error - fail encode": {
			code:     http.StatusOK,
			opts:     Options{Data: map[string]interface{}{"foo": make(chan int)}},
			wantErr:  true,
			wantCode: http.StatusOK,
			wantBody: "",
		},
		"ok": {
			code:     http.StatusOK,
			opts:     Options{Data: map[string]string{"foo": "bar"}},
			wantCode: http.StatusOK,
			wantBody: `{"foo":"bar"}`,
		},
		"ok; multi-key": {
			code:     http.StatusCreated,
			opts:     Options{Data: map[string]string{"baz": "qux", "foo": "bar"}},
			wantCode: http.StatusCreated,
			wantBody: `{"baz":"qux","foo":"bar"}`,
		},
	}
	jw := JSONWriter{}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := jw.Reply(w, c.code, c.opts)
			if (err != nil) != c.wantErr {
				t.Errorf(errorString, err, c.wantErr)
			}
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); !strings.Contains(got, c.wantBody) {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestJSONError(t *testing.T) {
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
	jw := JSONWriter{}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			jw.Error(w, http.StatusText(c.code), c.code)
			if got := w.Code; got != c.code {
				t.Errorf(errorString, got, c.code)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.want {
				t.Errorf(errorString, got, c.want)
			}
		})
	}
}
