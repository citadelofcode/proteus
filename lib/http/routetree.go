package http

import (
	"strings"
	"path"
)

// Contains the match information when a given route is matched with the prefix tree.
type MatchInfo struct {
	// List of all path parameter(s) fetched by comparing the given route with the one matched in the prefix tree.
	Segments Params
	// The matched route in the prefix tree.
	MatchedPath string
	// Route instance associated with the given path.
	MatchedRoute *Route
}

// Structure to represent each individual node of the prefix tree (trie tree).
type PrefixTreeNode struct {
	// Child elements to the current node stored as a map.
	Children map[string]*PrefixTreeNode
	// Route instance mapped to the current node. Default value is nil.
	Route *Route
}

// Creates and returns pointer to a new node in the prefix tree.
func NewPrefixTreeNode() *PrefixTreeNode {
	newNode := new(PrefixTreeNode)
	newNode.Route = nil
	newNode.Children = make(map[string]*PrefixTreeNode)
	return newNode
}

// Structure to represent the prefix tree to be built
type PrefixTree struct {
	// Root node of the prefix tree
	Root *PrefixTreeNode
}

// Creates an empty prefix tree and returns a pointer to the root node of the tree. An empty prefix tree contains only the root node with an empty string assigned as its value.
func EmptyPrefixTree() *PrefixTree {
	pt := new(PrefixTree)
	pt.Root = NewPrefixTreeNode()
	return pt
}

// Inserts the given route path to the prefix tree.
func (pt *PrefixTree) Insert(RoutePath string, MappedRoute *Route) {
	RouteParts := NormalizeRoute(RoutePath)
	if len(RouteParts) == 0 {
		pt.Root.Route = MappedRoute
		return
	}
	Current := pt.Root
	for _, part := range RouteParts {
		if _, exists := Current.Children[part]; !exists {
			Current.Children[part] = NewPrefixTreeNode()
		}
		Current = Current.Children[part]
	}
	Current.Route = MappedRoute
}

// Get all the routes available in the prefix tree.
func (pt *PrefixTree) GetAllRoutes() []string {
	routes := make([]string, 0)
	var traverse func(*PrefixTreeNode, string)
	traverse = func(CurrentNode *PrefixTreeNode, RoutePath string) {
		if CurrentNode.Route != nil {
			routes = append(routes, RoutePath)
		}
		for part, NextNode := range CurrentNode.Children {
			traverse(NextNode, path.Join(RoutePath, part))
		}
	}
	traverse(pt.Root, "/")
	for index := range routes {
		routes[index] = CleanRoute(routes[index])
	}
	return routes
}

// Find a match for the given route in the prefix tree.
func (pt *PrefixTree) Match(RoutePath string) *MatchInfo {
	MatchedRouteInfo := new(MatchInfo)
	MatchedRouteInfo.Segments = make(Params)
	ipRouteParts := NormalizeRoute(RoutePath)
	if len(ipRouteParts) == 0 {
		MatchedRouteInfo.MatchedPath = ROUTE_SEPERATOR
		MatchedRouteInfo.MatchedRoute = pt.Root.Route
		return MatchedRouteInfo
	}
	opRouteParts := make([]string, 0)
	Current := pt.Root
	for _, part := range ipRouteParts {
		Next, exists := Current.Children[part]
		if !exists {
			hasBeenFound := false
			for key, nextNode := range Current.Children {
				paramName, isFound := strings.CutPrefix(key, ":")
				if isFound {
					MatchedRouteInfo.Segments.Add(paramName, []string { part })
					opRouteParts = append(opRouteParts, key)
					Current = nextNode
					hasBeenFound = true
					break
				}
			}
			if !hasBeenFound {
				MatchedRouteInfo.MatchedRoute = nil
				MatchedRouteInfo.MatchedPath = ""
				return MatchedRouteInfo
			}
		} else {
			opRouteParts = append(opRouteParts, part)
			Current = Next
		}
	}
	MatchedRouteInfo.MatchedRoute = Current.Route
	MatchedRouteInfo.MatchedPath = CleanRoute(path.Join(opRouteParts...))
	return MatchedRouteInfo
}

// Normalizes the given route path into a slice of route parts present in the path.
// This function also removes any leading or trailing space and '/' before getting the route parts.
func NormalizeRoute(RoutePath string) []string {
	RoutePath = CleanRoute(RoutePath)
	RoutePath = strings.TrimPrefix(RoutePath, ROUTE_SEPERATOR)
	if strings.EqualFold(RoutePath, "") {
		return make([]string, 0)
	}
	RouteParts := strings.Split(RoutePath, ROUTE_SEPERATOR)
	return RouteParts
}
