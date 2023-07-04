package utils

import (
	"net/http"
)

const (
	TestOrigin = "TestOrigin"
	TestMethod = "RequestMethod"
)

var (
	AllowCors_TestHeaderStart = map[string][]string{
		"Origin":                        {TestOrigin},
		"Access-Control-Request-Method": {TestMethod},
	}
	AllowCors_TestHeaderEnd = map[string][]string{
		"Origin":                        {TestOrigin},
		"Access-Control-Request-Method": {TestMethod},
		"Access-Control-Allow-Origin":   {"origin"},
		"Access-Control-Allow-Headers":  {"Content-Type,Accept,Authorization"},
		"Access-Control-Allow-Methods":  {"GET,HEAD,POST,PUT,DELETE"},
	}
	AllowCors_TestHttpHandlerEmpty        = TestHttpHandler{}
	AllowCors_TestHttpResponseWriter      = TestHttpResponseWriter{TestHeader: AllowCors_TestHeaderStart}
	AllowCors_TestHttpResponseWriterEmpty = TestHttpResponseWriter{}
	AllowCors_TestHttpRequest             = http.Request{
		Header: map[string][]string{
			"Origin":                        {"origin"},
			"Access-Control-Request-Method": {"test"},
		},
		Method: "OPTIONS",
	}
	AllowCors_TestHttpRequestEmpty = http.Request{}
)

type TestHttpHandler struct {
	Served bool
}

func (t *TestHttpHandler) ServeHTTP(http.ResponseWriter, *http.Request) { t.Served = true }

type TestHttpResponseWriter struct {
	TestHeader
	int
	error
}
type TestHeader http.Header

func (t TestHeader) Get(s string) string {
	//if len(t) == 0 {
	//	return ""
	//}
	//if len(t[s]) == 0 {
	//	return ""
	//}
	return t[s][0]
}

func (t TestHeader) Set(s, ss string) {
	//if t == nil {
	//	t = make(map[string][]string)
	//}
	t[s] = append(t[s], ss)
}

func (t *TestHttpResponseWriter) Header() http.Header       { return (http.Header)(t.TestHeader) }
func (t *TestHttpResponseWriter) Write([]byte) (int, error) { return t.int, t.error }
func (t *TestHttpResponseWriter) WriteHeader(int)           {}
