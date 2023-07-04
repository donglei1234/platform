package utils_test

import (
	"github.com/appscode/go/strings"
	"net/http"
	"runtime/debug"
	"testing"

	"github.com/donglei1234/platform/services/common/utils"
)

func TestAllowCORS(t *testing.T) {
	testcases := map[string]struct {
		utils.TestHttpHandler
		utils.TestHttpResponseWriter
		http.Request
		expHeader http.Header
		expPanic  bool
		expServe  bool
	}{"SetHeadersWithoutServing": {
		TestHttpHandler:        utils.AllowCors_TestHttpHandlerEmpty,
		TestHttpResponseWriter: utils.AllowCors_TestHttpResponseWriter,
		Request:                utils.AllowCors_TestHttpRequest,
		expHeader:              utils.AllowCors_TestHeaderEnd,
	}, "ServeWithoutSettingHeaders": {
		TestHttpHandler:        utils.AllowCors_TestHttpHandlerEmpty,
		TestHttpResponseWriter: utils.AllowCors_TestHttpResponseWriterEmpty,
		Request:                utils.AllowCors_TestHttpRequestEmpty,
		expServe:               true,
	}}
	for i, tc := range testcases {
		func() {
			defer func() {
				if r := recover(); r != nil && !tc.expPanic {
					t.Errorf("Test %s: Encountered unexpected panic: %v\n%s", i, r, debug.Stack())
				}
			}()
			res := utils.AllowCORS(&tc.TestHttpHandler)
			if res == nil {
				t.Errorf("Test %s: Received nil result", i)
				return
			}
			res.ServeHTTP(&tc.TestHttpResponseWriter, &tc.Request)
			if tc.TestHttpHandler.Served && !tc.expServe {
				t.Errorf("Test %s: ServeHttp was unexpectedly invoked", i)
			} else if !tc.TestHttpHandler.Served && tc.expServe {
				t.Errorf("Test %s: ServeHttp was unexpectedly not invoked", i)
			}
			if !mapsEqual(tc.TestHttpResponseWriter.TestHeader, tc.expHeader) {
				t.Errorf(
					"Test %s: ResponseWriter had Header %#v; expected Header %#v",
					i,
					tc.TestHttpResponseWriter.TestHeader,
					tc.expHeader,
				)
			}
			if tc.expPanic {
				t.Errorf("Test %s: Did not encounter expected panic", i)
			}
		}()
	}
}

func mapsEqual(map1, map2 map[string][]string) bool {
	for i, strings1 := range map1 {
		if strings2, ok := map2[i]; !ok {
			return false
		} else {
			for _, s := range strings1 {
				if !strings.Contains(strings2, s) {
					return false
				}
			}
		}
	}
	return true
}
