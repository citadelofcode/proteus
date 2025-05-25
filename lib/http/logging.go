package http

import (
	"fmt"
	"log"
)

// Custom logger to record logs of requests being processed by the web server instance.
type Sonar struct {
	// Underlying logger instance to be used to record logs.
	logger *log.Logger
	// Name of the web server instance.
	serverName string
	// Format of the log message to be recorded for each request processed.
	logFormat string
}

// Logs an error message to the log file. If logger is not initialized, the error message is printed to stdout.
func (sn *Sonar) customError(Msg string) {
	if sn.logger != nil {
		sn.logger.Printf("%s  ERROR  %s", sn.serverName, Msg)
	} else {
		fmt.Printf("%s  %s  ERROR  %s", getRfc1123Time(), sn.serverName, Msg)
	}
}

// Logs a message to the log file. If logger is not initialized, the message is printed to stdout.
func (sn *Sonar) customInfo(Msg string) {
	if sn.logger != nil {
		sn.logger.Printf("%s  INFO  %s", sn.serverName, Msg)
	} else {
		fmt.Printf("%s  %s  INFO  %s", getRfc1123Time(), sn.serverName, Msg)
	}
}
