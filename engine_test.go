package reply

import (
	"fmt"
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
			reply:    Engine{TemplateWriter{}},
			wantCode: http.StatusBadRequest,
			wantBody: http.StatusText(http.StatusBadRequest),
		},
		"json writer": {
			reply:    Engine{JSONWriter{}},
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"Bad Request"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.BadRequest(w)
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
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
			reply:    Engine{TemplateWriter{}},
			wantCode: http.StatusUnauthorized,
			wantBody: http.StatusText(http.StatusUnauthorized),
		},
		"json writer": {
			reply:    Engine{JSONWriter{}},
			wantCode: http.StatusUnauthorized,
			wantBody: `{"error":"Unauthorized"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.Unauthorized(w)
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
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
			reply:    Engine{TemplateWriter{}},
			wantCode: http.StatusForbidden,
			wantBody: http.StatusText(http.StatusForbidden),
		},
		"json writer": {
			reply:    Engine{JSONWriter{}},
			wantCode: http.StatusForbidden,
			wantBody: `{"error":"Forbidden"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.Forbidden(w)
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
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
			reply:    Engine{TemplateWriter{}},
			wantCode: http.StatusNotFound,
			wantBody: http.StatusText(http.StatusNotFound),
		},
		"json writer": {
			reply:    Engine{JSONWriter{}},
			wantCode: http.StatusNotFound,
			wantBody: `{"error":"Not Found"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.NotFound(w)
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
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
			reply:     Engine{TemplateWriter{}},
			allow:     []string{http.MethodGet},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  http.StatusText(http.StatusMethodNotAllowed),
			wantAllow: http.MethodGet,
		},
		"template writer; allow multiple": {
			reply:     Engine{TemplateWriter{}},
			allow:     []string{http.MethodGet, http.MethodPost},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  http.StatusText(http.StatusMethodNotAllowed),
			wantAllow: strings.Join([]string{http.MethodGet, http.MethodPost}, ", "),
		},
		"json writer; allow one": {
			reply:     Engine{JSONWriter{}},
			allow:     []string{http.MethodGet},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  `{"error":"Method Not Allowed"}`,
			wantAllow: http.MethodGet,
		},
		"json writer; allow multiple": {
			reply:     Engine{JSONWriter{}},
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
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
			}
			if got := w.Header().Get("Allow"); got != c.wantAllow {
				t.Fatalf(errorString, got, c.wantAllow)
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
			reply:    Engine{TemplateWriter{}},
			wantCode: http.StatusInternalServerError,
			wantBody: http.StatusText(http.StatusInternalServerError),
		},
		"json writer": {
			reply:    Engine{JSONWriter{}},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"Internal Server Error"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.InternalServerError(w, fmt.Errorf("sample error"))
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
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
		"template writer": {
			reply:    Engine{TemplateWriter{}},
			wantCode: http.StatusOK,
			wantBody: "",
		},
		"json writer": {
			reply:    Engine{JSONWriter{}},
			wantCode: http.StatusOK,
			wantBody: "null",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.OK(w, Options{})
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
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
		"template writer": {
			reply:    Engine{TemplateWriter{}},
			wantCode: http.StatusCreated,
			wantBody: "",
		},
		"json writer": {
			reply:    Engine{JSONWriter{}},
			wantCode: http.StatusCreated,
			wantBody: "null",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.Created(w, Options{})
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
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
		"template writer": {
			reply:    Engine{TemplateWriter{}},
			wantCode: http.StatusNoContent,
			wantBody: "",
		},
		"json writer": {
			reply:    Engine{JSONWriter{}},
			wantCode: http.StatusNoContent,
			wantBody: "null",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.NoContent(w)
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
			}
		})
	}
}
