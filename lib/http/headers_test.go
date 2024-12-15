package http

import (
	"strings"
	"testing"
)

// Test case to validate the working of adding new header to the Headers collection.
func Test_Headers_Add(t *testing.T) {
	testHeaders := make(Headers)
	testCases := []struct {
		Name string
		HdrKey string
		HdrValue string
		ExpHdrCount int
	} {
		{ "Adding first parameter", "Name", "Proteus", 1 },
		{ "Adding values to first parameter", "Name", "WebServer", 1 },
		{ "Adding second parameter", "Age", "18", 2 },
		{ "Adding third parameter", "Value", "Proteus Web Server", 3 },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			testHeaders.Add(testCase.HdrKey, testCase.HdrValue)
			if len(testHeaders) == testCase.ExpHdrCount {
				tt.Logf("The expected parameter count [%d] matches the actual parameter count [%d].", testCase.ExpHdrCount, len(testHeaders))
			} else {
				tt.Errorf("The expected parameter count [%d] does not match the actual parameter count [%d].", testCase.ExpHdrCount, len(testHeaders))
			}
		})
	}
}

// Test case to validate the working of fetching the values for a 'key' from the Headers collection.
func Test_Headers_Get(t *testing.T) {
	testHeaders := make(Headers)
	testHeaders.Add("Name", "proteus")
	testHeaders.Add("Server", "WebServer")
	testHeaders.Add("Server", "HTTP-Compliant")
	testCases := []struct {
		Name string
		HdrKey string
		ExpHdrValue string
	} {
		{ "Fetching parameterwith a single value in the collection", "Name", "proteus" },
		{ "Fetching parameter not in the collection", "Age", "" },
		{ "Fetching parameter with multiple values in the collection", "Server", "WebServer,HTTP-Compliant" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			value, ok := testHeaders.Get(testCase.HdrKey)
			if ok {
				if strings.EqualFold(value, testCase.ExpHdrValue) {
					tt.Logf("The expected value [%s] matches the returned value [%s]", testCase.ExpHdrValue, value)
				} else {
					tt.Errorf("The expected value [%s] does not match the returned value [%s]", testCase.ExpHdrValue, value)
				}
			} else {
				if strings.EqualFold(testCase.ExpHdrValue, "") {
					tt.Logf("As expected, the key [%s] was not found in the Headers collection", testCase.HdrKey)
				} else {
					tt.Errorf("The key [%s] was supposed to be absent in the Headers collection, but value returned [%s] shows it was present.", testCase.HdrKey, value)
				}
			}
		})
	}
}