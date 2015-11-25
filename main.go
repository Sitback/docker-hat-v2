// Command helm allows for faster development using Docker and Docker Compose
// on Mac OS X.
package main

import (
	"fmt"
	"os"

	"github.com/Sitback/helm/host"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Name    string
	Version string
	Build   string
)

func main() {
	var err error

	pp, err := docker.NewProject(&docker.Context{
		Context: project.Context{
			ComposeFile: "docker-compose.yml",
			ProjectName: "helm",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	pp.Parse()
	// pp.Up()

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: could not determine your current directory:", err.Error())
		os.Exit(1)
	}

	var (
		verbose = kingpin.
			Flag("verbose", "Verbose mode.").
			Short('v').
			Bool()
	)
	// Clean up help.
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.CommandLine.Help = "Fast, Docker Compose-based development for Mac OS X."
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

	cmdHost := kingpin.Command("host", "Control and configure the Helm Docker Machine.")
	hostInit := cmdHost.Command("init", "Initialise the Helm Docker Machine for use.")
	hostInitForce := hostInit.Flag("force", "Force re-initialisation").Short('f').Bool()
	hostStart := cmdHost.Command("start", "Start the Helm Docker Machine.").Alias("up")
	hostStop := cmdHost.Command("stop", "Stop the Helm Docker Machine.").Alias("down")
	hostRestart := cmdHost.Command("restart", "Restart the Helm Docker Machine.")
	hostDestroy := cmdHost.Command("destroy", "Stop and remove the Helm Docker Machine.")
	hostStatus := cmdHost.Command("status", "Show status of the Helm Docker Machine.")

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
		_, err := host.NewHost(true, *hostInitForce)
		if err != nil {
			log.Fatal("The Helm host already exists, run `helm host destroy` first if you wish to recreate it. Alternatively, run init with the force flag: `helm host init --force`.")
		}
		log.Info("host init!")

	case hostStart.FullCommand():
		helmHost, err := host.NewHost(false, false)
		if err != nil {
			log.Fatal("Could not start host, it might not exist. Try running `helm host init`.")
		}
		err = helmHost.Start()
		if err != nil {
			log.Fatal(err)
		}
		log.Info("host start!")

	case hostStop.FullCommand():
		helmHost, err := host.NewHost(false, false)
		if err != nil {
			log.Fatal("Could not stop host, it might not exist. Try running `helm host init`.")
		}
		err = helmHost.Stop()
		if err != nil {
			log.Fatal(err)
		}
		log.Info("host stop!")

	case hostRestart.FullCommand():
		helmHost, err := host.NewHost(false, false)
		if err != nil {
			log.Fatal("Could not restart host, it might not exist. Try running `helm host init`.")
		}
		err = helmHost.Restart()
		if err != nil {
			log.Fatal(err)
		}
		log.Info("host restart!")

	case hostDestroy.FullCommand():
		helmHost, err := host.NewHost(false, false)
		if err != nil {
			log.Fatal("Could not destroy host, it might not exist. Try running `helm host init`.")
		}
		err = helmHost.Destroy()
		if err != nil {
			log.Fatal(err)
		}
		log.Info("host destory!")

	case hostStatus.FullCommand():
		helmHost, err := host.NewHost(false, false)
		if err != nil {
			log.Fatal("Host doesn't exist. Try running `helm host init`.")
		}

		status, err := helmHost.Host.Driver.GetState()
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Host status: %v", status.String())
	}

	if *verbose {
		log.Info("Hello, this is verbose mode.")
		log.Info("The current directory is:", pwd)
	}
}
