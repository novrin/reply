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
	// errorHTML is the HTML content for the error.html template. It is
	// used in NewTemplateWriter as a default if the given templates map does
	// not contain the key "error.html".
	errorHTML string = `<p>{{.Error}}</p>`

	// NoContentHTML is blank HTML content for the no_content.html template. It
	// is used in NewTemplateWriter to supply a default "no_content.html"
	// template if absent in the templates map.
	NoContentHTML string = ""
)

// Template writer implements Writer for template responses.
type TemplateWriter struct {
	Templates map[string]*template.Template
	buffer    *bytes.Buffer
}

// Execute applies the template mapped to key to the given data object, writing
// the output to tw's buffer. If an error occurs executing the template or
// writing its output, execution stops, the buffer is reset, and the error is
// returned. It is called in Reply to prevent partial HTML responses if an
// error occurs.
func (tw *TemplateWriter) Execute(key string, data any) error {
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

// ExecuteTemplate applies the template mapped to key that has the given name to
// the given data object and writes the output to tw's buffer. If an error
// occurs executing the template or writing its output, execution stops, the
// buffer is reset, and the error is returned. It is called in Reply to prevent
// partial HTML responses if an error occurs.
func (tw *TemplateWriter) ExecuteTemplate(key string, name string, data any) error {
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

// WriteTo writes data to w until tw's buffer is drained or an error occurs.
// Any values returned by the buffer's WriteTo are returned.
func (tw *TemplateWriter) WriteTo(w http.ResponseWriter) (int64, error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	return tw.buffer.WriteTo(w)
}

// Error replies to a request with tw's "error.html" template, given error
// message, and HTTP code. It does not otherwise end the request; the caller
// should ensure no further writes are done to w.
func (tw *TemplateWriter) Error(w http.ResponseWriter, error string, code int) {
	_ = tw.Execute("error.html", struct{ Error string }{Error: error})
	w.WriteHeader(code)
	_, _ = tw.WriteTo(w)
}

// Options represents fields used in Reply.
type Options struct {
	// Debug defines whether transparent error strings encountered in Reply are
	// sent in responses. If debug is false, the error message will simply be
	// the text representation of the error code.
	Debug bool

	// Key defines a lookup in an TemplateWriter's Templates. This is always
	// required for a TemplateWriter; if not supplied, its Reply will write an
	// Internal Server Error.
	Key string

	// Name defines an optional named template to execute.
	Name string

	// Data defines data for use in a TemplateWriter's Execution or JSON output.
	Data any
}

// Reply executes templates according to the given options. If an error occurs
// at any point in the process, it replies to the request with an Internal
// Server Error. The transparency of the errors are denoted by opts.Debug.
func (tw *TemplateWriter) Reply(w http.ResponseWriter, code int, opts Options) {
	var err error
	if opts.Name != "" {
		err = tw.ExecuteTemplate(opts.Key, opts.Name, opts.Data)
	} else {
		err = tw.Execute(opts.Key, opts.Data)
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
	_, _ = tw.WriteTo(w)
}

// NewTemplateWriter returns a new TemplateWriter with the given templates and
// an empty buffer. If no "error.html" or "no_content.html" are supplied in
// templates, the default HTML provided in vars errorHTML and NoContentHTML are
// parsed and used.
func NewTemplateWriter(templates map[string]*template.Template) *TemplateWriter {
	if _, ok := templates["error.html"]; !ok {
		templates["error.html"] = template.Must(template.New("error.html").Parse(errorHTML))
	}
	if _, ok := templates["no_content.html"]; !ok {
		templates["no_content.html"] = template.Must(template.New("no_content.html").Parse(NoContentHTML))
	}
	return &TemplateWriter{Templates: templates, buffer: new(bytes.Buffer)}
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
