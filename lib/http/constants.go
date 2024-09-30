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
	VALIDATE_ROUTE_PATTERN = "^[a-zA-z][a-zA-Z0-9_/:-]*$"
)

var SrvLogger *log.Logger = nil
var ServerInstance *HttpServer = nil
var DateHeaders []string
var AllowedContentTypes map[string]string
var ServerDefaults map[string]string
var Versions map[string][]string


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
				return "application/octet-stream", true
			}
		}
	}
	return "", false
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

func LogError(errorMsg string) {
	if SrvLogger != nil {
		SrvLogger.Printf("%s\n", errorMsg)
	} else {
		fmt.Printf("%s\n", errorMsg)
	}
}