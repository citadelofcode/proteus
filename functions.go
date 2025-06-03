package proteus

import (
	"github.com/citadelofcode/proteus/internal"
)

// Creates a new web server capable of accepting HTTP requests.
// It returns a reference to the newly instantiated server.
var CreateServer = internal.NewServer

// Creates a new router to declare endpoints and associated handlers.
// The created router instance must be mapped to a server instance for the route paths to be functional.
var CreateRouter = internal.NewRouter

// Cleans the given URL route path and returns the cleaned route.
var CleanRoute = internal.CleanRoute

// Cleans the given local file system path and returns the cleaned path. This function must be used only for file system paths.
var CleanPath = internal.CleanPath

// Gets the contents of the file available at the given path (points to the loca file system).
var GetFile = internal.GetFile
