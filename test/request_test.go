package test

import (
	"io"
	"strings"
	"testing"
)

// Test case to validate the HTTP request message read and parse functionality.
func Test_Request_Read(t *testing.T) {
	testServer := NewTestServer(t)
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
			testReq := NewTestRequest(tt, testServer, strings.NewReader(testCase.InputRequest))
			err := testReq.Read()
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

// Test case to validate the addHeader functionality of the HTTP Request instance.
func Test_Request_AddHeader(t *testing.T) {
	testServer := NewTestServer(t)
	testCases := []struct {
		Name string
		InputHeader string
		InputValue string
		HeaderExists bool
	} {
		{ "Valid non-date Header field and value", "Content-Type", "application/json", true },
		{ "Valid RFC1123 value for a date header", "Date", "Mon, 06 Jun 2025 18:45:00 GMT", true },
		{ "Valid ANSIC value for a date header", "Expires", "Fri Jun  6 18:45:00 2025", true },
		{ "Valid RFC850 value for a date header", "If-Modified-Since", "Friday, 06-Jun-25 18:45:00 GMT", true },
		{ "Date header with value of invalid format", "Last-Modified", "2025-06-06T18:45:00Z", false },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			testRequest := NewTestRequest(t, testServer, nil)
			testRequest.AddHeader(testCase.InputHeader, testCase.InputValue)
			_, headerExists := testRequest.Headers.Get(testCase.InputHeader)
			if headerExists == testCase.HeaderExists {
				if headerExists {
					tt.Logf("Input Header [%s] has been added successfully to the request as expected", testCase.InputHeader)
				} else {
					tt.Logf("Invalid Input Header [%s] was not added to the request headers collection as expected", testCase.InputHeader)
				}
			} else {
				if headerExists {
					tt.Errorf("An invalid header [%s] with value [%s] has been added to the request headers collection", testCase.InputHeader, testCase.InputValue)
				} else {
					tt.Errorf("A valid header [%s] with value [%s] has not been aded to the request headers collection", testCase.InputHeader, testCase.InputValue)
				}
			}
		})
	}
}
