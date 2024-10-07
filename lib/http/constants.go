package http

import (
	"strings"
	"strconv"
	"fmt"
	"slices"
	"path/filepath"
	"log"
	"github.com/maheshkumaarbalaji/proteus/lib/fs"
)

const (
	ERROR_MSG_CONTENT_TYPE = "text/html"
	HEADER_LINE_SEPERATOR = "\r\n"
	REQUEST_LINE_SEPERATOR = " "
	HEADER_KEY_VALUE_SEPERATOR = ":"
	VALIDATE_ROUTE_PATTERN = "^[a-zA-z][a-zA-Z0-9_/:-]*"
)

// Global logger instance to log messages generated while HTTP server processes incoming requests.
var SrvLogger *log.Logger = nil
// Global instance of HTTP web server.
var ServerInstance *HttpServer = nil
// Collection of headers supported by the server that has a date value.
var DateHeaders []string
// List of content types supported by the web server.
var AllowedContentTypes map[string]string
// A map containing all the default server configuration values.
var ServerDefaults map[string]string
// List of all versions of HTTP supported by the web server.
var Versions map[string][]string

// Structure to represent a response status code and its associated information.
type respStatus struct {
	// HTTP response status code.
	Code StatusCode
	// Short message for the corresponding status code.
	Message string
	// Error description for error status codes (>=400).
	ErrorDescription string
}

// List of response status codes and their associated information.
var ResponseStatusCodes []respStatus

// Returns the file media type for the given file path.
func GetContentType(CompleteFilePath string) (string, bool) {
	pathType, err := fs.GetPathType(CompleteFilePath)
	if err == nil {
		if pathType == fs.FILE_TYPE_PATH {
			fileExtension := filepath.Ext(CompleteFilePath)
			fileExtension = strings.TrimSpace(fileExtension)
			fileExtension = strings.ToLower(fileExtension)
			contentType, exists := AllowedContentTypes[fileExtension]
			if exists {
				return contentType, exists
			} else {
				return strings.TrimSpace(ServerDefaults["content_type"]), true
			}
		}
	}
	return "", false
}

// Returns the default hostname from the list of default configuration values.
func GetDefaultHostname() string {
	hostname := strings.TrimSpace(ServerDefaults["hostname"])
	return hostname
}

// Returns the default port number from the list of default configuration values.
func GetDefaultPort() int {
	portNumberValue := ServerDefaults["port"]
	portNumber, _ := strconv.Atoi(portNumberValue)
	return portNumber
}

// Returns the value for the given key from server default configuration values.
func GetServerDefaultsValue(key string) string {
	value := ServerDefaults[strings.TrimSpace(key)]
	value = strings.TrimSpace(value)
	return value
}

// Gets the highest version of HTTP supported by the web server.
func GetHighestVersion() string {
	var maxVersion float64 = 0.0
	for versionNo := range Versions {
		currentVersion, err := strconv.ParseFloat(versionNo, 64)
		if err == nil {
			if currentVersion > maxVersion {
				maxVersion = currentVersion
			}
		}
	}

	return fmt.Sprintf("%.1f", maxVersion)
}

// Gets an array of all the versions of HTTP supported by the web server.
func GetAllVersions() []string {
	vers := make([]string, 0)
	for versionNo := range Versions {
		tempVer := strings.TrimSpace(versionNo)
		vers = append(vers, tempVer)
	}

	return vers
}

// Gets the list of allowed HTTP methods supported by the web server for the given HTTP version.
func GetAllowedMethods(version string) string {
	for versionNo, AllowedMethods := range Versions {
		if strings.EqualFold(versionNo, version) {
			return strings.Join(AllowedMethods, ", ")
		}
	}

	return ""
}

// Checks if the given HTTP method is supported by the web server for the given version.
func IsMethodAllowed(version string, requestMethod string) bool {
	for versionNo, AllowedMethods := range Versions {
		if strings.EqualFold(versionNo, version) && slices.Contains(AllowedMethods, requestMethod) {
			return true
		}
	}

	return false
}

// Returns the HTTP response version for the given request version value.
func GetResponseVersion(requestVersion string) string {
	isCompatible := false

	for _, version := range GetAllVersions() {
		if strings.EqualFold(version, requestVersion) {
			isCompatible = true
			break
		}
	}

	if isCompatible {
		return requestVersion
	} else {
		return GetHighestVersion()
	}
}

// Logs an error message to the log file. If logger is not initialized, the error message is printed to stdout.
func LogError(Msg string) {
	Msg = strings.TrimSpace(Msg)
	templateString := fmt.Sprintf("%s  %s  ERROR  %s", getRfc1123Time(), GetServerDefaultsValue("server_name"), Msg)
	if SrvLogger != nil {
		SrvLogger.Printf("%s\n", templateString)
	} else {
		fmt.Printf("%s\n", templateString)
	}
}

// Logs a message to the log file. If logger is not initialized, the message is printed to stdout.
func LogInfo(Msg string) {
	Msg = strings.TrimSpace(Msg)
	templateString := fmt.Sprintf("%s  %s  INFO  %s", getRfc1123Time(), GetServerDefaultsValue("server_name"), Msg)
	if SrvLogger != nil {
		SrvLogger.Printf("%s\n", templateString)
	} else {
		fmt.Printf("%s\n", templateString)
	}
}