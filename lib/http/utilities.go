package http

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"regexp"
	"github.com/maheshkumaarbalaji/proteus/lib/config"
)

// Creates and returns a new instance of HTTP request.
func NewRequest(Connection net.Conn) *HttpRequest {
	var httpRequest HttpRequest
	httpRequest.initialize()
	reader := bufio.NewReader(Connection)
	httpRequest.setReader(reader)
	return &httpRequest
}

// Creates and returns a new instance of HTTP response.
func NewResponse(Connection net.Conn, request *HttpRequest) *HttpResponse {
	var httpResponse HttpResponse
	httpResponse.initialize(GetResponseVersion(request.Version))
	writer := bufio.NewWriter(Connection)
	httpResponse.setWriter(writer)
	return &httpResponse
}

// Creates and returns a new instance of HTTP web server.
func NewServer() (*HttpServer, error) {
	if SrvLogger == nil {
		SrvLogger = log.New(os.Stdout, "", log.Ldate | log.Ltime)
	}

	if ServerInstance == nil {
		var server HttpServer
		server.HostAddress = "";
		server.PortNumber = 0
		ServerConfig, err := config.GetConfig()
		if err != nil {
			return nil, err
		}
		server.innerRouter = Router{}
		server.innerRouter.StaticRoutes = make(map[string]string)
		DateHeaders = make([]string, 0)
		DateHeaders = append(DateHeaders, ServerConfig.DateHeaders...)
		AllowedContentTypes = ServerConfig.AllowedContentTypes
		ServerDefaults = ServerConfig.ServerDefaults
		Versions = ServerConfig.GetVersionMap()
		ResponseStatusCodes = make([]respStatus, 0)
		for _, stat := range ServerConfig.ResponseStatus {
			newStat := respStatus{
				Code: StatusCode(stat.Code),
				Message: stat.Message,
				ErrorDescription: stat.ErrorDescription,
			}
			ResponseStatusCodes = append(ResponseStatusCodes, newStat)
		}
		ServerInstance = &server

		return &server, nil
	}

	return ServerInstance, nil
}

// Returns the current UTC time in RFC 1123 format.
func getRfc1123Time() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC1123)
}

// Validates if a given route is syntactically correct.
func validateRoute(Route string) bool {
	if strings.HasPrefix(Route, "//") || !strings.HasPrefix(Route, "/") {
		return false
	}

	RouteName := strings.TrimPrefix(Route, "/")
	isRouteValid, err := regexp.MatchString(VALIDATE_ROUTE_PATTERN, RouteName)
	if err != nil {
		return false
	}

	if !isRouteValid {
		return false
	}
	
	return true
}