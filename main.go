// Command docker-hat allows for faster development using Docker and Docker
// Compose on Mac OS X.
package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	Name    string
	Version string
	Build   string
)

func main() {
	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		fmt.Println("Error: could not determine your current directory:", pwdErr.Error())
		os.Exit(1)
	}

	var (
		verbose = kingpin.
			Flag("verbose", "Verbose mode.").
			Short('v').
			Bool()
	)

	// Setup '-h' as an alias for the help flag.
	kingpin.CommandLine.HelpFlag.Short('h')

	// Setup version printing.
	kingpin.Flag("version", "Show version.").PreAction(kingpin.Action(func(*kingpin.ParseContext) error {
		fmt.Printf("%v version %v, build %v\n", Name, Version, Build)
		os.Exit(0)
		return nil
	})).Bool()

	kingpin.Parse()

	if *verbose {
		fmt.Println("Hello, this is verbose mode.")
		fmt.Println("The current directory is:", pwd)
	}
}
