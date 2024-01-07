package reply

import (
	"bytes"
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

const errorString = "\nGot:\t%#v\nWant:\t%#v\n"

func TestError(t *testing.T) {
	cases := map[string]int{
		"ok":          http.StatusOK,
		"bad Request": http.StatusBadRequest,
		"not found":   http.StatusNotFound,
	}
	tw := TemplateWriter{}
	for name, code := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tw.Error(w, code)
			if got := w.Code; got != code {
				t.Fatalf(errorString, got, code)
			}
		})
	}
}

// errorWriter wraps a response writter. It is used in this testing context to
// throw an intentional error when calling it's Write receiver.
type errorWriter struct {
	http.ResponseWriter
	Code int
	Body *bytes.Buffer
}

func (ew *errorWriter) WriteHeader(statusCode int) {
	ew.Code = statusCode
	ew.ResponseWriter.WriteHeader(statusCode)
}

func (ew *errorWriter) Write(b []byte) (int, error) {
	ew.Body.Reset()
	ew.Body.Write(b)
	return 0, errors.New("intentional error on write")
}

func TestWrite(t *testing.T) {
	templates := map[string]*template.Template{
		"foo": template.Must(template.New("foo").
			Parse(`{{define "base"}}Hello, {{.Name}}{{end}}`)),
	}
	cases := map[string]struct {
		tw       TemplateWriter
		code     int
		opts     Options
		wantCode int
		wantBody string
	}{
		"nil Templates; nil Options": {
			tw:       TemplateWriter{},
			code:     http.StatusOK,
			opts:     Options{},
			wantCode: http.StatusOK,
			wantBody: "",
		},
		"nil Options": {
			tw:       TemplateWriter{Templates: templates},
			code:     http.StatusOK,
			opts:     Options{},
			wantCode: http.StatusOK,
			wantBody: "",
		},
		"nil Templates": {
			tw:       TemplateWriter{},
			code:     http.StatusOK,
			opts:     Options{Template: "foo"},
			wantCode: http.StatusInternalServerError,
			wantBody: http.StatusText(http.StatusInternalServerError),
		},
		"template key not ok": {
			tw:       TemplateWriter{Templates: templates},
			code:     http.StatusOK,
			opts:     Options{Template: "boo"},
			wantCode: http.StatusInternalServerError,
			wantBody: http.StatusText(http.StatusInternalServerError),
		},
		"template key ok; no invoke": {
			tw:       TemplateWriter{Templates: templates},
			code:     http.StatusOK,
			opts:     Options{Template: "foo"},
			wantCode: http.StatusOK,
			wantBody: "",
		},
		"template key ok; bad invoke": {
			tw:       TemplateWriter{Templates: templates},
			code:     http.StatusOK,
			opts:     Options{Template: "foo", Invoke: "bass"},
			wantCode: http.StatusInternalServerError,
			wantBody: http.StatusText(http.StatusInternalServerError),
		},
		"template key ok; good invoke": {
			tw:       TemplateWriter{Templates: templates},
			code:     http.StatusOK,
			opts:     Options{Template: "foo", Invoke: "base"},
			wantCode: http.StatusOK,
			wantBody: "Hello,",
		},
		"template key ok; good invoke; bad data": {
			tw:   TemplateWriter{Templates: templates},
			code: http.StatusOK,
			opts: Options{
				Template: "foo",
				Invoke:   "base",
				Data:     struct{ Fame string }{Fame: "Stars"},
			},
			wantCode: http.StatusInternalServerError,
			wantBody: http.StatusText(http.StatusInternalServerError),
		},
		"template key ok; good invoke; good data": {
			tw:   TemplateWriter{Templates: templates},
			code: http.StatusOK,
			opts: Options{
				Template: "foo",
				Invoke:   "base",
				Data:     struct{ Name string }{Name: "Stars"},
			},
			wantCode: http.StatusOK,
			wantBody: "Hello, Stars",
		},
		"everything ok; error in buff": {
			tw:   TemplateWriter{Templates: templates},
			code: http.StatusOK,
			opts: Options{
				Template: "foo",
				Invoke:   "base",
				Data:     struct{ Name string }{Name: "Stars"},
			},
			wantCode: http.StatusInternalServerError,
			wantBody: http.StatusText(http.StatusInternalServerError),
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var gotCode int
			var gotBody string
			switch name {
			case "everything ok; error in buff":
				w := &errorWriter{ResponseWriter: httptest.NewRecorder(), Body: new(bytes.Buffer)}
				c.tw.Write(w, c.code, c.opts)
				gotCode = w.Code
				gotBody = w.Body.String()
			default:
				w := httptest.NewRecorder()
				c.tw.Write(w, c.code, c.opts)
				gotCode = w.Code
				gotBody = w.Body.String()
			}
			if gotCode != c.wantCode {
				t.Fatalf(errorString, gotCode, c.wantCode)
			}
			gotBody = strings.TrimSpace(gotBody)
			if gotBody != c.wantBody {
				t.Fatalf(errorString, gotBody, c.wantBody)
			}
		})
	}
}

func TestTemplateMap(t *testing.T) {
	templates := fstest.MapFS{
		"html/base.html":        {Data: []byte(`{{define "base"}}Base here. {{template "main" .}}{{end}}`)},
		"html/pages/hello.html": {Data: []byte(`{{define "main"}}Hello, {{.Name}}{{end}}`)},
		"html/pages/hey.html":   {Data: []byte(`{{define "main"}}Hey, {{.Name}}{{end}}`)},
	}
	badTemplates := fstest.MapFS{
		"html/base.html":        {Data: []byte(`{{define "base"}}Base here. {{template "main" .}}{{end}}`)},
		"html/pages/hello.html": {Data: []byte(`{{define "main"}`)}, // invalid template
	}
	funcs := template.FuncMap{
		"uppercase": func(s string) string { return strings.ToUpper(s) },
	}
	cases := map[string]struct {
		fsys     fs.FS
		src      string
		base     string
		wantErr  bool
		wantLen  int
		wantKeys []string
	}{
		"fs Glob error": {
			fsys:    templates,
			src:     "[",
			wantErr: true,
		},
		"template parse error": {
			fsys:    badTemplates,
			src:     "html/pages/*.html",
			base:    "html/base.html",
			wantErr: true,
		},
		"no error": {
			fsys:     templates,
			src:      "html/pages/*.html",
			base:     "html/base.html",
			wantErr:  false,
			wantLen:  2,
			wantKeys: []string{"hello.html", "hey.html"},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			pages, err := TemplateMap(c.fsys, c.src, c.base, funcs)
			if (err != nil) != c.wantErr {
				t.Fatalf(errorString, err, c.wantErr)
			}
			if got := len(pages); got != c.wantLen {
				t.Fatalf(errorString, got, c.wantLen)
			}
			for _, key := range c.wantKeys {
				if _, ok := pages[key]; !ok {
					t.Fatal("wanted 'hello.html' in pages")
				}
			}
		})
	}
}
