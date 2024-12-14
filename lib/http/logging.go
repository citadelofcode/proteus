package http

import (
	"fmt"
	"log"
	"strings"
)

// Logger to log messages to the console or to a file.
type Logger struct {
	// Pointer to the log.Logger instance created.
	srvLogger *log.Logger
}

// Logs an error message to the log file. If logger is not initialized, the error message is printed to stdout.
func (lg *Logger) LogError(Msg string) {
	Msg = strings.TrimSpace(Msg)
	ServerName := getServerDefaults("server_name")

	if lg.srvLogger != nil {
		lg.srvLogger.Printf("%s  ERROR  %s", ServerName, Msg)
	} else {
		fmt.Printf("%s  %s  ERROR  %s", getRfc1123Time(), ServerName, Msg)
	}
}

// Logs a message to the log file. If logger is not initialized, the message is printed to stdout.
func (lg *Logger)LogInfo(Msg string) {
	Msg = strings.TrimSpace(Msg)
	ServerName := getServerDefaults("server_name")
	
	if lg.srvLogger != nil {
		lg.srvLogger.Printf("%s  INFO  %s", ServerName, Msg)
	} else {
		fmt.Printf("%s  %s  INFO  %s", getRfc1123Time(), ServerName, Msg)
	}
}