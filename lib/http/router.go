package http

import (
	"errors"
	"path/filepath"
	"strings"
	"regexp"
	"github.com/maheshkumaarbalaji/proteus/lib/fs"
)

// Represents a handler function that is executed once any received request is parsed. You can define different handlers for different routes and HTTP methods.
type Handler func (*HttpRequest, *HttpResponse)

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
}

// Validates if a given route is syntactically correct.
func (rtr *Router) validateRoute(Route string) bool {
	if strings.HasPrefix(Route, "//") || !strings.HasPrefix(Route, "/") {
		return false
	}

	RouteName := strings.TrimPrefix(Route, "/")
	isRouteValid, err := regexp.MatchString(VALIDATE_ROUTE_PATTERN, RouteName)
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
	RoutePath = strings.TrimSpace(RoutePath)
	TargetPath = strings.TrimSpace(TargetPath)
	Method = strings.TrimSpace(Method)
	Method = strings.ToUpper(Method)
	isRouteValid := rtr.validateRoute(RoutePath)
	if !isRouteValid {
		return errors.New("route contains one or more invalid characters")
	}
	isAbsolutePath := filepath.IsAbs(TargetPath)
	if !isAbsolutePath {
		return errors.New("parameter 'TargetPath' must be an absolute path")
	}
	PathType, err := fs.GetPathType(TargetPath)
	if err != nil {
		return errors.New("error occurred while determining target path type: " + err.Error())
	}
	if PathType == fs.FILE_TYPE_PATH {
		return errors.New("target path should be a directory")
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
	return nil
}

/* func (rtr *Router) processRequest(request *HttpRequest, response *HttpResponse) {
	Method := strings.TrimSpace(request.Method)
	Method = strings.ToUpper(Method)

} */

// Adds a new dynamic route and its associated handler function to the collection of routes defined in the router instance.
func (rtr *Router) addDynamicRoute(Method string, RoutePath string, handlerFunc Handler) error {
	RoutePath = strings.TrimSpace(RoutePath)
	Method = strings.TrimSpace(Method)
	Method = strings.ToUpper(Method)

	isRouteValid := rtr.validateRoute(RoutePath)
	if !isRouteValid {
		return errors.New("route path contains one or more invalid characters")
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
	return nil
}