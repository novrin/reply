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

// Error replies to the request with the given HTTP status code and its text
// description in a JSON error key.
func (jw JSONWriter) Error(w http.ResponseWriter, error string, code int) {
	_ = jw.Encode(map[string]string{"error": error})
	w.WriteHeader(code)
	_, _ = jw.WriteTo(w)
}

// Write replies to a request with the given status code and opts Data encoded
// to JSON. The encoding is first written to a buffer. If an error occurs, it
// replies with an Internal Server Error. Otherwise, it writes the given status
// code and the encoded data.
func (jw JSONWriter) Reply(w http.ResponseWriter, code int, opts Options) {
	if err := jw.Encode(opts.Data); err != nil {
		message := err.Error()
		if !opts.Debug {
			message = http.StatusText(http.StatusInternalServerError)
		}
		jw.Error(w, message, http.StatusInternalServerError)
	}
	w.WriteHeader(code)
	_, _ = jw.buffer.WriteTo(w)
}

// NewJSONWriter returns a new JSONWriter with an empty buffer.
func NewJSONWriter() *JSONWriter {
	return &JSONWriter{buffer: new(bytes.Buffer)}
}
