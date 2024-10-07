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

// Structure to hold all the routes and the associated routing logic.
type Router struct {
	// Collection of all static routes defined by the user.
	StaticRoutes map[string]string
	// Collection of all dynamic routes and their corresponding handlers defined.
	DynamicRoutes map[string]Handler
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

// Added a new static route and target folder to the static routes collection.
func (rtr *Router) addStaticRoute(RoutePath string, TargetPath string) error {
	RoutePath = strings.TrimSpace(RoutePath)
	TargetPath = strings.TrimSpace(TargetPath)
	_, ok := rtr.StaticRoutes[RoutePath]
	if ok {
		return errors.New("static route already exists in router")
	}
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
	rtr.StaticRoutes[RoutePath] = TargetPath
	return nil
}

// Matches the request path from HTTP request to all the configured static routes and returns the matched route's target path. This target path will be combined with the remaining unmatched part of the request path and returned back to the calling function.
func (rtr *Router) matchStatic(RequestPath string) (string, bool) {
	for staticRoute, targetPath := range rtr.StaticRoutes {
		if strings.HasPrefix(RequestPath, staticRoute) {
			RequestPath, _ = strings.CutPrefix(RequestPath, staticRoute)
			TargetFilePath := filepath.Join(targetPath, RequestPath)
			return TargetFilePath, true
		}
	}

	return "", false
}

func (rtr *Router) addDynamicRoute(RouteKey string, handlerFunc Handler) error {
	RouteKey = strings.TrimSpace(RouteKey)
	_, ok := rtr.DynamicRoutes[RouteKey]
	if ok {
		return errors.New("dynamic route already exists in router")
	}

	_, RoutePath, found := strings.Cut(RouteKey, " ")
	if found {
		RoutePath = strings.TrimSpace(RoutePath)
		isRouteValid := rtr.validateRoute(RoutePath)
		if !isRouteValid {
			return errors.New("route path contains one or more invalid characters")
		}

		rtr.DynamicRoutes[RouteKey] = handlerFunc
	} else {
		return errors.New("invalid route key value provided")
	}

	return nil
}