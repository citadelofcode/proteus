package http

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

// Helper function to create and return a new test instance of HttpRequest.
func newTestRequest(t testing.TB) *HttpRequest {
	t.Helper()
	testReq := new(HttpRequest)
	testReq.initialize()
	return testReq
}

// Test case to validate the HTTP request message read and parse functionality.
func Test_Request_Read(t *testing.T) {
	testCases := []struct {
		Name string
		InputRequest string
		ExpHttpMethod string
		ExpHttpReqPath string
		ExpHttpVersion string
		ExpHeaderCount int
		ExpQpCount int
	} {
		{ "HTTP v0.9 GET Request", "GET /user/abc\r\n", "GET", "/user/abc", "0.9", 0, 0 },
		{ "HTTP v1.0 GET Request", "GET /user/abc HTTP/1.0\r\nHost: example.com\r\n\r\n", "GET", "/user/abc", "1.0", 1, 0 },
		{ "HTTP v1.0 GET Request with Query Params", "GET /user/abc?name=sample HTTP/1.0\r\nHost: example.com\r\n\r\n", "GET", "/user/abc", "1.0", 1, 1 },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			testReq := newTestRequest(tt)
			stringReader := strings.NewReader(testCase.InputRequest)
			testReq.setReader(bufio.NewReader(stringReader))
			err := testReq.read()
			if err != nil && err != io.EOF {
				tt.Errorf("The given request could not be parsed. Error :: %s", err.Error())
				return
			}

			if !strings.EqualFold(testReq.Method, testCase.ExpHttpMethod) {
				tt.Errorf("Expected request method was to be %s but got %s", testCase.ExpHttpMethod, testReq.Method)
			} else {
				tt.Logf("Expected requested method %s matches the returned value %s", testCase.ExpHttpMethod, testReq.Method)
			}

			if !strings.EqualFold(testReq.ResourcePath, testCase.ExpHttpReqPath) {
				tt.Errorf("Expected request path was to be %s but got %s", testCase.ExpHttpReqPath, testReq.ResourcePath)
			} else {
				tt.Logf("Expected request path %s matches the returned request path %s", testCase.ExpHttpReqPath, testReq.ResourcePath)
			}

			if !strings.EqualFold(testReq.Version, testCase.ExpHttpVersion) {
				tt.Errorf("Expected request version was to be %s but got %s", testCase.ExpHttpVersion, testReq.Version)
			} else {
				tt.Logf("Expected request version %s matches the returned request version %s", testCase.ExpHttpVersion, testReq.Version)
			}

			if testReq.Headers.Length() != testCase.ExpHeaderCount {
				tt.Errorf("Expected header count in the request was to be %d instead the request had %d headers", testCase.ExpHeaderCount, testReq.Headers.Length())
			} else {
				tt.Logf("Expected header count in the request %d matches the returned request header count %d", testCase.ExpHeaderCount, testReq.Headers.Length())
			}

			if testReq.Query.Length() != testCase.ExpQpCount {
				tt.Errorf("Expected query parameters count in the request was to be %d instead the request had %d query parameters", testCase.ExpQpCount, testReq.Query.Length())
			} else {
				tt.Logf("Expected query parameter count in the request %d matches the returned request query parameter count %d", testCase.ExpQpCount, testReq.Query.Length())
			}
		})
	}
}

// Test case to validate the addition of headers to a HTTP request message.
func Test_Request_AddHeader(t *testing.T) {
	testRequest := newTestRequest(t)
	testCases := []struct {
		Name string
		InputHeaderKey string
		InputHeaderValue string
		ExpHdrCount int
		ExpectedErrType string
	} {
		{ "A non-date header field", "Content-Type", "application/pdf", 1, "" },
		{ "A date header field with value in ANSIC format", "Date", "Sun Nov  6 08:49:37 1994", 2, "" },
		{ "A date header field with value in RFC 1123 format", "Last-Modified", "Mon, 30 Jun 2008 11:05:30 GMT", 3, "" },
		{ "A date header field with invalid date value", "If-Modified-Since", "2024-12-11T12:34:56Z" ,3, "RequestParseError" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			err := testRequest.addHeader(testCase.InputHeaderKey, testCase.InputHeaderValue)
			if testCase.ExpectedErrType == "" {
				if err != nil {
					tt.Errorf("Was not expecting an error and yet received one - %v", err)
					return
				}
			}

			if testCase.ExpectedErrType == "RequestParseError" {
				rpErr, ok := err.(*RequestParseError)
				if !ok {
					tt.Errorf("Was expecting a request parse error, but got this instead - %v", err)
				} else {
					tt.Logf("Received a request parse error - %v as expected", rpErr)
				}
				return
			}

			if testRequest.Headers.Length() != testCase.ExpHdrCount {
				tt.Errorf("The request header count - %d does not match the expected header count - %d", testRequest.Headers.Length(), testCase.ExpHdrCount)
			} else {
				tt.Logf("The request header count - %d matches the expected header count - %d", testRequest.Headers.Length(), testCase.ExpHdrCount)
			}
		})
	}
}