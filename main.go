package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/mkbworks/proteus/lib/http"
)

func main() {
	server := http.NewServer()
	CurrentDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Error occurred while getting current working directory: " + err.Error())
		os.Exit(1)
	}
	
	TargetDirectory := filepath.Join(CurrentDirectory, "Files")
	server.Static("/files", TargetDirectory)
	server.Listen(8080, "localhost")
}