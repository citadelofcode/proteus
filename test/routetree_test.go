package test

import (
	"slices"
	"strings"
	"testing"
	"github.com/citadelofcode/proteus/internal"
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
		{ "Root route path", "/", 0 },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			routeParts := internal.NormalizeRoute(testCase.RoutePath)
			if len(routeParts) != testCase.ExpArgCount {
				tt.Errorf("Number of route parts (%d) does not match the expected route part count (%d)", len(routeParts), testCase.ExpArgCount)
			} else {
				tt.Logf("Number of route parts (%d) matches the expected route part count (%d)", len(routeParts), testCase.ExpArgCount)
			}
		})
	}
}

// Test case to validate if an empty route tree is being created correctly.
func Test_RouteTree_Create(t *testing.T) {
	pt := internal.EmptyPrefixTree()
	if pt.Root.Routes != nil {
		t.Errorf("No route instance must be mapped to the root node of an empty prefix tree")
	} else {
		t.Logf("No route instance has been mapped to the root node of the empty prefix tree as expected")
	}

	if len(pt.Root.Children) > 0 {
		t.Errorf("The root node for an empty route tree cannot have children")
	} else {
		t.Logf("The root node does not have any children as expected for an empty route tree")
	}
}

// Test case to validate the addition of a new route to the route tree and fetching all the routes present in the tree.
func Test_RouteTree_AddNGetRoute(t *testing.T) {
	pt := internal.EmptyPrefixTree()
	testCases := []struct {
		Name string
		RoutePath string
		AddedRoutePath string
		ExpRouteCount int
	} {
		{ "Normal Route Path", "/abc/nbsdf/123", "/abc/nbsdf/123", 1 },
		{ "Route Path with a suffix", "/ahf/dggdg/sdfgdfg/", "/ahf/dggdg/sdfgdfg", 2 },
		{ "Route Path with two slashes in the middle", "/afsf/bfsdf//nfsdnf", "/afsf/bfsdf/nfsdnf", 3 },
		{ "Route Path with no prefix", "abc/ert/123", "/abc/ert/123", 4 },
		{ "Route path with no prefix and a suffix", "abc/dfgdf/sbfusd/124/", "/abc/dfgdf/sbfusd/124", 5 },
		{ "Route path with multiple slash characters as prefix", "//abc/xyz/pqr/123", "/abc/xyz/pqr/123", 6 },
		{ "Route path with multiple slash characters as suffix", "abc/fgbfdg/pqr/123//", "/abc/fgbfdg/pqr/123", 7 },
		{ "Route path with multiple slash characters as prefix and suffix", "//abc/bfghgf/pqr/123//", "/abc/bfghgf/pqr/123", 8 },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			pt.Insert(testCase.RoutePath, new(internal.Route))
			routes := pt.GetAllRoutes()
			if slices.Contains(routes, testCase.AddedRoutePath) {
				tt.Logf("The route %s was added successfully to the route tree.", testCase.RoutePath)
			} else {
				tt.Errorf("The route %s was not added to the route tree.", testCase.RoutePath)
			}

			if len(routes) == testCase.ExpRouteCount {
				tt.Logf("The total route count in tree %d matches the expected route count %d", len(routes), testCase.ExpRouteCount)
			} else {
				tt.Errorf("The total route count in tree %d does not match the expected route count %d", len(routes), testCase.ExpRouteCount)
			}
		})
	}
}

// Test case to validate if a request route path is being matched correctly against the routes present in the route tree.
func Test_RouteTree_MatchRoute(t *testing.T) {
	pt := internal.EmptyPrefixTree()
	pt.Insert("/users/list-all", new(internal.Route))
	pt.Insert("/users/:userId/get_name", new(internal.Route))

	testCases := []struct {
		Name string
		RequestRoute string
		MappedRoute string
		PathParamCount int
	} {
		{ "Request Route Path with no path parameters", "/users/list-all", "/users/list-all", 0 },
		{ "Request Route Path with a single path parameter", "/users/6/get_name", "/users/:userId/get_name", 1 },
		{ "Request route without a match in the prefix tree", "/favicon.ico", "", 0 },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			matchInfo := pt.Match(testCase.RequestRoute)
			if !strings.EqualFold(testCase.MappedRoute, matchInfo.MatchedPath) {
				tt.Errorf("The matched route [%s] returned does not match the expected route path [%s]", matchInfo.MatchedPath, testCase.MappedRoute)
			} else {
				tt.Logf("The matched route [%s] returned matches the expected route path [%s]", matchInfo.MatchedPath, testCase.MappedRoute)
			}

			if len(matchInfo.Segments) != testCase.PathParamCount {
				tt.Errorf("The number of path parameters returned (%d) does not match the expected parameter count (%d).", matchInfo.Segments.Length(), testCase.PathParamCount)
			} else {
				tt.Logf("The number of path parameters returned (%d) matches the expected parameter count (%d).", matchInfo.Segments.Length(), testCase.PathParamCount)
			}
		})
	}
}
