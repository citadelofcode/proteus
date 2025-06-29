package internal

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
// List of all versions of HTTP supported by the web server and the HTTP methods supported for each version.
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
		"txt": "text/plain",
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
		"1.0":  { "GET", "POST", "HEAD", "OPTIONS", "TRACE" },
		"1.1":  { "GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "OPTIONS", "CONNECT", "PATCH" },
	}

	ResponseStatusCodes = []HttpStatus {
		{ Code: Status100, Message: "Continue", ErrorDescription: "" },
		{ Code: Status101, Message: "Switching Protocols", ErrorDescription: "" },
		{ Code: Status102, Message: "Processing", ErrorDescription: "" },
		{ Code: Status103, Message: "Early Hints", ErrorDescription: "" },

		{ Code: Status200, Message: "OK", ErrorDescription: "" },
		{ Code: Status201, Message: "Created", ErrorDescription: "" },
		{ Code: Status202, Message: "Accepted", ErrorDescription: "" },
		{ Code: Status203, Message: "Non-Authoritative Information", ErrorDescription: "" },
		{ Code: Status204, Message: "No Content", ErrorDescription: "" },
		{ Code: Status205, Message: "Reset Content", ErrorDescription: "" },
		{ Code: Status206, Message: "Partial Content", ErrorDescription: "" },
		{ Code: Status207, Message: "Multi Status", ErrorDescription: "" },
		{ Code: Status208, Message: "Already Reported", ErrorDescription: "" },

		{ Code: Status300, Message: "Multiple Choices", ErrorDescription: "" },
		{ Code: Status301, Message: "Moved Permanently", ErrorDescription: "" },
		{ Code: Status302, Message: "Found", ErrorDescription: "" },
		{ Code: Status303, Message: "See Other", ErrorDescription: "" },
		{ Code: Status304, Message: "Not Modified", ErrorDescription: "" },
		{ Code: Status305, Message: "Use Proxy", ErrorDescription: "" },
		{ Code: Status307, Message: "Temporary Redirect", ErrorDescription: "" },
		{ Code: Status308, Message: "Permanent Redirect", ErrorDescription: "" },

		{ Code: Status400, Message: "Bad Request", ErrorDescription: "The server cannot or will not process the request due to an apparent client error (e.g., malformed request syntax, size too large, invalid request message framing, or deceptive request routing)." },
		{ Code: Status401, Message: "Unauthorized", ErrorDescription: "Similar to 403 Forbidden, but specifically for use when authentication is required and has failed or has not yet been provided. The response must include a WWW-Authenticate header field containing a challenge applicable to the requested resource." },
		{ Code: Status402, Message: "Payment Required", ErrorDescription: "Reserved for future use. The original intention was that this code might be used as part of some form of digital cash or micropayment scheme, as proposed, for example, by GNU Taler, but that has not yet happened, and this code is not widely used. Google Developers API uses this status if a particular developer has exceeded the daily limit on requests. Sipgate uses this code if an account does not have sufficient funds to start a call. Shopify uses this code when the store has not paid their fees and is temporarily disabled. Stripe uses this code for failed payments where parameters were correct, for example blocked fraudulent payments." },
		{ Code: Status403, Message: "Forbidden", ErrorDescription: "The request contained valid data and was understood by the server, but the server is refusing action. This may be due to the user not having the necessary permissions for a resource or needing an account of some sort, or attempting a prohibited action (e.g. creating a duplicate record where only one is allowed). This code is also typically used if the request provided authentication by answering the WWW-Authenticate header field challenge, but the server did not accept that authentication. The request should not be repeated." },
		{ Code: Status404, Message: "Not Found", ErrorDescription: "The requested resource could not be found but may be available in the future. Subsequent requests by the client are permissible." },
		{ Code: Status405, Message: "Method Not Allowed", ErrorDescription: "TA request method is not supported for the requested resource; for example, a GET request on a form that requires data to be presented via POST, or a PUT request on a read-only resource." },
		{ Code: Status406, Message: "None Acceptable", ErrorDescription: "The requested resource is capable of generating only content not acceptable according to the Accept headers sent in the request." },
		{ Code: Status407, Message: "Proxy Authentication Required", ErrorDescription: "The client must first authenticate itself with the proxy." },
		{ Code: Status408, Message: "Request Timeout", ErrorDescription: "The server timed out waiting for the request." },
		{ Code: Status409, Message: "Conflict", ErrorDescription: "Indicates that the request could not be processed because of conflict in the current state of the resource, such as an edit conflict between multiple simultaneous updates." },
		{ Code: Status410, Message: "Gone", ErrorDescription: "Indicates that the resource requested was previously in use but is no longer available and will not be available again. This should be used when a resource has been intentionally removed and the resource should be purged. Upon receiving a 410 status code, the client should not request the resource in the future. Clients such as search engines should remove the resource from their indices. Most use cases do not require clients and search engines to purge the resource, and a 404 Not Found may be used instead." },
		{ Code: Status411, Message: "Length Required", ErrorDescription: "The request did not specify the length of its content, which is required by the requested resource." },
		{ Code: Status412, Message: "Precondition Failed", ErrorDescription: "The server does not meet one of the preconditions that the requester put on the request header fields." },
		{ Code: Status413, Message: "Content Too Large", ErrorDescription: "The request is larger than the server is willing or able to process." },
		{ Code: Status414, Message: "URI Too Long", ErrorDescription: "The URI provided was too long for the server to process. Often the result of too much data being encoded as a query-string of a GET request, in which case it should be converted to a POST request. Called Request-URI Too Long previously." },
		{ Code: Status415, Message: "Unsupported Media Type", ErrorDescription: "The request entity has a media type which the server or resource does not support. For example, the client uploads an image as image/svg+xml, but the server requires that images use a different format." },
		{ Code: Status416, Message: "Range not satisfiable", ErrorDescription: "The client has asked for a portion of the file (byte serving), but the server cannot supply that portion. For example, if the client asked for a part of the file that lies beyond the end of the file. Called Requested Range Not Satisfiable previously." },
		{ Code: Status417, Message: "Expectation Failed", ErrorDescription: "The server cannot meet the requirements of the Expect request-header field." },
		{ Code: Status421, Message: "Misdirected Request", ErrorDescription: "The request was directed at a server that is not able to produce a response (for example because of connection reuse)." },
		{ Code: Status422, Message: "Unprocessable Content", ErrorDescription: "The request was well-formed (i.e., syntactically correct) but could not be processed." },
		{ Code: Status423, Message: "Locked", ErrorDescription: "The resource that is being accessed is locked." },
		{ Code: Status424, Message: "Failed Dependency", ErrorDescription: "The request failed because it depended on another request and that request failed." },
		{ Code: Status425, Message: "Too Early", ErrorDescription: "Indicates that the server is unwilling to risk processing a request that might be replayed." },
		{ Code: Status426, Message: "Upgrade Required", ErrorDescription: "The client should switch to a different protocol such as TLS/1.3, given in the Upgrade header field." },
		{ Code: Status428, Message: "Precondition Required", ErrorDescription: "The origin server requires the request to be conditional. Intended to prevent the 'lost update' problem, where a client GETs a resource's state, modifies it, and PUTs it back to the server, when meanwhile a third party has modified the state on the server, leading to a conflict." },
		{ Code: Status429, Message: "Too Many Requests", ErrorDescription: "The user has sent too many requests in a given amount of time. Intended for use with rate-limiting schemes." },
		{ Code: Status431, Message: "Request Header Fields Too Large", ErrorDescription: "The server is unwilling to process the request because either an individual header field, or all the header fields collectively, are too large." },

		{ Code: Status500, Message: "Internal Server Error", ErrorDescription: "A generic error message, given when an unexpected condition was encountered and no more specific message is suitable." },
		{ Code: Status501, Message: "Not Implemented", ErrorDescription: "The server either does not recognize the request method, or it lacks the ability to fulfil the request. Usually this implies future availability (e.g., a new feature of a web-service API)." },
		{ Code: Status502, Message: "Bad Gateway", ErrorDescription: "The server was acting as a gateway or proxy and received an invalid response from the upstream server." },
		{ Code: Status503, Message: "Service Unavailable", ErrorDescription: "The server cannot handle the request (because it is overloaded or down for maintenance). Generally, this is a temporary state." },
		{ Code: Status504, Message: "Gateway Timeout", ErrorDescription: "The server was acting as a gateway or proxy and did not receive a timely response from the upstream server." },
		{ Code: Status505, Message: "HTTP Version Not Supported", ErrorDescription: "The server does not support the HTTP version used in the request." },
		{ Code: Status506, Message: "Variant also Negotiates", ErrorDescription: "Transparent content negotiation for the request results in a circular reference." },
		{ Code: Status507, Message: "Insufficient Storage", ErrorDescription: "The server is unable to store the representation needed to complete the request." },
		{ Code: Status508, Message: "Loop Detected", ErrorDescription: "The server detected an infinite loop while processing the request." },
		{ Code: Status511, Message: "Network Authentication Required", ErrorDescription: "The client needs to authenticate to gain network access. Intended for use by intercepting proxies used to control access to the network." },
	}
}
