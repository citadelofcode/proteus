# Proteus

Proteus is a versatile web server framework written using Go. Currently, this framework supports only HTTP/1.0 requests.

## Running the project

To build and run the project, execute the following commands on project root directory.

```bash
# Below command builds the project and generates the output executable file.
# The name of the executable file is main.out by default. 
# If you want a different name, you can run the build command with the -o flag.
go build -o proteus.out

# Below command runs the executable file.
./proteus.out
```

## Compliance

The `proteus` web server is compliant with the below RFCs.

- [HTTP/1.0 - RFC 1945](https://datatracker.ietf.org/doc/html/rfc1945)