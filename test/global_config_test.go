package test

import (
	"testing"
	"github.com/citadelofcode/proteus/internal"
)

// Test case to valiodate if the global configuration information for "internal" package was loaded to the main memory as expected.
func Test_GlobalConfig_Init(t *testing.T) {
	hasError := false
	if internal.DateHeaders == nil || (internal.DateHeaders != nil && len(internal.DateHeaders) == 0){
		t.Error("The list of date value headers accepted by the HTTP server is not available as configuration.")
		hasError = true
	}

	if internal.AllowedContentTypes == nil || (internal.AllowedContentTypes != nil && len(internal.AllowedContentTypes) == 0){
		t.Error("The list of content types allowed by the HTTP server is not available as configuration.")
		hasError = true
	}

	if internal.ServerDefaults == nil || (internal.ServerDefaults != nil && len(internal.ServerDefaults) == 0){
		t.Error("The default values configured for the HTTP server instance is not available.")
		hasError = true
	}

	if internal.Versions == nil || (internal.Versions != nil && len(internal.Versions) == 0) {
		t.Error("The list of all HTTP versions supported by the server is not available as configuration.")
		hasError = true
	}

	if internal.ResponseStatusCodes == nil || (internal.ResponseStatusCodes != nil && len(internal.ResponseStatusCodes) == 0) {
		t.Error("The list of all HTTP response status codes supported by the web server is not available as configuration.")
		hasError = true
	}

	if !hasError {
		t.Log("All the configuration information needed, are available as expected.")
	}
}
