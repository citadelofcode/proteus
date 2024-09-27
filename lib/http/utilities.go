package http

import (
	"bufio"
	"fmt"
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

func NewResponse(Connection net.Conn) *HttpResponse {
	var httpResponse HttpResponse
	httpResponse.initialize()
	writer := bufio.NewWriter(Connection)
	httpResponse.setWriter(writer)
	return &httpResponse
}

func NewServer() (*HttpServer, error) {
	var server HttpServer
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
	server.SrvLogger = logger
	server.HostAddress = "";
	server.PortNumber = 0
	configObj, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	server.Config = configObj
	ServerName = strings.TrimSpace(configObj.ServerDefaults["servername"])
	DateHeaders = configObj.GetDateHeaders()
	DefaultHostname = configObj.GetDefaultHostname()
	DefaultPortNumber = configObj.GetDefaultPort()
	return &server, nil
}

func getRfc1123Time() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC1123)
}

func getW3CLogLine(req *HttpRequest, res *HttpResponse, ClientAddress string) string {
	return fmt.Sprintf("%s %s %s %d %s", ClientAddress, req.Method, req.ResourcePath, res.StatusCode, req.Version)
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