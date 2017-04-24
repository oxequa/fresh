package main

import (
	"flag"
	"os"
)

func main() {
	// Commands
	helpCmd := flag.NewFlagSet("", flag.ExitOnError)
	initCmd := flag.NewFlagSet("new", flag.ExitOnError)

	// Count subcommand flag pointers
	initG := initCmd.String("gateway", "", "Gateway pattern for microservices architecture.")
	initD := initCmd.String("default", "", "Web server with DDD pattern.")
	initM := initCmd.String("minimal", "", "Minimal web server starter project.")

	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) == 1 {
		// print commands list
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		helpCmd.Parse(os.Args[2:])
	case "count":
		initCmd.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
	if helpCmd.Parsed() {
		println(*initG, *initD, *initM)
	}
}
