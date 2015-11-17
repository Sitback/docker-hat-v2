// Command docker-hat allows for faster development using Docker and Docker
// Compose on Mac OS X.
package main

import (
	"fmt"
	// "github.com/docker/libcompose"
	log "github.com/Sirupsen/logrus"
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

	up := kingpin.Command("up", "Start up Docker containers defined in a Docker Compose file.")
	down := kingpin.Command("down", "Stop any Docker containers running from a Docker Compose file.")

	host := kingpin.Command("host", "Control and configure your Docker Machine.")
	hostInit := host.Command("init", "Initialises your Docker Machine for use.")
	hostStart := host.Command("start", "Starts your Docker Machine.").Alias("up")
	hostStop := host.Command("stop", "Stops your Docker Machine.").Alias("down")
	hostRestart := host.Command("restart", "Restarts your Docker Machine.")

	switch kingpin.Parse() {
	case up.FullCommand():
		log.Info("Up!")

	case down.FullCommand():
		log.Info("Down!")

	case hostInit.FullCommand():
		log.Info("host init!")

	case hostStart.FullCommand():
		log.Info("host start!")

	case hostStop.FullCommand():
		log.Info("host stop!")

	case hostRestart.FullCommand():
		log.Info("host restart!")
	}

	if *verbose {
		log.Info("Hello, this is verbose mode.")
		log.Info("The current directory is:", pwd)
	}
}
