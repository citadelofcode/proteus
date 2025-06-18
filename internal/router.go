package internal

import (
	"strings"
	"path/filepath"
)

// Structure to contain information about a single route declared in the Router.
type Route struct {
	// Handler function to be executed for the route paths.
	RouteHandler RouteHandler
	// HTTP method for which the route is defined
	Method string
	// List of all route level middlewares configured.
	Middlewares []Middleware
}

// Structure to hold all the routes and the associated routing logic.
type Router struct {
	// Prefix tree containing all the routes declared on the router.
	routeTree *PrefixTree
	// Collection of static paths configured for the router.
	staticRoutes map[string]string
	// To access the underlying filesystem and its files/folders.
	fs *FileSystem
}

// Adds a new static route and target folder to the static routes collection.
func (rtr *Router) Static(RoutePath string, TargetPath string) error {
	RoutePath = CleanRoute(RoutePath)
	isAbsolute := rtr.fs.IsAbsolute(TargetPath)
	if !isAbsolute {
		reError := new(RoutingError)
		reError.RoutePath = TargetPath
		reError.Message = "Target path must be absolute"
		return reError
	}
	isDirectory := rtr.fs.IsDirectory(TargetPath)
	if !isDirectory {
		reError := new(RoutingError)
		reError.RoutePath = TargetPath
		reError.Message = "Target path given should point to a directory in the local file system"
		return reError
	}

	rtr.staticRoutes[RoutePath] = TargetPath
	return nil
}

// Creates a new GET endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (rtr *Router) Get(RoutePath string, handlerFunc RouteHandler, middlewareList ...Middleware) error {
	RoutePath = CleanRoute(RoutePath)
	return rtr.addRoute("GET", RoutePath, handlerFunc, middlewareList)
}

// Creates a new HEAD endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (rtr *Router) Head(RoutePath string, handlerFunc RouteHandler, middlewareList ...Middleware) error {
	RoutePath = CleanRoute(RoutePath)
	return rtr.addRoute("HEAD", RoutePath, handlerFunc, middlewareList)
}

// Creates a new POST endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (rtr *Router) Post(RoutePath string, handlerFunc RouteHandler, middlewareList ...Middleware) error {
	RoutePath = CleanRoute(RoutePath)
	return rtr.addRoute("POST", RoutePath, handlerFunc, middlewareList)
}

// Creates a new PUT endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (rtr *Router) Put(RoutePath string, handlerFunc RouteHandler, middlewareList ...Middleware) error {
	RoutePath = CleanRoute(RoutePath)
	return rtr.addRoute("PUT", RoutePath, handlerFunc, middlewareList)
}

// Creates a new DELETE endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (rtr *Router) Delete(RoutePath string, handlerFunc RouteHandler, middlewareList ...Middleware) error {
	RoutePath = CleanRoute(RoutePath)
	return rtr.addRoute("DELETE", RoutePath, handlerFunc, middlewareList)
}

// Creates a new TRACE endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (rtr *Router) Trace(RoutePath string, handlerFunc RouteHandler, middlewareList ...Middleware) error {
	RoutePath = CleanRoute(RoutePath)
	return rtr.addRoute("TRACE", RoutePath, handlerFunc, middlewareList)
}

// Creates a new OPTIONS endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (rtr *Router) Options(RoutePath string, handlerFunc RouteHandler, middlewareList ...Middleware) error {
	RoutePath = CleanRoute(RoutePath)
	return rtr.addRoute("OPTIONS", RoutePath, handlerFunc, middlewareList)
}

// Creates a new CONNECT endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (rtr *Router) Connect(RoutePath string, handlerFunc RouteHandler, middlewareList ...Middleware) error {
	RoutePath = CleanRoute(RoutePath)
	return rtr.addRoute("CONNECT", RoutePath, handlerFunc, middlewareList)
}

// Adds a new dynamic route and its associated handler function to the collection of routes defined in the router instance.
func (rtr *Router) addRoute(Method string, RoutePath string, handlerFunc RouteHandler, middlewareList []Middleware) error {
	RoutePath = CleanRoute(RoutePath)
	Method = strings.TrimSpace(Method)
	Method = strings.ToUpper(Method)

	routeObj := Route{
		RouteHandler: handlerFunc,
		Method: Method,
		Middlewares: make([]Middleware, 0),
	}

	routeObj.Middlewares = append(routeObj.Middlewares, middlewareList...)
	rtr.routeTree.Insert(RoutePath, &routeObj)
	return nil
}

// Function that matches a given route with the route tree and fetches the matched route, uses this route to get the corresponding handler.
func (rtr *Router) Match(request *HttpRequest) (*Route, error) {
	routePath := CleanRoute(request.ResourcePath)
	if strings.EqualFold(request.Method, "GET") || strings.EqualFold(request.Method, "HEAD") {
		for routeKey, TargetPath := range rtr.staticRoutes {
			if strings.HasPrefix(routePath, routeKey) {
				RouteAfterPrefix := strings.TrimPrefix(routePath, routeKey)
				RouteAfterPrefix = CleanRoute(RouteAfterPrefix)
				FinalPath := filepath.Join(TargetPath, RouteAfterPrefix)
				if rtr.fs.Exists(FinalPath) {
					request.Locals["StaticFilePath"] = FinalPath
					finalRoute := new(Route)
					finalRoute.Method = request.Method
					finalRoute.RouteHandler = StaticFileHandler
					finalRoute.Middlewares = make([]Middleware, 0)
					return finalRoute, nil
				}
			}
		}
	}

	routeInfo := rtr.routeTree.Match(routePath)
	if routeInfo.MatchedRoutes == nil {
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

	var finalRoute *Route = nil
	for _, route := range routeInfo.MatchedRoutes {
		if strings.EqualFold(route.Method, request.Method) {
			finalRoute = route
			break
		}
	}

	if finalRoute == nil {
		reError := new(RoutingError)
		reError.RoutePath = routePath
		reError.Message = "matchRoute: A match was not for the HTTP method and route combination"
		return nil, reError
	}

	return finalRoute, nil
}

// Creates a new instance of Router and returns a reference to the instance.
func NewRouter() *Router {
	router := new(Router)
	router.routeTree = EmptyPrefixTree()
	router.staticRoutes = make(map[string]string)
	router.fs = new(FileSystem)
	return router
}
