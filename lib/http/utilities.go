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
	"io"
	"path/filepath"
	"regexp"
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

	return &server, nil
}

func getRfc1123Time() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC1123)
}

func getW3CLogLine(req *HttpRequest, res *HttpResponse, ClientAddress string) string {
	return fmt.Sprintf("%s %s %s %d %s", ClientAddress, req.Method, req.ResourcePath, res.StatusCode, req.Version)
}

func getResponseVersion(requestVersion string) string {
	isCompatible := false
	for _, version := range COMPATIBLE_VERSIONS {
		if strings.EqualFold(version, requestVersion) {
			isCompatible = true
			break
		}
	}

	if isCompatible {
		return requestVersion
	} else {
		return MAX_VERSION
	}
}

func GetPathType(TargetPath string) (string, error) {
	fileStat, err := os.Stat(TargetPath)
	if err != nil {
		return "", err
	}
	fileMode := fileStat.Mode()
	if fileMode.IsDir() {
		return FOLDER_TYPE_PATH, nil
	} else if fileMode.IsRegular() {
		return FILE_TYPE_PATH, nil
	} else {
		return "", errors.New("given path is neither a file nor a folder")
	}
}

func GetFileContents(CompleteFilePath string) ([]byte, error) {
	fileContents := make([]byte, 0)
	fileHandler, err := os.Open(CompleteFilePath)
	if err != nil {
		return nil, err
	}
	defer fileHandler.Close()
	reader := bufio.NewReader(fileHandler)
	for {
		chunk := make([]byte, CHUNK_SIZE)
		bytesRead, err := reader.Read(chunk)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		if bytesRead < CHUNK_SIZE {
			newChunk := chunk[0: bytesRead]
			fileContents = append(fileContents, newChunk...)
		} else {
			fileContents = append(fileContents, chunk...)
		}
	}

	return fileContents, nil
}

func GetContentType(CompleteFilePath string) (string, error) {
	pathType, err := GetPathType(CompleteFilePath)
	if err != nil {
		return "", err
	}
	if pathType != FILE_TYPE_PATH {
		return "", errors.New("path provided does not point to a file")
	}
	fileExtension := filepath.Ext(CompleteFilePath)
	if fileExtension == "" {
		return "", errors.New("file path provided does not contain a file extension")
	}
	fileExtension = strings.ToLower(fileExtension)
	fileMediaType := ""
	switch fileExtension {
	case "pdf": 
		fileMediaType = "application/pdf"
	case "htm", "html":
		fileMediaType = "text/html"
	case "css":
		fileMediaType = "text/css"
	case "js", "mjs":
		fileMediaType = "text/javascript"
	case "jpg", "jpeg":
		fileMediaType = "image/jpeg"
	case "gif":
		fileMediaType = "image/gif"
	case "png":
		fileMediaType = "image/png"
	case "apng":
		fileMediaType = "image/apng"
	case "svg", "svgz":
		fileMediaType = "image/svg+xml"
	case "tif", "tiff":
		fileMediaType = "image/tiff"
	case "aac":
		fileMediaType = "audio/aac"
	case "abw":
		fileMediaType = "application/x-abiword"
	case "arc":
		fileMediaType = "application/x-freearc"
	case "avif":
		fileMediaType = "image/avif"
	case "avi":
		fileMediaType = "video/x-msvideo"
	case "azw":
		fileMediaType = "application/vnd.amazon.ebook"
	case "bin":
		fileMediaType = "application/octet-stream"
	case "bmp":
		fileMediaType = "image/bmp"
	case "bz":
		fileMediaType = "application/x-bzip"
	case "bz2":
		fileMediaType = "application/x-bzip2"
	case "cda":
		fileMediaType = "application/x-cdf"
	case "csh":
		fileMediaType = "application/x-csh"
	case "csv":
		fileMediaType = "text/csv"
	case "doc":
		fileMediaType = "application/msword"
	case "docx":
		fileMediaType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "eot":
		fileMediaType = "application/vnd.ms-fontobject"
	case "epub":
		fileMediaType = "application/epub+zip"
	case "gz":
		fileMediaType = "application/gzip"
	case "ico":
		fileMediaType = "image/vnd.microsoft.icon"
	case "ics":
		fileMediaType = "text/calendar"
	case "jar":
		fileMediaType = "application/java-archive"
	case "json":
		fileMediaType = "application/json"
	case "jsonld":
		fileMediaType = "application/ld+json"
	case "odp":
		fileMediaType = "application/vnd.oasis.opendocument.presentation"
	case "odt":
		fileMediaType = "application/vnd.oasis.opendocument.text"
	case "ods":
		fileMediaType = "application/vnd.oasis.opendocument.spreadsheet"
	case "php":
		fileMediaType = "application/x-httpd-php"
	case "ppt":
		fileMediaType = "application/vnd.ms-powerpoint"
	case "pptx":
		fileMediaType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case "rtf":
		fileMediaType = "application/rtf"
	case "xhtml":
		fileMediaType = "application/xhtml+xml"
	case "xls":
		fileMediaType = "application/vnd.ms-excel"
	case "xlsx":
		fileMediaType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "xml":
		fileMediaType = "application/xml"
	case "zip":
		fileMediaType = "application/zip"
	default:
		fileMediaType = "application/octet-stream"
	}

	return fileMediaType, nil
}

func GetFile(CompleteFilePath string) (*File, error) {
	var file File
	contentType, err := GetContentType(CompleteFilePath)
	if err != nil {
		return nil, err
	}
	fileContents, err := GetFileContents(CompleteFilePath)
	if err != nil {
		return nil, err
	}
	file.Contents = fileContents
	file.ContentType = contentType
	return &file, nil
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