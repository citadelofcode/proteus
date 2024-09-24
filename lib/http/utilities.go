package http

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"regexp"
	"encoding/json"
	"runtime"
	"path/filepath"
)

func NewRequest(Connection net.Conn) *HttpRequest {
	var httpRequest HttpRequest
	httpRequest.Initialize()
	reader := bufio.NewReader(Connection)
	httpRequest.setReader(reader)
	return &httpRequest
}

func NewResponse(Connection net.Conn) *HttpResponse {
	var httpResponse HttpResponse
	httpResponse.Initialize()
	writer := bufio.NewWriter(Connection)
	httpResponse.setWriter(writer)
	return &httpResponse
}

func NewServer(ServerHost string) (*HttpServer, error) {
	var server HttpServer
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
	server.SrvLogger = logger
	if ServerHost == "" {
		return nil, errors.New("server host address cannot be empty")
	} else {
		server.HostAddress = ServerHost
	}
	
	server.PortNumber = 0
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("unable to access call stack to fetch current file being executed")
	}
	currentFilePath, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	currentDirectory := filepath.Dir(currentFilePath)
	contentTypesJson := filepath.Join(currentDirectory, "assets", "contenttypes.json")
	server.AllowedContentTypes = make(map[string]string)
	fileContents, err := readFileContents(contentTypesJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileContents, &server.AllowedContentTypes)
	if err != nil {
		return nil, err
	}

	compatibilityJson := filepath.Join(currentDirectory, "assets", "compatibility.json")
	server.HttpCompatibility = Compatibility{}
	fileContents, err = readFileContents(compatibilityJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileContents, &server.HttpCompatibility)
	if err != nil {
		return nil, err
	}
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