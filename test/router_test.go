package test

import (
	"strings"
	"testing"
	"path/filepath"
	"github.com/citadelofcode/proteus/internal"
)

// Test case to check the working of Static() function of Router instance.
func Test_Router_Static(t *testing.T) {
	testRouter := NewTestRouter(t)
	root := t.TempDir()
	err := CreateStaticAssets(t, root)
	if err != nil {
		t.Fatalf("Error occurred while creating static test assets: %s", err.Error())
		return
	}
	staticFolder := filepath.Join(root, "static")
	homeFile := filepath.Join(root, "index.html")
	testCases := []struct {
		Name string
		InputRoute string
		TargetPath string
		ExpectedErr string
	} {
		{ "Valid route with valid absolute target folder path", "/files/static", staticFolder, "" },
		{ "Valid route with a absolute target file path", "/files/staticone", homeFile, "RoutingError" },
		{ "Valid route with a relative target folder path", "/files/statictwo", "./statictwo", "RoutingError" },
		{ "Valid route with a relative target file path", "/files/staticthree", "./staticthree.html", "RoutingError" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			err := testRouter.Static(testCase.InputRoute, testCase.TargetPath)
			if err != nil {
				if strings.EqualFold(testCase.ExpectedErr, "RoutingError") {
					rtrError, ok := err.(*internal.RoutingError)
					if !ok {
						tt.Errorf("Expected a routing error while adding static route to router, but got this instead - %#v", err)
					} else {
						tt.Logf("Was expecting a routing error and got a routing error as well - %#v", rtrError)
					}
				} else {
					tt.Errorf("Was not expecting an error for adding static route to router and yet got this instead - %#v", err)
				}
				return
			}

			tt.Logf("As expected, route [%s] and associated path [%s] have been added to the router's list of static routes", testCase.InputRoute, testCase.TargetPath)
		})
	}
}


// Test case to validate the working of static route match for the Router structure.
func Test_Router_StaticMatch(t *testing.T) {
	testRouter := NewTestRouter(t)
	testServer := NewTestServer(t)
	root := t.TempDir()
	err := CreateStaticAssets(t, root)
	if err != nil {
		t.Fatalf("Error occurred while creating static test assets: %s", err.Error())
		return
	}
	staticFolder := filepath.Join(root, "static")
	FileInsideStatic := filepath.Join(staticFolder, "file-one.html")
	err = testRouter.Static("/public", staticFolder)
	if err != nil {
		t.Fatalf("Error occurred while setting up the necessary static routes: %s", err.Error())
		return
	}
	testCases := []struct {
		Name string
		RequestPath string
		RequestMethod string
		ExpStaticPath string
		ExpError string
	} {
		{ "Valid GET request path pointing to a file inside the static folder", "/public/file-one.html", "GET", FileInsideStatic, "" },
		{ "Invalid request path pointing to a file not inside the static folder", "/public/file-two.html", "GET", "", "RoutingError" },
		{ "A valid HEAD request path pointing to a file inside the static folder", "/public/file-one.html", "HEAD", FileInsideStatic, "" },
		{ "A request that is not a GET or HEAD request", "/public/file-one.html", "POST", "", "RoutingError" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			request := NewTestRequest(tt, testServer, nil)
			request.ResourcePath = testCase.RequestPath
			request.Method = testCase.RequestMethod
			route, err := testRouter.Match(request)
			if err != nil {
				if strings.EqualFold(testCase.ExpError, "") {
					tt.Errorf("Was not expecting an error, but yet got one - %#v", err)
				} else {
					routingErr, ok := err.(*internal.RoutingError)
					if !ok {
						tt.Errorf("Was expecting a routing error, but got this instead - %#v", err)
					} else {
						tt.Logf("Was expecting a routing error, and got one - %#v", routingErr)
					}
				}
				return
			}

			if strings.EqualFold(route.Method, testCase.RequestMethod) {
				tt.Logf("The request method [%s] matches the matched route instance's method [%s]", testCase.RequestMethod, route.Method)
			} else {
				tt.Errorf("The request method [%s] does not match the matched route instance's method [%s]", testCase.RequestMethod, route.Method)
			}

			staticFilePath, ok := request.Locals["StaticFilePath"].(string)
			if !ok {
				tt.Error("The static file path was supposed to be in the request locals object, but was missing instead")
			} else {
				if strings.EqualFold(staticFilePath, testCase.ExpStaticPath) {
					tt.Logf("The expected static file path [%s] matches the received static file path [%s].", testCase.ExpStaticPath, staticFilePath)
				} else {
					tt.Errorf("The expected static file path [%s] does not match the received static file path [%s].", testCase.ExpStaticPath, staticFilePath)
				}
			}
		})
	}
}

// Test case to validate the working of static route match for the Router structure.
func Test_Router_DynamicMatch(t *testing.T) {
	testRouter := NewTestRouter(t)
	testServer := NewTestServer(t)
	err := testRouter.Post("/link/:linkId", func(request *internal.HttpRequest, response *internal.HttpResponse) {
		request.Server.Log("New link has been added!", internal.INFO_LEVEL)
	}, func(request *internal.HttpRequest, response *internal.HttpResponse, stop internal.StopFunction) {
		request.Server.Log("Middleware processed for the route", internal.INFO_LEVEL)
	})
	if err != nil {
		t.Fatalf("Failed to setup POST route: %s", err.Error())
		return
	}

	err = testRouter.Delete("/link/:linkId", func(request *internal.HttpRequest, response *internal.HttpResponse) {
		request.Server.Log("Given link has been deleted!", internal.INFO_LEVEL)
	})
	if err != nil {
		t.Fatalf("Failed to setup DELETE route: %s", err.Error())
		return
	}

	err = testRouter.Get("/link/:linkId", func(request *internal.HttpRequest, response *internal.HttpResponse) {
		request.Server.Log("Request is being redirected to mapped url", internal.INFO_LEVEL)
	})
	if err != nil {
		t.Fatalf("Failed to setup GET route: %s", err.Error())
		return
	}

	err = testRouter.Post("/user/create", func(request *internal.HttpRequest, response *internal.HttpResponse) {
		request.Server.Log("New user has been added!", internal.INFO_LEVEL)
	})
	if err != nil {
		t.Fatalf("Failed to setup POST route: %s", err.Error())
		return
	}

	testCases := []struct {
		Name string
		RequestRoute string
		RequestMethod string
		ExpMiddlewareCount int
		ExpParamCount int
		ExpError string
	} {
		{ "A GET request with one path parameter", "/link/6", "GET", 0, 1, "" },
		{ "A DELETE request with one path parameter", "/link/7", "DELETE", 0, 1, "" },
		{ "A POST request with no path parameters", "/user/create", "POST", 0, 0, "" },
		{ "A invalid route not configured in the router", "/user/list-all", "GET", 0, 0, "RoutingError" },
		{ "A valid route with invalid HTTP method", "/link/6", "PUT", 0, 1, "RoutingError" },
		{ "A valid route with middleware", "/link/6", "POST", 1, 1, "" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			request := NewTestRequest(tt, testServer, nil)
			request.Method = testCase.RequestMethod
			request.ResourcePath = testCase.RequestRoute
			route, err := testRouter.Match(request)
			if err != nil {
				if strings.EqualFold(testCase.ExpError, "") {
					tt.Errorf("Was not expecting an error, but yet got one - %#v", err)
				} else {
					routingErr, ok := err.(*internal.RoutingError)
					if !ok {
						tt.Errorf("Was expecting a routing error, but got this instead - %#v", err)
					} else {
						tt.Logf("Was expecting a routing error, and got one - %#v", routingErr)
					}
				}
				return
			}

			if strings.EqualFold(route.Method, testCase.RequestMethod) {
				tt.Logf("The request method [%s] matches the matched route instance's method [%s]", testCase.RequestMethod, route.Method)
			} else {
				tt.Errorf("The request method [%s] does not match the matched route instance's method [%s]", testCase.RequestMethod, route.Method)
			}

			if len(request.Segments) == testCase.ExpParamCount {
				tt.Logf("The expected path parameter count [%d] matches the received parameter count [%d].", testCase.ExpParamCount, len(request.Segments))
			} else {
				tt.Errorf("The expected path parameter count [%d] does not match the received parameter count [%d].", testCase.ExpParamCount, len(request.Segments))
			}

			if len(route.Middlewares) == testCase.ExpMiddlewareCount {
				tt.Logf("The expected middleware count [%d] matches the received middleware count [%d].", testCase.ExpMiddlewareCount, len(route.Middlewares))
			} else {
				tt.Errorf("The expected middleware count [%d] does not match the received middleware count [%d].", testCase.ExpMiddlewareCount, len(route.Middlewares))
			}
		})
	}
}
