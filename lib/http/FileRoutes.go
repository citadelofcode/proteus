package http

import (
	"errors"
	"path/filepath"
	"strings"
	"github.com/maheshkumaarbalaji/proteus/lib/fs"
)

// Holds all the static routes and the mapped target folder in the local file system.
type FileRoutes map[string]string

// Added a new static route and target folder to the collection of file routes.
func (sr FileRoutes) Add(RoutePath string, TargetPath string) error {
	RoutePath = strings.TrimSpace(RoutePath)
	TargetPath = strings.TrimSpace(TargetPath)
	_, ok := sr[RoutePath]
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
	sr[RoutePath] = TargetPath
	return nil
}

// Gets the target folder path mapped to the given static route.
func (sr FileRoutes) Get(Route string) (string, bool) {
	targetPath, ok := sr[Route]
	return targetPath, ok
}

// Matches the request path from HTTP request to all the configured static routes and returns the matched route's target path. This target path will be combined with the remaining unmatched part of the request path and returned back to the calling function.
func (sr FileRoutes) Match(RequestPath string) (string, bool) {
	for staticRoute, targetPath := range sr {
		if strings.HasPrefix(RequestPath, staticRoute) {
			RequestPath, _ = strings.CutPrefix(RequestPath, staticRoute)
			TargetFilePath := filepath.Join(targetPath, RequestPath)
			return TargetFilePath, true
		}
	}

	return "", false
}