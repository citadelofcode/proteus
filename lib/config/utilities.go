package config

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"fmt"
	"github.com/mkbworks/proteus/lib/fs"
)

// Returns the configuration information imported from "config.json".
func GetConfig() (*Configuration ,error) {
	var ServerConfig Configuration
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		ce := new(ConfigError)
		ce.Message = "Unable to access call stack to fetch the path for the current file"
		return nil, ce
	}
	currentFilePath, err := filepath.Abs(file)
	if err != nil {
		ce := new(ConfigError)
		ce.Message = fmt.Sprintf("Error while fetching absolute path: %s", err.Error())
		return nil, ce
	}
	currentDirectory := filepath.Dir(currentFilePath)
	configFilePath := filepath.Join(currentDirectory, "config.json")
	fileContents, err := fs.ReadFileContents(configFilePath)
	if err != nil {
		ce := new(ConfigError)
		ce.Message = err.Error()
		return nil, ce
	}

	err = json.Unmarshal(fileContents, &ServerConfig)
	if err != nil {
		ce := new(ConfigError)
		ce.Message = fmt.Sprintf("Error occurred while unmarshalling config file contents: %s", err.Error())
		return nil, ce
	}

	return &ServerConfig, nil
}