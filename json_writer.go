package reply

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// JSONWriter implements Writer for JSON responses.
type JSONWriter struct{}

// Reply sends an HTTP status response header with the given status code
// and writes encoded JSON to w using the opts provided. If an error occurs
// at encoding, the function exits and does not write to w.
func (jw JSONWriter) Reply(w http.ResponseWriter, code int, opts Options) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(opts.Data); err != nil {
		return err
	}
	w.WriteHeader(code)
	_, _ = buf.WriteTo(w)
	return nil
}

// Error sends an HTTP response header with the given status code
// and writes an encoded JSON error to w.
func (jw JSONWriter) Error(w http.ResponseWriter, error string, code int) {
	_ = jw.Reply(w, code, Options{Data: map[string]string{"error": error}})
}
