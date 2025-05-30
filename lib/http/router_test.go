package http

import (
	"testing"
	"path/filepath"
	"runtime"
)

// Test case to check the working of addStaticRoute() function of Router instance.
func Test_Router_AddStaticRoute(t *testing.T) {
	testRouter := NewRouter()
	testCases := []struct {
		Name string
		InputMethod string
		InputRoute string
		TargetPath string
		ExpectedErr string
	} {
		{ "Valid route with valid target folder path", "GET", "/files/static", "../assets", "" },
		{ "Valid route with a target file path", "GET", "/files/staticone", "../assets/index.html", "RoutingError" },
	}

	_, CurrentFilePath, _, _ := runtime.Caller(0)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			testCaseTargetPath := testCase.TargetPath
			isAbsolutePath := filepath.IsAbs(testCaseTargetPath)
			if !isAbsolutePath {
				testCaseTargetPath = filepath.Join(filepath.Dir(CurrentFilePath), testCaseTargetPath)
			}
			err := testRouter.addStaticRoute(testCase.InputMethod, testCase.InputRoute, testCaseTargetPath)
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
