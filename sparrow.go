package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/maheshkumaarbalaji/project-sparrow/lib/http"
)

func main() {
	server, err := http.NewServer("localhost")
	if err != nil {
		fmt.Printf("Error occurred while creating web server: %s", err.Error())
		os.Exit(1)
	}
	CurrentDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Error occurred while getting current working directory: " + err.Error())
		os.Exit(1)
	}
	TargetDirectory := filepath.Join(CurrentDirectory, "Files")
	server.Static("/files", TargetDirectory)
	server.Listen(8080)
}