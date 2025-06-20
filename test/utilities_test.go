package test

import (
	"strings"
	"testing"
	"slices"
	"github.com/citadelofcode/proteus/internal"
)

// Test case to validate the working of the getHighestVersion utility function.
func Test_GetHighestVersion(t *testing.T) {
	testCases := []struct {
		Name string
		InputVersion string
		ExpVersion string
	} {
		{ "HTTP v2.0", "2.0", "1.1" },
		{ "HTTP v3.0", "3.0", "1.1" },
		{ "HTTP v0.9", "0.9", "0.9" },
		{ "HTTP v1.0", "1.0", "1.0"},
		{ "HTTP v1.1", "1.1", "1.1" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			maxVersion := internal.GetHighestVersion(testCase.InputVersion)
			if strings.EqualFold(maxVersion, testCase.ExpVersion) {
				tt.Logf("The expected version [%s] matches the version returned by the utility function [%s].", testCase.ExpVersion, maxVersion)
			} else {
				tt.Errorf(internal.TextColor.Red("The expected version [%s] does not match the version returned by the utility function [%s]."), testCase.ExpVersion, maxVersion)
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
		{ "A version (HTTP/0.9) compatible with the server", "0.9", "0.9" },
		{ "A version (HTTP/1.0) compatible with the server", "1.0", "1.0" },
		{ "A version (HTTP/1.1) compatible with the server", "1.1", "1.1" },
		{ "A version not compatible with the server", "2.0", "1.1" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			responseVersion := internal.GetResponseVersion(testCase.InputVersion)
			if strings.EqualFold(responseVersion, testCase.ExpVersion) {
				tt.Logf("The expected response version [%s] matches the version returned by the utility function [%s].", testCase.ExpVersion, responseVersion)
			} else {
				tt.Errorf(internal.TextColor.Red("The expected response version [%s] does not match the version returned by the utility function [%s]."), testCase.ExpVersion, responseVersion)
			}
		})
	}
}

// Test case to validate the IsHttpDate utility function.
func Test_IsHttpDate(t *testing.T) {
	testCases := []struct {
		Name string
		IpDateString string
		ExpOutput bool
	} {
		{ "RFC 1123 date string", "Mon, 06 Jun 2025 18:45:00 GMT", true },
		{ "ANSIC date string", "Fri Jun  6 18:45:00 2025", true },
		{ "RFC 850 date string", "Friday, 06-Jun-25 18:45:00 GMT", true },
		{ "RFC 3339 date string", "2025-06-06T18:45:00Z", false },
		{ "Ruby date formatted string", "Mon Jan 02 15:04:05 -0700 2006", false },
		{ "Unix date formatted string", "Mon Jan _2 15:04:05 MST 2006", false },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			got, _ := internal.IsHttpDate(testCase.IpDateString)
			if got == testCase.ExpOutput {
				if got {
					tt.Logf("The given date string [%s] has been correctly classified as HTTP date", testCase.IpDateString)
				} else {
					tt.Logf("The given date string [%s] has been correctly classified as not a HTTP date", testCase.IpDateString)
				}
			} else {
				if got {
					tt.Errorf(internal.TextColor.Red("The given date string [%s] has been incorrectly classified as HTTP date"), testCase.IpDateString)
				} else {
					tt.Errorf(internal.TextColor.Red("The given date string [%s] has been incorrectly classified as not a HTTP date"), testCase.IpDateString)
				}
			}
		})
	}
}

// Test case to validate the CleanRoute utility function.
func Test_CleanRoute(t *testing.T) {
	testCases := []struct {
		Name string
		IpRoute string
		ExpRoute string
	} {
		{ "An already clean route path", "/user/register", "/user/register" },
		{ "A route with multiple leading slashes", "//user/register", "/user/register" },
		{ "A route with multiple trailing slashes", "/user/register//", "/user/register" },
		{ "A route with a single trailing slash", "/user/register/", "/user/register" },
		{ "A route with multiple slash as route part seperator", "/user//register", "/user/register" },
		{ "A route path with no leading slash", "user/register", "/user/register" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			got := internal.CleanRoute(testCase.IpRoute)
			if strings.EqualFold(got, testCase.ExpRoute) {
				tt.Logf("The input route [%s] matches the expected route [%s]", testCase.IpRoute, got)
			} else {
				tt.Errorf(internal.TextColor.Red("The input route [%s] does not match the expected route [%s]"), testCase.IpRoute, got)
			}
		})
	}
}

// Test case to validate the GetAllVersions() utility function.
func Test_GetAllVersions(t *testing.T) {
	versions := internal.GetAllVersions()
	if len(versions) != 3 {
		t.Errorf(internal.TextColor.Red("Expected 3 compatible versions, but got %d versions instead"), len(versions))
		return
	} else {
		t.Log("There are 3 versions of HTTP supported by the server as expected")
	}

	allVersionsFound := true
	if !slices.Contains(versions, "0.9") {
		t.Error(internal.TextColor.Red("HTTP/0.9 was not found among the list of compatible versions"))
		allVersionsFound = false
	}

	if !slices.Contains(versions, "1.0") {
		t.Error(internal.TextColor.Red("HTTP/1.0 was not found among the list of compatible versions"))
		allVersionsFound = false
	}

	if !slices.Contains(versions, "1.1") {
		t.Error(internal.TextColor.Red("HTTP/1.1 was not found among the list of compatible versions"))
		allVersionsFound = false
	}

	if allVersionsFound {
		t.Log("All 3 versions of HTTP [0.9, 1.0, 1.1] were found among the list of compatible versions.")
	}
}

// Test case to validate the GetAllowedMethods() utility function.
func Test_GetAllowedMethods(t *testing.T) {
	testCases := []struct {
		Name string
		IpVersion string
		OpMethods string
	} {
		{ "HTTP version 0.9 - Compatible", "0.9", "GET" },
		{ "HTTP version 1.0 - Compatible", "1.0", "GET, POST, HEAD, OPTIONS, TRACE" },
		{ "HTTP version 1.1 - Compatible", "1.1", "GET, HEAD, POST, PUT, DELETE, TRACE, OPTIONS, CONNECT, PATCH" },
		{ "HTTP version 2.0 - Not Compatible", "2.0", "" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			allowedMethods := internal.GetAllowedMethods(testCase.IpVersion)
			if strings.EqualFold(allowedMethods, testCase.OpMethods) {
				tt.Logf("The list of allowed methods [%s] received for version [%s] matches the expected list of methods", allowedMethods, testCase.IpVersion)
			} else {
				tt.Errorf(internal.TextColor.Red("The list of allowed methods [%s] received for version [%s] does not match the expected list of methods"), allowedMethods, testCase.IpVersion)
			}
		})
	}
}

// Test case to validate the working of the IsMethodAllowed() utility function.
func Test_IsMethodAllowed(t *testing.T) {
	testCases := []struct {
		Name string
		IpVersion string
		IpMethod string
		IsAllowed bool
	} {
		{ "Allowed method for HTTP/0.9", "0.9", "GET", true },
		{ "Not allowed method for HTTP/0.9", "0.9", "POST", false },
		{ "Allowed method for HTTP/1.0", "1.0", "HEAD", true },
		{ "Not Allowed method for HTTP/1.0", "1.0", "PATCH", false },
		{ "Allowed method for HTTP/1.1", "1.1", "PATCH", true },
		{ "Not allowed method for HTTP/1.1", "1.1", "PURGE", false },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			isAllowed := internal.IsMethodAllowed(testCase.IpVersion, testCase.IpMethod)
			if isAllowed == testCase.IsAllowed {
				if isAllowed {
					tt.Logf("The given method [%s] for version [%s] was correctly determined as allowed.", testCase.IpMethod, testCase.IpVersion)
				} else {
					tt.Logf("The given method [%s] for version [%s] was correctly determined as not allowed.", testCase.IpMethod, testCase.IpVersion)
				}
			} else {
				if isAllowed {
					tt.Errorf(internal.TextColor.Red("The given method [%s] for version [%s] was incorrectly determined as allowed."), testCase.IpMethod, testCase.IpVersion)
				} else {
					tt.Errorf(internal.TextColor.Red("The given method [%s] for version [%s] was incorrectly determined as not allowed."), testCase.IpMethod, testCase.IpVersion)
				}
			}
		})
	}
}
