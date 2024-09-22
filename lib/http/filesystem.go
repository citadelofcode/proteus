package http

import (
	"time"
	"os"
	"errors"
	"io"
	"bufio"
)

type File struct {
	Contents []byte
	ContentType string
	LastModifiedAt time.Time
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

func GetFileContents(CompleteFilePath string) ([]byte, error) {
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