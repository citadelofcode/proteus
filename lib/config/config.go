package config

import (
	"strings"
)

type HttpVersion struct {
	VersionNumber string `json:"versionNumber"`
	AllowedMethods []string `json:"allowedMethods"`
	HighestSupported bool `json:"highestSupported"`
}

type Configuration struct {
	AllowedContentTypes map[string]string `json:"content_types"`
	Versions []HttpVersion `json:"versions"`
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
	for _, ver := range cy.Versions {
		if ver.HighestSupported {
			return strings.TrimSpace(ver.VersionNumber)
		}
	}

	return "";
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