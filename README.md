# Proteus

<img src="https://citadel-of-code.s3.us-west-1.amazonaws.com/static/project-proteus/proteus-logo.jpeg" style="border-radius:50%" align="right" width="159px" alt="Proteus logo">

Proteus is a versatile web server framework written using Go. 

This is my solution to the challenge available at [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-webserver) to create my own web server.

## Running the project

To build and run the project, execute the following commands on the project's root directory.

```bash
# The below command makes the bash script file as an executable.
chmod +x ./run.sh

# Below command runs the shell script that compiles and runs the server
./run.sh
```

## Testing

To run all the test scripts available in the module, execute the following commands.

```bash
# Below command gives executable permissions to the shell script file.
chmod +x ./test.sh

# Execute the shell script to run the test cases.
./test.sh
```

## Example Usage

To work with creating a HTTP server and process incoming requests, add the below import statement at the top of your Go file.

```go
import "github.com/mkbworks/proteus/lib/http"
```

Once the import statement is included, use the below statement to create a new instance of a web server to handle incoming HTTP requests.

```go
server := http.NewServer()
```

Please note that, the above statement merely creates an instance of the web server. To make it listen for incoming requests, use the **Listen()** method of the server instance, as given below.

```go
server.Listen(8080, "localhost")
```

The **Listen()** method accepts two arguments - the port number where the server will listen for incoming requests and the hostname of the machine where the server instance is running.

To create static directory in the web server instance, use the following code.

```go
server.Static("/files/static", **TargetDirectoryPath**)
```

To declare a custom route and its associated handler function, refer to the following code snippet.

```go
server.Get("/user/:name", func(req *http.HttpRequest, res *http.HttpResponse) error {
    names, _ := req.Segments.Get("name")
    server.LogInfo(fmt.Sprintf("The name value in the path is %s\n", strings.Join(names, ",")))
    res.Status(http.StatusOK)
    res.Send("The name parameter value received is: " + strings.Join(names, ", "))
    return nil
})
```

The handler function must accept two parameters - pointer to a HttpRequest and pointer to a HttpResponse instance. It must return an error.

## HTTP Version Compatibility

The `proteus` web server supports the below HTTP versions.

- [HTTP/0.9 & HTTP/1.0 - RFC 1945](https://datatracker.ietf.org/doc/html/rfc1945)
- [HTTP/1.1 - RFC 2616](https://datatracker.ietf.org/doc/html/rfc2616#autoid-45)