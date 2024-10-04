package http

import (
	"net/textproto"
	"strings"
)

// Represents a collection of headers (request or response).
type Headers map[string][]string

// Add a new key-value pair to the collection of headers.
func (headers Headers) Add(key string, value string) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	valueParts := strings.Split(value, ",")
	_, ok := headers[key]
	if ok {
		headers[key] = append(headers[key], valueParts...)
	} else {
		headers[key] = valueParts
	}
}

// Gets the value for a given header key from the collection of headers. The function also returns a boolean value to indicate if the key was found in the collection.
func (headers Headers) Get(key string) (string, bool) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	valueParts, ok := headers[key]
	if ok {
		return strings.Join(valueParts, ","), true
	} else {
		return "", false
	}
}