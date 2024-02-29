package reply

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
)

// Template writer implements Writer for template responses.
type TemplateWriter struct {
	Templates map[string]*template.Template
}

// Options represents fields used in Reply.
type Options struct {
	// TemplateKey defines a lookup in an TemplateWriter's Templates.
	// This is always required for a TemplateWriter; if not supplied,
	// its Reply will write an Internal Server Error.
	TemplateKey string

	// TemplateName defines an optional named template to execute.
	// This optional even for a TemplateWriter.
	TemplateName string

	// Data defines data for use in a reply.
	Data any
}

// Reply sends an HTTP status response header with the given status code and
// writes an executed template to w using the opts provided. If an error occurs
// at template execution, the function exits and does not write to w.
func (tw *TemplateWriter) Reply(w http.ResponseWriter, code int, opts Options) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	tmpl, ok := tw.Templates[opts.TemplateKey]
	if !ok {
		return fmt.Errorf("no such template '%s'", opts.TemplateKey)
	}
	buf := new(bytes.Buffer)
	if opts.TemplateName != "" {
		if err := tmpl.ExecuteTemplate(buf, opts.TemplateName, opts.Data); err != nil {
			return err
		}
	} else {
		if err := tmpl.Execute(buf, opts.Data); err != nil {
			return err
		}
	}
	w.WriteHeader(code)
	_, _ = buf.WriteTo(w)
	return nil
}

// Error sends an HTTP response header with the given status code and writes
// tw's executed "error.html" template to w. It does not otherwise end the
// request; the caller should ensure no further writes are done to w.
func (tw *TemplateWriter) Error(w http.ResponseWriter, error string, code int) {
	_ = tw.Reply(w, code, Options{
		TemplateKey: "error.html",
		Data:        struct{ Error string }{Error: error},
	})
}

// NewTemplateWriter returns a new TemplateWriter with the given templates and
// an empty buffer. If no "error.html" or "no_content.html" are supplied in
// templates, defaults are parsed and used.
func NewTemplateWriter(templates map[string]*template.Template) *TemplateWriter {
	if _, ok := templates["error.html"]; !ok {
		templates["error.html"] = template.
			Must(template.New("error.html").Parse("<p>{{.Error}}</p>"))
	}
	if _, ok := templates["no_content.html"]; !ok {
		templates["no_content.html"] = template.
			Must(template.New("no_content.html").Parse(""))
	}
	return &TemplateWriter{Templates: templates}
}

// TemplateMap returns a map of string to HTML template using fsys as its source.
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
