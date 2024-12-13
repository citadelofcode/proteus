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
				tt.Errorf("Number of route parts (%d) does not match the expected route part count (%d)", len(routeParts), testCase.ExpArgCount)
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
	} else {
		t.Logf("The route part for the root node is an empty string as expected for a route tree")
	}

	if len(root.Children) > 0 {
		t.Errorf("The root node for an empty route tree cannot have child nodes")
	} else {
		t.Logf("The root node does not have any children as expected for an empty route tree")
	}
}

// Test case to validate the addition of a new route to the route tree and fetching all the routes present in the tree.
func Test_RouteTree_AddNGetRoute(t *testing.T) {
	root := createTree()
	testCases := []struct {
		Name string
		RoutePath string
		AddedRoutePath string
		ExpRouteCount int
	} {
		{ "Normal Route Path", "/abc/nbsdf/123", "abc/nbsdf/123", 1 },
		{ "Route Path with a suffix", "/ahf/dggdg/sdfgdfg/", "ahf/dggdg/sdfgdfg", 2 },
		{ "Route Path with two slashes in the middle", "/afsf/bfsdf//nfsdnf", "afsf/bfsdf/nfsdnf", 3 }, 
		{ "Route Path with no prefix", "abc/ert/123", "abc/ert/123", 4 },
		{ "Route path with no prefix and a suffix", "abc/dfgdf/sbfusd/124/", "abc/dfgdf/sbfusd/124", 5 },
		{ "Route path with multiple slash characters as prefix", "//abc/xyz/pqr/123", "abc/xyz/pqr/123", 6 }, 
		{ "Route path with multiple slash characters as suffix", "abc/fgbfdg/pqr/123//", "abc/fgbfdg/pqr/123", 7 },
		{ "Route path with multiple slash characters as prefix and suffix", "//abc/bfghgf/pqr/123//", "abc/bfghgf/pqr/123", 8 },
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
				tt.Errorf("The matched route [%s] returned does not match the expected route path [%s]", matchInfo.RoutePath, testCase.MappedRoute)
			} else {
				tt.Logf("The matched route [%s] returned matches the expected route path [%s]", matchInfo.RoutePath, testCase.MappedRoute)
			}

			if len(matchInfo.Segments) != testCase.PathParamCount {
				tt.Errorf("The number of path parameters returned (%d) does not match the expected parameter count (%d).", matchInfo.Segments.Length(), testCase.PathParamCount)
			} else {
				tt.Logf("The number of path parameters returned (%d) matches the expected parameter count (%d).", matchInfo.Segments.Length(), testCase.PathParamCount)
			}
		})
	}
}