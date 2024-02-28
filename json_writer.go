package reply

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// JSONWriter implements Writer for JSON responses.
type JSONWriter struct {
	buffer *bytes.Buffer
}

// Encode writes the JSON encoding of data followed by a newline character to
// jw's buffer. If an error occurs encoding the data or writing its output,
// execution stops, the buffer is reset, and the error is returned.
func (jw JSONWriter) Encode(data any) error {
	if err := json.NewEncoder(jw.buffer).Encode(data); err != nil {
		jw.buffer.Reset()
		return err
	}
	return nil
}

// WriteTo writes data to w until jw's buffer is drained or an error occurs.
// Any values returned by the buffer's WriteTo are returned.
func (jw JSONWriter) WriteTo(w http.ResponseWriter) (int64, error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	return jw.buffer.WriteTo(w)
}

// Error sends an HTTP response header with the given status code and writes an
// encoded JSON error to w.
func (jw JSONWriter) Error(w http.ResponseWriter, error string, code int) {
	_ = jw.Encode(map[string]string{"error": error})
	w.WriteHeader(code)
	_, _ = jw.WriteTo(w)
}

// Reply sends an HTTP status response header with the given status code and
// writes encoded JSON to w using the opts provided. If an error occurs at
// encoding, the function exits and does not write to w.
func (jw JSONWriter) Reply(w http.ResponseWriter, code int, opts Options) error {
	if err := jw.Encode(opts.Data); err != nil {
		return err
	}
	w.WriteHeader(code)
	_, _ = jw.WriteTo(w)
	return nil
}

// NewJSONWriter returns a new JSONWriter with an empty buffer.
func NewJSONWriter() *JSONWriter {
	return &JSONWriter{buffer: new(bytes.Buffer)}
}
