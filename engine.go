package reply

import (
	"net/http"
)

// Writer is used by an Engine to construct replies to HTTP server requests.
type Writer interface {
	WriteTo(w http.ResponseWriter) (int64, error)
	Error(w http.ResponseWriter, error string, code int)
	Reply(w http.ResponseWriter, code int, opts Options) error
}

// Engine provides convenience reply methods by wrapping its embedded Writer's
// Error and Reply.
type Engine struct {
	// Debug defines whether error strings encountered in the Writer's Reply are
	// sent in responses. If debug is false, the error string will simply be the
	// plain text representation of the error code.
	Debug bool

	// Writer is an interface used to construct replies to HTTP server requests.
	Writer
}

// BadRequest replies with HTTP Status 400 Bad Request.
func (e Engine) BadRequest(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

// Unauthorized replies with HTTP Status 401 Unauthorized.
func (e Engine) Unauthorized(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// Forbidden replies with HTTP Status 403 Forbidden.
func (e Engine) Forbidden(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

// NotFound replies with HTTP Status 404 Not Found.
func (e Engine) NotFound(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// MethodNotAllowed replies with HTTP Status 405 Method Not Allowed.
func (e Engine) MethodNotAllowed(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// NotAcceptable replies with HTTP Status 406 Not Acceptable.
func (e Engine) NotAcceptable(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
}

// RequestTimeout replies with HTTP Status 408 Request Timeout.
func (e Engine) RequestTimeout(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
}

// Conflict replies with HTTP Status 409 Conflict.
func (e Engine) Conflict(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
}

// Gone replies with HTTP Status 410 Gone.
func (e Engine) Gone(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
}

// UnprocessableEntity replies with HTTP Status 422 Unprocessable Entity.
func (e Engine) UnprocessableEntity(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
}

// TooManyRequests replies with HTTP Status 429 Too Many Requests.
func (e Engine) TooManyRequests(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
}

// InternalServerError replies with HTTP Status 500 Internal Server Error.
func (e Engine) InternalServerError(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// ReplyOrError wraps Reply with error debugging. If an error is encountered in
// Reply, the Writer's Error function is triggered. Error essages are replaced
// with 'Internal Server Error' if e.Debug is false.
func (e Engine) ReplyOrError(w http.ResponseWriter, code int, opts Options) {
	if err := e.Reply(w, code, opts); err != nil {
		code = http.StatusInternalServerError
		if !e.Debug {
			e.Error(w, http.StatusText(code), code)
		} else {
			e.Error(w, err.Error(), code)
		}
	}
}

// OK replies with HTTP 200 Status OK.
func (e Engine) OK(w http.ResponseWriter, opts Options) {
	e.ReplyOrError(w, http.StatusOK, opts)
}

// Created replies with HTTP 201 Status Created.
func (e Engine) Created(w http.ResponseWriter, opts Options) {
	e.ReplyOrError(w, http.StatusCreated, opts)
}

// NoContent replies with HTTP Status 204 No Content.
func (e Engine) NoContent(w http.ResponseWriter) {
	e.ReplyOrError(w, http.StatusNoContent, Options{Key: "no_content.html"})
}
