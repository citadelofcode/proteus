package http

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"path/filepath"
)

const (
	// Size in bytes for each chunk of data being read from a file.
	CHUNK_SIZE = 1024
	// Type value for paths pointing to a folder in the file system.
	FOLDER_TYPE_PATH = "Folder"
	// Type value for paths pointing to a file in the file system.
	FILE_TYPE_PATH = "File"
)

// Structure to represent a file in the local file system.
type File struct {
	// Contents of the file as a stream of bytes.
	Contents []byte
	// Media type of the file.
	ContentType string
	// Time at which the file was last modified.
	LastModifiedAt time.Time
	// Base name of the file.
	Name string
	// Size of the file in bytes.
	Size int64
}

// Returns the type of the given path i.e., file or folder. An error is returned if the given path is neither a file nor a folder.
func GetPathType(TargetPath string) (string, error) {
	TargetPath = CleanPath(TargetPath)
	fileStat, err := os.Stat(TargetPath)
	if err != nil {
		fsfErr := new(FileSystemError)
		fsfErr.TargetPath = TargetPath
		fsfErr.Message = fmt.Sprintf("GetPathType: Error occurred while fetching file stats: %s", err.Error())
		return "", fsfErr
	}
	fileMode := fileStat.Mode()
	if fileMode.IsDir() {
		return FOLDER_TYPE_PATH, nil
	} else if fileMode.IsRegular() {
		return FILE_TYPE_PATH, nil
	} else {
		nfErr := new(FileSystemError)
		nfErr.TargetPath = TargetPath
		nfErr.Message = "Given path points neither to a file nor to a folder"
		return "", nfErr
	}
}

// Reads the contents of the file available at the given path and returns it as a byte slice.
func ReadFileContents(CompleteFilePath string) ([]byte, error) {
	CompleteFilePath = CleanPath(CompleteFilePath)
	fileContents := make([]byte, 0)
	fileHandler, err := os.Open(CompleteFilePath)
	if err != nil {
		fsfErr := new(FileSystemError)
		fsfErr.TargetPath = CompleteFilePath
		fsfErr.Message = fmt.Sprintf("Error occurred while reading file contents: %s", err.Error())
		return nil, fsfErr
	}
	defer fileHandler.Close()
	reader := bufio.NewReader(fileHandler)
	for {
		chunk := make([]byte, CHUNK_SIZE)
		bytesRead, err := reader.Read(chunk)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		if bytesRead < CHUNK_SIZE {
			newChunk := chunk[0: bytesRead]
			fileContents = append(fileContents, newChunk...)
		} else {
			fileContents = append(fileContents, chunk...)
		}
	}

	return fileContents, nil
}

// Returns pointer to a FILE object that contains metadata for file available at the given path.
// The metadata include file contents, last modified time, base name and size in bytes. If the given path does not point to a file, then an error is returned.
func GetFile(CompleteFilePath string, ContentType string, OnlyMetadata bool) (*File, error) {
	var file File
	pathType, err := GetPathType(CompleteFilePath)
	if err != nil {
		return nil, err
	}
	if pathType == FILE_TYPE_PATH {
		file.ContentType = strings.TrimSpace(ContentType)
		if !OnlyMetadata {
			fileContents, err := ReadFileContents(CompleteFilePath)
			if err != nil {
				return nil, err
			}

			file.Contents = fileContents
		}

		fileStats, _ := os.Stat(CompleteFilePath)
		file.LastModifiedAt = fileStats.ModTime()
		file.Name = fileStats.Name()
		file.Size = fileStats.Size()
		return &file, nil
	} else {
		fsfErr := new(FileSystemError)
		fsfErr.TargetPath = CompleteFilePath
		fsfErr.Message = "Given path does not point to a file"
		return nil, fsfErr
	}
}

// Cleans the path by replacing multiple seperators with a single seperator.
// It also removes any trailing seperators in the given path.
func CleanPath(Path string) string {
	Path = strings.TrimSpace(Path)
	Path = filepath.Clean(Path)
	return Path
}

// Returns a boolean value indicating if the given path is an absolute path.
func IsAbsolute(FilePath string) bool {
	return filepath.IsAbs(CleanPath(FilePath))
}

// Gets the file extension of the given file path without the period (".") precending it.
func GetFileExtension(CompleteFilePath string) string {
	CompleteFilePath = CleanPath(CompleteFilePath)
	fileExtension := filepath.Ext(CompleteFilePath)
	fileExtension = strings.TrimPrefix(fileExtension, ".")
	return fileExtension
}

// Returns the file media type for the given file path.
func GetContentType(CompleteFilePath string) (string, error) {
	pathType, err := GetPathType(CompleteFilePath)
	if err != nil {
		return "", err
	}

	if pathType == FILE_TYPE_PATH {
		fileExtension := strings.ToLower(GetFileExtension(CompleteFilePath))
		contentType, exists := AllowedContentTypes[fileExtension]
		if exists {
			return contentType, nil
		} else {
			defaultContentType := getServerDefaults("content_type").(string)
			return strings.TrimSpace(defaultContentType), nil
		}
	}

	nfErr := new(FileSystemError)
	nfErr.TargetPath = CompleteFilePath
	nfErr.Message = "Given path does not point to a file"
	return "", nfErr
}
