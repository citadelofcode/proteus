package http

import (
	"bufio"
	"fmt"
	"net/textproto"
	"slices"
	"strconv"
	"strings"
	"time"
	"github.com/mkbworks/proteus/lib/fs"
)

// Structure to represent a response status code and its associated information.
type respStatus struct {
	// HTTP response status code.
	Code StatusCode
	// Short message for the corresponding status code.
	Message string
	// Error description for error status codes (>=400).
	ErrorDescription string
}

// Structure to represent a HTTP response sent back by the server to the client.
type HttpResponse struct {
	// Status code of the response being sent back to the client like 200, 203, 404 etc.
	StatusCode int
	// Status message associated with the response status code.
	StatusMessage string
	// HTTP version of the response being sent back.
	Version string
	// Collection of all response headers being sent by the server.
	Headers Headers
	// Complete contents of the response body.
	Body []byte
	// Streamed writer instance to write the response bytes to the network stream.
	writer *bufio.Writer
	// Boolean value to indicate if the response created is a test object.
	isTest bool
}

// // Initializes the instance of HttpResponse with default values for all its fields.
func (res *HttpResponse) initialize(version string, isTest bool) {
	version = strings.TrimSpace(version)
	if version == "" {
		res.Version = "0.9"
	} else {
		res.Version = version
	}
	res.isTest = isTest
	res.Headers = make(Headers)
	res.addGeneralHeaders()
	res.addResponseHeaders()
}

// // Assigns the stream writer field of HttpResponse with a valid response stream.
func (res *HttpResponse) setWriter(writer *bufio.Writer) {
	res.writer = writer
}

// Adds all the general HTTP headers to the HttpResponse instance.
// Headers are added only if the given HttpResponse object is not a test instance and the response version is not HTTP/0.9.
func (res *HttpResponse) addGeneralHeaders() {
	if !strings.EqualFold(res.Version, "0.9") && !res.isTest {
		res.Headers.Add("Date", getRfc1123Time())
	}
}

// Adds all the default response HTTP headers to the HttpResponse instance. 
// Headers are added only if the given HttpResponse object is not a test instance and the response version is not HTTP/0.9.
func (res *HttpResponse) addResponseHeaders() {
	if !strings.EqualFold(res.Version, "0.9") && !res.isTest {
		res.Headers.Add("Server", getServerDefaults("server_name"))
	}
}

// Writes bytes of data to response byte stream from the HttpResponse instance.
func (res *HttpResponse) write() error {
	if res.writer == nil {
		resErr := new(ResponseError)
		resErr.Section = "RespWrite"
		resErr.Value = ""
		resErr.Message = "Writer object not initialized"
		return resErr
	}

	var err error
	if !strings.EqualFold(res.Version, "0.9") {
		err = res.writeStatusLine()
		if err != nil {
			return err
		}
	}

	if !strings.EqualFold(res.Version, "0.9") {
		err = res.writeHeaders()
		if err != nil {
			return err
		}
	}

	err = res.writeBody()
	if err != nil {
		return err
	}

	err = res.writer.Flush()
	if err != nil {
		resErr := new(ResponseError)
		resErr.Section = "RespWrite"
		resErr.Value = ""
		resErr.Message = "Writer object could not be flushed"
		return resErr
	}

	return nil
}

// Writes the HTTP response status line to the response byte stream.
func (res *HttpResponse) writeStatusLine() error {
	if res.StatusCode == 0 {
		resErr := new(ResponseError)
		resErr.Section = "StatusLine"
		resErr.Value = ""
		resErr.Message = "Status code for the response cannot be zero"
		return resErr
	}

	if res.Version == "" {
		resErr := new(ResponseError)
		resErr.Section = "StatusLine"
		resErr.Value = ""
		resErr.Message = "Response version cannot be empty"
		return resErr
	}

	_, err := res.writer.WriteString(fmt.Sprintf("HTTP/%s %d %s%s", res.Version, res.StatusCode, res.StatusMessage, HEADER_LINE_SEPERATOR))
	if err != nil {
		resErr := new(ResponseError)
		resErr.Section = "StatusLine"
		resErr.Value = ""
		resErr.Message = fmt.Sprintf("Error while writing response status line :: %s", err.Error())
		return resErr
	}

	return nil
}

// Writes the HTTP response headers to the response byte stream.
func (res *HttpResponse) writeHeaders() error {
	for key, values := range res.Headers {
		value := strings.Join(values, ",")
		_, err := res.writer.WriteString(fmt.Sprintf("%s: %s%s", key, value, HEADER_LINE_SEPERATOR))
		if err != nil {
			resErr := new(ResponseError)
			resErr.Section = "Header"
			resErr.Value = fmt.Sprintf("%s: %s", key, value)
			resErr.Message = fmt.Sprintf("Error while writing response header :: %s", err.Error())
			return resErr
		}
	}

	_, err := res.writer.WriteString(HEADER_LINE_SEPERATOR)
	if err != nil {
		resErr := new(ResponseError)
		resErr.Section = "Header"
		resErr.Value = HEADER_LINE_SEPERATOR
		resErr.Message = fmt.Sprintf("Error while writing response header :: %s", err.Error())
		return resErr
	}

	return nil
}

// Writes the response body to the response byte stream.
func (res *HttpResponse) writeBody() error {
	if len(res.Body) > 0 {
		ContentType, exists := res.Headers.Get("Content-Type")
		if exists {
			ContentType = strings.TrimSpace(ContentType)
			ContentType = strings.ToLower(ContentType)
			if strings.HasPrefix(ContentType, "text") {
				_, err := res.writer.WriteString(string(res.Body))
				if err != nil {
					resErr := new(ResponseError)
					resErr.Section = "Body"
					resErr.Value = ContentType
					resErr.Message = fmt.Sprintf("Error while writing response body :: %s", err.Error())
					return resErr
				}
			} else {
				_, err := res.writer.Write(res.Body)
				if err != nil {
					resErr := new(ResponseError)
					resErr.Section = "Body"
					resErr.Value = ContentType
					resErr.Message = fmt.Sprintf("Error while writing response body :: %s", err.Error())
					return resErr
				}
			}
		} else {
			_, err := res.writer.Write(res.Body)
			if err != nil {
				resErr := new(ResponseError)
				resErr.Section = "Body"
				resErr.Value = "Content type Not Available"
				resErr.Message = fmt.Sprintf("Error while writing response body :: %s", err.Error())
				return resErr
			}
		}
	}

	return nil
}

// Adds a new key-value pair to the request headers collection.
func (res *HttpResponse) AddHeader(HeaderKey string, HeaderValue string) error {
	if slices.Contains(DateHeaders, textproto.CanonicalMIMEHeaderKey(HeaderKey)) {
		isValid, _ := isHttpDate(HeaderValue)
		if isValid {
			res.Headers.Add(HeaderKey, HeaderValue)
		} else {
			resErr := new(ResponseError)
			resErr.Section = "Header"
			resErr.Value = fmt.Sprintf("%s: %s", HeaderKey, HeaderValue)
			resErr.Message = "Date string must conform to one of these formats - RFC1123 or ANSIC"
			return resErr
		}
	} else {
		res.Headers.Add(HeaderKey, HeaderValue)
	}

	return nil
}

// Sets the status of the HTTP response instance.
func (res *HttpResponse) Status(status StatusCode) {
	res.StatusCode = int(status)
	res.StatusMessage = status.GetStatusMessage()
}

// Send the given file from the local file system as the HTTP response.
func (res *HttpResponse) SendFile(CompleteFilePath string, OnlyMetadata bool) error {
	fileMediaType, exists := getContentType(CompleteFilePath)
	if exists {
		file, err := fs.GetFile(CompleteFilePath, fileMediaType, OnlyMetadata)
		if err == nil {
			res.Headers.Add("Content-Type", fileMediaType)
			res.Headers.Add("Content-Length", strconv.FormatInt(file.Size, 10))
			res.Headers.Add("Last-Modified", file.LastModifiedAt.Format(time.RFC1123))
			if !OnlyMetadata {
				res.Body = file.Contents
			}
			
			err := res.write()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Sends a the given error content as response back to the client.
func (res *HttpResponse) SendError(Content string) error {
	responseContent := []byte(Content)
	res.Headers.Add("Content-Type", ERROR_MSG_CONTENT_TYPE)
	res.Headers.Add("Content-Length", strconv.Itoa(len(responseContent)))
	res.Body = responseContent
	err := res.write()
	if err != nil {
		return err
	}

	return nil
}

// Send the given string as response back to the client.
func (res *HttpResponse) Send(content string) error {
	content = strings.TrimSpace(content)
	contentBuffer := []byte(content)
	res.Headers.Add("Content-Type", "text/plain")
	res.Headers.Add("Content-Length", strconv.Itoa(len(contentBuffer)))
	res.Body = contentBuffer
	err := res.write()
	return err
}