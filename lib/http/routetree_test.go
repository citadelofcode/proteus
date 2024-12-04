package http

import (
	"slices"
	"testing"
)

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

func Test_RouteTree(t *testing.T) {
	root := createTree()
	routes := root.getAllRoutes()
	if len(routes) != 1 {
		t.Fatalf("Empty route tree was not created as expected")
		return
	}

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
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			addRouteToTree(root, testCase.RoutePath)
			routes = root.getAllRoutes()
			if slices.Contains(routes, testCase.AddedRoutePath) {
				tt.Logf("The route %s was added successfully to the route tree.", testCase.RoutePath)
			} else {
				tt.Errorf("The route %s was not added to the route tree.", testCase.RoutePath)
			}
		})
	}
}