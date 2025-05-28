package http

// Function used to update if subsequent middlewares in the stack have to be processed before a response can be sent back to the client.
type StopFunction func()
// A function to execute operations like validations and transformations before before performing the backend tasks.
type Middleware func(*HttpRequest, *HttpResponse, StopFunction) error

// Structure to hold and process one or more middlewares.
// Middlewares are executed in the order in which they were defined.
// They can either be associated with a route or to a server instance.
// Miidlewares associated with a server instance are applicable to all requests directed to the server instance.
type Middlewares struct {
	// Flag to denote if the next middleware in the stack will be executed.
	// the default value for this flag is true.
	// At the end of execution of each middleware this flag can be updated by using the "StopFunction" to end the execution of request with that middleware.
	ProcessNext bool
	// List of all middlewares to be processed as part of the stack for the associated instance.
	Stack []Middleware
}

// Stops further execution of middlewares in the stack by marking the "ProcessNext" flag as false.
func (mds *Middlewares) Stop() {
	mds.ProcessNext = false
}

// Create a new instance of "Middlewares" and return a reference to the instance.
func CreateMiddlewares(middlewareList []Middleware) *Middlewares {
	middlewares := new(Middlewares)
	middlewares.ProcessNext = true
	middlewares.Stack = make([]Middleware, 0)
	middlewares.Stack = append(middlewares.Stack, middlewareList...)
	return middlewares
}
