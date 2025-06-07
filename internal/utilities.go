package internal

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"path"
)

// Returns the value for the given key from server default configuration values.
func GetServerDefaults(key string) any {
	key = strings.TrimSpace(key)
	value := ServerDefaults[key]
	return value
}

// Gets the highest version of HTTP supported by the web server that is less than the given request version.
func GetHighestVersion(requestVersion string) string {
	var maxVersion float64 = 0.0
	requestVersion = strings.TrimSpace(requestVersion)
	reqVer, _ := strconv.ParseFloat(requestVersion, 64)
	for versionNo := range Versions {
		currentVersion, err := strconv.ParseFloat(versionNo, 64)
		if err == nil {
			if currentVersion > maxVersion && currentVersion < reqVer {
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
		return GetHighestVersion(requestVersion)
	}
}

// Returns the current UTC time in RFC 1123 format.
func GetRfc1123Time() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC1123)
}

// Checks if the given date time value corresponds to a valid HTTP date and returns two values.
// First returned is a boolean value which indicates if the given date value conforms to a valid format.
// Second returned is a time.Time value corresponding to the given string and if its invalid, returns the zero time.
func IsHttpDate(value string) (bool, time.Time) {
	rfc1123Time, err := time.Parse(time.RFC1123, value)
	ansicTime, errOne := time.Parse(time.ANSIC, value)
	rfc850Time, errTwo := time.Parse(time.RFC850, value)

	if err == nil {
		return true, rfc1123Time
	} else if errOne == nil {
		return true, ansicTime
	} else if errTwo == nil {
		return true, rfc850Time
	} else {
		return false, time.Time{}
	}
}

// Removes all but one leading '/' and all the trailing '/' from the given route path and returns the cleaned value.
// Replaces all instances of "//" with "/". This function is used for cleaning URL route paths.
// If the given route path is an empty string, it is returned as-is.
func CleanRoute(RoutePath string) string {
	RoutePath = strings.TrimSpace(RoutePath)
	if !strings.EqualFold(RoutePath, "") {
		RoutePath = path.Clean(RoutePath)
		if !strings.HasPrefix(RoutePath, ROUTE_SEPERATOR) {
			RoutePath = path.Join(ROUTE_SEPERATOR, RoutePath)
		}
	}
	return RoutePath
}

// Creates a new instance of HTTP web server and binds it to the given hostname and port number.
// If the hostname is empty, the web server instance is bound to the locahost.
// If the port number is zero, the web server instance is bound to port - 8080.
func NewServer(HostAddress string, PortNumber int) *HttpServer {
	server := new(HttpServer)
	if strings.TrimSpace(HostAddress) == "" {
		defaultHost := GetServerDefaults("hostname").(string)
		server.HostAddress = strings.TrimSpace(defaultHost)
	} else {
		server.HostAddress = strings.TrimSpace(HostAddress)
	}

	if PortNumber == 0 {
		defaultPort := GetServerDefaults("port").(int)
		server.PortNumber = defaultPort
	} else {
		server.PortNumber = PortNumber
	}

	server.Router = NewRouter()
	server.requestLogger = log.New(os.Stdout, "", 0)
	server.logFormat = COMMON_LOGGER
	server.middlewares = make([]Middleware, 0)

	return server
}
