package http

import (
	"testing"
)

// Test case to check the working of th route validation logic.
func Test_Router_ValidateRoute(t *testing.T) {
	testRouter := newRouter()
	testCases := []struct {
		Name string
		InputRoute string
		ExpectedOp bool
	} {
		{ "Valid route containing alphabets and numbers", "/abc/xyz/123", true },
		{ "Valid route containing hyphen and underscore", "/abc/xyz_123", true },
		{ "Valid route containing path parameters", "/abc/:name", true },
		{ "Invalid route containing multiple slashes as prefix", "//pqr/abc/123", false },
		{ "Invalid route containing multiple slashes as prefix", "/pqr/abc/123/", false },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			isValid := testRouter.validateRoute(testCase.InputRoute)
			if testCase.ExpectedOp != isValid {
				if testCase.ExpectedOp {
					tt.Errorf("The route (%s) is a valid route, but was deemed invalid.", testCase.InputRoute)
				} else {
					tt.Errorf("The route (%s) is an invalid route, but was deemed valid.", testCase.InputRoute)
				}
			} else {
				if testCase.ExpectedOp {
					tt.Logf("The route - %s was established correctly to be valid.", testCase.InputRoute)
				} else {
					tt.Logf("The route - %s was established correctly to be invalid.", testCase.InputRoute)
				}
			}
		})
	}
}

// Test case to check the working of addStaticRoute() function of Router instance.
func Test_Router_AddStaticRoute(t *testing.T) {
	testRouter := newRouter()
	testCases := []struct {
		Name string
		InputMethod string
		InputRoute string
		TargetFilePath string
		ExpectedErr string
	} {
		{ "Valid route with valid target folder path", "GET", "/files/static", "/Users/maheshkumaarbalaji/Projects/proteus/Files/", "" },
		{ "Valid route with a target file path", "GET", "/files/staticone", "/Users/maheshkumaarbalaji/Projects/proteus/Files/home.html", "RoutingError" },
		{ "Valid route with a relative target path", "GET", "/files/statictwo", "./../proteus/Files/", "RoutingError" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			err := testRouter.addStaticRoute(testCase.InputMethod, testCase.InputRoute, testCase.TargetFilePath)
			if testCase.ExpectedErr == "" {
				if err != nil {
					tt.Errorf("Was not expecting an error for adding static route to router and yet got this instead - %v", err)
					return
				}
			}

			if testCase.ExpectedErr == "RoutingError" {
				rtrError, ok := err.(*RoutingError)
				if !ok {
					tt.Errorf("Expected a routing error while adding static route to router, but got this instead - %v", err)
				} else {
					tt.Logf("Was expecting a routing error and got a routing error as well - %v", rtrError)
				}
			}
		})
	}
}