package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"github.com/citadelofcode/proteus/lib/fs"
)

// Structure to represent a HTTP request received by the web server.
type HttpRequest struct {
	// HTTP request method like GET, POST, PUT etc.
	Method string
	// Resource path requested by the client.
	ResourcePath string
	// HTTP version that the request complies with. It is of format <major>.<minor> which refers to the major and minor versions respectively.
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
	// Collection of all query parameters stored as key-values pair.
	Query Params
	// Collection of all path parameter values stored as key-value pair.
	Segments Params
	// The IP address and port number of the client who made the request to the server
	ClientAddress string
}

// Initializes the instance of HttpRequest with default values for all its fields.
func (req *HttpRequest) initialize() {
	req.Body = make([]byte, 0)
	req.Headers = make(Headers)
	req.Version = "0.9"
	req.staticFilePath = ""
	req.Query = make(Params)
	req.Segments = make(Params)
}

// Assigns the stream reader field of HttpRequest with a valid request stream.
func (req *HttpRequest) setReader(reader *bufio.Reader) {
	req.reader = reader
}

// Reads bytes of data from request byte stream and stores it in individual fields of HttpRequest instance.
func (req *HttpRequest) read() error {
	err := req.readHeader()
	if err != nil {
		return err
	}

	err = req.parseQueryParams()
	if err != nil {
		return err
	}

	clength, ok := req.Headers.Get("Content-Length")
	if ok {
		req.ContentLength, err = strconv.Atoi(clength)
		if err != nil {
			reqError := new(RequestParseError)
			reqError.Section = "Header"
			reqError.Message = err.Error()
			reqError.Value = fmt.Sprintf("Content Length parsing error for value - %s", strings.TrimSpace(clength))
			return reqError
		}

		err = req.readBody()
		if err != nil {
			return err
		}
	}

	return nil
}

// Reads the values for all request headers and stores them in the HttpRequest instance.
func (req *HttpRequest) readHeader() error {
	RequestLineProcessed := false
	HeaderProcessingCompleted := false

	for {
		message, err := req.reader.ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return &ReadTimeoutError{}
			} else if len(message) == 0 && err != io.EOF {
				reqError := new(RequestParseError)
				reqError.Section = "Header"
				reqError.Message = err.Error()
				reqError.Value = strings.TrimSpace(message)
				return reqError
			} else if len(message) == 0 && err == io.EOF {
				return err
			}
		}

		message = strings.TrimSuffix(message, HEADER_LINE_SEPERATOR)
		if len(message) == 0 && !HeaderProcessingCompleted {
			HeaderProcessingCompleted = true
			break
		} else if !RequestLineProcessed {
			RequestLineParts := strings.Split(message, REQUEST_LINE_SEPERATOR)
			if len(RequestLineParts) != 2 && len(RequestLineParts) != 3 {
				reqError := new(RequestParseError)
				reqError.Section = "Header"
				reqError.Message = "Request line should contain either 2 or 3 values, seperated by a single whitespace"
				reqError.Value = strings.TrimSpace(message)
				return reqError
			}

			tempVersion := ""
			if len(RequestLineParts) == 2 || len(RequestLineParts) == 3 {
				req.Method = strings.TrimSpace(RequestLineParts[0])
				req.ResourcePath = strings.TrimSpace(RequestLineParts[1])
			}

			if len(RequestLineParts) == 2 {
				tempVersion = "HTTP/0.9"
			}

			if len(RequestLineParts) == 3 {
				tempVersion = strings.TrimSpace(RequestLineParts[2])
			}

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
			err := req.addHeader(HeaderKey, HeaderValue)
			if err != nil {
				return err
			}
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
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return &ReadTimeoutError{}
			} else if err != nil {
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
// Once the parsing is done, it removes the query parameters string from the Resource Path field.
func (req *HttpRequest) parseQueryParams() error {
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

	if len(req.Query) > 0 {
		req.ResourcePath, _, _ = strings.Cut(req.ResourcePath, "?")
	}

	return nil
}

// Checks if the given HTTP GET request made is a CONDITIONAL GET request.
func (req *HttpRequest) isConditionalGet(CompleteFilePath string) (bool, error) {
	if !strings.EqualFold(req.Method, "GET") {
		return false, nil
	}

	fileMediaType, err := getContentType(CompleteFilePath)
	if err != nil {
		return false, err
	}

	file, err := fs.GetFile(CompleteFilePath, fileMediaType, true)
	if err != nil {
		return false, err
	}

	LastModifiedString, ok := req.Headers.Get("If-Modified-Since")
	if !ok {
		return false, nil
	}

	LastModifiedString = strings.TrimSpace(LastModifiedString)
	isValid, LastModifiedSince := isHttpDate(LastModifiedString)
	if !isValid {
		reqError := new(RequestParseError)
		reqError.Section = "Header"
		reqError.Value = LastModifiedString
		reqError.Message = "The given datetime string value must either conform to ANSIC or RFC 1123 format"
		return false, reqError
	}

	if file.LastModifiedAt.After(LastModifiedSince) {
		return false, nil
	}

	return true, nil
}

// Adds a new key-value pair to the request headers collection.
func (req *HttpRequest) addHeader(HeaderKey string, HeaderValue string) error {
	if slices.Contains(DateHeaders, textproto.CanonicalMIMEHeaderKey(HeaderKey)) {
		isValid, _ := isHttpDate(HeaderValue)
		if isValid {
			req.Headers.Add(HeaderKey, HeaderValue)
		} else {
			reqError := new(RequestParseError)
			reqError.Section = "Header"
			reqError.Value = fmt.Sprintf("%s: %s", HeaderKey, HeaderValue)
			reqError.Message = "The given date header value should be one of either of these formats - RFC 1123 or ANSIC"
			return reqError
		}
	} else {
		req.Headers.Add(HeaderKey, HeaderValue)
	}

	return nil
}
