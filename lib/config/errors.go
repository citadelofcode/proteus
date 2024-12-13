package config

import (
	"fmt"
)

// Custom error to track the errors raised while reading config information from the JSON file.
type ConfigError struct {
	// Complete error message for the ConfigError raised.
	Message string
}

// This function returns a formatted error message string associated with the ConfigError instance.
func (ce *ConfigError) Error() string {
	return fmt.Sprintf("ConfigError :: %s", ce.Message)
}