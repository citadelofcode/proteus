package internal

import (
	"bufio"
	"fmt"
	"net/textproto"
	"slices"
	"strconv"
	"strings"
	"time"
	"io"
)

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
	// Complete contents of the response body as a stream of bytes.
	BodyBytes []byte
	// Streamed writer instance to write the response bytes to the network stream.
	writer *bufio.Writer
	// The server instance processing this response.
	Server *HttpServer
	// key-value pairs to hold variables available during the entire response lifecycle.
	Locals map[string]any
	// FileSystem instance to access the local file system.
	fs *FileSystem
}

// // Initializes the instance of HttpResponse with default values for all its fields.
func (res *HttpResponse) Initialize(version string, writer io.Writer) {
	version = strings.TrimSpace(version)
	if version == "" {
		res.Version = "0.9"
	} else {
		res.Version = version
	}
	res.Headers = make(Headers)
	res.Locals = make(map[string]any)
	res.addGeneralHeaders()
	res.addResponseHeaders()
	res.writer = bufio.NewWriter(writer)
	res.Server = nil
	res.fs = new(FileSystem)
}

// Sets the server field to the given server instance reference.
func (res *HttpResponse) SetServer(serverRef *HttpServer) {
	res.Server = serverRef
}

// Adds all the general HTTP headers to the HttpResponse instance.
// Headers are added only if the given HttpResponse object is not a test instance and the response version is not HTTP/0.9.
func (res *HttpResponse) addGeneralHeaders() {
	if !strings.EqualFold(res.Version, "0.9") {
		res.Headers.Add("Date", GetRfc1123Time())
	}
}

// Adds all the default response HTTP headers to the HttpResponse instance.
// Headers are added only if the given HttpResponse object is not a test instance and the response version is not HTTP/0.9.
func (res *HttpResponse) addResponseHeaders() {
	if !strings.EqualFold(res.Version, "0.9") {
		res.Headers.Add("Server", GetServerDefaults("server_name").(string))
	}
}

// Writes bytes of data to response byte stream from the HttpResponse instance.
func (res *HttpResponse) Write() error {
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
	if len(res.BodyBytes) > 0 {
		ContentType, exists := res.Headers.Get("Content-Type")
		if exists {
			ContentType = strings.TrimSpace(ContentType)
			ContentType = strings.ToLower(ContentType)
			if strings.HasPrefix(ContentType, "text") {
				_, err := res.writer.WriteString(string(res.BodyBytes))
				if err != nil {
					resErr := new(ResponseError)
					resErr.Section = "Body"
					resErr.Value = ContentType
					resErr.Message = fmt.Sprintf("Error while writing response body :: %s", err.Error())
					return resErr
				}
			} else {
				_, err := res.writer.Write(res.BodyBytes)
				if err != nil {
					resErr := new(ResponseError)
					resErr.Section = "Body"
					resErr.Value = ContentType
					resErr.Message = fmt.Sprintf("Error while writing response body :: %s", err.Error())
					return resErr
				}
			}
		} else {
			_, err := res.writer.Write(res.BodyBytes)
			if err != nil {
				resErr := new(ResponseError)
				resErr.Section = "Body"
				resErr.Value = "Body Write without Content Type"
				resErr.Message = fmt.Sprintf("Error while writing response body :: %s", err.Error())
				return resErr
			}
		}
	}

	return nil
}

// Adds a new key-value pair to the request headers collection.
func (res *HttpResponse) AddHeader(HeaderKey string, HeaderValue string) {
	if slices.Contains(DateHeaders, textproto.CanonicalMIMEHeaderKey(HeaderKey)) {
		isValid, _ := IsHttpDate(HeaderValue)
		if isValid {
			res.Headers.Add(HeaderKey, HeaderValue)
		} else {
			res.Server.Log(fmt.Sprintf("Error while adding header - [%s] :: Date string must conform to one of these formats - RFC1123 or ANSIC", HeaderKey), ERROR_LEVEL)
		}
	} else {
		res.Headers.Add(HeaderKey, HeaderValue)
	}
}

// Sets the status of the HTTP response instance.
func (res *HttpResponse) Status(status StatusCode) {
	res.StatusCode = int(status)
	res.StatusMessage = status.GetStatusMessage()
}

// Send the given file from the local file system as the HTTP response.
func (res *HttpResponse) SendFile(CompleteFilePath string, OnlyMetadata bool) error {
	file, err := res.fs.GetFile(CompleteFilePath)
	if err != nil {
		return err
	}

	res.Headers.Add("Content-Length", strconv.FormatInt(file.Size(), 10))
	res.Headers.Add("Last-Modified", file.LastModified().Format(time.RFC1123))
	_, ok := res.Headers.Get("Content-Type")
	if !ok {
		res.Headers.Add("Content-Type", file.MediaType())
	}

	if !OnlyMetadata {
		contents, err := file.Contents()
		if err != nil {
			return err
		}
		res.BodyBytes = contents
	}

	return res.Write()
}

// Sends a the given error content as response back to the client.
func (res *HttpResponse) SendError(Content string) error {
	responseContent := []byte(Content)
	res.Headers.Add("Content-Type", ERROR_MSG_CONTENT_TYPE)
	res.Headers.Add("Content-Length", strconv.Itoa(len(responseContent)))
	res.BodyBytes = responseContent
	return res.Write()
}

// Send the given string as response back to the client.
func (res *HttpResponse) Send(content string) error {
	_, ok := res.Headers.Get("Content-Type")
	if !ok {
		fileMediaType := GetServerDefaults("content_type").(string)
		res.Headers.Add("Content-Type", fileMediaType)
	}

	content = strings.TrimSpace(content)
	contentBuffer := []byte(content)
	res.Headers.Add("Content-Length", strconv.Itoa(len(contentBuffer)))
	res.BodyBytes = contentBuffer
	return res.Write()
}
