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

// ProxyAuthRequired replies with HTTP Status 407 Proxy Authentication Required.
func (e Engine) ProxyAuthRequired(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusProxyAuthRequired), http.StatusProxyAuthRequired)
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

// LengthRequired replies with HTTP Status 411 Length Required.
func (e Engine) LengthRequired(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusLengthRequired), http.StatusLengthRequired)
}

// PreconditionFailed replies with HTTP Status 412 Precondition Failed.
func (e Engine) PreconditionFailed(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusPreconditionFailed), http.StatusPreconditionFailed)
}

// RequestEntityTooLarge replies with HTTP Status 413 Request Entity Too Large.
func (e Engine) RequestEntityTooLarge(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
}

// RequestURITooLong replies with HTTP Status 414 Request URI Too Long.
func (e Engine) RequestURITooLong(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusRequestURITooLong), http.StatusRequestURITooLong)
}

// UnsupportedMediaType replies with HTTP Status 415 Unsupported Media Type.
func (e Engine) UnsupportedMediaType(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
}

// RequestedRangeNotSatisfiable replies with HTTP Status 416 Requested Range Not Satisfiable.
func (e Engine) RequestedRangeNotSatisfiable(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusRequestedRangeNotSatisfiable), http.StatusRequestedRangeNotSatisfiable)
}

// ExpectationFailed replies with HTTP Status 417 Expectation Failed.
func (e Engine) ExpectationFailed(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusExpectationFailed), http.StatusExpectationFailed)
}

// Teapot replies with HTTP Status 418 I'm a teapot.
func (e Engine) Teapot(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
}

// MisdirectedRequest replies with HTTP Status 421 Misdirected Request.
func (e Engine) MisdirectedRequest(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusMisdirectedRequest), http.StatusMisdirectedRequest)
}

// UnprocessableEntity replies with HTTP Status 422 Unprocessable Entity.
func (e Engine) UnprocessableEntity(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
}

// Locked replies with HTTP Status 423 Locked.
func (e Engine) Locked(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusLocked), http.StatusLocked)
}

// FailedDependency replies with HTTP Status 424 Failed Dependency.
func (e Engine) FailedDependency(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusFailedDependency), http.StatusFailedDependency)
}

// TooEarly replies with HTTP Status 425 Too Early.
func (e Engine) TooEarly(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusTooEarly), http.StatusTooEarly)
}

// UpgradeRequired replies with HTTP Status 426 Upgrade Required.
func (e Engine) UpgradeRequired(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusUpgradeRequired), http.StatusUpgradeRequired)
}

// PreconditionRequired replies with HTTP Status 428 Precondition Required.
func (e Engine) PreconditionRequired(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusPreconditionRequired), http.StatusPreconditionRequired)
}

// TooManyRequests replies with HTTP Status 429 Too Many Requests.
func (e Engine) TooManyRequests(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
}

// RequestHeaderFieldsTooLarge replies with HTTP Status 431 Request Header Fields Too Large.
func (e Engine) RequestHeaderFieldsTooLarge(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusRequestHeaderFieldsTooLarge), http.StatusRequestHeaderFieldsTooLarge)
}

// UnavailableForLegalReasons replies with HTTP Status 451 Unavailable For Legal Reasons.
func (e Engine) UnavailableForLegalReasons(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusUnavailableForLegalReasons), http.StatusUnavailableForLegalReasons)
}

// InternalServerError replies with HTTP Status 500 Internal Server Error.
func (e Engine) InternalServerError(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// NotImplemented replies with HTTP Status 501 Not Implemented.
func (e Engine) NotImplemented(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// BadGateway replies with HTTP Status 502 Bad Gateway.
func (e Engine) BadGateway(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
}

// ServiceUnavailable replies with HTTP Status 503 Service Unavailable.
func (e Engine) ServiceUnavailable(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

// GatewayTimeout replies with HTTP Status 504 Gateway Timeout.
func (e Engine) GatewayTimeout(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
}

// HTTPVersionNotSupported replies with HTTP Status 505 HTTP Version Not Supported.
func (e Engine) HTTPVersionNotSupported(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusHTTPVersionNotSupported), http.StatusHTTPVersionNotSupported)
}

// VariantAlsoNegotiates replies with HTTP Status 506 Variant Also Negotiates.
func (e Engine) VariantAlsoNegotiates(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusVariantAlsoNegotiates), http.StatusVariantAlsoNegotiates)
}

// InsufficientStorage replies with HTTP Status 507 Insufficient Storage.
func (e Engine) InsufficientStorage(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusInsufficientStorage), http.StatusInsufficientStorage)
}

// LoopDetected replies with HTTP Status 508 Loop Detected.
func (e Engine) LoopDetected(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusLoopDetected), http.StatusLoopDetected)
}

// NetworkAuthenticationRequired replies with HTTP Status 511 Network Authentication Required
func (e Engine) NetworkAuthenticationRequired(w http.ResponseWriter) {
	e.Error(w, http.StatusText(http.StatusNetworkAuthenticationRequired), http.StatusNetworkAuthenticationRequired)
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
