# Proteus

Proteus is a versatile web server framework written using Go. ***Currently, this framework supports only HTTP/0.9 and HTTP/1.0 requests.***

## Running the project

To build and run the project, execute the following commands on the project's root directory.

```bash
# Below command builds the project and generates the output executable file.
# The name of the executable file is main.out by default. 
# If you want a different name, you can run the build command with the -o flag.
go build -o proteus.out

# Below command runs the executable file.
./proteus.out
```

## Example Usage

To work with creating a HTTP server and process incoming requests, add the below import statement at the top of your Go file.

```
import "github.com/maheshkumaarbalaji/proteus/lib/http"
```

Once the import statement is included, use the below statement to create a new instance of a web server to handle incoming HTTP requests.

```
server := http.NewServer()
```

Please note that, the above statement merely creates an instance of the web server. To make it listen for incoming requests, use the **Listen()** method of the server instance, as given below.

```
server.Listen(8080, "localhost")
```

The **Listen()** method accepts two arguments - the port number where the server will listen for incoming requests and the hostname of the machine where the server instance is running.

To create static directory in the web server instance, use the following code.

```
server.Static("/files/static", **TargetDirectoryPath**)
```

## Testing

Each package in the module contains unit test scripts which can be identified by the "_test.go" suffix present in the files. To run all test scripts in the module, execute the following command.

```bash
go test ./lib/... -v -cover
```

- The **-v** command line option enables the verbose logs to be printed as the test scripts are being executed.
- The **-cover** command line option displays the code coverage information for the test cases executed.

To run test cases for a specific package, run the following command.

```bash
go test ./lib/http -v -cover
```

This command runs the test cases defined for the entire `http` package. To run for a package of your choice, replace `http` with the name of the package for which you want to run the test cases.

## Compliance

The `proteus` web server is compliant with the below RFCs.

- [HTTP/0.9 & HTTP/1.0 - RFC 1945](https://datatracker.ietf.org/doc/html/rfc1945)