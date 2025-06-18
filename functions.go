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
