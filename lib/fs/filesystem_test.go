package fs

import (
	"strings"
	"testing"
)

// Test case to validate the working of the GetPathType() function to fetch the path type for a given file system path.
func Test_GetPathType(t *testing.T) {
	testCases := []struct {
		Name string
		testPath string
		ExpectedPathType string
		ExpectedErrorType string
	} {
		{ "Path pointing to a folder", "/Users/maheshkumaarbalaji/Downloads", FOLDER_TYPE_PATH, "" },
		{ "Path pointing to a file", "/Users/maheshkumaarbalaji/Projects/proteus/Files/home.html", FILE_TYPE_PATH, "" },
		{ "Path pointing to neither a file nor a folder", "https://www.google.com", "", "FileSystemError" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			PathType, err := GetPathType(testCase.testPath)
			if testCase.ExpectedErrorType == "" {
				if err != nil {
					tt.Errorf("Was not expecting an error, and yet received one - %v", err)
					return
				}
			}

			if testCase.ExpectedErrorType == "FileSystemError" {
				fsErr, ok := err.(*FileSystemError)
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