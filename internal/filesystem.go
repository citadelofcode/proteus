package internal

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
)

// Structure to represent a file in the local file system.
type File struct {
	// Base name of the file.
	Name string
	// Complete Path of the file in the local file system.
	Path string
	// Stats interface associated with the given file. If the value is nil, it implies the path points to a file that does not exist.
	stats os.FileInfo
}

// Reads the contents of the file available at the given path and returns it as a byte slice.
func (file *File) Contents() ([]byte, error) {
	CompleteFilePath := file.Path
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

// Gets the file extension of the given file path without the period (".") preceding it and in lowercase.
func (file *File) Extension() string {
	CompleteFilePath := file.Path
	fileExtension := filepath.Ext(CompleteFilePath)
	fileExtension = strings.TrimPrefix(fileExtension, ".")
	fileExtension = strings.TrimSpace(fileExtension)
	fileExtension = strings.ToLower(fileExtension)
	return fileExtension
}

// Returns the media type for the given file path.
func (file *File) MediaType() string {
	fileExtension := file.Extension()
	contentType, exists := AllowedContentTypes[fileExtension]
	if exists {
		return contentType
	} else {
		defaultContentType := GetServerDefaults("content_type").(string)
		return strings.TrimSpace(defaultContentType)
	}
}

// Returns the total size of the file in bytes. If the file does not existsd, it returns zero.
func (file *File) Size() int64 {
	if file.stats == nil {
		return 0
	} else {
		return file.stats.Size()
	}
}

// Returns the last modified time for the file. If the target file does not exist, it returns the zero value for the "time.Time" type.
func (file *File) LastModified() time.Time {
	if file.stats == nil {
		return time.Time{}
	} else {
		return file.stats.ModTime()
	}
}

// Structure to connect to the local file system and access files/folders.
type FileSystem struct {}

// Cleans the path by replacing multiple seperators with a single seperator.
// It also removes any trailing seperators in the given path.
func (fs *FileSystem) CleanPath(Path string) string {
	Path = strings.TrimSpace(Path)
	Path = filepath.Clean(Path)
	return Path
}

// Returns pointer to a FILE object that contains metadata for file available at the given path.
// The metadata include file contents, last modified time, base name and size in bytes. If the given path does not point to a file, then an error is returned.
func (fs *FileSystem) GetFile(CompleteFilePath string) (*File, error) {
	CompleteFilePath = fs.CleanPath(CompleteFilePath)
	fileStat, err := os.Stat(CompleteFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fsfErr := new(FileSystemError)
			fsfErr.TargetPath = CompleteFilePath
			fsfErr.Message = "GetFile :: File or Directory referenced by the given path does not exist in the file system"
			return nil, fsfErr
		}
		fsfErr := new(FileSystemError)
		fsfErr.TargetPath = CompleteFilePath
		fsfErr.Message = fmt.Sprintf("GetFile: Error occurred while fetching file stats: %s", err.Error())
		return nil, fsfErr
	}
	fileMode := fileStat.Mode()
	if fileMode.IsRegular() {
		file := new(File)
		file.Path = CompleteFilePath
		file.Name = filepath.Base(file.Path)
		file.stats = fileStat
		return file, nil
	} else {
		fsfErr := new(FileSystemError)
		fsfErr.TargetPath = CompleteFilePath
		fsfErr.Message = "Given path does not point to a file"
		return nil, fsfErr
	}
}

// Returns a boolean value indicating if the given path is an absolute path.
func (fs *FileSystem) IsAbsolute(CompleteFilePath string) bool {
	CompleteFilePath = fs.CleanPath(CompleteFilePath)
	return filepath.IsAbs(CompleteFilePath)
}

// Returns a boolean value indicating if the given path points to a directory in the file system.
// Itn returns a false if the path points to a folder that does not exist or if the program does not have access to the file system.
func (fs *FileSystem) IsDirectory(CompletePath string) bool {
	CompletePath = fs.CleanPath(CompletePath)
	stats, err := os.Stat(CompletePath)
	if err != nil {
		return false
	}
	mode := stats.Mode()
	if mode.IsDir() {
		return true
	} else {
		return false
	}
}

// Returns a boolean value indicating if the file or folder represented by the given path exists in the file system.
func (fs *FileSystem) Exists(CompletePath string) bool {
	CompletePath = fs.CleanPath(CompletePath)
	_, err := os.Stat(CompletePath)
	return err == nil
}
