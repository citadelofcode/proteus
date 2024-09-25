package http

import (
	"errors"
	"path/filepath"
	"strings"
	"github.com/maheshkumaarbalaji/proteus/lib/fs"
)

type FileRoutes map[string]string

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

func (sr FileRoutes) Get(Route string) (string, bool) {
	targetPath, ok := sr[Route]
	return targetPath, ok
}

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