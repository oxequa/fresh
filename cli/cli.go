package main

import "fmt"

// Commands list
func list() {
	fmt.Println("Fresh commands:")
	fmt.Println("", "new", "\t\t", "Create a new project from an available boilerplate.")
	fmt.Println("", "list", "\t\t", "List the boilerplates available.")
	fmt.Println("", "run", "\t\t", "Run a Fresh server with Hot reload.")
}
