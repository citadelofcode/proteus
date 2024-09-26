package config

import (
	"errors"
	"path/filepath"
	"runtime"
	"encoding/json"
	"github.com/maheshkumaarbalaji/proteus/lib/fs"
)

func GetConfig() (*Configuration, error) {
	var config Configuration
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("unable to access call stack to fetch current file being executed")
	}
	currentFilePath, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	currentDirectory := filepath.Dir(currentFilePath)
	configFilePath := filepath.Join(currentDirectory, "config.json")
	fileContents, err := fs.ReadFileContents(configFilePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileContents, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}