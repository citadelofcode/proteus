package http

import (
	"strings"
	"strconv"
	"fmt"
	"slices"
)

const (
	ERROR_MSG_CONTENT_TYPE = "text/html"
	HEADER_LINE_SEPERATOR = "\r\n"
	REQUEST_LINE_SEPERATOR = " "
	HEADER_KEY_VALUE_SEPERATOR = ":"
	VALIDATE_ROUTE_PATTERN = "^[a-zA-z][a-zA-Z0-9_/:-]*$"
)

var DateHeaders []string
var AllowedContentTypes map[string]string
var ServerDefaults map[string]string
var Versions map[string][]string


func GetContentType(fileExtension string) (string, bool) {
	fileExtension = strings.ToLower(fileExtension)
	contentType, exists := AllowedContentTypes[fileExtension]
	return contentType, exists
}

func GetDefaultHostname() string {
	hostname := strings.TrimSpace(ServerDefaults["hostname"])
	return hostname
}

func GetDefaultPort() int {
	portNumberValue := ServerDefaults["port"]
	portNumber, _ := strconv.Atoi(portNumberValue)
	return portNumber
}

func GetServerName() string {
	serverName := ServerDefaults["serverName"]
	serverName = strings.TrimSpace(serverName)
	return serverName
}

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

func GetAllVersions() []string {
	vers := make([]string, 0)
	for versionNo := range Versions {
		tempVer := strings.TrimSpace(versionNo)
		vers = append(vers, tempVer)
	}

	return vers
}

func GetAllowedMethods(version string) string {
	for versionNo, AllowedMethods := range Versions {
		if strings.EqualFold(versionNo, version) {
			return strings.Join(AllowedMethods, ", ")
		}
	}

	return ""
}

func IsMethodAllowed(version string, requestMethod string) bool {
	for versionNo, AllowedMethods := range Versions {
		if strings.EqualFold(versionNo, version) && slices.Contains(AllowedMethods, requestMethod) {
			return true
		}
	}

	return false
}