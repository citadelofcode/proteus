package http

import (
	"fmt"
	"log"
)

// Logger to log messages to the console or to a file.
type logger struct {
	// Pointer to the log.Logger instance created.
	srvLogger *log.Logger
	// Name of the server instance for which logs sre being recorded.
	serverName string
}

// Logs an error message to the log file. If logger is not initialized, the error message is printed to stdout.
func (lg *logger) logError(Msg string) {
	if lg.srvLogger != nil {
		lg.srvLogger.Printf("%s  ERROR  %s", lg.serverName, Msg)
	} else {
		fmt.Printf("%s  %s  ERROR  %s", getRfc1123Time(), lg.serverName, Msg)
	}
}

// Logs a message to the log file. If logger is not initialized, the message is printed to stdout.
func (lg *logger)logInfo(Msg string) {
	if lg.srvLogger != nil {
		lg.srvLogger.Printf("%s  INFO  %s", lg.serverName, Msg)
	} else {
		fmt.Printf("%s  %s  INFO  %s", getRfc1123Time(), lg.serverName, Msg)
	}
}