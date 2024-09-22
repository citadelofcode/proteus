package http

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type HttpResponse struct {
	StatusCode int
	StatusMessage string
	Version string
	Headers Headers
	Body []byte
	writer *bufio.Writer
	ContentType string
}

func (res *HttpResponse) Initialize() {
	res.SetVersion(MAX_VERSION)
	res.Headers = make(Headers)
	res.AddGeneralHeaders()
	res.AddResponseHeaders()
}

func (res *HttpResponse) setWriter(writer *bufio.Writer) {
	res.writer = writer
}

func (res *HttpResponse) SetVersion(version string) {
	res.Version = strings.TrimSpace(version)
}

func (res *HttpResponse) Set(status StatusCode, version string, contentType string, contents []byte) {
	res.SetStatus(status)
	if version != "" {
		res.SetVersion(version)
	}
	if contentType != "" {
		res.SetContentType(contentType)
	}
	if len(contents) > 0 {
		res.SetContents(contents)
	}
}

func (res *HttpResponse) SetStatus(status StatusCode) {
	res.StatusCode = int(status)
	res.StatusMessage = status.GetStatusMessage()
}

func (res *HttpResponse) SetContents(contents []byte) {
	res.Body = contents
	contentLength := strconv.Itoa(len(contents))
	res.Headers.Add("Content-Length", contentLength)
}

func (res *HttpResponse) SetContentType(ContentType string) {
	res.ContentType = ContentType
	res.Headers.Add("Content-Type", ContentType)
}

func (res *HttpResponse) AddGeneralHeaders() {
	res.Headers.Add("Date", getRfc1123Time())
}

func (res *HttpResponse) AddResponseHeaders() {
	res.Headers.Add("Server", SERVER_NAME)
}

func (res *HttpResponse) Write() error {
	err := res.writeStatusLine()
	if err != nil {
		return err
	}

	err = res.writeHeaders()
	if err != nil {
		return err
	}

	err = res.writeBody()
	if err != nil {
		return err
	}

	err = res.writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (res *HttpResponse) writeStatusLine() error {
	if res.writer == nil {
		return errors.New("error occurred while writing response status line: writer object not initialized")
	}

	if res.StatusCode == 0 {
		return errors.New("error occurred while writing response status line: status code cannot be zero")
	}

	if res.Version == "" {
		return errors.New("error occurred while writing response status line: Protocol version not set")
	}

	_, err := res.writer.WriteString(fmt.Sprintf("HTTP/%s %d %s%s", res.Version, res.StatusCode, res.StatusMessage, HEADER_LINE_SEPERATOR))
	if err != nil {
		return errors.New("error occurred while writing response status line: " + err.Error())
	}

	return nil
}

func (res *HttpResponse) writeHeaders() error {
	if res.writer == nil {
		return errors.New("error occurred while writing response headers: writer object not initialized")
	}

	for key, values := range res.Headers {
		value := strings.Join(values, ",")
		_, err := res.writer.WriteString(fmt.Sprintf("%s: %s%s", key, value, HEADER_LINE_SEPERATOR))
		if err != nil {
			return errors.New("error occurred while writing response headers: " + err.Error())
		}
	}
	res.writer.WriteString(HEADER_LINE_SEPERATOR)

	return nil
}

func (res *HttpResponse) writeBody() error {
	if res.writer == nil {
		return errors.New("error occurred while writing response body: writer object not initialized")
	}

	if len(res.Body) > 0 {
		if strings.HasPrefix(res.ContentType, "text") {
			_, err := res.writer.WriteString(string(res.Body))
			if err != nil {
				return errors.New("error occurred while writing response body: " + err.Error())
			}
		} else {
			_, err := res.writer.Write(res.Body)
			if err != nil {
				return errors.New("error occurred while writing response body: " + err.Error())
			}
		}
	}

	return nil
}