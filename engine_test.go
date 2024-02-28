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
		andErr    bool
		methodErr func(http.ResponseWriter, error)
		method    func(http.ResponseWriter)
		wantCode  int
		wantBody  string
	}{
		"bad request - tw": {
			method:   etw.BadRequest,
			wantCode: http.StatusBadRequest,
			wantBody: errorTemplateBody(http.StatusBadRequest),
		},
		"bad request - jw": {
			method:   ejw.BadRequest,
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"Bad Request"}`,
		},
		"unauthorized - tw": {
			method:   etw.Unauthorized,
			wantCode: http.StatusUnauthorized,
			wantBody: errorTemplateBody(http.StatusUnauthorized),
		},
		"unauthorized - jw": {
			method:   ejw.Unauthorized,
			wantCode: http.StatusUnauthorized,
			wantBody: `{"error":"Unauthorized"}`,
		},
		"forbidden - tw": {
			method:   etw.Forbidden,
			wantCode: http.StatusForbidden,
			wantBody: errorTemplateBody(http.StatusForbidden),
		},
		"forbidden - kw": {
			method:   ejw.Forbidden,
			wantCode: http.StatusForbidden,
			wantBody: `{"error":"Forbidden"}`,
		},
		"not found - tw": {
			method:   etw.NotFound,
			wantCode: http.StatusNotFound,
			wantBody: errorTemplateBody(http.StatusNotFound),
		},
		"not found - jw": {
			method:   ejw.NotFound,
			wantCode: http.StatusNotFound,
			wantBody: `{"error":"Not Found"}`,
		},
		"not acceptable - template writer": {
			method:   etw.NotAcceptable,
			wantCode: http.StatusNotAcceptable,
			wantBody: errorTemplateBody(http.StatusNotAcceptable),
		},
		"not acceptable - json writer": {
			method:   ejw.NotAcceptable,
			wantCode: http.StatusNotAcceptable,
			wantBody: `{"error":"Not Acceptable"}`,
		},
		"request timeout - tw": {
			method:   etw.RequestTimeout,
			wantCode: http.StatusRequestTimeout,
			wantBody: errorTemplateBody(http.StatusRequestTimeout),
		},
		"request timeout - jw": {
			method:   ejw.RequestTimeout,
			wantCode: http.StatusRequestTimeout,
			wantBody: `{"error":"Request Timeout"}`,
		},
		"conflict - tw": {
			method:   etw.Conflict,
			wantCode: http.StatusConflict,
			wantBody: errorTemplateBody(http.StatusConflict),
		},
		"conflict - jw": {
			method:   ejw.Conflict,
			wantCode: http.StatusConflict,
			wantBody: `{"error":"Conflict"}`,
		},
		"gone - tw": {
			method:   etw.Gone,
			wantCode: http.StatusGone,
			wantBody: errorTemplateBody(http.StatusGone),
		},
		"gone - jw": {
			method:   ejw.Gone,
			wantCode: http.StatusGone,
			wantBody: `{"error":"Gone"}`,
		},
		"unprocessable entity - tw": {
			method:   etw.UnprocessableEntity,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: errorTemplateBody(http.StatusUnprocessableEntity),
		},
		"unprocessable entity - jw": {
			method:   ejw.UnprocessableEntity,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"error":"Unprocessable Entity"}`,
		},
		"too many requests - tw": {
			method:   etw.TooManyRequests,
			wantCode: http.StatusTooManyRequests,
			wantBody: errorTemplateBody(http.StatusTooManyRequests),
		},
		"too many requests - jw": {
			method:   ejw.TooManyRequests,
			wantCode: http.StatusTooManyRequests,
			wantBody: `{"error":"Too Many Requests"}`,
		},
		"internal server error - tw": {
			andErr:    true,
			methodErr: etw.InternalServerError,
			wantCode:  http.StatusInternalServerError,
			wantBody:  errorTemplateBody(http.StatusInternalServerError),
		},
		"internal server error - jw": {
			andErr:    true,
			methodErr: ejw.InternalServerError,
			wantCode:  http.StatusInternalServerError,
			wantBody:  `{"error":"Internal Server Error"}`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if !c.andErr {
				c.method(w)
			} else {
				c.methodErr(w, fmt.Errorf("sample error"))
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

func TestMethodNotAllowed(t *testing.T) {
	etw := Engine{Writer: NewTemplateWriter(map[string]*template.Template{})}
	ejw := Engine{Writer: NewJSONWriter()}
	cases := map[string]struct {
		reply     Engine
		allow     []string
		wantCode  int
		wantBody  string
		wantAllow string
	}{
		"allow one - tw": {
			reply:     etw,
			allow:     []string{http.MethodGet},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  errorTemplateBody(http.StatusMethodNotAllowed),
			wantAllow: http.MethodGet,
		},
		"allow one - jw": {
			reply:     ejw,
			allow:     []string{http.MethodGet},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  `{"error":"Method Not Allowed"}`,
			wantAllow: http.MethodGet,
		},
		"allow multiple - tw": {
			reply:     etw,
			allow:     []string{http.MethodGet, http.MethodPost},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  errorTemplateBody(http.StatusMethodNotAllowed),
			wantAllow: strings.Join([]string{http.MethodGet, http.MethodPost}, ", "),
		},
		"allow multiple - jw": {
			reply:     ejw,
			allow:     []string{http.MethodGet, http.MethodPost},
			wantCode:  http.StatusMethodNotAllowed,
			wantBody:  `{"error":"Method Not Allowed"}`,
			wantAllow: strings.Join([]string{http.MethodGet, http.MethodPost}, ", "),
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.reply.MethodNotAllowed(w, c.allow...)
			if got := w.Code; got != c.wantCode {
				t.Errorf(errorString, got, c.wantCode)
			}
			if got := strings.TrimSpace(w.Body.String()); got != c.wantBody {
				t.Errorf(errorString, got, c.wantBody)
			}
			if got := w.Header().Get("Allow"); got != c.wantAllow {
				t.Errorf(errorString, got, c.wantAllow)
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
		"ok jw": {
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
		"ok - tw": {
			methodOpts: rtw.OK,
			wantCode:   http.StatusOK,
			wantBody:   "Hello, Sherlock",
		},
		"ok - jw": {
			methodOpts: ejw.OK,
			wantCode:   http.StatusOK,
			wantBody:   `{"name":"Sherlock"}`,
		},
		"created - tw": {
			methodOpts: rtw.Created,
			wantCode:   http.StatusCreated,
			wantBody:   "Hello, Sherlock",
		},
		"created - jw": {
			methodOpts: ejw.Created,
			wantCode:   http.StatusCreated,
			wantBody:   `{"name":"Sherlock"}`,
		},
		"no content - tw": {
			noOpts:   true,
			method:   rtw.NoContent,
			wantCode: http.StatusNoContent,
			wantBody: "",
		},
		"no content - jw": {
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
