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

func NewRequest(Connection net.Conn) *HttpRequest {
	var httpRequest HttpRequest
	httpRequest.initialize()
	reader := bufio.NewReader(Connection)
	httpRequest.setReader(reader)
	return &httpRequest
}

func NewResponse(Connection net.Conn, request *HttpRequest) *HttpResponse {
	var httpResponse HttpResponse
	httpResponse.initialize(GetResponseVersion(request.Version))
	writer := bufio.NewWriter(Connection)
	httpResponse.setWriter(writer)
	return &httpResponse
}

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

		DateHeaders = make([]string, 0)
		DateHeaders = append(DateHeaders, ServerConfig.DateHeaders...)
		AllowedContentTypes = ServerConfig.AllowedContentTypes
		ServerDefaults = ServerConfig.ServerDefaults
		Versions = ServerConfig.GetVersionMap()
		ServerInstance = &server

		return &server, nil
	}

	return ServerInstance, nil
}

func getRfc1123Time() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC1123)
}

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