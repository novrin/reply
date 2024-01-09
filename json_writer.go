package reply

import (
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

// Write sets the Content-Type header with "application/json". If it fails to
// marshal the Data provided in opts, it writes a Internal Server Error.
// Otherwise it write the given status code and the marshalled data.
func (jw JSONWriter) Write(w http.ResponseWriter, statusCode int, opts Options) {
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(opts.Data)
	if err != nil {
		statusCode = http.StatusInternalServerError
		j = []byte(fmt.Sprint("failed to marshal ", err))
	}
	w.WriteHeader(statusCode)
	w.Write(j)
}
