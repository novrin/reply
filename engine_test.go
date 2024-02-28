package reply

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBadRequest(t *testing.T) {
	cases := map[string]struct {
		reply    Engine
		wantCode int
		wantBody string
	}{
		"template writer": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{})},
			wantCode: http.StatusBadRequest,
			wantBody: fmt.Sprintf("<p>%s</p>", http.StatusText(http.StatusBadRequest)),
		},
		"json writer": {
			reply:    Engine{Writer: NewJSONWriter()},
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"Bad Request"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.BadRequest(w)
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestUnauthorized(t *testing.T) {
	cases := map[string]struct {
		reply    Engine
		wantCode int
		wantBody string
	}{
		"template writer": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{})},
			wantCode: http.StatusUnauthorized,
			wantBody: fmt.Sprintf("<p>%s</p>", http.StatusText(http.StatusUnauthorized)),
		},
		"json writer": {
			reply:    Engine{Writer: NewJSONWriter()},
			wantCode: http.StatusUnauthorized,
			wantBody: `{"error":"Unauthorized"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.Unauthorized(w)
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestForbidden(t *testing.T) {
	cases := map[string]struct {
		reply    Engine
		wantCode int
		wantBody string
	}{
		"template writer": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{})},
			wantCode: http.StatusForbidden,
			wantBody: fmt.Sprintf("<p>%s</p>", http.StatusText(http.StatusForbidden)),
		},
		"json writer": {
			reply:    Engine{Writer: NewJSONWriter()},
			wantCode: http.StatusForbidden,
			wantBody: `{"error":"Forbidden"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.Forbidden(w)
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	cases := map[string]struct {
		reply    Engine
		wantCode int
		wantBody string
	}{
		"template writer": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{})},
			wantCode: http.StatusNotFound,
			wantBody: fmt.Sprintf("<p>%s</p>", http.StatusText(http.StatusNotFound)),
		},
		"json writer": {
			reply:    Engine{Writer: NewJSONWriter()},
			wantCode: http.StatusNotFound,
			wantBody: `{"error":"Not Found"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.NotFound(w)
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	cases := map[string]struct {
		reply     Engine
		allow     []string
		wantCode  int
		wantBody  string
		wantAllow string
	}{
		"template writer; allow one": {
			reply:     Engine{Writer: NewTemplateWriter(map[string]*template.Template{})},
			allow:     []string{http.MethodGet},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  fmt.Sprintf("<p>%s</p>", http.StatusText(http.StatusMethodNotAllowed)),
			wantAllow: http.MethodGet,
		},
		"template writer; allow multiple": {
			reply:     Engine{Writer: NewTemplateWriter(map[string]*template.Template{})},
			allow:     []string{http.MethodGet, http.MethodPost},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  fmt.Sprintf("<p>%s</p>", http.StatusText(http.StatusMethodNotAllowed)),
			wantAllow: strings.Join([]string{http.MethodGet, http.MethodPost}, ", "),
		},
		"json writer; allow one": {
			reply:     Engine{Writer: NewJSONWriter()},
			allow:     []string{http.MethodGet},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  `{"error":"Method Not Allowed"}`,
			wantAllow: http.MethodGet,
		},
		"json writer; allow multiple": {
			reply:     Engine{Writer: NewJSONWriter()},
			allow:     []string{http.MethodGet, http.MethodPost},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  `{"error":"Method Not Allowed"}`,
			wantAllow: strings.Join([]string{http.MethodGet, http.MethodPost}, ", "),
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.MethodNotAllowed(w, c.allow...)
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
			if got := w.Header().Get("Allow"); got != c.wantAllow {
				t.Errorf(errorString, got, c.wantAllow)
			}
		})
	}
}

func TestInternalServerError(t *testing.T) {
	cases := map[string]struct {
		reply    Engine
		wantCode int
		wantBody string
	}{
		"template writer": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{})},
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf("<p>%s</p>", http.StatusText(http.StatusInternalServerError)),
		},
		"json writer": {
			reply:    Engine{Writer: NewJSONWriter()},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"Internal Server Error"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.InternalServerError(w, fmt.Errorf("sample error"))
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestReplyOrError(t *testing.T) {
	cases := map[string]struct {
		reply    Engine
		code     int
		opts     Options
		wantCode int
		wantBody string
	}{
		"template writer, no such template, debug false": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{})},
			code:     http.StatusOK,
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf("<p>%s</p>", http.StatusText(http.StatusInternalServerError)),
		},
		"template writer, no such template, debug true": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{}), Debug: true},
			code:     http.StatusOK,
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf("<p>%s</p>", "no such template &#39;foo&#39;"),
		},
		"template writer, no template": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{"foo": foo})},
			code:     http.StatusOK,
			wantCode: http.StatusOK,
			wantBody: "Hello, Sherlock",
		},
		"json writer, error - fail encode, debug false": {
			reply:    Engine{Writer: NewJSONWriter()},
			code:     http.StatusOK,
			opts:     Options{Data: map[string]interface{}{"foo": make(chan int)}},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"Internal Server Error"}`,
		},
		"json writer, error - fail encode debug true": {
			reply:    Engine{Writer: NewJSONWriter(), Debug: true},
			code:     http.StatusOK,
			opts:     Options{Data: map[string]interface{}{"foo": make(chan int)}},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"json: unsupported type: chan int"}`,
		},
		"json writer": {
			reply:    Engine{Writer: NewJSONWriter()},
			code:     http.StatusOK,
			wantCode: http.StatusOK,
			wantBody: `{"name":"Sherlock"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if c.opts == (Options{}) {
				c.opts = Options{
					Key:  "foo",
					Name: "base",
					Data: struct {
						Name string `json:"name"`
					}{
						Name: "Sherlock",
					},
				}
			}
			c.reply.ReplyOrError(w, c.code, c.opts)
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestOK(t *testing.T) {
	cases := map[string]struct {
		reply    Engine
		wantCode int
		wantBody string
	}{
		"template writer, no template": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{"foo": foo})},
			wantCode: http.StatusOK,
			wantBody: "Hello, Sherlock",
		},
		"json writer": {
			reply:    Engine{Writer: NewJSONWriter()},
			wantCode: http.StatusOK,
			wantBody: `{"name":"Sherlock"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.OK(w, Options{
				Key:  "foo",
				Name: "base",
				Data: struct {
					Name string `json:"name"`
				}{
					Name: "Sherlock",
				},
			})
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestCreated(t *testing.T) {
	cases := map[string]struct {
		reply    Engine
		wantCode int
		wantBody string
	}{
		"template writer; no template": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{"foo": foo})},
			wantCode: http.StatusCreated,
			wantBody: "Hello, Sherlock",
		},
		"json writer": {
			reply:    Engine{Writer: NewJSONWriter()},
			wantCode: http.StatusCreated,
			wantBody: `{"name":"Sherlock"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.Created(w, Options{
				Key:  "foo",
				Name: "base",
				Data: struct {
					Name string `json:"name"`
				}{
					Name: "Sherlock",
				},
			})
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestNoContent(t *testing.T) {
	cases := map[string]struct {
		reply    Engine
		wantCode int
		wantBody string
	}{
		"template writer; no template": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{"quux": quux})},
			wantCode: http.StatusNoContent,
			wantBody: "",
		},
		"json writer": {
			reply:    Engine{Writer: NewJSONWriter()},
			wantCode: http.StatusNoContent,
			wantBody: "null",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.NoContent(w)
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}
