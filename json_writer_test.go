package reply

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewJSONWriterHasEmptyBuffer(t *testing.T) {
	w := httptest.NewRecorder()
	jw := NewJSONWriter()
	jw.WriteTo(w)
	if got, want := w.Body.String(), ""; got != want {
		t.Fatalf(errorString, got, want)
	}
}

func TestEncode(t *testing.T) {
	cases := map[string]struct {
		data     any
		wantErr  bool
		wantBody string
	}{
		"error - fail encode": {
			data:     make(chan int),
			wantErr:  true,
			wantBody: "",
		},
		"ok": {
			data:     map[string]string{"foo": "bar"},
			wantErr:  false,
			wantBody: `{"foo":"bar"}`,
		},
	}
	jw := NewJSONWriter()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := jw.Encode(c.data)
			if err != nil && !c.wantErr {
				t.Fatalf("got unwanted error - %s", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("wanted an error but didn't get one")
			}
			jw.WriteTo(w)
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
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
	jw := NewJSONWriter()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			jw.Error(w, http.StatusText(c.code), c.code)
			if got := w.Code; got != c.code {
				t.Fatalf(errorString, got, c.code)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.want {
				t.Fatalf(errorString, got, c.want)
			}
		})
	}
}

func TestJSONReply(t *testing.T) {
	cases := map[string]struct {
		code     int
		opts     Options
		wantCode int
		wantBody string
	}{
		"error - fail encode": {
			code:     http.StatusOK,
			opts:     Options{Data: map[string]interface{}{"foo": make(chan int)}},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"Internal Server Error"}`,
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
	jw := NewJSONWriter()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			jw.Reply(w, c.code, c.opts)
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); !strings.Contains(got, c.wantBody) {
				t.Fatalf(errorString, got, c.wantBody)
			}
		})
	}
}
