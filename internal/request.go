package internal

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
	"time"
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
	// Represents the complete contents of the request body as a stream of bytes.
	BodyBytes []byte
	// key-value pairs to hold variables available during the entire request lifecycle.
	Locals map[string]any
	// Streamed reader instance to read the HTTP request from the network stream.
	reader *bufio.Reader
	// Collection of all query parameters stored as key-values pair.
	Query Params
	// Collection of all path parameter values stored as key-value pair.
	Segments Params
	// The IP address and port number of the client who made the request to the server
	ClientAddress string
	// The server instance processing this request.
	Server *HttpServer
	// The actual content contained by the request. The type of the data is determined at run time depending on the data sent as part of the request.
	Body any
	// FileSystem instance to access the local file system.
	fs *FileSystem
}

// Initializes the instance of HttpRequest with default values for all its fields.
func (req *HttpRequest) Initialize(reader io.Reader) {
	req.BodyBytes = make([]byte, 0)
	req.Headers = make(Headers)
	req.Version = "0.9"
	req.Locals = make(map[string]any)
	req.Query = make(Params)
	req.Segments = make(Params)
	req.reader = bufio.NewReader(reader)
	req.Server = nil
	req.Locals["Started"] = time.Time{}
	req.Locals["ContentLength"] = 0
	req.fs = new(FileSystem)
}

// Reads bytes of data from request byte stream and stores it in individual fields of HttpRequest instance.
func (req *HttpRequest) Read() error {
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
		reqContentLength, err := strconv.Atoi(clength)
		if err != nil {
			reqError := new(RequestParseError)
			reqError.Section = "Header"
			reqError.Message = err.Error()
			reqError.Value = fmt.Sprintf("Content Length parsing error for value - %s", strings.TrimSpace(clength))
			return reqError
		}

		req.Locals["ContentLength"] = reqContentLength
		err = req.readBody()
		if err != nil {
			return err
		}
	}

	return nil
}

// Gets the time elapsed since request processing started (in milliseconds).
// If start time is not available, it returns zero.
func (req *HttpRequest) ProcessingTime() int64 {
	Started := req.Locals["Started"].(time.Time)
	if Started.IsZero() {
		return 0
	}
	timeSinceStart := time.Since(Started)
	return timeSinceStart.Milliseconds()
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
			req.AddHeader(HeaderKey, HeaderValue)
		}
	}

	return nil
}

// Reads the body from request byte stream and stores them in the HttpRequest instance.
func (req *HttpRequest) readBody() error {
	reqContentLength := req.Locals["ContentLength"].(int)
	if reqContentLength > 0 {
		req.BodyBytes = make([]byte, reqContentLength)
		for index := range reqContentLength {
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
			req.BodyBytes[index] = bodyByte
		}
	}

	return nil
}

// Parses all the query paramaters from the request URL and stores in the HttpRequest instance.
// Once the parsing is done, it removes the query parameters string from the Resource Path field.
func (req *HttpRequest) parseQueryParams() error {
	CleanedPath := CleanRoute(req.ResourcePath)
	parsedUrl, err := url.Parse(CleanedPath)
	if err != nil {
		reqError := new(RequestParseError)
		reqError.Section = "QueryParams"
		reqError.Value = CleanedPath
		reqError.Message = err.Error()
		return reqError
	}

	queryParams := parsedUrl.Query()
	for paramName, paramValues := range queryParams {
		req.Query.Add(paramName, paramValues)
	}

	if len(req.Query) > 0 {
		req.ResourcePath, _, _ = strings.Cut(CleanedPath, "?")
	}

	return nil
}

// Checks if the given HTTP GET request made is a CONDITIONAL GET request.
func (req *HttpRequest) IsConditionalGet(CompleteFilePath string) (bool, error) {
	if !strings.EqualFold(req.Method, "GET") {
		return false, nil
	}

	IfModifiedSinceString, ok := req.Headers.Get("If-Modified-Since")
	if !ok {
		return false, nil
	}

	IfModifiedSinceString = strings.TrimSpace(IfModifiedSinceString)
	isValid, IfModifiedSince := IsHttpDate(IfModifiedSinceString)
	if !isValid {
		reqError := new(RequestParseError)
		reqError.Section = "Header"
		reqError.Value = IfModifiedSinceString
		reqError.Message = "The given datetime string value must either conform to ANSIC, RFC 1123 or  RFC 850 format"
		return false, reqError
	}

	file, err := req.fs.GetFile(CompleteFilePath)
	if err != nil {
		return false, err
	}

	LastModifiedAt := file.LastModified()
	if LastModifiedAt.After(IfModifiedSince) {
		return false, nil
	}

	return true, nil
}

// Adds a new key-value pair to the request headers collection.
func (req *HttpRequest) AddHeader(HeaderKey string, HeaderValue string) {
	if slices.Contains(DateHeaders, textproto.CanonicalMIMEHeaderKey(HeaderKey)) {
		isValid, _ := IsHttpDate(HeaderValue)
		if isValid {
			req.Headers.Add(HeaderKey, HeaderValue)
		} else {
			req.Server.Log(fmt.Sprintf("Error while adding header - [%s] :: Date string must conform to one of these formats - RFC1123 or ANSIC", HeaderKey), ERROR_LEVEL)
		}
	} else {
		req.Headers.Add(HeaderKey, HeaderValue)
	}
}
