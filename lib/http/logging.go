package http

import (
	"strings"
	"fmt"
)

// Logs an error message to the log file. If logger is not initialized, the error message is printed to stdout.
func LogError(Msg string) {
	Msg = strings.TrimSpace(Msg)
	ServerName := getServerDefaults("server_name")

	if SrvLogger != nil {
		SrvLogger.Printf("%s  ERROR  %s", ServerName, Msg)
	} else {
		fmt.Printf("%s  %s  ERROR  %s", getRfc1123Time(), ServerName, Msg)
	}
}

// Logs a message to the log file. If logger is not initialized, the message is printed to stdout.
func LogInfo(Msg string) {
	Msg = strings.TrimSpace(Msg)
	ServerName := getServerDefaults("server_name")
	
	if SrvLogger != nil {
		SrvLogger.Printf("%s  INFO  %s", ServerName, Msg)
	} else {
		fmt.Printf("%s  %s  INFO  %s", getRfc1123Time(), ServerName, Msg)
	}
}