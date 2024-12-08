package http

import (
	"bufio"
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
	// HTTP version that the request complies with.
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

// Initializes the instance of HttpRequest with default values for all its fields. 
func (req *HttpRequest) initialize() {
	req.Body = make([]byte, 0)
	req.Headers = make(Headers)
	req.Version = getHighestVersion()
	req.staticFilePath = ""
	req.Query = nil
	req.Segments = nil
}

// Assigns the stream reader field of HttpRequest with a valid request stream.
func (req *HttpRequest) setReader(reader *bufio.Reader) {
	req.reader = reader
}

// Reads bytes of data from request byte stream and stores it in individual fields of HttpRequest instance.
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

// Reads the values for all request headers and stores them in the HttpRequest instance.
func (req *HttpRequest) readHeader() error {
	RequestLineProcessed := false
	HeaderProcessingCompleted := false

	for {
		message, err := req.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			
			reqError := new(RequestParseError)
			reqError.Section = "Header"
			reqError.Message = err.Error()
			reqError.Value = strings.TrimSpace(message)
			return reqError
		}

		message = strings.TrimSuffix(message, HEADER_LINE_SEPERATOR)
		if len(message) == 0 && !HeaderProcessingCompleted {
			HeaderProcessingCompleted = true
			break
		} else if !RequestLineProcessed {
			RequestLineParts := strings.Split(message, REQUEST_LINE_SEPERATOR)
			if len(RequestLineParts) != 3 {
				reqError := new(RequestParseError)
				reqError.Section = "Header"
				reqError.Message = "Request line should contain exactly three values seperated by a single whitespace"
				reqError.Value = strings.TrimSpace(message)
				return reqError
			}
			req.Method = strings.TrimSpace(RequestLineParts[0])
			req.ResourcePath = strings.TrimSpace(RequestLineParts[1])
			tempVersion := strings.TrimSpace(RequestLineParts[2])
			tempVersion, found := strings.CutPrefix(tempVersion, "HTTP/")
			if !found {
				reqError := new(RequestParseError)
				reqError.Section = "Header"
				reqError.Value = strings.TrimSpace(tempVersion)
				reqError.Message = "Invalid HTTP Version found in header"
				return reqError
			}
			req.Version = strings.TrimSpace(tempVersion) 
			RequestLineProcessed = true
		} else {
			HeaderKey, HeaderValue, found := strings.Cut(message, HEADER_KEY_VALUE_SEPERATOR)
			if !found {
				reqError := new(RequestParseError)
				reqError.Section = "Header"
				reqError.Value = strings.TrimSpace(message)
				reqError.Message = "Invalid header string found among request headers"
				return reqError
			}

			HeaderKey = strings.TrimSpace(HeaderKey)
			HeaderValue = strings.TrimSpace(HeaderValue)
			req.AddHeader(HeaderKey, HeaderValue)
		}
	}

	return nil
}

// Reads the body from request byte stream and stores them in the HttpRequest instance.
func (req *HttpRequest) readBody() error {
	if req.ContentLength > 0 {
		req.Body = make([]byte, req.ContentLength)
		for index := 0; index < req.ContentLength; index++ {
			bodyByte, err := req.reader.ReadByte()
			if err != nil {
				reqError := new(RequestParseError)
				reqError.Section = "Body"
				reqError.Value = "Request Body"
				reqError.Message = err.Error()
				return reqError
			}
			req.Body[index] = bodyByte
		}
	}

	return nil
}

// Parses all the query paramaters from the request URL and stores in the HttpRequest instance.
func (req *HttpRequest) parseQueryParams() error {
	req.Query = make(Params)
	parsedUrl, err := url.Parse(req.ResourcePath)
	if err != nil {
		reqError := new(RequestParseError)
		reqError.Section = "QueryParams"
		reqError.Value = req.ResourcePath
		reqError.Message = err.Error()
		return reqError
	}

	queryParams := parsedUrl.Query()
	for paramName, paramValues := range queryParams {
		req.Query.Add(paramName, paramValues)
	}

	return nil
}

// Checks if the given HTTP GET request made is a CONDITIONAL GET request.
func (req *HttpRequest) isConditionalGet(CompleteFilePath string) bool {
	if !strings.EqualFold(req.Method, "GET") {
		return false
	}

	fileMediaType, exists := getContentType(CompleteFilePath)
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