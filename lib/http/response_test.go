package http

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"bufio"
)

// Helper function to create and return a new test instance of HttpResponse.
func newTestResponse(t testing.TB, version string) *HttpResponse {
	t.Helper()
	testRes := new(HttpResponse)
	testRes.initialize(version, true)
	return testRes
}

// Test case to validate the addition of headers to a HTTP response message.
func Test_Response_AddHeader(t *testing.T) {
	testResponse := newTestResponse(t, "")
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
		{ "A date header field with invalid date value", "If-Modified-Since", "2024-12-11T12:34:56Z" ,3, "ResponseError" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			err := testResponse.AddHeader(testCase.InputHeaderKey, testCase.InputHeaderValue)
			if testCase.ExpectedErrType == "" {
				if err != nil {
					tt.Errorf("Was not expecting an error and yet received one - %v", err)
					return
				}
			}

			if testCase.ExpectedErrType == "ResponseError" {
				rpErr, ok := err.(*ResponseError)
				if !ok {
					tt.Errorf("Was expecting a response error, but got this instead - %v", err)
				} else {
					tt.Logf("Received a response error - %v as expected", rpErr)
				}
				return
			}

			if testResponse.Headers.Length() != testCase.ExpHdrCount {
				tt.Errorf("The response header count - %d does not match the expected header count - %d", testResponse.Headers.Length(), testCase.ExpHdrCount)
			} else {
				tt.Logf("The response header count - %d matches the expected header count - %d", testResponse.Headers.Length(), testCase.ExpHdrCount)
			}
		})
	}
}

// Test case to validate the working of the response write function.
func Test_Response_Write(t *testing.T) {
	testCases := []struct {
		Name string
		IpVersion string
		IpContentType string
		IpContent string
		IpStatus StatusCode
		ExpErr string
		ExpResponse string
	} {
		{ "Simple v0.9 Response", "0.9", "", "Hello, this is a simple response from Proteus!", StatusOK, "", "Hello, this is a simple response from Proteus!" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			res := newTestResponse(tt, testCase.IpVersion)
			var opBuffer bytes.Buffer
			writer := bufio.NewWriter(&opBuffer)
			res.setWriter(writer)
			res.bodyBytes = []byte(testCase.IpContent)
			if !strings.EqualFold(testCase.IpVersion, "0.9") {
				res.AddHeader("Content-Type", testCase.IpContentType)
				res.AddHeader("Content-Length", strconv.Itoa(len(res.bodyBytes)))
				res.Status(testCase.IpStatus)
			}

			err := res.write()

			if testCase.ExpErr == "" {
				if err != nil {
					tt.Errorf("Was not expecting an error and yet got this error - %v", err)
					return
				}
			}

			if testCase.ExpErr == "ResponseError" {
				respErr, ok := err.(*ResponseError)
				if !ok {
					tt.Errorf("Was expecting a response error, but got this error instead - %v", err)
				} else {
					tt.Logf("Was expecting a response error and got one - %v", respErr)
				}

				return
			}

			opString := opBuffer.String()
			if strings.EqualFold(opString, testCase.ExpResponse) {
				tt.Logf("The expected response [%s] matches the response written by the write function [%s].", testCase.ExpResponse, opString)
			} else {
				tt.Errorf("The expected response [%s] does not match the response written by the write function [%s].", testCase.ExpResponse, opString)
			}
		})
	}
}
