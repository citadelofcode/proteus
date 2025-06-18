package internal

import (
	"strings"
	"encoding/json"
	"fmt"
	"net/url"
)

// Returns a middleware function that checks for requests with JSON payloads and parses them from source byte stream.
// If the request payload contains a JSON object, the request body will contain the parsed value as map[string]any instance.
// If the request payload contains a JSON array, the request body will contain the parsed value as []any{ map[string]any } instance.
func JsonParser() Middleware {
	return func (request *HttpRequest, response *HttpResponse, stop StopFunction) {
		reqContentType, ok := request.Headers.Get("Content-Type")
		if ok {
			reqContentType = strings.TrimSpace(reqContentType)
			if strings.EqualFold(reqContentType, "application/json") {
				sourceBytes := request.BodyBytes
				err := json.Unmarshal(sourceBytes, &request.Body)
				if err != nil {
					request.Server.Log(fmt.Sprintf("Error occurred while parsing JSON request payload: %s", err.Error()), ERROR_LEVEL)
					request.Server.Log(fmt.Sprintf("HTTP Request :: %s %s", request.Method, request.ResourcePath), ERROR_LEVEL)
					return
				}
			}
		}
	}
}

// Returns a middleware function that checks for request payloads of type "application/x-www-form-urlencoded" and fetches the parsed values. The parsed values are stored in an instance of type - map[string][]string.
func UrlEncoded() Middleware {
	return func (request *HttpRequest, response *HttpResponse, stop StopFunction) {
		reqContentType, ok := request.Headers.Get("Content-Type")
		if ok {
			reqContentType = strings.TrimSpace(reqContentType)
			if strings.EqualFold(reqContentType, "application/x-www-form-urlencoded") {
				payloadString := string(request.BodyBytes)
				payloadString = strings.TrimSpace(payloadString)
				parsedParams, err := url.ParseQuery(payloadString)
				if err != nil {
					request.Server.Log(fmt.Sprintf("Error occurred while parsing url encoded request payload: %s", err.Error()), ERROR_LEVEL)
					request.Server.Log(fmt.Sprintf("HTTP Request :: %s %s", request.Method, request.ResourcePath), ERROR_LEVEL)
					return
				}

				request.Body = map[string][]string(parsedParams)
			}
		}
	}
}
