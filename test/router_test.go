package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"github.com/citadelofcode/proteus/internal"
)

// Test case to check the working of Static() function of Router instance.
func Test_Router_Static(t *testing.T) {
	testRouter := internal.NewRouter()
	root := t.TempDir()
	staticFolder := filepath.Join(root, "static")
	homeFile := filepath.Join(root, "index.html")
	err := os.Mkdir(staticFolder, 0755)
	if err != nil {
		t.Fatalf("Error occurred while creating static test folder: %s", err.Error())
		return
	}
	err = os.WriteFile(homeFile, []byte("<p>This is a sample home page!</p>"), 0644)
	if err != nil {
		t.Fatalf("Error occurred while creating index.html: %s", err.Error())
		return
	}
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
