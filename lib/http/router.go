package http

import (
	"regexp"
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
	// Route path being defined for the router
	RoutePath string
	// List of all route level middlewares configured.
	middlewares []Middleware
}

// Structure to hold all the routes and the associated routing logic.
type Router struct {
	// Collection of all routes defined in the router.
	Routes []Route
	// Contains the last sequence number generated for a route defined for the router.
	LastSequenceNumber int
	// Prefix tree containing all the routes declared on the router.
	RouteTree *prefixTreeNode
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
	Method = strings.TrimSpace(strings.ToUpper(Method))
	isRouteValid := rtr.validateRoute(RoutePath)
	if !isRouteValid {
		reError := new(RoutingError)
		reError.RoutePath = RoutePath
		reError.Message = "addStaticRoute: Route contains one or more invalid characters"
		return reError
	}
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
		RoutePath: RoutePath,
		middlewares: make([]Middleware, 0),
	}

	rtr.Routes = append(rtr.Routes, routeObj)
	addRouteToTree(rtr.RouteTree, RoutePath)
	return nil
}

// Adds a new dynamic route and its associated handler function to the collection of routes defined in the router instance.
func (rtr *Router) addDynamicRoute(Method string, RoutePath string, handlerFunc RouteHandler, middlewareList []Middleware) error {
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
		middlewares: make([]Middleware, 0),
	}

	routeObj.middlewares = append(routeObj.middlewares, middlewareList...)
	rtr.Routes = append(rtr.Routes, routeObj)
	addRouteToTree(rtr.RouteTree, RoutePath)
	return nil
}

// Function that matches a given route with the route tree and fetches the matched route, uses this route to get the corresponding handler (static or dynamic).
func (rtr *Router) matchRoute(request *HttpRequest) (*Route, error) {
	routePath := request.ResourcePath
	routeInfo := matchRouteInTree(rtr.RouteTree, routePath)
	if routeInfo.routePath == "" {
		reError := new(RoutingError)
		reError.RoutePath = routePath
		reError.Message = "matchRoute: A match was not found in the router prefix tree"
		return nil, reError
	}

	if routeInfo.segments.Length() > 0 {
		for key, values := range routeInfo.segments {
			request.Segments.Add(key, values)
		}
	}

	var finalRoute *Route
	for _, route := range rtr.Routes {
		if strings.EqualFold(routeInfo.routePath, route.RoutePath) {
			finalRoute = &route
			if route.IsStatic {
				request.staticFilePath = strings.Replace(request.ResourcePath, routeInfo.routePath, route.StaticFolderPath, 1)
			}
			break
		}
	}

	return finalRoute, nil
}

// Creates a new instance of Router and returns a reference to the instance.
func NewRouter() *Router {
	router := new(Router)
	router.Routes = make([]Route, 0)
	router.RouteTree = createTree()
	return router
}
