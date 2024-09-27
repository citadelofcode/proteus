package fs

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"time"
)

const (
	CHUNK_SIZE = 1024
	FOLDER_TYPE_PATH = "Folder"
	FILE_TYPE_PATH = "File"
)

type File struct {
	Contents []byte
	ContentType string
	LastModifiedAt time.Time
	Name string
	Size int64
}

func GetPathType(TargetPath string) (string, error) {
	fileStat, err := os.Stat(TargetPath)
	if err != nil {
		return "", err
	}
	fileMode := fileStat.Mode()
	if fileMode.IsDir() {
		return FOLDER_TYPE_PATH, nil
	} else if fileMode.IsRegular() {
		return FILE_TYPE_PATH, nil
	} else {
		return "", errors.New("given path is neither a file nor a folder")
	}
}

func ReadFileContents(CompleteFilePath string) ([]byte, error) {
	fileContents := make([]byte, 0)
	fileHandler, err := os.Open(CompleteFilePath)
	if err != nil {
		return nil, err
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

func GetFile(CompleteFilePath string, ContentType string) (*File, error) {
	var file File
	fileStat, err := os.Stat(CompleteFilePath)
	if err != nil {
		return nil, err
	}
	Mode := fileStat.Mode()
	if Mode.IsRegular() {
		file.ContentType = strings.TrimSpace(ContentType)
		fileContents, err := ReadFileContents(CompleteFilePath)
		if err != nil {
			return nil, err
		}

		file.Contents = fileContents
		file.LastModifiedAt = fileStat.ModTime()
		file.Name = fileStat.Name()
		file.Size = fileStat.Size()
		return &file, nil
	} else {
		return nil, errors.New("given path is not a file")
	}
}