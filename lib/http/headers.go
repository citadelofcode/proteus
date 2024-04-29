package http

import (
	"net/textproto"
	"strings"
)

type Headers map[string][]string

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

func (headers Headers) Get(key string) (string, bool) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	valueParts, ok := headers[key]
	if ok {
		return strings.Join(valueParts, ","), true
	} else {
		return "", false
	}
}

func (headers Headers) Contains(key string) bool {
	key = textproto.CanonicalMIMEHeaderKey(key)
	_, ok := headers[key]
	return ok
}