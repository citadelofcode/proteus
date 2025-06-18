package test

import (
	"path/filepath"
	"strings"
	"testing"
	"github.com/citadelofcode/proteus/internal"
)

// Test case to validate the working of the IsAbsolute() function of FileSystem.
func Test_FileSystem_IsAbsolute(t *testing.T) {
	root := t.TempDir()
	AbsFolderExists := filepath.Join(root, "abs-exists")
	AbsFolderNoExists := filepath.Join(root, "abs-noexists")
	err := CreateDirectories(t, root, []string{ "abs-exists" })
	if err != nil {
		t.Fatalf("Error occurred while creating test folders: %s", err.Error())
		return
	}

	fs := new(internal.FileSystem)
	testCases := []struct {
		Name string
		IpPath string
		ExpOp bool
	} {
		{ "A valid absolute path that exists in the local file system", AbsFolderExists, true },
		{ "A valid absolute path that does not exist in the local file system", AbsFolderNoExists, true },
		{ "A relative path that exists in the file system", "./rel-exists", false },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			isAbolsute := fs.IsAbsolute(testCase.IpPath)
			if isAbolsute == testCase.ExpOp {
				tt.Logf("The result returned by isAbsolute() for path [%s] matches the expected output", testCase.IpPath)
			} else {
				tt.Errorf("The result returned by isAbsolute() for path [%s] does not match the expected output", testCase.IpPath)
			}
		})
	}
}

// Test case to validate the working of the IsDirectory() method of FileSystem.
func Test_FileSystem_IsDirectory(t *testing.T) {
	root := t.TempDir()
	ExistsFolder := filepath.Join(root, "folderone")
	NotExistsFolder := filepath.Join(root, "foldertwo")
	err := CreateDirectories(t, root, []string{ "folderone" })
	if err != nil {
		t.Fatalf("Error occurred while creating test folders: %s", err.Error())
		return
	}

	ExistsTextFile := filepath.Join(root, "exists.txt")
	err = CreateFiles(t, root, map[string][]byte{
		"exists.txt": []byte("Hello, this is a sample text file"),
	})
	if err != nil {
		t.Fatalf("Error occurred while creating test file: %s", err.Error())
		return
	}

	NotExistsTextFile := filepath.Join(root, "not-exists.txt")
	fs := new(internal.FileSystem)
	testCases := []struct {
		Name string
		IpPath string
		ExpOp bool
	} {
		{ "A folder that exists in the file system", ExistsFolder, true },
		{ "A folder that does not exist in the file system", NotExistsFolder, false },
		{ "A file that exists in the file system", ExistsTextFile, false },
		{ "A file that does not exist in the file system", NotExistsTextFile, false },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			isDir := fs.IsDirectory(testCase.IpPath)
			if isDir == testCase.ExpOp {
				tt.Logf("The result returned by isDirectory() for path [%s] matches the expected output", testCase.IpPath)
			} else {
				tt.Errorf("The result returned by isDirectory() for path [%s] does not match the expected output", testCase.IpPath)
			}
		})
	}
}

// Test case to validate the working of GetFile() method of FileSystem.
func Test_FileSystem_GetFile(t *testing.T) {
	root := t.TempDir()
	Folder := filepath.Join(root, "folder-tc")
	ExistsFile := filepath.Join(root, "exists.txt")
	NotExistsFile := filepath.Join(root, "not-exists.txt")
	err := CreateDirectories(t, root, []string{ "folder-tc" })
	if err != nil {
		t.Fatalf("Error occurred while creating test folder: %s", err.Error())
		return
	}
	err = CreateFiles(t, root, map[string][]byte{
		"exists.txt": []byte("This is a sample text!"),
	})
	if err != nil {
		t.Fatalf("Error occurred while creating test file: %s", err.Error())
		return
	}
	fs := new(internal.FileSystem)
	testCases := []struct {
		Name string
		IpPath string
		ExpSize int64
		ExpExtension string
		ExpMediaType string
		ExpError string
	} {
		{ "A valid text file available in the file system", ExistsFile, 22, "txt", "text/plain", "" },
		{ "A valid text file not available in the file system", NotExistsFile , 0, "txt", "text/plain", "FileSystemError" },
		{ "A path pointing to a folder in the file system", Folder, 0, "", "text/plain", "FileSystemError" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			file, err := fs.GetFile(testCase.IpPath)
			if err != nil {
				if strings.EqualFold(testCase.ExpError, "FileSystemError") {
					fsErr, ok := err.(*internal.FileSystemError)
					if !ok {
						tt.Errorf("Expected a FileSystemError, but received something else: %#v", err)
					} else {
						tt.Logf("Expected a FileSystemError, and received an error of same type - %#v", *fsErr)
					}
				} else {
					tt.Errorf("An error was not expected, but yet received one - %s", err.Error())
				}
				return
			}

			if file.Size() == testCase.ExpSize {
				tt.Logf("The expected file size [%d] matches the computed file size [%d]", testCase.ExpSize, file.Size())
			} else {
				tt.Errorf("The expected file size [%d] does not match the computed file size [%d]", testCase.ExpSize, file.Size())
			}

			if strings.EqualFold(file.MediaType(), testCase.ExpMediaType) {
				tt.Logf("The expected media type [%s] matches the fetched file's media type [%s]", testCase.ExpMediaType, file.MediaType())
			} else {
				tt.Errorf("The expected media type [%s] does not match the fetched file's media type [%s]", testCase.ExpMediaType, file.MediaType())
			}

			if strings.EqualFold(file.Extension(), testCase.ExpExtension) {
				tt.Logf("The expected file extension [%s] matches the fetched file's extension [%s]", testCase.ExpExtension, file.Extension())
			} else {
				tt.Errorf("The expected file extension [%s] does not match the fetched file's extension [%s]", testCase.ExpExtension, file.Extension())
			}
		})
	}
}

// Test case to validate the working of method that fetches the contents of a given file from the underlying file system.
func Test_FileSystem_FetchContents(t *testing.T) {
	root := t.TempDir()
	FileContentsToBeWritten := "This is a sample text\nThis is another line\nAnd another line"
	err := CreateFiles(t, root, map[string][]byte{
		"file-to-be-read.txt": []byte(FileContentsToBeWritten),
	})
	if err != nil {
		t.Fatalf("Error occurred while creating test file: %s", err.Error())
		return
	}
	FileToBeRead := filepath.Join(root, "file-to-be-read.txt")
	fs := new(internal.FileSystem)
	file, err := fs.GetFile(FileToBeRead)
	if err != nil {
		t.Errorf("Error occurred while getting file properties, instead of successful parse opertion: %s", err.Error())
		return
	}
	fileContentBytes, err := file.Contents()
	if err != nil {
		t.Errorf("Error occurred while reading the file contents, instead of successful read opertion: %s", err.Error())
		return
	}

	fileContents := string(fileContentBytes)
	if strings.EqualFold(FileContentsToBeWritten, fileContents) {
		t.Logf("The contents of file [%s] read from the file system matches the expected content", FileToBeRead)
	} else {
		t.Errorf("The contents of file [%s] read from the file system does not match the expected content", FileToBeRead)
	}
}

// Test case to validate the Exists() function of FileSystem.
func Test_FileSystem_Exists(t *testing.T) {
	root := t.TempDir()
	err := CreateDirectories(t, root, []string{ "static" })
	if err != nil {
		t.Fatalf("Error occurred while creating temporary folders for testing: %s", err.Error())
		return
	}
	err = CreateFiles(t, root, map[string][]byte{
		"home.html": []byte("<p>Hello, World!</p>"),
	})
	if err != nil {
		t.Fatalf("Error occurred while creating temporary files for testing: %s", err.Error())
		return
	}
	fs := new(internal.FileSystem)
	testCases := []struct {
		Name string
		IpPath string
		OpValue bool
	} {
		{ "A folder that exists in the file system", filepath.Join(root, "static"), true },
		{ "A folder that does not exist in the file system", filepath.Join(root, "public"), false },
		{ "A file that exists in the file system", filepath.Join(root, "home.html"), true },
		{ "A file that does not exist in the file system", filepath.Join(root, "index.html"), false },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			isExists := fs.Exists(testCase.IpPath)
			if isExists == testCase.OpValue {
				tt.Logf("The expected value [%t] matches the returned value [%t].", testCase.OpValue, isExists)
			} else {
				tt.Errorf("The expected value [%t] does not match the returned value [%t].", testCase.OpValue, isExists)
			}
		})
	}
}
