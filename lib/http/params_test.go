package http

import (
	"testing"
	"slices"
)

// Test case to validate the working of adding new parameter to the Params collection.
func Test_Params_Add(t *testing.T) {
	testParams := make(Params)
	testCases := []struct {
		Name string
		ParamKey string
		ParamValues []string
		ExpParamCount int
	} {
		{ "Adding first parameter", "Name", []string{ "Proteus" }, 1 },
		{ "Adding values to first parameter", "Name", []string{ "WebServer" }, 1 },
		{ "Adding second parameter", "Age", []string{ "18" }, 2 },
		{ "Adding third parameter", "Value", []string{ "Proteus Web Server" }, 3 },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			testParams.Add(testCase.ParamKey, testCase.ParamValues)
			if len(testParams) == testCase.ExpParamCount {
				tt.Logf("The expected parameter count [%d] matches the actual parameter count [%d].", testCase.ExpParamCount, len(testParams))
			} else {
				tt.Errorf("The expected parameter count [%d] does not match the actual parameter count [%d].", testCase.ExpParamCount, len(testParams))
			}
		})
	}
}

// Test case to validate the working of fetching the values for a 'key' from the params collection.
func Test_Params_Get(t *testing.T) {
	testParams := make(Params)
	testParams.Add("Name", []string{ "proteus" })
	testCases := []struct {
		Name string
		ParamKey string
		ExpParamValues []string
	} {
		{ "Fetching parameter in the collection", "Name", []string{ "proteus" } },
		{ "Fetching parameter not in the collection", "Age", nil },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			values, _ := testParams.Get(testCase.ParamKey)
			if testCase.ExpParamValues == nil {
				if values != nil {
					tt.Errorf("Expected returned values array to be nil, but got a non-nil value instead - %v", values)
				} else {
					tt.Logf("Expected a nil slice, and got exactly that.")
				}

				return
			}

			if slices.Equal(testCase.ExpParamValues, values) {
				tt.Logf("The returned slice of values [%v], matches the expected slice of values [%v]", values, testCase.ExpParamValues)
			} else {
				tt.Errorf("The returned slice of values [%v], does not match the expected slice of values [%v]", values, testCase.ExpParamValues)
			}
		})
	}
}