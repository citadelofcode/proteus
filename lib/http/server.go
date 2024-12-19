package http

import (
	"fmt"
	"net"
	"strconv"
	"strings"
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
	innerRouter *Router
	// Logger instance associated with the Server instance.
	eventLogger *logger
}

// Define a static route and map to a static file or folder in the file system.
func (srv *HttpServer) Static(Route string, TargetPath string) error {
	err := srv.innerRouter.addStaticRoute("GET", Route, TargetPath)
	if err != nil {
		return err
	}

	err = srv.innerRouter.addStaticRoute("HEAD", Route, TargetPath)
	if err != nil {
		return err
	}

	return nil
}

// Setup the web server instance to listen for incoming HTTP requests at the given hostname and port number.
func (srv * HttpServer) Listen(PortNumber int, HostAddress string) {
	if PortNumber == 0 {
		srv.PortNumber = getDefaultPort()
	} else {
		srv.PortNumber = PortNumber
	}

	if HostAddress == "" {
		srv.HostAddress = getServerDefaults("hostname")
	} else {
		srv.HostAddress = strings.TrimSpace(HostAddress)
	}

	serverAddress := srv.HostAddress + ":" + strconv.Itoa(srv.PortNumber)
	server, err := net.Listen("tcp", serverAddress)
	if err != nil {
		srv.LogError(fmt.Sprintf("Error occurred while setting up listener socket: %s", err.Error()))
		return
	}

	srv.Socket = server
	defer srv.Socket.Close()
	srv.LogInfo(fmt.Sprintf("Web server is listening at http://%s", serverAddress))

	for {
		clientConnection, err := srv.Socket.Accept()
		if err != nil {
			srv.LogError(fmt.Sprintf("Error occurred while accepting a new client: %s", err.Error()))
			continue
		}

		srv.LogInfo(fmt.Sprintf("A new client - %s has connected to the server", clientConnection.RemoteAddr().String()))
		go srv.handleClient(clientConnection)
	}
}

// Handles incoming HTTP requests sent from each individual client trying to connect to the web server instance.
func (srv *HttpServer) handleClient(ClientConnection net.Conn) {
	defer ClientConnection.Close()
	httpRequest := newRequest(ClientConnection)
	err := httpRequest.read()
	if err != nil {
		srv.LogError(err.Error())
		return
	}

	httpResponse := newResponse(ClientConnection, httpRequest)

	if !isMethodAllowed(httpResponse.Version, strings.ToUpper(strings.TrimSpace(httpRequest.Method))) {
		httpResponse.Status(StatusMethodNotAllowed)
		err = ErrorHandler(httpRequest, httpResponse)
		if err != nil {
			srv.LogError(err.Error())
		}
	} else {
		routeHandler, err := srv.innerRouter.matchRoute(httpRequest)
		if err != nil {
			srv.LogError(err.Error())
			httpResponse.Status(StatusNotFound)
			err = ErrorHandler(httpRequest, httpResponse)
			if err != nil {
				srv.LogError(err.Error())
			}
		} else {
			err = routeHandler(httpRequest, httpResponse)
			if err != nil {
				srv.LogError(err.Error())
			}
		}
	}
}

// Creates a new GET endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Get(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("GET", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new HEAD endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Head(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("HEAD", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new POST endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Post(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("POST", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new PUT endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Put(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("PUT", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new DELETE endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Delete(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("DELETE", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new TRACE endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Trace(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("TRACE", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new OPTIONS endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Options(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("OPTIONS", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new CONNECT endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Connect(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("CONNECT", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}

// Logs the given message as an error in the server logs.
func (srv *HttpServer) LogError(message string) {
	message = strings.TrimSpace(message)
	srv.eventLogger.logError(message)
}

// Logs the given message as an information in the server logs.
func (srv *HttpServer) LogInfo(message string) {
	message = strings.TrimSpace(message)
	srv.eventLogger.logInfo(message)
}

// Logs the status for a HTTP request to the server logger.
func (srv *HttpServer) Log(request *HttpRequest, response *HttpResponse) {
	logMsg := fmt.Sprintf("  %s  %s  %s  HTTP/%s  %d  %s", request.ClientAddress, request.Method, request.ResourcePath, request.Version, response.StatusCode, response.StatusMessage)
	if response.StatusCode < 400 {
		srv.eventLogger.logInfo(logMsg)
	} else {
		srv.eventLogger.logError(logMsg)
	}
}