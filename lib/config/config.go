package config

import (
	"strings"
) 

// Structure to contain the compatibility information for the web server instance. This includes the different versions of HTTP and the corresponding methods supported.
type HttpVersion struct {
	// HTTP version number supported by the server.
	VersionNumber string `json:"versionNumber"`
	// List of HTTP methods allowed for the version number.
	AllowedMethods []string `json:"allowed_methods"`
}

// Structure to represent individial HTTP Status Code information.
type HttpStatus struct {
	// HTTP status code. For example - 200, 201 etc.
	Code int
	// HTTP status message. For example - OK, Created etc.
	Message string
	// Error description associated with the status code. This is applicableonly for error status codes (>= 400). For other status codes, the field contains an empty string.
	ErrorDescription string
}

// Structure to hold the configuration information exported from "config.json" file.
type Configuration struct {
	// List of content types supported by the server. It is represented as a map with the key value being a file extension and the value pointing to the media type corresponding to the file extension.
	AllowedContentTypes map[string]string `json:"content_types"`
	// List of HTTP Versions supported by the server.
	Versions []HttpVersion `json:"versions"`
	// List of default configurations for the server instance.
	ServerDefaults map[string]string `json:"server_defaults"`
	// List of date headers processed by the server instance.
	DateHeaders []string `json:"date_headers"`
	// List of response status codes
	ResponseStatus []HttpStatus `json:"status_codes"`
}

// Returns the versions array in server configuration as a map with version number as its key and the array of allowed methods as its value.
func (cy *Configuration) GetVersionMap() map[string][]string {
	versionMap := make(map[string][]string)
	for _, ver := range cy.Versions {
		versionNo := strings.TrimSpace(ver.VersionNumber)
		versionMap[versionNo] = ver.AllowedMethods
	}

	return versionMap
}