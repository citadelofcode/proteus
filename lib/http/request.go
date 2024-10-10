package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"net/textproto"
	"slices"
	"net/url"
	"github.com/maheshkumaarbalaji/proteus/lib/fs"
)

// Structure to represent a HTTP request received by the web server.
type HttpRequest struct {
	// HTTP request method like GET, POST, PUT etc.
	Method string
	// Resource path requested by the client.
	ResourcePath string
	// HTTP version that the request received complies with.
	Version string
	// Collection of all the request headers received.
	Headers Headers
	// Represents the complete contents of the request body.
	Body []byte
	// Total length of the request body (in bytes).
	ContentLength int
	// Streamed reader instance to read the HTTP request from the network stream.
	reader *bufio.Reader
	// Contains the target file path in case the request is for a static file.
	staticFilePath string
	// Collection of all query parameters
	Query Params
	// Collection of all path parameter values
	Segments Params
}

func (req *HttpRequest) initialize() {
	req.Body = make([]byte, 0)
	req.Headers = make(Headers)
	req.Version = GetHighestVersion()
	req.staticFilePath = ""
	req.Query = nil
	req.Segments = nil
}

func (req *HttpRequest) setReader(reader *bufio.Reader) {
	req.reader = reader
}

func (req *HttpRequest) read() {
	err := req.readHeader()
	if err != nil {
		LogError(fmt.Sprintf("Error while reading request headers: %s\n", err.Error()))
		return
	}

	err = req.parseQueryParams()
	if err != nil {
		LogError(fmt.Sprintf("Error while parsing query parameters: %s\n", err.Error()))
		return
	}

	clength, ok := req.Headers.Get("Content-Length")
	if ok {
		req.ContentLength, err = strconv.Atoi(clength)
		if err != nil {
			LogError(fmt.Sprintf("Error while reading value of Content-Length: %s\n", err.Error()))
			return
		}

		err = req.readBody()
		if err != nil {
			LogError(fmt.Sprintf("Error while reading request body: %s\n", err.Error()))
			return
		}
	}
}

func (req *HttpRequest) readHeader() error {
	RequestLineProcessed := false
	HeaderProcessingCompleted := false

	for {
		message, err := req.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		message = strings.TrimSuffix(message, HEADER_LINE_SEPERATOR)
		if len(message) == 0 && !HeaderProcessingCompleted {
			HeaderProcessingCompleted = true
			break
		} else if !RequestLineProcessed {
			RequestLineParts := strings.Split(message, REQUEST_LINE_SEPERATOR)
			if len(RequestLineParts) != 3 {
				return errors.New("request line should contain exactly three values seperated by a single whitespace")
			}
			req.Method = strings.TrimSpace(RequestLineParts[0])
			req.ResourcePath = strings.TrimSpace(RequestLineParts[1])
			tempVersion := strings.TrimSpace(RequestLineParts[2])
			tempVersion, found := strings.CutPrefix(tempVersion, "HTTP/")
			if !found {
				return errors.New("invalid value for HTTP Version in request line")
			}
			req.Version = strings.TrimSpace(tempVersion) 
			RequestLineProcessed = true
		} else {
			HeaderKey, HeaderValue, found := strings.Cut(message, HEADER_KEY_VALUE_SEPERATOR)
			if !found {
				errorString := fmt.Sprintf("error while processing header: %s :: Semicolon is missing", message)
				return errors.New(errorString)
			}

			HeaderKey = strings.TrimSpace(HeaderKey)
			HeaderValue = strings.TrimSpace(HeaderValue)
			req.AddHeader(HeaderKey, HeaderValue)
		}
	}

	return nil
}

func (req *HttpRequest) readBody() error {
	if req.ContentLength > 0 {
		req.Body = make([]byte, req.ContentLength)
		for index := 0; index < req.ContentLength; index++ {
			bodyByte, err := req.reader.ReadByte()
			if err != nil {
				return errors.New("unexpected error occurred. Unable to read request body")
			}
			req.Body[index] = bodyByte
		}
	}

	return nil
}

func (req *HttpRequest) parseQueryParams() error {
	req.Query = make(Params)
	parsedUrl, err := url.Parse(req.ResourcePath)
	if err != nil {
		return err
	}

	queryParams := parsedUrl.Query()
	for paramName, paramValues := range queryParams {
		req.Query.Add(paramName, paramValues)
	}

	return nil
}

func (req *HttpRequest) isConditionalGet(CompleteFilePath string) bool {
	if !strings.EqualFold(req.Method, "GET") {
		return false
	}

	fileMediaType, exists := GetContentType(CompleteFilePath)
	if !exists {
		return false
	}

	file, err := fs.GetFile(CompleteFilePath, fileMediaType, true)
	if err != nil {
		return false
	}

	LastModifiedString, ok := req.Headers.Get("If-Modified-Since")
	if !ok {
		return false
	}
	LastModifiedString = strings.TrimSpace(LastModifiedString)
	LastModifiedSince, err := time.Parse(time.RFC1123, LastModifiedString)
	if err != nil {
		LastModifiedSince, err = time.Parse(time.ANSIC, LastModifiedString)
		if err != nil {
			return false
		}
	}
	if file.LastModifiedAt.After(LastModifiedSince) {
		return false
	}
	return true
}

// Adds a new key-value pair to the request headers collection.
func (req *HttpRequest) AddHeader(HeaderKey string, HeaderValue string) {
	if slices.Contains(DateHeaders, textproto.CanonicalMIMEHeaderKey(HeaderKey)) {
		_, err := time.Parse(time.RFC1123, HeaderValue)
		_, errOne := time.Parse(time.ANSIC, HeaderValue)

		if err == nil || errOne == nil {
			req.Headers.Add(HeaderKey, HeaderValue)
		}
	} else {
		req.Headers.Add(HeaderKey, HeaderValue)
	}
}