package reply

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
)

var (
	errorHTML     string = `<p>{{.Error}}</p>`
	NoContentHTML string = ""
)

func NewTemplateWriter(templates map[string]*template.Template) *TemplateWriter {
	if _, ok := templates["error.html"]; !ok {
		templates["error.html"] = template.Must(template.New("error.html").Parse(errorHTML))
		templates["no_content.html"] = template.Must(template.New("no_content.html").Parse(NoContentHTML))
	}
	return &TemplateWriter{Templates: templates, buffer: new(bytes.Buffer)}
}

// Template writer implements Writer for template responses.
type TemplateWriter struct {
	Templates map[string]*template.Template
	buffer    *bytes.Buffer
}

func (tw *TemplateWriter) BufferExecute(key string, data any) error {
	template, ok := tw.Templates[key]
	if !ok {
		return fmt.Errorf("no such template '%s'", key)
	}
	if err := template.Execute(tw.buffer, data); err != nil {
		tw.buffer.Reset()
		return err
	}
	return nil
}

func (tw *TemplateWriter) BufferExecuteTemplate(key string, name string, data any) error {
	template, ok := tw.Templates[key]
	if !ok {
		return fmt.Errorf("no such template '%s'", key)
	}
	if err := template.ExecuteTemplate(tw.buffer, name, data); err != nil {
		tw.buffer.Reset()
		return err
	}
	return nil
}

func (tw *TemplateWriter) WriteTo(w http.ResponseWriter) (int64, error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	return tw.buffer.WriteTo(w)
}

// Error replies to the request with the given HTTP status code and its text
// description in plain text.
func (tw *TemplateWriter) Error(w http.ResponseWriter, error string, code int) {
	_ = tw.BufferExecute("error.html", struct{ Error string }{Error: error})
	w.WriteHeader(code)
	_, _ = tw.WriteTo(w)
}

// Options represents fields used in Write.
//   - Debug defines whether transparent error strings are sent in responses
//   - Key defines a lookup in an TemplateWriter's Templates
//   - Name defines a named template to invoke
//   - Data defines data for use in a template or JSON output
type Options struct {
	Debug bool
	Key   string
	Name  string
	Data  any
}

func (tw *TemplateWriter) Reply(w http.ResponseWriter, code int, opts Options) {
	var err error
	if opts.Name != "" {
		err = tw.BufferExecuteTemplate(opts.Key, opts.Name, opts.Data)
	} else {
		err = tw.BufferExecute(opts.Key, opts.Data)
	}
	if err != nil {
		message := err.Error()
		if !opts.Debug {
			message = http.StatusText(http.StatusInternalServerError)
		}
		tw.Error(w, message, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	_, _ = tw.buffer.WriteTo(w)
}

// TemplateMap returns a map of string to HTML template.
// It uses fsys as the source for files.
func TemplateMap(fsys fs.FS, src string, base string, funcs template.FuncMap) (map[string]*template.Template, error) {
	sources, err := fs.Glob(fsys, src)
	if err != nil {
		return nil, err
	}
	cache := map[string]*template.Template{}
	for _, s := range sources {
		name := filepath.Base(s)
		files := []string{s}
		if base != "" {
			files = append([]string{base}, files...)
		}
		tmpl, err := template.New(name).Funcs(funcs).ParseFS(fsys, files...)
		if err != nil {
			return nil, err
		}
		cache[name] = tmpl
	}
	return cache, nil
}
