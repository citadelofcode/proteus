package http

import (
	"strings"
	"fmt"
)

// Structure to represent each individual node of the route tree (trie tree).
type routeTreeNode struct {
	// Part of the route present between 2 '/'s.
	RoutePart string
	// A slice containing all the child nodes for the current node in the tree.
	Children []*routeTreeNode
}

// Creates and returns pointer to a new node in the route tree.
func newRouteTreeNode(RoutePart string) *routeTreeNode {
	newNode := new(routeTreeNode)
	newNode.RoutePart = strings.TrimSpace(RoutePart)
	newNode.Children = make([]*routeTreeNode, 0)
	return newNode
}

// Creates an empty route tree and returns a pointer to the root node of the tree. An empty route tree contains only the root node with an empty string assigned as its routepart.
func createTree() *routeTreeNode {
	return newRouteTreeNode("")
}

// Normalizes the given route path into a slice of route parts present in the path. This function also removes any leading or trailing space and '/' before getting the route parts.
func normalizeRoute(RoutePath string) []string {
	RoutePath = strings.TrimSpace(RoutePath)
	RoutePath = strings.ToLower(RoutePath)
	RoutePath = strings.TrimSuffix(RoutePath, "/")
	RoutePath = strings.TrimPrefix(RoutePath, "/")
	RouteParts := strings.Split(RoutePath, "/")
	NormalizedParts := make([]string, 0)
	for _, routePart := range RouteParts {
		routePart = strings.TrimSpace(routePart)
		if routePart != "" {
			NormalizedParts = append(NormalizedParts, routePart)
		}
	}

	return NormalizedParts
}

// Inserts the given route path in the route tree.
func addRouteToTree(RouteTree *routeTreeNode, RoutePath string) {
	RouteParts := normalizeRoute(RoutePath)
	RouteTree.insert(RouteParts)
}

// Recursively adds the route parts to the route tree by creating nodes in the tree for individual route parts.
func (rtn *routeTreeNode) insert(RouteParts []string) {
	if len(rtn.Children) == 0 {
		// If the route node does not have any child nodes of its own
		newNode := newRouteTreeNode(RouteParts[0])
		rtn.Children = append(rtn.Children, newNode)
		if len(RouteParts) > 1 {
			newNode.insert(RouteParts[1:])
		}
	} else {
		// If the root node has one or more child nodes
		cnFound := false
		var rtnNode *routeTreeNode
		for _, cl := range rtn.Children {
			if strings.EqualFold(RouteParts[0], cl.RoutePart) {
				cnFound = true
				rtnNode = cl
				break
			}
		}

		if !cnFound {
			// If none of the child nodes of the root node had the first route part of the given route.
			rtnNode = newRouteTreeNode(RouteParts[0])
			rtn.Children = append(rtn.Children, rtnNode)
			if len(RouteParts) > 1 {
				rtnNode.insert(RouteParts[1:])
			}
		} else {
			// If one of the child nodes of the root node contained the first route part of the given route.
			if len(RouteParts) > 1 {
				rtnNode.insert(RouteParts[1:])
			}
		}
	}
}

// Gets the list of all routes from the route tree node to all the leaf nodes in the tree.
func (rtn *routeTreeNode) getAllRoutes() []string {
	routeParts := make([]string, 0)
	if len(rtn.Children) != 0 {
		for _, cn := range rtn.Children {
			childParts := cn.getAllRoutes()
			for _, part := range childParts {
				if rtn.RoutePart != "" {
					routeParts = append(routeParts, fmt.Sprintf("%s/%s", rtn.RoutePart, part))
				} else {
					routeParts = append(routeParts, part)
				}
			}
		}
	} else {
		routeParts = append(routeParts, rtn.RoutePart)
	}

	return routeParts
}