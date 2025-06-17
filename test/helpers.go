package test

import (
	"testing"
	"github.com/citadelofcode/proteus/internal"
)

// Helper function to return a new instance of HTTP request and response.
func TCParams(t testing.TB) (*internal.HttpRequest, *internal.HttpResponse) {
	t.Helper()
	request := new(internal.HttpRequest)
	request.Initialize(nil)
	request.Version = "1.1"
	response := new(internal.HttpResponse)
	response.Initialize(request.Version, nil)
	return request, response
}

// Test stop function to be passed as the third argument to the middleware.
func Stop() {
	// This is an empty function and does nothing.
}
