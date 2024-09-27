package config

import (
	"fmt"
	"strconv"
	"strings"
)

// Structure to contain the compatibility information for the web server instance. This includes the different versions of HTTP and the corresponding methods supported.
type HttpVersion struct {
	VersionNumber string `json:"versionNumber"`
	AllowedMethods []string `json:"allowed_methods"`
}

// Structure to hold the configuration information exported from "config.json" file.
type Configuration struct {
	AllowedContentTypes map[string]string `json:"content_types"`
	Versions []HttpVersion `json:"versions"`
	ServerDefaults map[string]string `json:"server_defaults"`
	DateHeaders []string `json:"date_headers"`
}

func (cy *Configuration) GetAllVersions() []string {
	versions := make([]string, 0)
	for _, ver := range cy.Versions {
		tempVer := strings.TrimSpace(ver.VersionNumber)
		versions = append(versions, tempVer)
	}

	return versions
}

func (cy *Configuration) GetHighestVersion() string {
	var maxVersion float64 = 0.0
	for _, ver := range cy.Versions {
		currentVersion, err := strconv.ParseFloat(ver.VersionNumber, 64)
		if err == nil {
			if currentVersion > maxVersion {
				maxVersion = currentVersion
			}
		}
	}

	return fmt.Sprintf("%.1f", maxVersion)
}

func (cy *Configuration) GetAllowedMethods(version string) string {
	for _, ver := range cy.Versions {
		if strings.EqualFold(ver.VersionNumber, version) {
			return strings.Join(ver.AllowedMethods, ", ")
		}
	}

	return ""
}

func (cy *Configuration) IsMethodAllowed(version string, reqMethod string) bool {
	for _, ver := range cy.Versions {
		if strings.EqualFold(ver.VersionNumber, version) {
			for _, method := range ver.AllowedMethods {
				if strings.EqualFold(method, reqMethod) {
					return true
				}
			}
		}
	}

	return false
}

func (cy *Configuration) GetContentType(fileExtension string) (string, bool) {
	contentType, exists := cy.AllowedContentTypes[fileExtension]
	return contentType, exists
}

func (cy *Configuration) GetDefaultHostname() string {
	hostname := strings.TrimSpace(cy.ServerDefaults["hostname"])
	return hostname
}

func (cy *Configuration) GetDefaultPort() int {
	portNumberValue := cy.ServerDefaults["port"]
	portNumber, _ := strconv.Atoi(portNumberValue)
	return portNumber
}

func (cy *Configuration) GetDateHeaders() []string {
	return cy.DateHeaders
}