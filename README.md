![Proteus Logo](https://static.citadelofcode.com/proteus/logo.png)

![Go Version](https://img.shields.io/badge/go-1.24.0-blue) ![Static Badge](https://img.shields.io/badge/powered_by-Citadel_of_Code-orange)

`Proteus` is a versatile web server framework written using Go.

## Testing

To run all the test scripts available in the module, execute the below `Makefile` command.

```bash
make test
```

To run a specific test script file, run the following command.

```bash
make testfile FILE=./test/body_parser_test.go
```

Replace `./test/body_parser_test.go` with the path of the test script file you want to execute.

## Example Usage

To work with creating a HTTP server and process incoming requests, add the below import statement at the top of your Go file.

```go
import "github.com/citadelofcode/proteus"
```

Once the import statement is included, use the below statement to create a new instance of a web server to handle incoming HTTP requests.

```go
server := proteus.NewServer()
```

Please note that, the above statement merely creates an instance of the web server. To make it listen for incoming requests, use the **Listen()** method of the server instance, as given below.

```go
server.Listen(8080, "localhost")
```

The **Listen()** method accepts two arguments - the port number where the server will listen for incoming requests and the hostname of the machine where the server instance is running.

## HTTP Version Compatibility

The `proteus` web server supports the below HTTP versions.

- [HTTP/0.9 & HTTP/1.0 - RFC 1945](https://datatracker.ietf.org/doc/html/rfc1945)
- [HTTP/1.1 - RFC 2616](https://datatracker.ietf.org/doc/html/rfc2616#autoid-45)

## HTTP Response Status Codes

The list of status codes supported by the `proteus` web server are compliant with the list of codes registered with `IANA`.

[IANA HTTP Status Code](https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml)
