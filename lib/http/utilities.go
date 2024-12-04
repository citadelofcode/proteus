package http

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"
	"github.com/maheshkumaarbalaji/proteus/lib/config"
)

// Initializes the global variables used in the 'http' package.
func init() {
	ServerConfig, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

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
}

// Creates and returns pointer to a new instance of HTTP request.
func newRequest(Connection net.Conn) *HttpRequest {
	var httpRequest HttpRequest
	httpRequest.initialize()
	reader := bufio.NewReader(Connection)
	httpRequest.setReader(reader)
	return &httpRequest
}

// Creates and returns pointer to a new instance of HTTP response.
func newResponse(Connection net.Conn, request *HttpRequest) *HttpResponse {
	var httpResponse HttpResponse
	httpResponse.initialize(GetResponseVersion(request.Version))
	writer := bufio.NewWriter(Connection)
	httpResponse.setWriter(writer)
	return &httpResponse
}

// Creates and returns pointer to a new instance of Router.
func newRouter() *Router {
	router := new(Router)
	router.Routes = make([]Route, 0)
	router.RouteTree = createTree()
	return router
}

// Returns an instance of HTTP web server.
func NewServer() (*HttpServer, error) {
	if SrvLogger == nil {
		SrvLogger = log.New(os.Stdout, "", log.Ldate | log.Ltime)
	}

	if ServerInstance == nil {
		var server HttpServer
		server.HostAddress = "";
		server.PortNumber = 0
		server.innerRouter = newRouter()
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