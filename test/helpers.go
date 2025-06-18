package test

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"github.com/citadelofcode/proteus/internal"
)

// Test stop function to be passed as the third argument to the middleware.
func Stop() {
	// This is an empty function and does nothing.
}

// Helper function to create a new test server instance.
func NewTestServer(t testing.TB) *internal.HttpServer {
	t.Helper()
	return internal.NewServer("", 0)
}

// Helper function to create a new test request for the given server instance.
func NewTestRequest(t testing.TB, server *internal.HttpServer, reader io.Reader) *internal.HttpRequest {
	t.Helper()
	request := new(internal.HttpRequest)
	request.Initialize(reader)
	request.Server = server
	return request
}

// Helper function to create a new test response for the given server instance.
func NewTestResponse(t testing.TB, version string, server *internal.HttpServer, writer io.Writer) *internal.HttpResponse {
	t.Helper()
	response := new(internal.HttpResponse)
	response.Initialize(version, writer)
	response.Server = server
	return response
}

// Helper function to create a new empty test router with no static or dynamic rotues configured.
func NewTestRouter(t testing.TB) *internal.Router {
	t.Helper()
	return internal.NewRouter()
}

// Helper method to create static assets for testing in the given root directory. This function creates a "static" folder and "index.html" file in the given root directory.
func CreateStaticAssets(t testing.TB, root string) error {
	t.Helper()
	err := CreateDirectories(t, root, []string{ "static" })
	if err != nil {
		return err
	}

	err = CreateFiles(t, root, map[string][]byte {
		"index.html": []byte("<p>This is a sample home page in the root directory!</p>"),
	})
	if err != nil {
		return err
	}

	staticFolder := filepath.Join(root, "static")
	staticFolder = filepath.Clean(staticFolder)
	err = CreateFiles(t, staticFolder, map[string][]byte {
		"file-one.html": []byte("<p>This is a sample page in the static folder of root directory!</p>"),
	})
	if err != nil {
		return err
	}

	return nil
}

// Helper function to create the given list of folders in the root directory provided.
func CreateDirectories(t testing.TB, root string, folders []string) error {
	t.Helper()
	for _, folder := range folders {
		completePath := filepath.Join(root, folder)
		completePath = filepath.Clean(completePath)
		err := os.Mkdir(completePath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper function to create the given list of files in the root directory provided.
func CreateFiles(t testing.TB, root string, files map[string][]byte) error {
	t.Helper()
	for file, contents := range files {
		completePath := filepath.Join(root, file)
		completePath = filepath.Clean(completePath)
		err := os.WriteFile(completePath, contents, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
