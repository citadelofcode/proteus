package proteus

import (
	"github.com/citadelofcode/proteus/internal"
)

// Creates a new web server capable of accepting HTTP requests.
// It returns a reference to the newly instantiated server.
var CreateServer = internal.NewServer
