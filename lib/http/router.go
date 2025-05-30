package http

import (
	"strings"
)

// Structure to contain information about a single route declared in the Router.
type Route struct {
	// Is true if the route path being defined points to a static folder path. Is false if the route path has a dynamic handler declared.
	IsStatic bool
	// Defined only for static routes. This field contains the target folder path mapped to the given route path. It is assigned an empty string for dynamic routes.
	StaticFolderPath string
	// Handler function to be executed for the route paths.
	RouteHandler RouteHandler
	// Represents the order in which the route was defined by the users. This also determines the priority of a path being chosen when a request is being processed.
	SequenceNumber int
	// HTTP method for which the route is defined
	Method string
	// List of all route level middlewares configured.
	middlewares []Middleware
}

// Structure to hold all the routes and the associated routing logic.
type Router struct {
	// Contains the last sequence number generated for a route defined for the router.
	LastSequenceNumber int
	// Prefix tree containing all the routes declared on the router.
	RouteTree *PrefixTree
}

// Adds a new static route and target folder to the static routes collection.
func (rtr *Router) addStaticRoute(Method string, RoutePath string, TargetPath string) error {
	RoutePath = CleanRoute(RoutePath)
	Method = strings.TrimSpace(strings.ToUpper(Method))
	isAbsolutePath := IsAbsolute(TargetPath)
	if !isAbsolutePath {
		reError := new(RoutingError)
		reError.RoutePath = TargetPath
		reError.Message = "Target path must be absolute"
		return reError
	}
	PathType, err := GetPathType(TargetPath)
	if err != nil {
		return err
	}
	if PathType == FILE_TYPE_PATH {
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
		middlewares: make([]Middleware, 0),
	}

	rtr.RouteTree.Insert(RoutePath, &routeObj)
	return nil
}

// Adds a new dynamic route and its associated handler function to the collection of routes defined in the router instance.
func (rtr *Router) addDynamicRoute(Method string, RoutePath string, handlerFunc RouteHandler, middlewareList []Middleware) error {
	RoutePath = CleanRoute(RoutePath)
	Method = strings.TrimSpace(Method)
	Method = strings.ToUpper(Method)
	rtr.LastSequenceNumber++
	routeObj := Route{
		IsStatic: false,
		StaticFolderPath: "",
		RouteHandler: handlerFunc,
		SequenceNumber: rtr.LastSequenceNumber,
		Method: Method,
		middlewares: make([]Middleware, 0),
	}

	routeObj.middlewares = append(routeObj.middlewares, middlewareList...)
	rtr.RouteTree.Insert(RoutePath, &routeObj)
	return nil
}

// Function that matches a given route with the route tree and fetches the matched route, uses this route to get the corresponding handler (static or dynamic).
func (rtr *Router) matchRoute(request *HttpRequest) (*Route, error) {
	routePath := request.ResourcePath
	routeInfo := rtr.RouteTree.Match(routePath)
	if routeInfo.MatchedRoute == nil {
		reError := new(RoutingError)
		reError.RoutePath = routePath
		reError.Message = "matchRoute: A match was not found in the router's prefix tree"
		return nil, reError
	}

	if routeInfo.Segments.Length() > 0 {
		for key, values := range routeInfo.Segments {
			request.Segments.Add(key, values)
		}
	}

	if routeInfo.MatchedRoute.IsStatic {
		request.staticFilePath = strings.Replace(request.ResourcePath, routeInfo.MatchedPath, routeInfo.MatchedRoute.StaticFolderPath, 1)
	}

	return routeInfo.MatchedRoute, nil
}

// Creates a new instance of Router and returns a reference to the instance.
func NewRouter() *Router {
	router := new(Router)
	router.LastSequenceNumber = 0
	router.RouteTree = EmptyPrefixTree()
	return router
}
