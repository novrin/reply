package reply

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

const errorString = "\nGot:\t%#v\nWant:\t%#v\n"

var foo = template.Must(template.New("foo").Parse(`{{define "base"}}Hello, {{.Name}}{{end}}`))
var bar = template.Must(template.New("bar").Parse(`{{define "base"}}Hiya{{end}}`))
var baz = template.Must(template.New("baz").Parse(`Hiya, {{.Name}}`))
var quux = template.Must(template.New("baz").Parse(`HELLO`))

func TestNewTemplateWriter(t *testing.T) {
	cases := map[string]struct {
		templates map[string]*template.Template
		want      []string
	}{
		"empty templates has default error.html and no_content.html": {
			templates: map[string]*template.Template{},
			want:      []string{"error.html", "no_content.html"},
		},
		"one template": {
			templates: map[string]*template.Template{"foo": foo},
			want:      []string{"error.html", "no_content.html", "foo"},
		},
		"many templates": {
			templates: map[string]*template.Template{"foo": foo, "bar": bar},
			want:      []string{"error.html", "no_content.html", "foo", "bar"},
		},
	}
	for name, c := range cases {
		tw := NewTemplateWriter(c.templates)
		t.Run(name, func(t *testing.T) {
			if got := len(tw.Templates); got != len(c.want) {
				t.Fatalf(errorString, got, c.want)
			}
			for _, key := range c.want {
				if _, ok := tw.Templates[key]; !ok {
					t.Fatalf("absent key '%s' in writer templates", key)
				}
			}
		})
	}
}

func TestExecute(t *testing.T) {
	cases := map[string]struct {
		key      string
		wantErr  bool
		wantBody string
	}{
		"error - no such template": {
			key:      "foo",
			wantErr:  true,
			wantBody: "",
		},
		"error - buffer execute failed": {
			key:      "baz",
			wantErr:  true,
			wantBody: "",
		},
		"ok": {
			key:      "error.html",
			wantErr:  false,
			wantBody: "<p>qux</p>",
		},
	}
	tw := NewTemplateWriter(map[string]*template.Template{"bar": bar, "baz": baz})
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := tw.Execute(c.key, struct{ Error string }{Error: "qux"})
			if err != nil && !c.wantErr {
				t.Fatalf("got unwanted error - %s", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("wanted an error but didn't get one")
			}
			tw.WriteTo(w)
			if got := w.Body.String(); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestExecuteTemplate(t *testing.T) {
	cases := map[string]struct {
		key      string
		wantErr  bool
		wantBody string
	}{
		"error - no such template": {
			key:      "foo",
			wantErr:  true,
			wantBody: "",
		},
		"error - buffer execute failed": {
			key:      "error.html",
			wantErr:  true,
			wantBody: "",
		},
		"ok": {
			key:      "bar",
			wantErr:  false,
			wantBody: "Hiya",
		},
	}
	tw := NewTemplateWriter(map[string]*template.Template{"bar": bar, "baz": baz})
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := tw.ExecuteTemplate(c.key, "base", struct{}{})
			if err != nil && !c.wantErr {
				t.Fatalf("got unwanted error - %s", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("wanted an error but didn't get one")
			}
			tw.WriteTo(w)
			if got := w.Body.String(); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestTemplateError(t *testing.T) {
	cases := map[string]int{
		"bad Request": http.StatusBadRequest,
		"not found":   http.StatusNotFound,
	}
	tw := NewTemplateWriter(map[string]*template.Template{})
	for name, code := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := http.StatusText(code)
			tw.Error(w, body, code)
			if got := w.Code; got != code {
				t.Fatalf(errorString, got, code)
			}
			want := fmt.Sprintf("<p>%s</p>", body)
			if got := w.Body.String(); got != want {
				t.Fatalf(errorString, got, "")
			}
		})
	}
}

func TestTemplateReply(t *testing.T) {
	cases := map[string]struct {
		writer   Writer
		code     int
		opts     Options
		wantCode int
		wantBody string
	}{
		"Options nil": {
			writer:   NewTemplateWriter(map[string]*template.Template{}),
			code:     http.StatusOK,
			opts:     Options{},
			wantCode: http.StatusInternalServerError,
			wantBody: "<p>Internal Server Error</p>",
		},
		"error - no such template": {
			writer:   NewTemplateWriter(map[string]*template.Template{}),
			code:     http.StatusOK,
			opts:     Options{Key: "foo"},
			wantCode: http.StatusInternalServerError,
			wantBody: "<p>Internal Server Error</p>",
		},
		"template key ok; name nil": {
			writer:   NewTemplateWriter(map[string]*template.Template{"quux": quux}),
			code:     http.StatusOK,
			opts:     Options{Key: "quux"},
			wantCode: http.StatusOK,
			wantBody: "HELLO",
		},
		"template key ok; name not ok": {
			writer:   NewTemplateWriter(map[string]*template.Template{"foo": foo}),
			code:     http.StatusOK,
			opts:     Options{Key: "foo", Name: "bass"},
			wantCode: http.StatusInternalServerError,
			wantBody: "<p>Internal Server Error</p>",
		},
		"template key ok; name ok": {
			writer: NewTemplateWriter(map[string]*template.Template{"foo": foo}),
			code:   http.StatusOK,
			opts: Options{
				Key:  "foo",
				Name: "base",
				Data: struct{ Name string }{Name: "Sherlock"},
			},
			wantCode: http.StatusOK,
			wantBody: "Hello, Sherlock",
		},
		"template key ok; name ok; data not ok": {
			writer: NewTemplateWriter(map[string]*template.Template{"foo": foo}),
			code:   http.StatusOK,
			opts: Options{
				Key:  "foo",
				Name: "base",
				Data: struct{ Mame string }{Mame: "Sherlock"},
			},
			wantCode: http.StatusInternalServerError,
			wantBody: "<p>Internal Server Error</p>",
		},
		"template key ok; name ok; data ok": {
			writer: NewTemplateWriter(map[string]*template.Template{"foo": foo}),
			code:   http.StatusOK,
			opts: Options{
				Key:  "foo",
				Name: "base",
				Data: struct{ Name string }{Name: "Sherlock"},
			},
			wantCode: http.StatusOK,
			wantBody: "Hello, Sherlock",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.writer.Reply(w, c.code, c.opts)
			if got := w.Code; got != c.wantCode {
				t.Fatalf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Fatalf(errorString, got, c.wantBody)
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
		"error - fs Glob ": {
			fsys:    templates,
			src:     "[",
			wantErr: true,
		},
		"error - template parse": {
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
