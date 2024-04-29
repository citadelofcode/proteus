package http

import (
	"log"
	"net"
	"strconv"
	"github.com/maheshkumaarbalaji/project-sparrow/lib/fs"
	"github.com/maheshkumaarbalaji/project-sparrow/lib/router"
)

type HttpServer struct {
	HostAddress string
	PortNumber int
	Socket net.Listener
	SrvLogger *log.Logger
	StaticRouter router.FileRoutes
}

func (srv *HttpServer) Static(Route string, TargetPath string) {
	if srv.StaticRouter == nil {
		srv.StaticRouter = make(router.FileRoutes)
	}
	err := srv.StaticRouter.Add(Route, TargetPath)
	if err != nil {
		srv.SrvLogger.Printf("AddStaticRoute() :: %s", err.Error())
		return
	}
}

func (srv * HttpServer) Listen(PortNumber int) {
	if PortNumber == 0 {
		srv.PortNumber = getRandomPort()
	} else {
		srv.PortNumber = PortNumber
	}
	serverAddress := srv.HostAddress + ":" + strconv.Itoa(srv.PortNumber)
	server, err := net.Listen(SERVER_TYPE, serverAddress)
	if err != nil {
		srv.SrvLogger.Printf("Error occurred while setting up listener socket: %s", err.Error())
		return
	}
	srv.Socket = server
	defer srv.Socket.Close()
	srv.SrvLogger.Printf("Web server is listening at http://%s", serverAddress)

	for {
		clientConnection, err := srv.Socket.Accept()
		if err != nil {
			srv.SrvLogger.Printf("Error occurred while accepting a new client: %s", err.Error())
			continue
		}
		srv.SrvLogger.Printf("A new client - %s has connected to the server", clientConnection.RemoteAddr().String())
		go srv.handleClient(clientConnection)
	}
}

func (srv *HttpServer) handleClient(ClientConnection net.Conn) {
	defer ClientConnection.Close()
	httpRequest := NewRequest(ClientConnection)
	httpResponse := NewResponse(ClientConnection)
	err := httpRequest.Read()
	if err != nil {
		srv.logError(err.Error())
		httpResponse.Set(StatusInternalServerError, "", "", StatusInternalServerError.GetErrorContent())
		srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
		return
	}

	responseVersion := getResponseVersion(httpRequest.Version)

	switch httpRequest.Method {
	case GET_METHOD:
		TargetFilePath, staticRouteExists := srv.StaticRouter.Match(httpRequest.ResourcePath)
		if !staticRouteExists {
			srv.logError("A Static route matching the given resource does  not exist")
			httpResponse.Set(StatusNotFound, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusNotFound.GetErrorContent())
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
			return
		}
		file, err := fs.GetFile(TargetFilePath)
		if err != nil {
			srv.logError(err.Error())
			httpResponse.Set(StatusInternalServerError, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusInternalServerError.GetErrorContent())
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
			return
		}
		if !httpRequest.isConditionalGet(file.LastModifiedAt) {
			httpResponse.Set(StatusOK, responseVersion, file.ContentType, file.Contents)
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
		} else {
			httpResponse.Set(StatusNotModified, responseVersion, file.ContentType, make([]byte, 0))
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
		}
	case HEAD_METHOD:
		TargetFilePath, staticRouteExists := srv.StaticRouter.Match(httpRequest.ResourcePath)
		if !staticRouteExists {
			srv.logError("A Static route matching the given resource does  not exist")
			httpResponse.Set(StatusNotFound, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusNotFound.GetErrorContent())
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
			return
		}
		file, err := fs.GetFile(TargetFilePath)
		if err != nil {
			srv.logError(err.Error())
			httpResponse.Set(StatusInternalServerError, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusInternalServerError.GetErrorContent())
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
			return
		}
		httpResponse.Set(StatusOK, responseVersion, file.ContentType, make([]byte, 0))
		srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
	default:
		srv.logError("The HTTP method is not allowed by the server. Allowed Methods are - " + ALLOWED_METHODS)
		httpResponse.Set(StatusMethodNotAllowed, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusMethodNotAllowed.GetErrorContent())
		httpResponse.Headers.Add("Allow", ALLOWED_METHODS)
		srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
	}
}

func (srv *HttpServer) sendResponse(httpRequest *HttpRequest, httpResponse *HttpResponse, ClientAddress string) {
	err := httpResponse.Write()
	if err != nil {
		srv.SrvLogger.Printf("Error occurred while sending response to client (%s): %s", ClientAddress, err.Error())
	} else {
		srv.SrvLogger.Print(getW3CLogLine(httpRequest, httpResponse, ClientAddress))
	}
}

func (srv *HttpServer) logError(errorString string) {
	srv.SrvLogger.Print(errorString)
}