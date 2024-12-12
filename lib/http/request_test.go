package http

import (
	"testing"
	"strings"
	"bufio"
)

func newTestRequest(t testing.TB) *HttpRequest {
	t.Helper()
	testReq := new(HttpRequest)
	testReq.initialize()
	return testReq
}

func Test_RequestRead(t *testing.T) {
	testCases := []struct {
		Name string
		InputRequest string
		ExpHttpMethod string
		ExpHttpReqPath string
		ExpHttpVersion string
		ExpHeaderCount int
	} {
		{"HTTP v0.9 GET Request", "GET /user/abc", "GET", "/user/abc", "0.9", 0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			testReq := newTestRequest(tt)
			stringReader := strings.NewReader(testCase.InputRequest)
			testReq.setReader(bufio.NewReader(stringReader))
			testReq.read()

			if !strings.EqualFold(testReq.Method, testCase.ExpHttpMethod) {
				tt.Errorf("Expected request method to be %s but got %s", testCase.ExpHttpMethod, testReq.Method)
			}

			if !strings.EqualFold(testReq.ResourcePath, testCase.ExpHttpReqPath) {
				tt.Errorf("Expected request path to be %s but got %s", testCase.ExpHttpReqPath, testReq.ResourcePath)
			}

			if !strings.EqualFold(testReq.Version, testCase.ExpHttpVersion) {
				tt.Errorf("Expected request version to be %s but got %s", testCase.ExpHttpVersion, testReq.Version)
			}

			if len(testReq.Headers) != testCase.ExpHeaderCount {
				tt.Errorf("Expected header count in the request to be %d instead the request had %d headers", testCase.ExpHeaderCount, len(testReq.Headers))
			}
		})
	}
}