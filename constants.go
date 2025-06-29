package proteus

import (
	"github.com/citadelofcode/proteus/internal"
)

// List of HTTP response codes supported by the server.
const (
	Status100 StatusCode = internal.Status100
	Status101 StatusCode = internal.Status101
	Status102 StatusCode = internal.Status102
	Status103 StatusCode = internal.Status102

	Status200 StatusCode = internal.Status200
	Status201 StatusCode = internal.Status201
	Status202 StatusCode = internal.Status202
	Status203 StatusCode = internal.Status203
	Status204 StatusCode = internal.Status204
	Status205 StatusCode = internal.Status205
	Status206 StatusCode = internal.Status206
	Status207 StatusCode = internal.Status207
	Status208 StatusCode = internal.Status208

	Status300 StatusCode = internal.Status300
	Status301 StatusCode = internal.Status301
	Status302 StatusCode = internal.Status302
	Status303 StatusCode = internal.Status303
	Status304 StatusCode = internal.Status304
	Status305 StatusCode = internal.Status305
	Status307 StatusCode = internal.Status307
	Status308 StatusCode = internal.Status308

	Status400 StatusCode = internal.Status400
	Status401 StatusCode = internal.Status401
	Status402 StatusCode = internal.Status402
	Status403 StatusCode = internal.Status403
	Status404 StatusCode = internal.Status404
	Status405 StatusCode = internal.Status405
	Status406 StatusCode = internal.Status406
	Status407 StatusCode = internal.Status407
	Status408 StatusCode = internal.Status408
	Status409 StatusCode = internal.Status409
	Status410 StatusCode = internal.Status410
	Status411 StatusCode = internal.Status411
	Status412 StatusCode = internal.Status412
	Status413 StatusCode = internal.Status413
	Status414 StatusCode = internal.Status414
	Status415 StatusCode = internal.Status415
	Status416 StatusCode = internal.Status416
	Status417 StatusCode = internal.Status417
	Status421 StatusCode = internal.Status421
	Status422 StatusCode = internal.Status422
	Status423 StatusCode = internal.Status423
	Status424 StatusCode = internal.Status424
	Status425 StatusCode = internal.Status425
	Status426 StatusCode = internal.Status426
	Status428 StatusCode = internal.Status428
	Status429 StatusCode = internal.Status429
	Status431 StatusCode = internal.Status431

	Status500 StatusCode = internal.Status500
	Status501 StatusCode = internal.Status501
	Status502 StatusCode = internal.Status502
	Status503 StatusCode = internal.Status503
	Status504 StatusCode = internal.Status504
	Status505 StatusCode = internal.Status505
	Status506 StatusCode = internal.Status506
	Status507 StatusCode = internal.Status507
	Status508 StatusCode = internal.Status508
	Status511 StatusCode = internal.Status511
)

// Logging levels available for server logs.
const (
	// Informational data logged to the terminal.
	INFO_LEVEL = internal.INFO_LEVEL
	// Error data logged to the terminal.
	ERROR_LEVEL = internal.ERROR_LEVEL
	// Warning(s) logged to the terminal.
	WARN_LEVEL = internal.WARN_LEVEL
)

// Format(s) in which the request processing logs can be generated.
const (
	// :remote-addr [:date[clf]] ":method :url HTTP/:http-version" :status :res[content-length]
	COMMON_LOGGER = internal.COMMON_LOGGER
	// :method :url :status :response-time ms - :res[content-length]
	DEV_LOGGER = internal.DEV_LOGGER
	// :method :url :status :res[content-length] - :response-time ms
	TINY_LOGGER = internal.TINY_LOGGER
	// :remote-addr :method :url HTTP/:http-version :status :res[content-length] - :response-time ms
	SHORT_LOGGER = internal.SHORT_LOGGER
)

// Exposes member functions to apply colors for texts before being logged to any ANSI-supported terminals.
var TextColor = internal.TextColor
