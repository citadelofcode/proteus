package config

import (
	"fmt"
	"strconv"
	"strings"
)

type HttpHeader struct {
	Type string `json:"type"`
	Category string `json:"category"`
	Compatibility []string `json:"compatibility"`
}

type HttpVersion struct {
	VersionNumber string `json:"versionNumber"`
	AllowedMethods []string `json:"allowed_methods"`
}

type Configuration struct {
	AllowedContentTypes map[string]string `json:"content_types"`
	Versions []HttpVersion `json:"versions"`
	ServerDefaults map[string]string `json:"server_defaults"`
	AllowedHeaders map[string]HttpHeader `json:"allowed_headers"`
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

func (cy *Configuration) GetAllowedHeaders() map[string]HttpHeader {
	return cy.AllowedHeaders
}