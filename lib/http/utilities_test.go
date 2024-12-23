package http

import (
	"strings"
	"testing"
)

// Test case to validate the working of the getHighestVersion utility function.
func Test_GetHighestVersion(t *testing.T) {
	testCases := []struct {
		Name string
		InputVersion string
		ExpVersion string
	} {
		{ "HTTP v2", "2.0", "1.1" },
		{ "HTTP v3", "3.0", "1.1" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			maxVersion := getHighestVersion(testCase.InputVersion)
			if strings.EqualFold(maxVersion, testCase.ExpVersion) {
				tt.Logf("The expected version [%s] matches the version returned by the utility function [%s].", testCase.ExpVersion, maxVersion)
			} else {
				tt.Errorf("The expected version [%s] does not match the version returned by the utility function [%s].", testCase.ExpVersion, maxVersion)
			}
		})
	}
}

// Test case to validate the working of the getResponseVersion utility function.
func Test_GetResponseVersion(t *testing.T) {
	testCases := []struct {
		Name string
		InputVersion string
		ExpVersion string
	} {
		{ "A version compatible with the server", "1.1", "1.1" },
		{ "A version not compatible with the server", "2.0", "1.1" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			responseVersion := getResponseVersion(testCase.InputVersion)
			if strings.EqualFold(responseVersion, testCase.ExpVersion) {
				tt.Logf("The expected response version [%s] matches the version returned by the utility function [%s].", testCase.ExpVersion, responseVersion)
			} else {
				tt.Errorf("The expected response version [%s] does not match the version returned by the utility function [%s].", testCase.ExpVersion, responseVersion)
			}
		})
	}
}