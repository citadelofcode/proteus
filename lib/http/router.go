package http

import (
	"errors"
	"path/filepath"
	"strings"
	"github.com/maheshkumaarbalaji/proteus/lib/fs"
)

// Represents a handler function that is executed once any received request is parsed. You can define different handlers for different routes and HTTP methods.
type Handler func (*HttpRequest, *HttpResponse)

// Structure to hold all the routes and the associated routing logic.
type Router struct {
	StaticRoutes map[string]string
}

// Added a new static route and target folder to the static routes collection.
func (rtr *Router) AddStatic(RoutePath string, TargetPath string) error {
	RoutePath = strings.TrimSpace(RoutePath)
	TargetPath = strings.TrimSpace(TargetPath)
	_, ok := rtr.StaticRoutes[RoutePath]
	if ok {
		return errors.New("static route already exists in server")
	}
	isRouteValid := validateRoute(RoutePath)
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

// Gets the target folder path mapped to the given static route.
func (rtr *Router) GetStatic(Route string) (string, bool) {
	targetPath, ok := rtr.StaticRoutes[Route]
	return targetPath, ok
}

// Matches the request path from HTTP request to all the configured static routes and returns the matched route's target path. This target path will be combined with the remaining unmatched part of the request path and returned back to the calling function.
func (rtr *Router) MatchStatic(RequestPath string) (string, bool) {
	for staticRoute, targetPath := range rtr.StaticRoutes {
		if strings.HasPrefix(RequestPath, staticRoute) {
			RequestPath, _ = strings.CutPrefix(RequestPath, staticRoute)
			TargetFilePath := filepath.Join(targetPath, RequestPath)
			return TargetFilePath, true
		}
	}

	return "", false
}