package reply

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
)

// Template writer implements Writer for template responses.
type TemplateWriter struct {
	Templates map[string]*template.Template
}

// Error replies to the request with the given HTTP status code and its text
// description in plain text.
func (tw TemplateWriter) Error(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

// Options represents fields used in Write.
//   - Template defines the lookup in an Engine's Templates map
//   - Invoke defines an optional named template to invoke
//   - Data defines the data for use in a template
type Options struct {
	Template string
	Invoke   string
	Data     interface{}
}

// Write searches tw's Templates map using the Template key provided in opts.
// If a key is not provided, it simple writes the given status code. If a key
// given but does not exist, it throws an error. Otherwise, the template is
// applied and written on a buffer. The buffer then attempts to write to the
// writer which succeeds if and only if the attempt yields no errors.
func (tw TemplateWriter) Write(w http.ResponseWriter, statusCode int, opts Options) {
	if opts.Template == "" {
		w.WriteHeader(statusCode)
		return
	}
	tmpl, ok := tw.Templates[opts.Template]
	if !ok {
		tw.Error(w, http.StatusInternalServerError)
		return
	}
	buf := new(bytes.Buffer)
	name := opts.Template
	if opts.Invoke != "" {
		name = opts.Invoke
	}
	if err := tmpl.ExecuteTemplate(buf, name, opts.Data); err != nil {
		tw.Error(w, http.StatusInternalServerError)
		return
	}
	if _, err := buf.WriteTo(w); err != nil {
		tw.Error(w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
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
