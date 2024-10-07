package http

import (
	"net"
	"strconv"
	"strings"
	"fmt"
)

// Structure to create an instance of a web server.
type HttpServer struct {
	// Hostname of the web server instance.
	HostAddress string
	// Port number where web server instance is listening for incoming requests.
	PortNumber int
	// Server socket created and bound to the port number.
	Socket net.Listener
	// Router instance that contains all the routes and their associated handlers.
	innerRouter Router
}

// Define a static route and map to a static file or folder in the file system.
func (srv *HttpServer) Static(Route string, TargetPath string) {
	err := srv.innerRouter.addStaticRoute(Route, TargetPath)
	if err != nil {
		LogError(fmt.Sprintf("Adding Static Route to Server :: %s", err.Error()))
		return
	}
}

// Setup the web server instance to listen for incoming HTTP requests at the given hostname and port number.
func (srv * HttpServer) Listen(PortNumber int, HostAddress string) {
	if PortNumber == 0 {
		srv.PortNumber = GetDefaultPort()
	} else {
		srv.PortNumber = PortNumber
	}

	if HostAddress == "" {
		srv.HostAddress = GetDefaultHostname()
	} else {
		srv.HostAddress = strings.TrimSpace(HostAddress)
	}

	serverAddress := srv.HostAddress + ":" + strconv.Itoa(srv.PortNumber)
	server, err := net.Listen("tcp", serverAddress)
	if err != nil {
		LogError(fmt.Sprintf("Error occurred while setting up listener socket: %s", err.Error()))
		return
	}
	srv.Socket = server
	defer srv.Socket.Close()
	LogInfo(fmt.Sprintf("Web server is listening at http://%s", serverAddress))

	for {
		clientConnection, err := srv.Socket.Accept()
		if err != nil {
			LogError(fmt.Sprintf("Error occurred while accepting a new client: %s", err.Error()))
			continue
		}
		LogInfo(fmt.Sprintf("A new client - %s has connected to the server", clientConnection.RemoteAddr().String()))
		go srv.handleClient(clientConnection)
	}
}

// Handles incoming HTTP requests sent from each individual client trying to connect to the web server instance.
func (srv *HttpServer) handleClient(ClientConnection net.Conn) {
	defer ClientConnection.Close()
	httpRequest := NewRequest(ClientConnection)
	httpRequest.read()
	httpResponse := NewResponse(ClientConnection, httpRequest)
	
	switch httpRequest.Method {
	case "GET":
		if !IsMethodAllowed(httpResponse.Version, "GET") {
			httpResponse.Status(StatusMethodNotAllowed)
			ErrorHandler(httpRequest, httpResponse)
			return
		}
		TargetFilePath, staticRouteExists := srv.innerRouter.matchStatic(httpRequest.ResourcePath)
		if staticRouteExists {
			httpRequest.staticFilePath = TargetFilePath
			StaticFileHandler(httpRequest, httpResponse)
			return	
		}
	case "HEAD":
		if !IsMethodAllowed(httpResponse.Version, "HEAD") {
			httpResponse.Status(StatusMethodNotAllowed)
			ErrorHandler(httpRequest, httpResponse)
			return
		}
		TargetFilePath, staticRouteExists := srv.innerRouter.matchStatic(httpRequest.ResourcePath)
		if !staticRouteExists {
			LogError("A Static route matching the given resource does  not exist")
			httpResponse.Status(StatusNotFound)
			ErrorHandler(httpRequest, httpResponse)
			return
		}

		httpRequest.staticFilePath = TargetFilePath
		StaticFileHandler(httpRequest, httpResponse)
	case "POST":
		if !IsMethodAllowed(httpResponse.Version, "POST") {
			httpResponse.Status(StatusMethodNotAllowed)
			ErrorHandler(httpRequest, httpResponse)
			return
		}
	default:
		httpResponse.Status(StatusMethodNotAllowed)
		ErrorHandler(httpRequest, httpResponse)
	}
}

// Creates a new GET endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Get(routePath string, handlerFunc Handler) {
	routePath = strings.TrimSpace(routePath)
	key := fmt.Sprintf("GET %s", routePath)
	err := srv.innerRouter.addDynamicRoute(key, handlerFunc)
	if err != nil {
		LogError(fmt.Sprintf("Adding Dynamic Route to Server :: %s", err.Error()))
		return
	}
}

// Creates a new HEAD endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Head(routePath string, handlerFunc Handler) {
	routePath = strings.TrimSpace(routePath)
	key := fmt.Sprintf("HEAD %s", routePath)
	err := srv.innerRouter.addDynamicRoute(key, handlerFunc)
	if err != nil {
		LogError(fmt.Sprintf("Adding Dynamic Route to Server :: %s", err.Error()))
		return
	}
}

// Creates a new POST endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Post(routePath string, handlerFunc Handler) {
	routePath = strings.TrimSpace(routePath)
	key := fmt.Sprintf("POST %s", routePath)
	err := srv.innerRouter.addDynamicRoute(key, handlerFunc)
	if err != nil {
		LogError(fmt.Sprintf("Adding Dynamic Route to Server :: %s", err.Error()))
		return
	}
}