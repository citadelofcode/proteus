package internal

import (
	"fmt"
)

// Structure to apply ANSI escape sequences to color the text written to the terminal. These color codes are supported on all the major MacOS and Linux terminals and all the latest windows powershell terminal, windows terminals.
type textColor struct {}

// Function to color the given text as "black" and return the colored text.
func (tc *textColor) Black(value string) string {
	return fmt.Sprintf("\033[30m%s\033[0m", value)
}

// Function to color the given text as "red" and return the colored text.
func (tc *textColor) Red(value string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", value)
}

// Function to color the given text as "green" and return the colored text.
func (tc *textColor) Green(value string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", value)
}

// Function to color the given text as "yellow" and return the colored text.
func (tc *textColor) Yellow(value string) string {
	return fmt.Sprintf("\033[33m%s\033[0m", value)
}

// Function to color the given text as "blue" and return the colored text.
func (tc *textColor) Blue(value string) string {
	return fmt.Sprintf("\033[34m%s\033[0m", value)
}

// Function to color the given text as "magenta" and return the colored text.
func (tc *textColor) Magenta(value string) string {
	return fmt.Sprintf("\033[35m%s\033[0m", value)
}

// Function to color the given text as "cyan" and return the colored text.
func (tc *textColor) Cyan(value string) string {
	return fmt.Sprintf("\033[36m%s\033[0m", value)
}

// Function to color the given text as "white" and return the colored text.
func (tc *textColor) White(value string) string {
	return fmt.Sprintf("\033[37m%s\033[0m", value)
}


// Global instance to apply colors to text printed on ANSI supported terminals.
var TextColor = textColor{}
