package http

const (
	ERROR_MSG_CONTENT_TYPE = "text/html"
	HEADER_LINE_SEPERATOR = "\r\n"
	REQUEST_LINE_SEPERATOR = " "
	HEADER_KEY_VALUE_SEPERATOR = ":"
	ROUTE_SEPERATOR = "/"

	// Informational data logged to the terminal.
	INFO_LEVEL = "INFO"
	// Error data logged to the terminal.
	ERROR_LEVEL = "ERROR"
	// Warning(s) logged to the terminal.
	WARN_LEVEL = "WARNING"

	// This is the default format in which request status is logged to the terminal.
	// This follows the common Apache log format which is as follows:
	//
	// :remote-addr [:date[clf]] ":method :url HTTP/:http-version" :status :res[content-length]
	COMMON_LOGGER = "common"
	// Concise output colored by response status for development use.
	// The :status token will be colored green for success codes, red for server error codes, yellow for client error codes, cyan for redirection codes, and uncolored for information codes.
	//
	// :method :url :status :response-time ms - :res[content-length]
	DEV_LOGGER = "dev"
	// The minimal output.
	//
	// :method :url :status :res[content-length] - :response-time ms
	TINY_LOGGER = "tiny"
	// Shorter than default, also including response time.
	//
	// :remote-addr :method :url HTTP/:http-version :status :res[content-length] - :response-time ms
	SHORT_LOGGER = "short"
)

// Collection of headers supported by the server that has a date value.
var DateHeaders []string
// List of content types supported by the web server.
var AllowedContentTypes map[string]string
// A map containing all the default server configuration values.
var ServerDefaults map[string]any
// List of all versions of HTTP supported by the web server.
var Versions map[string][]string
// List of response status codes and their associated information.
var ResponseStatusCodes []HttpStatus

// Initializes the global variables used in this package.
func init() {
	DateHeaders = []string{"Date", "Expires", "If-Modified-Since", "Last-Modified"}
	AllowedContentTypes = map[string]string{
		"pdf": "application/pdf",
        "htm": "text/html",
        "html": "text/html",
        "css": "text/css",
        "js": "text/javascript",
        "mjs": "text/javascript",
        "cjs": "text/javascript",
        "jpg": "image/jpeg",
        "jpeg": "image/jpeg",
        "gif": "image/gif",
        "png": "image/png",
        "apng": "image/apng",
        "svg": "image/svg+xml",
        "svgz": "image/svg+xml",
        "tif": "image/tiff",
        "tiff": "image/tiff",
        "aac": "audio/aac",
        "abw": "application/x-abiword",
        "arc": "application/x-freearc",
        "avif": "image/avif",
        "avi": "video/x-msvideo",
        "azw": "application/vnd.amazon.ebook",
        "bin": "application/octet-stream",
        "bmp": "image/bmp",
        "bz": "application/x-bzip",
        "bz2": "application/x-bzip2",
        "cda": "application/x-cdf",
        "csh": "application/x-csh",
        "csv": "text/csv",
        "doc": "application/msword",
        "docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
        "eot": "application/vnd.ms-fontobject",
        "epub": "application/epub+zip",
        "gz": "application/gzip",
        "ico": "image/vnd.microsoft.icon",
        "ics": "text/calendar",
        "jar": "application/java-archive",
        "json": "application/json",
        "jsonld": "application/ld+json",
        "odp": "application/vnd.oasis.opendocument.presentation",
        "odt": "application/vnd.oasis.opendocument.text",
        "ods": "application/vnd.oasis.opendocument.spreadsheet",
        "php": "application/x-httpd-php",
        "ppt": "application/vnd.ms-powerpoint",
        "pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
        "rtf": "application/rtf",
        "xhtml": "application/xhtml+xml",
        "xls": "application/vnd.ms-excel",
        "xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        "xml": "application/xml",
        "zip": "application/zip",
	}

	ServerDefaults = map[string]any {
		"hostname": "localhost",
		"port": 8080,
		"server_name": "proteus",
		"content_type": "application/octet-stream",
		"shutdown_timeout": 60,
	}

	Versions = map[string][]string {
		"0.9":  { "GET" },
		"1.0":  { "GET", "POST", "HEAD" },
		"1.1":  { "GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "OPTIONS", "CONNECT" },
	}

	ResponseStatusCodes = []HttpStatus {
		{ Code: StatusOK, Message: "OK", ErrorDescription: "" },
		{ Code: StatusCreated, Message: "Created", ErrorDescription: "" },
		{ Code: StatusAccepted, Message: "Accepted", ErrorDescription: "" },
		{ Code: StatusNonAuthoritative, Message: "Non-Authoritative Information", ErrorDescription: "" },
		{ Code: StatusNoContent, Message: "No Content", ErrorDescription: "" },
		{ Code: StatusResetContent, Message: "Reset Content", ErrorDescription: "" },
		{ Code: StatusPartialContent, Message: "Partial Content", ErrorDescription: "" },
		{ Code: StatusMultipleChoices, Message: "Multiple Choices", ErrorDescription: "" },
		{ Code: StatusMovedPermanently, Message: "Moved Permanently", ErrorDescription: "" },
		{ Code: StatusMovedTemporarily, Message: "Moved Temporarily", ErrorDescription: "" },
		{ Code: StatusSeeOther, Message: "See Other", ErrorDescription: "" },
		{ Code: StatusNotModified, Message: "Not Modified", ErrorDescription: "" },
		{ Code: StatusUseProxy, Message: "Use Proxy", ErrorDescription: "" },
		{ Code: StatusTemporaryRedirect, Message: "Temporary Redirect", ErrorDescription: "" },
		{ Code: StatusBadRequest, Message: "Bad Request", ErrorDescription: "HTTP request made was not valid. Please check the request and send again." },
		{ Code: StatusUnauthorized, Message: "Unauthorized", ErrorDescription: "" },
		{ Code: StatusPaymentRequired, Message: "Payment Required", ErrorDescription: "" },
		{ Code: StatusForbidden, Message: "Forbidden", ErrorDescription: "" },
		{ Code: StatusNotFound, Message: "Not Found", ErrorDescription: "The resource you requested has not been found at the specified address. Please check the spelling of the address." },
		{ Code: StatusMethodNotAllowed, Message: "Method Not Allowed", ErrorDescription: "The operation requested is not allowed." },
		{ Code: StatusNoneAcceptable, Message: "None Acceptable", ErrorDescription: "Resource with the requested filters are not available." },
		{ Code: StatusProxyAuth, Message: "Proxy Authentication Required", ErrorDescription: "" },
		{ Code: StatusRequestTimeout, Message: "Request Timeout", ErrorDescription: "The transmission was not received quickly enough. Check internet connectivity and please try again." },
		{ Code: StatusConflict, Message: "Conflict", ErrorDescription: "This resource has been marked read-only. Please try changing the status and modify again." },
		{ Code: StatusGone, Message: "Gone", ErrorDescription: "The resource requested has expired and is no longer relevant." },
		{ Code: StatusLengthMissing, Message: "Length Required", ErrorDescription: "The request sent does not contain 'Content-Length' value." },
		{ Code: StatusInternalServerError, Message: "Internal Server Error", ErrorDescription: "Your request cannot be completed due to a server error" },
		{ Code: StatusNotImplemented, Message: "Not Implemented", ErrorDescription: "Your request can not be completed because this functionality is currently under development." },
		{ Code: StatusBadGateway, Message: "Bad Gateway", ErrorDescription: "The server is unreachable at this time." },
		{ Code: StatusServiceUnavailable, Message: "Service Unavailable", ErrorDescription: "The server is currently unavailable." },
		{ Code: StatusGatewayTimeout, Message: "Gateway Timeout", ErrorDescription: "The server is not responding." },
		{ Code: StatusHTTPVersionNotSupported, Message: "HTTP Version Not Supported", ErrorDescription: "Received HTTP version is not supported by the server." },
	}
}
