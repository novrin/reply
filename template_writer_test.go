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
		t.Run(name, func(t *testing.T) {
			tw := NewTemplateWriter(c.templates)
			if got := len(tw.Templates); got != len(c.want) {
				t.Errorf(errorString, got, c.want)
			}
			for _, key := range c.want {
				if _, ok := tw.Templates[key]; !ok {
					t.Errorf("absent key '%s' in writer templates", key)
				}
			}
			w := httptest.NewRecorder()
			tw.Reply(w, http.StatusOK, Options{})
			if got, want := w.Body.String(), ""; got != want {
				t.Errorf(errorString, got, want)
			}
		})
	}
}

func TestTemplateReply(t *testing.T) {
	cases := map[string]struct {
		writer   Writer
		code     int
		opts     Options
		wantErr  bool
		wantCode int
		wantBody string
	}{
		"Options nil": {
			writer:   NewTemplateWriter(map[string]*template.Template{}),
			code:     http.StatusOK,
			opts:     Options{},
			wantErr:  true,
			wantCode: http.StatusOK,
			wantBody: "",
		},
		"error - no such template": {
			writer:   NewTemplateWriter(map[string]*template.Template{}),
			code:     http.StatusOK,
			opts:     Options{TemplateKey: "foo"},
			wantErr:  true,
			wantCode: http.StatusOK,
			wantBody: "",
		},
		"template key ok; name nil": {
			writer:   NewTemplateWriter(map[string]*template.Template{"quux": quux}),
			code:     http.StatusOK,
			opts:     Options{TemplateKey: "quux"},
			wantCode: http.StatusOK,
			wantBody: "HELLO",
		},
		"template key ok; name nil; data not ok": {
			writer: NewTemplateWriter(map[string]*template.Template{"baz": baz}),
			code:   http.StatusOK,
			opts: Options{
				TemplateKey: "baz",
				Data:        struct{ Mame string }{Mame: "Sherlock"},
			},
			wantErr:  true,
			wantCode: http.StatusOK,
			wantBody: "",
		},
		"template key ok; name not ok": {
			writer:   NewTemplateWriter(map[string]*template.Template{"foo": foo}),
			code:     http.StatusOK,
			opts:     Options{TemplateKey: "foo", TemplateName: "bass"},
			wantErr:  true,
			wantCode: http.StatusOK,
			wantBody: "",
		},
		"template key ok; name ok": {
			writer: NewTemplateWriter(map[string]*template.Template{"foo": foo}),
			code:   http.StatusOK,
			opts: Options{
				TemplateKey:  "foo",
				TemplateName: "base",
				Data:         struct{ Name string }{Name: "Sherlock"},
			},
			wantCode: http.StatusOK,
			wantBody: "Hello, Sherlock",
		},
		"template key ok; name ok; data not ok": {
			writer: NewTemplateWriter(map[string]*template.Template{"foo": foo}),
			code:   http.StatusOK,
			opts: Options{
				TemplateKey:  "foo",
				TemplateName: "base",
				Data:         struct{ Mame string }{Mame: "Sherlock"},
			},
			wantErr:  true,
			wantCode: http.StatusOK,
			wantBody: "",
		},
		"template key ok; name ok; data ok": {
			writer: NewTemplateWriter(map[string]*template.Template{"foo": foo}),
			code:   http.StatusOK,
			opts: Options{
				TemplateKey:  "foo",
				TemplateName: "base",
				Data:         struct{ Name string }{Name: "Sherlock"},
			},
			wantCode: http.StatusOK,
			wantBody: "Hello, Sherlock",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := c.writer.Reply(w, c.code, c.opts)
			if (err != nil) != c.wantErr {
				t.Errorf(errorString, err, c.wantErr)
			}
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
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
				t.Errorf(errorString, got, code)
			}
			want := fmt.Sprintf("<p>%s</p>", body)
			if got := w.Body.String(); got != want {
				t.Errorf(errorString, got, "")
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
		fsys             fs.FS
		src              string
		base             string
		wantErr          bool
		wantLen          int
		wantTemplateKeys []string
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
			fsys:             templates,
			src:              "html/pages/*.html",
			base:             "html/base.html",
			wantErr:          false,
			wantLen:          2,
			wantTemplateKeys: []string{"hello.html", "hey.html"},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			pages, err := TemplateMap(c.fsys, c.src, c.base, funcs)
			if (err != nil) != c.wantErr {
				t.Errorf(errorString, err, c.wantErr)
			}
			if got := len(pages); got != c.wantLen {
				t.Errorf(errorString, got, c.wantLen)
			}
			for _, key := range c.wantTemplateKeys {
				if _, ok := pages[key]; !ok {
					t.Error("wanted 'hello.html' in pages")
				}
			}
		})
	}
}
