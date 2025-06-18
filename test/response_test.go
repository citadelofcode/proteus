package test

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"bufio"
	"github.com/citadelofcode/proteus/internal"
)

// Test case to validate the working of the response write function.
func Test_Response_Write(t *testing.T) {
	testServer := NewTestServer(t)
	testCases := []struct {
		Name string
		IpVersion string
		IpContentType string
		IpContent string
		IpStatus internal.StatusCode
		ExpErr string
		ExpResponse string
	} {
		{ "Simple v0.9 Response", "0.9", "", "Hello, this is a simple response from Proteus!", internal.StatusOK, "", "Hello, this is a simple response from Proteus!" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			var opBuffer bytes.Buffer
			res := NewTestResponse(tt, testCase.IpVersion, testServer, bufio.NewWriter(&opBuffer))
			res.BodyBytes = []byte(testCase.IpContent)
			if !strings.EqualFold(testCase.IpVersion, "0.9") {
				res.AddHeader("Content-Type", testCase.IpContentType)
				res.AddHeader("Content-Length", strconv.Itoa(len(res.BodyBytes)))
				res.Status(testCase.IpStatus)
			}

			err := res.Write()

			if testCase.ExpErr == "" {
				if err != nil {
					tt.Errorf("Was not expecting an error and yet got this error - %v", err)
					return
				}
			}

			if testCase.ExpErr == "ResponseError" {
				respErr, ok := err.(*internal.ResponseError)
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

// Test case to validate the default headers for response instance for different versions of HTTP.
func Test_Response_DefaultHeaders(t *testing.T) {
	testServer := NewTestServer(t)
	testCases := []struct {
		Name string
		HttpVersion string
		ExpHeaderCount int
	} {
		{ "For HTTP/0.9 response instance", "0.9",  0},
		{ "For HTTP/1 response instance", "1", 2 },
		{ "For HTTP/1.1 response instance", "1.1", 2 },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			testResponse := NewTestResponse(tt, testCase.HttpVersion, testServer, nil)
			if testResponse.Headers.Length() == testCase.ExpHeaderCount {
				tt.Logf("The expected header count [%d] matches the actual header count [%d] for HTTP/[%s] response instance", testCase.ExpHeaderCount, testResponse.Headers.Length(), testCase.HttpVersion)
			} else {
				tt.Errorf("The expected header count [%d] does not match the actual header count [%d] for HTTP/[%s] response instance", testCase.ExpHeaderCount, testResponse.Headers.Length(), testCase.HttpVersion)
			}
		})
	}
}

// Test case to validate the addHeader functionality of the HTTP Response instance.
func Test_Response_AddHeader(t *testing.T) {
	testServer := NewTestServer(t)
	testCases := []struct {
		Name string
		InputHeader string
		InputValue string
		HeaderExists bool
	} {
		{ "Valid non-date Header field and value", "Content-Type", "application/json", true },
		{ "Valid RFC1123 value for a date header", "Expires", "Mon, 06 Jun 2025 18:45:00 GMT", true },
		{ "Valid ANSIC value for a date header", "Expires", "Fri Jun  6 18:45:00 2025", true },
		{ "Valid RFC850 value for a date header", "If-Modified-Since", "Friday, 06-Jun-25 18:45:00 GMT", true },
		{ "Date header with value of invalid format", "Last-Modified", "2025-06-06T18:45:00Z", false },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			testResponse := NewTestResponse(t, "1.0", testServer, nil)
			testResponse.AddHeader(testCase.InputHeader, testCase.InputValue)
			_, headerExists := testResponse.Headers.Get(testCase.InputHeader)
			if headerExists == testCase.HeaderExists {
				if headerExists {
					tt.Logf("Input Header [%s] has been added successfully to the response as expected", testCase.InputHeader)
				} else {
					tt.Logf("Invalid Input Header [%s] was not added to the response headers collection as expected", testCase.InputHeader)
				}
			} else {
				if headerExists {
					tt.Errorf("An invalid header [%s] with value [%s] has been added to the response headers collection", testCase.InputHeader, testCase.InputValue)
				} else {
					tt.Errorf("A valid header [%s] with value [%s] has not been aded to the response headers collection", testCase.InputHeader, testCase.InputValue)
				}
			}
		})
	}
}
