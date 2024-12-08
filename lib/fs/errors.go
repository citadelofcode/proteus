package fs

import (
	"fmt"
)

// A custom error to track file system related errors raised.
type FileSystemError struct {
	// The target file path that is causing the error.
	TargetPath string
	// The actual error message raised by the program.
	Message string
}

// Returns a customized error message associated with the instance of FileSystemError.
func (fsf *FileSystemError) Error() string {
	return fmt.Sprintf("Error occurred for file [%s] :: %s", fsf.TargetPath, fsf.Message)
}