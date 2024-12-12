package http

import (
	"testing"
)

// Helper method to create an instance of a router for testing.
func getTestRouter(t testing.TB) *Router {
	t.Helper()
	testRouter := new(Router)
	testRouter.Routes = make([]Route, 0)
	testRouter.RouteTree = createTree()
	return testRouter
}

// Test case to check the working of th route validation logic.
func Test_ValidateRoute(t *testing.T) {
	testRouter := getTestRouter(t)
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
			}
		})
	}
}