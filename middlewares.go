package proteus

import (
	"github.com/citadelofcode/proteus/internal"
)

// Middleware to parse JSON payloads as request body and store them in the "Body" attribute of the request body.
// To enable JSON parsing of request payloads, call this function with the use() function on the router instance.
var JsonParser = internal.JsonParser

// Middleware to parse URL-Encoded payloads as request body and store them in the "Body" attribute of the request body.
// To enable parsing of url encoded values for request body, call this function with the use() function on the router instance.
var UrlEncoded = internal.UrlEncoded
