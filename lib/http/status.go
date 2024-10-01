package http

type StatusCode int

const (
	StatusOK StatusCode = 200
	StatusCreated StatusCode = 201
	StatusAccepted StatusCode = 202
	StatusNonAuthoritative StatusCode = 203
	StatusNoContent StatusCode = 204
	StatusMultipleChoices StatusCode = 300
	StatusMovedPermanently StatusCode = 301
	StatusMovedTemporarily StatusCode = 302
	StatusSeeOther StatusCode = 303
	StatusNotModified StatusCode = 304
	StatusBadRequest StatusCode = 400
	StatusUnauthorized StatusCode = 401
	StatusPaymentRequired StatusCode = 402
	StatusForbidden StatusCode = 403
	StatusNotFound StatusCode = 404
	StatusMethodNotAllowed StatusCode = 405
	StatusNoneAcceptable StatusCode = 406
	StatusProxyAuth StatusCode = 407
	StatusRequestTimeout StatusCode = 408
	StatusConflict StatusCode = 409
	StatusGone StatusCode = 410
	StatusLengthMissing StatusCode = 411
	StatusInternalServerError StatusCode = 500
	StatusNotImplemented StatusCode = 501
	StatusBadGateway StatusCode = 502
	StatusServiceUnavailable StatusCode = 503
	StatusGatewayTimeout StatusCode = 504
)

// Gets the minified message assosciated with a HTTP status code.
func (code StatusCode) GetStatusMessage() string {
	switch code {
		case StatusOK: 
			return "OK"
		case StatusCreated:
			return "Created"
		case StatusAccepted:
			return "Accepted"
		case StatusNonAuthoritative:
			return "Non-Authoritative Information"
		case StatusNoContent:
			return "No Content"
		case StatusMultipleChoices:
			return "Multiple Choices"
		case StatusMovedPermanently:
			return "Moved Permanently"
		case StatusMovedTemporarily:
			return "Moved Temporarily"
		case StatusSeeOther:
			return "See Other"
		case StatusNotModified:
			return "Not Modified"
		case StatusBadRequest:
			return "Bad Request"
		case StatusUnauthorized:
			return "Unauthorized"
		case StatusPaymentRequired:
			return "Payment Required"
		case StatusForbidden:
			return "Forbidden"
		case StatusNotFound:
			return "Not Found"
		case StatusMethodNotAllowed:
			return "Method Not Allowed"
		case StatusNoneAcceptable:
			return "None Acceptable"
		case StatusProxyAuth:
			return "Proxy Authentication Required"
		case StatusRequestTimeout:
			return "Request Timeout"
		case StatusConflict:
			return "Conflict"
		case StatusGone:
			return "Gone"
		case StatusLengthMissing:
			return "Length Required"
		case StatusInternalServerError:
			return "Internal Server Error"
		case StatusNotImplemented:
			return "Not Implemented"
		case StatusBadGateway:
			return "Bad Gateway"
		case StatusServiceUnavailable:
			return "Service Unavailable"
		case StatusGatewayTimeout:
			return "Gateway Timeout"
		default:
			return ""
	}
}

// Gets the default error content for a HTTP status code.
func (code StatusCode) GetErrorContent() string {
	responseString := ""
	switch code {
	case StatusBadRequest:
		responseString = `<html>
					<head>
					<title>400 Bad Request<\title>
					</head>
					<body>
					<h1>400 - Bad Request</h1>
					<p>HTTP request made was not valid. Please check the request and send again.</p>
					</body>
				</html>`
	case StatusNotFound:
		responseString = `<html>
					<head>
					<title>404 Resource Not Found</title>
					</head>
					<body>
					<h1>401 - Resource Not Found</h1>
					<p>The resource you requested has not been found at the specified address. Please check the spelling of the address.</p>
					</body>
				</html>`
	case StatusMethodNotAllowed:
		responseString = `<html>
					<head>
					<title>405 Operation Not Permitted</title>
					</head>
					<body>
					<h1>405 - Ooperation Not Permitted</h1>
					<p>The operation requested is not allowed.</p>
					</body>
				</html>`
	case StatusNoneAcceptable:
		responseString = `<html>
					<head>
					<title>406 - Not Acceptable</title>
					</head>
					<body>
					<h1>405 - Not Acceptable</h1>
					<p>Resource with the requested filters are not available.</p>
					</body>
				</html>`
	case StatusRequestTimeout:
		responseString = `<html>
					<head>
					<title>408 Request Timeout</title>
					</head>
					<body>
					<h1>408 - Request Timeout</h1>
					<p>The transmission was not received quickly enough. Check internet connectivity and please try again.</p>
					</body>
				</html>`
	case StatusConflict:
		responseString = `<html>
					<head>
					<title>409 Resource Conflict</title>
					</head>
					<body>
					<h1>409 - Resource Conflict</h1>
					<p>This resource has been marked read-only. Please try changing the status and modify again.</p>
					</body>
				</html>`
	case StatusGone:
		responseString = `<html>
					<head>
					<title>410 - Resource Expired</title>
					</head>
					<body>
					<h1>410 - Resource Expired</h1>
					<p>The resource requested has expired and is no longer relevant.</p>
					</body>
				</html>`
	case StatusLengthMissing:
		responseString = `<html>
					<head>
					<title>411 - Length Missing</title>
					</head>
					<body>
					<h1>411 - Length Missing</h1>
					<p>The request sent does not contain 'Content-Length' value.</p>
					</body>
				</html>`
	case StatusInternalServerError:
		responseString = `<html>
					<head>
					<title>500 - Request Failed<\title>
					</head>
					<body>
					<h1>500 - Internal Server Error</h1>
					<p>Your request cannot be completed due to a server error.</p>
					</body>
				</html>`
	case StatusNotImplemented:
		responseString = `<html>
					<head>
					<title>501 - Function Not Implemented<\title>
					</head>
					<body>
					<h1>501 - Not Implemented</h1>
					<p>Your request can not be completed because this functionality is currently under development.</p>
					</body>
				</html>`
	case StatusBadGateway:
		responseString = `<html>
					<head>
					<title>502 - Bad Gateway<\title>
					</head>
					<body>
					<h1>502 - Bad Gateway</h1>
					<p>The server is unreachable at this time.</p>
					</body>
				</html>`
	case StatusServiceUnavailable:
		responseString = `<html>
					<head>
					<title>503 - Resource Busy<\title>
					</head>
					<body>
					<h1>503 - Resource Busy</h1>
					<p>Your request cannot be completed at this time. Please try again after 30 minutes.</p>
					</body>
				</html>`
	case StatusGatewayTimeout:
		responseString = `<html>
					<head>
					<title>504 - Gateway Timeout<\title>
					</head>
					<body>
					<h1>504 - Gateway Timeout</h1>
					<p>The server is not responding.</p>
					</body>
				</html>`
	default:
		responseString = ""
	}

	return responseString
}