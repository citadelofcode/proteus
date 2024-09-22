package http

import (
	"strings"
)

type HttpVersion struct {
	VersionNumber string `json:"versionNumber"`
	AllowedMethods []string `json:"allowedMethods"`
	HighestSupported bool `json:"highestSupported"`
}

type Compatibility struct {
	Versions []HttpVersion `json:"versions"`
}

func (cy *Compatibility) getAllVersions() []string {
	versions := make([]string, 0)
	for _, ver := range cy.Versions {
		tempVer := strings.TrimSpace(ver.VersionNumber)
		versions = append(versions, tempVer)
	}

	return versions
}

func (cy *Compatibility) getHighestVersion() string {
	for _, ver := range cy.Versions {
		if ver.HighestSupported {
			return strings.TrimSpace(ver.VersionNumber)
		}
	}

	return "";
}