package http

import (
	"path/filepath"
	"regexp"
	"strings"
	"github.com/mkbworks/proteus/lib/fs"
)

// Structure to contain information about a single route declared in the Router.
type Route struct {
	// Is true if the route path being defined points to a static folder path. Is false if the route path has a dynamic handler declared.
	IsStatic bool
	// Defined only for static routes. This field contains the target folder path mapped to the given route path. It is assigned an empty string for dynamic routes.
	StaticFolderPath string
	// Handler function to be executed for the route paths.
	RouteHandler Handler
	// Represents the order in which the route was defined by the users. This also determines the priority of a path being chosen when a request is being processed.
	SequenceNumber int
	// HTTP method for which the route is defined
	Method string
	// Route path being defined for the router
	RoutePath string
}

// Structure to hold all the routes and the associated routing logic.
type Router struct {
	// Collection of all routes defined in the router.
	Routes []Route
	// Contains the last sequence number generated for a route defined for the router.
	LastSequenceNumber int
	// Contains the prefix tree representation of all the routes
	RouteTree *routeTreeNode
}

// Validates if a given route path is syntactically correct.
func (rtr *Router) validateRoute(routePath string) bool {
	isRouteValid, err := regexp.MatchString("^/[a-zA-z][a-zA-Z0-9_/:-]*[a-zA-Z0-9]$", routePath)
	if err != nil {
		return false
	}

	if !isRouteValid {
		return false
	}
	
	return true
}

// Adds a new static route and target folder to the static routes collection.
func (rtr *Router) addStaticRoute(Method string, RoutePath string, TargetPath string) error {
	RoutePath = cleanRoute(RoutePath)
	TargetPath = strings.TrimSpace(TargetPath)
	Method = strings.TrimSpace(Method)
	Method = strings.ToUpper(Method)
	isRouteValid := rtr.validateRoute(RoutePath)
	if !isRouteValid {
		reError := new(RoutingError)
		reError.RoutePath = RoutePath
		reError.Message = "addStaticRoute: Route contains one or more invalid characters"
		return reError
	}
	isAbsolutePath := filepath.IsAbs(TargetPath)
	if !isAbsolutePath {
		reError := new(RoutingError)
		reError.RoutePath = TargetPath
		reError.Message = "addStaticRoute: Given target folder path is not an absolute path"
		return reError
	}
	PathType, err := fs.GetPathType(TargetPath)
	if err != nil {
		return err
	}
	if PathType == fs.FILE_TYPE_PATH {
		reError := new(RoutingError)
		reError.RoutePath = TargetPath
		reError.Message = "Target path given should point to a directory not a file"
		return reError
	}
	rtr.LastSequenceNumber++
	routeObj := Route{
		IsStatic: true,
		StaticFolderPath: TargetPath,
		RouteHandler: StaticFileHandler,
		SequenceNumber: rtr.LastSequenceNumber,
		Method: Method,
		RoutePath: RoutePath,
	}
	
	rtr.Routes = append(rtr.Routes, routeObj)
	addRouteToTree(rtr.RouteTree, RoutePath)
	return nil
}

// Adds a new dynamic route and its associated handler function to the collection of routes defined in the router instance.
func (rtr *Router) addDynamicRoute(Method string, RoutePath string, handlerFunc Handler) error {
	RoutePath = cleanRoute(RoutePath)
	Method = strings.TrimSpace(Method)
	Method = strings.ToUpper(Method)

	isRouteValid := rtr.validateRoute(RoutePath)
	if !isRouteValid {
		reError := new(RoutingError)
		reError.RoutePath = RoutePath
		reError.Message = "addDynamicRoute: Route contains one or more invalid characters"
		return reError
	}

	rtr.LastSequenceNumber++
	routeObj := Route{
		IsStatic: false,
		StaticFolderPath: "",
		RouteHandler: handlerFunc,
		SequenceNumber: rtr.LastSequenceNumber,
		Method: Method,
		RoutePath: RoutePath,
	}
	
	rtr.Routes = append(rtr.Routes, routeObj)
	addRouteToTree(rtr.RouteTree, RoutePath)
	return nil
}

// Function that matches a given route with the route tree and fetches the matched route, uses this route to get the corresponding handler (static or dynamic).
func (rtr *Router) matchRoute(request *HttpRequest) (Handler, error) {
	routePath := request.ResourcePath
	routeInfo := matchRouteInTree(rtr.RouteTree, routePath)
	if routeInfo.RoutePath == "" {
		reError := new(RoutingError)
		reError.RoutePath = routePath
		reError.Message = "matchRoute: A match was not found in the router route tree"
		return nil, reError
	}

	if routeInfo.Segments.Length() > 0 {
		for key, values := range routeInfo.Segments {
			request.Segments.Add(key, values)
		}
	}

	var handler Handler
	for _, route := range rtr.Routes {
		if strings.EqualFold(routeInfo.RoutePath, route.RoutePath) {
			handler = route.RouteHandler
			if route.IsStatic {
				request.staticFilePath = strings.Replace(request.ResourcePath, routeInfo.RoutePath, route.StaticFolderPath, 1)
			}
			break
		}
	}

	return handler, nil
}