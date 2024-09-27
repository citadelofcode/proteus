package http

import (
	"log"
	"net"
	"strconv"
	"path/filepath"
	"errors"
	"strings"
	"github.com/maheshkumaarbalaji/proteus/lib/fs"
)

type HttpServer struct {
	HostAddress string
	PortNumber int
	Socket net.Listener
	SrvLogger *log.Logger
	StaticRouter FileRoutes
}

func (srv *HttpServer) Static(Route string, TargetPath string) {
	if srv.StaticRouter == nil {
		srv.StaticRouter = make(FileRoutes)
	}
	err := srv.StaticRouter.Add(Route, TargetPath)
	if err != nil {
		srv.SrvLogger.Printf("AddStaticRoute() :: %s", err.Error())
		return
	}
}

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
	err := httpRequest.read()
	if err != nil {
		srv.logError(err.Error())
		httpResponse.set(StatusInternalServerError, "", "", StatusInternalServerError.GetErrorContent())
		srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
		return
	}

	responseVersion := srv.getResponseVersion(httpRequest.Version)
	
	switch httpRequest.Method {
	case "GET":
		if !IsMethodAllowed(responseVersion, "GET") {
			srv.logError("'GET' method is not allowed in HTTP version " + responseVersion)
			httpResponse.set(StatusMethodNotAllowed, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusMethodNotAllowed.GetErrorContent())
			httpResponse.Headers.Add("Allow", GetAllowedMethods(responseVersion))
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
			return
		}
		TargetFilePath, staticRouteExists := srv.StaticRouter.Match(httpRequest.ResourcePath)
		if !staticRouteExists {
			srv.logError("A Static route matching the given resource does  not exist")
			httpResponse.set(StatusNotFound, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusNotFound.GetErrorContent())
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
			return
		}
		file, err := srv.getFile(TargetFilePath)
		if err != nil {
			srv.logError(err.Error())
			httpResponse.set(StatusInternalServerError, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusInternalServerError.GetErrorContent())
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
			return
		}
		if !httpRequest.isConditionalGet(file.LastModifiedAt) {
			httpResponse.set(StatusOK, responseVersion, file.ContentType, file.Contents)
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
		} else {
			httpResponse.set(StatusNotModified, responseVersion, file.ContentType, make([]byte, 0))
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
		}
	case "HEAD":
		if !IsMethodAllowed(responseVersion, "HEAD") {
			srv.logError("'HEAD' method is not allowed in HTTP version " + responseVersion)
			httpResponse.set(StatusMethodNotAllowed, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusMethodNotAllowed.GetErrorContent())
			httpResponse.Headers.Add("Allow", GetAllowedMethods(responseVersion))
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
			return
		}
		TargetFilePath, staticRouteExists := srv.StaticRouter.Match(httpRequest.ResourcePath)
		if !staticRouteExists {
			srv.logError("A Static route matching the given resource does  not exist")
			httpResponse.set(StatusNotFound, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusNotFound.GetErrorContent())
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
			return
		}
		file, err := srv.getFile(TargetFilePath)
		if err != nil {
			srv.logError(err.Error())
			httpResponse.set(StatusInternalServerError, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusInternalServerError.GetErrorContent())
			srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
			return
		}
		httpResponse.set(StatusOK, responseVersion, file.ContentType, make([]byte, 0))
		srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
	default:
		srv.logError("The HTTP method is not supported by the server. Allowed Methods are - " + GetAllowedMethods(responseVersion))
		httpResponse.set(StatusMethodNotAllowed, responseVersion, ERROR_MSG_CONTENT_TYPE, StatusMethodNotAllowed.GetErrorContent())
		httpResponse.Headers.Add("Allow", GetAllowedMethods(responseVersion))
		srv.sendResponse(httpRequest, httpResponse, ClientConnection.RemoteAddr().String())
	}
}

func (srv *HttpServer) sendResponse(httpRequest *HttpRequest, httpResponse *HttpResponse, ClientAddress string) {
	err := httpResponse.write()
	if err != nil {
		srv.SrvLogger.Printf("Error occurred while sending response to client (%s): %s", ClientAddress, err.Error())
	} else {
		srv.SrvLogger.Print(getW3CLogLine(httpRequest, httpResponse, ClientAddress))
	}
}

func (srv *HttpServer) logError(errorString string) {
	srv.SrvLogger.Print(errorString)
}

func (srv *HttpServer) getContentType(CompleteFilePath string) (string, error) {
	pathType, err := fs.GetPathType(CompleteFilePath)
	if err != nil {
		return "", err
	}
	if pathType != fs.FILE_TYPE_PATH {
		return "", errors.New("path provided does not point to a file")
	}
	fileExtension := filepath.Ext(CompleteFilePath)
	if fileExtension == "" {
		return "", errors.New("file path provided does not contain a file extension")
	}
	fileExtension = strings.ToLower(fileExtension)
	fileMediaType, exists := GetContentType(fileExtension)
	if !exists {
		fileMediaType = "application/octet-stream"
	}
	
	return fileMediaType, nil
}

func (srv *HttpServer) getFile(CompleteFilePath string) (*fs.File, error) {
	contentType, err := srv.getContentType(CompleteFilePath)
	if err != nil {
		return nil, err
	}
	
	file, err := fs.GetFile(CompleteFilePath, contentType)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (srv *HttpServer) getResponseVersion(requestVersion string) string {
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
		return GetHighestVersion()
	}
}