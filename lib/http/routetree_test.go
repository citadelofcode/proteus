package http

import (
	"slices"
	"strings"
	"testing"
)

// Test case to validate if a route path is being normalized into a slice of route parts correctly.
func Test_RoutePathNormalization(t *testing.T) {
	testCases := []struct {
		Name string
		RoutePath string
		ExpArgCount int
	} {
		{ "Normal Route Path", "/abc/nbsdf/123", 3 },
		{ "Route Path with a suffix", "/ahf/dggdg/sdfgdfg/", 3 },
		{ "Route Path with two slashes in the middle", "/afsf/bfsdf//nfsdnf", 3 }, 
		{ "Route Path with no prefix", "abc/ert/123", 3 },
		{ "Route path with no prefix and a suffix", "abc/dfgdf/sbfusd/124/", 4},
		{ "Route path with multiple slash characters as prefix", "//abc/xyz/pqr/123", 4 }, 
		{ "Route path with multiple slash characters as suffix", "abc/xyz/pqr/123//", 4 },
		{ "Route path with multiple slash characters as prefix and suffix", "//abc/xyz/pqr/123//", 4 },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			routeParts := normalizeRoute(testCase.RoutePath)
			if len(routeParts) != testCase.ExpArgCount {
				tt.Fatalf("Number of route parts (%d) does not match the expected route part count (%d)", len(routeParts), testCase.ExpArgCount)
			} else {
				tt.Logf("Number of route parts (%d) matches the expected route part count (%d)", len(routeParts), testCase.ExpArgCount)
			}
		})
	}
}

// Test case to validate if an empty route tree is being created correctly.
func Test_RouteTree_Create(t *testing.T) {
	root := createTree()
	if root.RoutePart != "" {
		t.Errorf("The route part for the root node of the route tree is expected to be an empty string")
	}

	if len(root.Children) > 0 {
		t.Errorf("The root node for an empty route tree cannot have child nodes")
	}
}

// Test case to validate the GetAllRoutes() method of the route tree.
func Test_RouteTree_GetAllRoutes(t *testing.T) {
	testCases := []struct {
		Name string
		Routes []string
		ExpectedResult int
	} {
		{ "Add a normal route", []string{ "/abc/xyz" }, 1 },
		{ "Add one normal route and one route with path parameters", []string { "/pqr/123/jdsbfds", "/xyz/user/:name" }, 2 },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			root := createTree()
			for _, route := range testCase.Routes {
				addRouteToTree(root, route)
			}

			addedRoutes := root.getAllRoutes()
			if len(addedRoutes) != testCase.ExpectedResult {
				tt.Errorf("The number of routes in the route tree (%d) does not match the expected route count (%d).", len(addedRoutes), testCase.ExpectedResult)
			}
		})
	}
}

// Test case to validate the addition of a new route to the route tree.
func Test_RouteTree_InsertRoute(t *testing.T) {
	root := createTree()
	testCases := []struct {
		Name string
		RoutePath string
		AddedRoutePath string
	} {
		{ "Normal Route Path", "/abc/nbsdf/123", "abc/nbsdf/123" },
		{ "Route Path with a suffix", "/ahf/dggdg/sdfgdfg/", "ahf/dggdg/sdfgdfg" },
		{ "Route Path with two slashes in the middle", "/afsf/bfsdf//nfsdnf", "afsf/bfsdf/nfsdnf" }, 
		{ "Route Path with no prefix", "abc/ert/123", "abc/ert/123" },
		{ "Route path with no prefix and a suffix", "abc/dfgdf/sbfusd/124/", "abc/dfgdf/sbfusd/124" },
		{ "Route path with multiple slash characters as prefix", "//abc/xyz/pqr/123", "abc/xyz/pqr/123" }, 
		{ "Route path with multiple slash characters as suffix", "abc/xyz/pqr/123//", "abc/xyz/pqr/123" },
		{ "Route path with multiple slash characters as prefix and suffix", "//abc/xyz/pqr/123//", "abc/xyz/pqr/123" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			addRouteToTree(root, testCase.RoutePath)
			routes := root.getAllRoutes()
			if slices.Contains(routes, testCase.AddedRoutePath) {
				tt.Logf("The route %s was added successfully to the route tree.", testCase.RoutePath)
			} else {
				tt.Errorf("The route %s was not added to the route tree.", testCase.RoutePath)
			}
		})
	}
}

// Test case to validate if a request route path is being matched correctly against the routes present in the route tree.
func Test_RouteTree_MatchRoute(t *testing.T) {
	root := createTree()
	addRouteToTree(root, "/users/list-all")
	addRouteToTree(root, "/users/:userId/get_name")
	testCases := []struct {
		Name string
		RequestRoute string
		MappedRoute string
		PathParamCount int
	} {
		{ "Request Route Path with no path parameters", "/users/list-all", "users/list-all", 0 },
		{ "Request Route Path with a single path parameter", "/users/6/get_name", "users/:userId/get_name", 1 },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			matchInfo := matchRouteInTree(root, testCase.RequestRoute)
			if !strings.EqualFold(testCase.MappedRoute, matchInfo.RoutePath) {
				tt.Errorf("The matched route [%s] returned does not match the expected route path [%s].", matchInfo.RoutePath, testCase.MappedRoute)
			}

			if len(matchInfo.Segments) != testCase.PathParamCount {
				tt.Errorf("The number of path parameters returned (%d) does not match the expected parameter count (%d).", len(matchInfo.Segments), testCase.PathParamCount)
			}
		})
	}
}