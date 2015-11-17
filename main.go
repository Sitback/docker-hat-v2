// Command helm allows for faster development using Docker and Docker Compose
// on Mac OS X.
package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	// "github.com/docker/libcompose"
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

	up := kingpin.Command("up", "Start up services defined in a Docker Compose file.")
	upFile := up.Flag("file", "Docker Compose file to use").Short('f').Default("docker-compose.yml").String()
	upDetached := up.Flag("detached", "Specify a Detached mode: Run containers in the background.").Short('d').Bool()

	down := kingpin.Command("down", "Stop any running services.")
	downFile := down.Flag("file", "Specify a Docker Compose file to use").Short('f').Default("docker-compose.yml").String()

	host := kingpin.Command("host", "Control and configure your Docker Machine.")
	hostInit := host.Command("init", "Initialise your Docker Machine for use.")
	hostStart := host.Command("start", "Start your Docker Machine.").Alias("up")
	hostStop := host.Command("stop", "Stop your Docker Machine.").Alias("down")
	hostRestart := host.Command("restart", "Restart your Docker Machine.")

	switch kingpin.Parse() {
	case up.FullCommand():
		log.Info("Up!")
		log.Info("File:", *upFile)
		if *upDetached {
			log.Info("Detaching...")
		}

	case down.FullCommand():
		log.Info("Down!")
		log.Info("File:", *downFile)

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
