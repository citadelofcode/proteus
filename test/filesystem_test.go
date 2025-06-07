package test

import (
	"strings"
	"testing"
	"runtime"
	"path/filepath"
	"github.com/citadelofcode/proteus/internal"
)

// Test case to validate the working of the GetPathType() function to fetch the path type for a given file system path.
func Test_GetPathType(t *testing.T) {
	testCases := []struct {
		Name string
		testPath string
		ExpectedPathType string
		ExpectedErrorType string
	} {
		{ "Path pointing to a folder", "../assets", internal.FOLDER_TYPE_PATH, "" },
		{ "Path pointing to a file", "../assets/index.html", internal.FILE_TYPE_PATH, "" },
		{ "Path pointing to neither a file nor a folder", "https://www.google.com", "", "FileSystemError" },
	}
	_, CurrentFilePath, _, _ := runtime.Caller(0)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			testCasePath := testCase.testPath
			isAbsolutePath := filepath.IsAbs(testCasePath)
			if !isAbsolutePath {
				testCasePath = filepath.Join(filepath.Dir(CurrentFilePath), testCasePath)
			}
			PathType, err := internal.GetPathType(testCasePath)
			if testCase.ExpectedErrorType == "" {
				if err != nil {
					tt.Errorf("Was not expecting an error, and yet received one - %v", err)
					return
				}
			}

			if testCase.ExpectedErrorType == "FileSystemError" {
				fsErr, ok := err.(*internal.FileSystemError)
				if !ok {
					tt.Errorf("Expected a FileSystemError, but got %v instead", err)
				} else {
					tt.Logf("Received a FileSystemError as expected - %v", fsErr)
				}
				return
			}

			if !strings.EqualFold(PathType, testCase.ExpectedPathType) {
				tt.Errorf("Computed path type (%s) does not match the expected path type (%s)", PathType, testCase.ExpectedPathType)
			} else {
				tt.Logf("Computed path type (%s) matches the expected path type (%s)", PathType, testCase.ExpectedPathType)
			}
		})
	}
}
