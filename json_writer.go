package reply

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONWriter implement Writer for JSON responses.
type JSONWriter struct{}

// Error replies to the request with the given HTTP status code and its text
// description in a JSON error key.
func (jw JSONWriter) Error(w http.ResponseWriter, statusCode int) {
	jw.Write(w, statusCode, Options{
		Data: map[string]string{"error": http.StatusText(statusCode)},
	})
}

// Write replies to a request with the given status code and opts Data encoded
// to JSON. The encoding is first written to a buffer. If an error occurs, it
// replies with an Internal Server Error. Otherwise, it writes the given status
// code and the encoded data.
func (jw JSONWriter) Write(w http.ResponseWriter, statusCode int, opts Options) {
	w.Header().Set("Content-Type", "application/json")
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(opts.Data); err != nil {
		statusCode = http.StatusInternalServerError
		buffer.Reset()
		e, _ := json.Marshal(map[string]string{"error": fmt.Sprintf("failed to marshal %v", err)})
		buffer.Write(e)
	}
	w.WriteHeader(statusCode)
	_, _ = buffer.WriteTo(w)
}
