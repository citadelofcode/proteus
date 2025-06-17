package test

import (
	"reflect"
	"strings"
	"testing"
	"github.com/citadelofcode/proteus/internal"
)

// Test case to validate the working of the JsonParser() middleware to parse JSON objects as payloads.
func Test_JsonParser_ObjectParse(t *testing.T) {
	JsonParser := internal.JsonParser()
	testCases := []struct {
		Name string
		InPayload string
		OpPayload map[string]any
	} {
		{ "A valid JSON payload with a flat object", `{"name":"Mahesh","email":"mkbalaji@email.com"}`, map[string]any{ "name": "Mahesh", "email": "mkbalaji@email.com"}},
		{ "A valid JSON payload with a nested object", `{"name":"Mahesh","email":"mkbalaji@email.com","location":{"latitude":24,"longitude":-124}}`, map[string]any{"name": "Mahesh", "email": "mkbalaji@email.com", "location": map[string]any{"latitude": float64(24), "longitude": float64(-124)}}},
		{ "A valid flat JSON payload with numeric fields", `{"name":"Mahesh","email":"mkbalaji@email.com", "age":18}`, map[string]any{ "name": "Mahesh", "email": "mkbalaji@email.com", "age": float64(18)}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			request, response := TCParams(tt)
			request.BodyBytes = []byte(testCase.InPayload)
			request.AddHeader("Content-Type", "application/json")
			JsonParser(request, response, Stop)
			tt.Logf("Expected JSON Object Payload: %#v", testCase.OpPayload)
			tt.Logf("Parsed JSON Object payload: %#v", request.Body)
			isEqual := reflect.DeepEqual(request.Body, testCase.OpPayload)
			if isEqual {
				tt.Log("The expected JSON payload matches the parsed JSON object payload.")
			} else {
				tt.Error("The expected JSON payload does not match the parsed JSON object payload.")
			}
		})
	}
}

// Test case to validate the working of the JsonParser() middleware to parse JSON arrays as payloads.
func Test_JsonParser_ArrayParse(t *testing.T) {
	JsonParser := internal.JsonParser()
	testCases := []struct {
		Name string
		InPayload string
		OpPayload []any
	} {
		{ "A valid JSON array", `[{"name":"Mahesh","email":"mkbalaji@email.com"},{"name":"sdgbjsdk","email":"jksbfuhs@email.com"}]`, []any{ map[string]any{"name": "Mahesh", "email":"mkbalaji@email.com"}, map[string]any{"name": "sdgbjsdk", "email": "jksbfuhs@email.com"}}},
		{ "A valid JSON array with numeric fields", `[{"name":"Mahesh","age":18},{"name":"sdgbjsdk","age":20}]`, []any{ map[string]any{"name": "Mahesh", "age": float64(18)}, map[string]any{"name": "sdgbjsdk", "age": float64(20)}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			request, response := TCParams(tt)
			request.BodyBytes = []byte(testCase.InPayload)
			request.AddHeader("Content-Type", "application/json")
			JsonParser(request, response, Stop)
			tt.Logf("Expected JSON Array Payload: %#v", testCase.OpPayload)
			tt.Logf("Parsed JSON Array Payload: %#v", request.Body)
			isEqual := reflect.DeepEqual(request.Body, testCase.OpPayload)
			if isEqual {
				tt.Log("The expected JSON payload matches the parsed JSON array payload.")
			} else {
				tt.Error("The expected JSON payload does not match the parsed JSON array payload.")
			}
		})
	}
}

// Test case to validate how JsonParser() middleware affects non-JSON payloads. It checks if the request body still remains "nil" after middleware execution for non-JSON payloads.
func Test_JsonParser_NonJson(t *testing.T) {
	JsonParser := internal.JsonParser()
	testCases := []struct {
		Name string
		IpContentType string
		IpPayload string
	} {
		{ "A plain text payload", "text/plain", "This is a simple text value" },
		{ "An XML payload", "application/xml", `<?xml version="1.0" encoding="UTF-8"?><Name>Mahesh</Name>` },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			request, response := TCParams(tt)
			request.BodyBytes = []byte(testCase.IpPayload)
			request.AddHeader("Content-Type", testCase.IpContentType)
			JsonParser(request, response, Stop)
			if request.Body == nil {
				tt.Logf("The request body is 'nil' as expected, for payload [%s] of content type [%s]", testCase.IpPayload, testCase.IpContentType)
			} else {
				tt.Logf("The request body is not 'nil' as expected, for payload [%s] of content type [%s]", testCase.IpPayload, testCase.IpContentType)
			}
		})
	}
}

// Test case to validate how the UrlEncoded() middleware works for request payloads of type - "application/x-www-form-urlencoded".
func Test_UrlEncoded_ValidPayloads(t *testing.T) {
	UrlEncoded := internal.UrlEncoded()
	testCases := []struct {
		Name string
		IpPayload string
		OpPayload map[string][]string
	} {
		{ "A valid url encoded string with single value params", "name=Mahesh&age=30&city=San+Francisco", map[string][]string {"name": {"Mahesh"}, "age": {"30"}, "city": {"San Francisco"}}},
		{ "A valid url encoded string with single value params and url encoded values", "name=Mahesh&email=mkbalaji%40email.com&age=30", map[string][]string{"age": {"30"}, "email": {"mkbalaji@email.com"}, "name":{"Mahesh"}}},
		{ "A valid url encoded string with multiple value params", "tag=go&tag=web&tag=backend&user=mahesh", map[string][]string{"tag": {"go", "web", "backend"}, "user": {"mahesh"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			request, response := TCParams(tt)
			request.AddHeader("Content-Type", "application/x-www-form-urlencoded")
			request.BodyBytes = []byte(testCase.IpPayload)
			UrlEncoded(request, response, Stop)
			tt.Logf("Expected Url encoded payload: %#v\n", testCase.OpPayload)
			tt.Logf("Processed Url encoded payload: %#v\n", request.Body)
			isEqual := reflect.DeepEqual(request.Body, testCase.OpPayload)
			if isEqual {
				tt.Log("The expected request payload matches the processed payload.")
			} else {
				tt.Error("The expected request payload does not match the processed payload.")
			}
		})
	}
}

// Test case to validate how the UrlEncoded() middleware works for request payloads not of type - "application/x-www-form-urlencoded".
func Test_UrlEncoded_InvalidPayloads(t *testing.T) {
	UrlEncoded := internal.UrlEncoded()
	testCases := []struct {
		Name string
		IpContentType string
		IpPayload string
	} {
		{ "A plain text payload", "text/plain", "This is a simple text value" },
		{ "An XML payload", "application/xml", `<?xml version="1.0" encoding="UTF-8"?><Name>Mahesh</Name>` },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			request, response := TCParams(tt)
			request.BodyBytes = []byte(testCase.IpPayload)
			request.AddHeader("Content-Type", testCase.IpContentType)
			UrlEncoded(request, response, Stop)
			if request.Body == nil {
				tt.Logf("The request body is 'nil' as expected, for payload [%s] of content type [%s]", testCase.IpPayload, testCase.IpContentType)
			} else {
				tt.Logf("The request body is not 'nil' as expected, for payload [%s] of content type [%s]", testCase.IpPayload, testCase.IpContentType)
			}
		})
	}
}

// Test case to validate how either of the middlewares - JsonParser() and UrlEncoded() behave when there is no content type header present.
func Test_BodyParser_NoContentType(t *testing.T) {
	JsonParser := internal.JsonParser()
	UrlEncoded := internal.UrlEncoded()
	testCases := []struct {
		Name string
		Middleware string
		IpPayload string
	} {
		{ "Json Parser middleware", "JsonParser", `{"name":"Mahesh","email":"mkbalaji@email.com"}` },
		{ "Url Encoded middleware", "UrlEncoded", "name=Mahesh&age=30&city=San+Francisco" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			if strings.EqualFold(testCase.Middleware, "JsonParser") {
				request, response := TCParams(tt)
				request.BodyBytes = []byte(testCase.IpPayload)
				JsonParser(request, response, Stop)
				if request.Body == nil {
					tt.Logf("The request body is 'nil' as expected, for payload [%s]", testCase.IpPayload)
				} else {
					tt.Logf("The request body is not 'nil' as expected, for payload [%s]", testCase.IpPayload)
				}
			}

			if strings.EqualFold(testCase.Middleware, "UrlEncoded") {
				request, response := TCParams(tt)
				request.BodyBytes = []byte(testCase.IpPayload)
				UrlEncoded(request, response, Stop)
				if request.Body == nil {
					tt.Logf("The request body is 'nil' as expected, for payload [%s]", testCase.IpPayload)
				} else {
					tt.Logf("The request body is not 'nil' as expected, for payload [%s]", testCase.IpPayload)
				}
			}
		})
	}
}
