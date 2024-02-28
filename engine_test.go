package reply

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func errorTemplateBody(code int) string {
	return fmt.Sprintf("<p>%s</p>", http.StatusText(code))
}

func TestGenericErrors(t *testing.T) {
	etw := Engine{Writer: NewTemplateWriter(map[string]*template.Template{})}
	ejw := Engine{Writer: NewJSONWriter()}
	cases := map[string]struct {
		method   func(http.ResponseWriter)
		wantCode int
		wantBody string
	}{
		"400 bad request - tw": {
			method:   etw.BadRequest,
			wantCode: http.StatusBadRequest,
			wantBody: errorTemplateBody(http.StatusBadRequest),
		},
		"400 bad request - jw": {
			method:   ejw.BadRequest,
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"Bad Request"}`,
		},
		"401 unauthorized - tw": {
			method:   etw.Unauthorized,
			wantCode: http.StatusUnauthorized,
			wantBody: errorTemplateBody(http.StatusUnauthorized),
		},
		"401 unauthorized - jw": {
			method:   ejw.Unauthorized,
			wantCode: http.StatusUnauthorized,
			wantBody: `{"error":"Unauthorized"}`,
		},
		"403 forbidden - tw": {
			method:   etw.Forbidden,
			wantCode: http.StatusForbidden,
			wantBody: errorTemplateBody(http.StatusForbidden),
		},
		"403 forbidden - kw": {
			method:   ejw.Forbidden,
			wantCode: http.StatusForbidden,
			wantBody: `{"error":"Forbidden"}`,
		},
		"404 not found - tw": {
			method:   etw.NotFound,
			wantCode: http.StatusNotFound,
			wantBody: errorTemplateBody(http.StatusNotFound),
		},
		"404 not found - jw": {
			method:   ejw.NotFound,
			wantCode: http.StatusNotFound,
			wantBody: `{"error":"Not Found"}`,
		},
		"405 method not allowed - tw": {
			method:   etw.MethodNotAllowed,
			wantCode: http.StatusMethodNotAllowed,
			wantBody: errorTemplateBody(http.StatusMethodNotAllowed),
		},
		"405 method not allowed - jw": {
			method:   ejw.MethodNotAllowed,
			wantCode: http.StatusMethodNotAllowed,
			wantBody: `{"error":"Method Not Allowed"}`,
		},
		"406 not acceptable - tw": {
			method:   etw.NotAcceptable,
			wantCode: http.StatusNotAcceptable,
			wantBody: errorTemplateBody(http.StatusNotAcceptable),
		},
		"406 not acceptable - jw": {
			method:   ejw.NotAcceptable,
			wantCode: http.StatusNotAcceptable,
			wantBody: `{"error":"Not Acceptable"}`,
		},
		"407 proxy authentication required - tw": {
			method:   etw.ProxyAuthRequired,
			wantCode: http.StatusProxyAuthRequired,
			wantBody: errorTemplateBody(http.StatusProxyAuthRequired),
		},
		"407 proxy authentication required - jw": {
			method:   ejw.ProxyAuthRequired,
			wantCode: http.StatusProxyAuthRequired,
			wantBody: `{"error":"Proxy Authentication Required"}`,
		},
		"408 request timeout - tw": {
			method:   etw.RequestTimeout,
			wantCode: http.StatusRequestTimeout,
			wantBody: errorTemplateBody(http.StatusRequestTimeout),
		},
		"408 request timeout - jw": {
			method:   ejw.RequestTimeout,
			wantCode: http.StatusRequestTimeout,
			wantBody: `{"error":"Request Timeout"}`,
		},
		"409 conflict - tw": {
			method:   etw.Conflict,
			wantCode: http.StatusConflict,
			wantBody: errorTemplateBody(http.StatusConflict),
		},
		"409 conflict - jw": {
			method:   ejw.Conflict,
			wantCode: http.StatusConflict,
			wantBody: `{"error":"Conflict"}`,
		},
		"410 gone - tw": {
			method:   etw.Gone,
			wantCode: http.StatusGone,
			wantBody: errorTemplateBody(http.StatusGone),
		},
		"410 gone - jw": {
			method:   ejw.Gone,
			wantCode: http.StatusGone,
			wantBody: `{"error":"Gone"}`,
		},
		"411 length required - tw": {
			method:   etw.LengthRequired,
			wantCode: http.StatusLengthRequired,
			wantBody: errorTemplateBody(http.StatusLengthRequired),
		},
		"411 length required - jw": {
			method:   ejw.LengthRequired,
			wantCode: http.StatusLengthRequired,
			wantBody: `{"error":"Length Required"}`,
		},
		"412 precondition failed - tw": {
			method:   etw.PreconditionFailed,
			wantCode: http.StatusPreconditionFailed,
			wantBody: errorTemplateBody(http.StatusPreconditionFailed),
		},
		"412 precondition failed - jw": {
			method:   ejw.PreconditionFailed,
			wantCode: http.StatusPreconditionFailed,
			wantBody: `{"error":"Precondition Failed"}`,
		},
		"413 request entity too large - tw": {
			method:   etw.RequestEntityTooLarge,
			wantCode: http.StatusRequestEntityTooLarge,
			wantBody: errorTemplateBody(http.StatusRequestEntityTooLarge),
		},
		"413 request entity too large - jw": {
			method:   ejw.RequestEntityTooLarge,
			wantCode: http.StatusRequestEntityTooLarge,
			wantBody: `{"error":"Request Entity Too Large"}`,
		},
		"414 request uri too long - tw": {
			method:   etw.RequestURITooLong,
			wantCode: http.StatusRequestURITooLong,
			wantBody: errorTemplateBody(http.StatusRequestURITooLong),
		},
		"414 request uri too long - jw": {
			method:   ejw.RequestURITooLong,
			wantCode: http.StatusRequestURITooLong,
			wantBody: `{"error":"Request URI Too Long"}`,
		},
		"415 unsupported media type - tw": {
			method:   etw.UnsupportedMediaType,
			wantCode: http.StatusUnsupportedMediaType,
			wantBody: errorTemplateBody(http.StatusUnsupportedMediaType),
		},
		"415 unsupported media type - jw": {
			method:   ejw.UnsupportedMediaType,
			wantCode: http.StatusUnsupportedMediaType,
			wantBody: `{"error":"Unsupported Media Type"}`,
		},
		"416 requested range not satisfiable - tw": {
			method:   etw.RequestedRangeNotSatisfiable,
			wantCode: http.StatusRequestedRangeNotSatisfiable,
			wantBody: errorTemplateBody(http.StatusRequestedRangeNotSatisfiable),
		},
		"416 requested range not satisfiable - jw": {
			method:   ejw.RequestedRangeNotSatisfiable,
			wantCode: http.StatusRequestedRangeNotSatisfiable,
			wantBody: `{"error":"Requested Range Not Satisfiable"}`,
		},
		"417 expectation failed - tw": {
			method:   etw.ExpectationFailed,
			wantCode: http.StatusExpectationFailed,
			wantBody: errorTemplateBody(http.StatusExpectationFailed),
		},
		"417 expectation failed - jw": {
			method:   ejw.ExpectationFailed,
			wantCode: http.StatusExpectationFailed,
			wantBody: `{"error":"Expectation Failed"}`,
		},
		"418 teapot - tw": {
			method:   etw.Teapot,
			wantCode: http.StatusTeapot,
			wantBody: "<p>I&#39;m a teapot</p>",
		},
		"418 teapot - jw": {
			method:   ejw.Teapot,
			wantCode: http.StatusTeapot,
			wantBody: `{"error":"I'm a teapot"}`,
		},
		"421 misdirected request - tw": {
			method:   etw.MisdirectedRequest,
			wantCode: http.StatusMisdirectedRequest,
			wantBody: errorTemplateBody(http.StatusMisdirectedRequest),
		},
		"421 misdirected request - jw": {
			method:   ejw.MisdirectedRequest,
			wantCode: http.StatusMisdirectedRequest,
			wantBody: `{"error":"Misdirected Request"}`,
		},
		"422 unprocessable entity - tw": {
			method:   etw.UnprocessableEntity,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: errorTemplateBody(http.StatusUnprocessableEntity),
		},
		"422 unprocessable entity - jw": {
			method:   ejw.UnprocessableEntity,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"error":"Unprocessable Entity"}`,
		},
		"423 locked - tw": {
			method:   etw.Locked,
			wantCode: http.StatusLocked,
			wantBody: errorTemplateBody(http.StatusLocked),
		},
		"423 locked - jw": {
			method:   ejw.Locked,
			wantCode: http.StatusLocked,
			wantBody: `{"error":"Locked"}`,
		},
		"424 failed dependency - tw": {
			method:   etw.FailedDependency,
			wantCode: http.StatusFailedDependency,
			wantBody: errorTemplateBody(http.StatusFailedDependency),
		},
		"424 failed dependency - jw": {
			method:   ejw.FailedDependency,
			wantCode: http.StatusFailedDependency,
			wantBody: `{"error":"Failed Dependency"}`,
		},
		"425 too early - tw": {
			method:   etw.TooEarly,
			wantCode: http.StatusTooEarly,
			wantBody: errorTemplateBody(http.StatusTooEarly),
		},
		"425 too early - jw": {
			method:   ejw.TooEarly,
			wantCode: http.StatusTooEarly,
			wantBody: `{"error":"Too Early"}`,
		},
		"426 upgrade required - tw": {
			method:   etw.UpgradeRequired,
			wantCode: http.StatusUpgradeRequired,
			wantBody: errorTemplateBody(http.StatusUpgradeRequired),
		},
		"426 upgrade required - jw": {
			method:   ejw.UpgradeRequired,
			wantCode: http.StatusUpgradeRequired,
			wantBody: `{"error":"Upgrade Required"}`,
		},
		"428 precondition required - tw": {
			method:   etw.PreconditionRequired,
			wantCode: http.StatusPreconditionRequired,
			wantBody: errorTemplateBody(http.StatusPreconditionRequired),
		},
		"428 precondition required - jw": {
			method:   ejw.PreconditionRequired,
			wantCode: http.StatusPreconditionRequired,
			wantBody: `{"error":"Precondition Required"}`,
		},
		"429 too many requests - tw": {
			method:   etw.TooManyRequests,
			wantCode: http.StatusTooManyRequests,
			wantBody: errorTemplateBody(http.StatusTooManyRequests),
		},
		"429 too many requests - jw": {
			method:   ejw.TooManyRequests,
			wantCode: http.StatusTooManyRequests,
			wantBody: `{"error":"Too Many Requests"}`,
		},
		"431 request header fields too large - tw": {
			method:   etw.RequestHeaderFieldsTooLarge,
			wantCode: http.StatusRequestHeaderFieldsTooLarge,
			wantBody: errorTemplateBody(http.StatusRequestHeaderFieldsTooLarge),
		},
		"431 request header fields too large - jw": {
			method:   ejw.RequestHeaderFieldsTooLarge,
			wantCode: http.StatusRequestHeaderFieldsTooLarge,
			wantBody: `{"error":"Request Header Fields Too Large"}`,
		},
		"451 unavailable for legal reasons - tw": {
			method:   etw.UnavailableForLegalReasons,
			wantCode: http.StatusUnavailableForLegalReasons,
			wantBody: errorTemplateBody(http.StatusUnavailableForLegalReasons),
		},
		"451 unavailable for legal reasons - jw": {
			method:   ejw.UnavailableForLegalReasons,
			wantCode: http.StatusUnavailableForLegalReasons,
			wantBody: `{"error":"Unavailable For Legal Reasons"}`,
		},
		"500 internal server error - tw": {
			method:   etw.InternalServerError,
			wantCode: http.StatusInternalServerError,
			wantBody: errorTemplateBody(http.StatusInternalServerError),
		},
		"500 internal server error - jw": {
			method:   ejw.InternalServerError,
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"Internal Server Error"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.method(w)
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestReplyOrError(t *testing.T) {
	etw := Engine{Writer: NewTemplateWriter(map[string]*template.Template{})}
	ejw := Engine{Writer: NewJSONWriter()}
	cases := map[string]struct {
		reply    Engine
		code     int
		opts     Options
		wantCode int
		wantBody string
	}{
		"error no such template, debug false - tw": {
			reply:    etw,
			code:     http.StatusOK,
			wantCode: http.StatusInternalServerError,
			wantBody: errorTemplateBody(http.StatusInternalServerError),
		},
		"error no such template, debug true - tw": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{}), Debug: true},
			code:     http.StatusOK,
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf("<p>%s</p>", "no such template &#39;foo&#39;"),
		},
		"error no template name - tw": {
			reply:    Engine{Writer: NewTemplateWriter(map[string]*template.Template{"foo": foo})},
			code:     http.StatusOK,
			wantCode: http.StatusOK,
			wantBody: "Hello, Sherlock",
		},
		"error fail encode, debug false - jw": {
			reply:    ejw,
			code:     http.StatusOK,
			opts:     Options{Data: map[string]interface{}{"foo": make(chan int)}},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"Internal Server Error"}`,
		},
		"error fail encode, debug true - jw": {
			reply:    Engine{Writer: NewJSONWriter(), Debug: true},
			code:     http.StatusOK,
			opts:     Options{Data: map[string]interface{}{"foo": make(chan int)}},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"json: unsupported type: chan int"}`,
		},
		"ok - jw": {
			reply:    Engine{Writer: NewJSONWriter()},
			code:     http.StatusOK,
			wantCode: http.StatusOK,
			wantBody: `{"name":"Sherlock"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if c.opts == (Options{}) {
				c.opts = Options{
					Key:  "foo",
					Name: "base",
					Data: struct {
						Name string `json:"name"`
					}{
						Name: "Sherlock",
					},
				}
			}
			c.reply.ReplyOrError(w, c.code, c.opts)
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}

func TestGenericReplies(t *testing.T) {
	rtw := Engine{Writer: NewTemplateWriter(map[string]*template.Template{"foo": foo})}
	ejw := Engine{Writer: NewJSONWriter()}
	cases := map[string]struct {
		noOpts     bool
		methodOpts func(http.ResponseWriter, Options)
		method     func(http.ResponseWriter)
		wantCode   int
		wantBody   string
	}{
		"200 ok - tw": {
			methodOpts: rtw.OK,
			wantCode:   http.StatusOK,
			wantBody:   "Hello, Sherlock",
		},
		"200 ok - jw": {
			methodOpts: ejw.OK,
			wantCode:   http.StatusOK,
			wantBody:   `{"name":"Sherlock"}`,
		},
		"201 created - tw": {
			methodOpts: rtw.Created,
			wantCode:   http.StatusCreated,
			wantBody:   "Hello, Sherlock",
		},
		"201 created - jw": {
			methodOpts: ejw.Created,
			wantCode:   http.StatusCreated,
			wantBody:   `{"name":"Sherlock"}`,
		},
		"204 no content - tw": {
			noOpts:   true,
			method:   rtw.NoContent,
			wantCode: http.StatusNoContent,
			wantBody: "",
		},
		"204 no content - jw": {
			noOpts:   true,
			method:   ejw.NoContent,
			wantCode: http.StatusNoContent,
			wantBody: "null",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if c.noOpts {
				c.method(w)
			} else {
				c.methodOpts(w, Options{
					Key:  "foo",
					Name: "base",
					Data: struct {
						Name string `json:"name"`
					}{
						Name: "Sherlock",
					},
				})
			}
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
		})
	}
}
