package host

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Sitback/helm/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/drivers/virtualbox"
	"github.com/docker/machine/libmachine"
	mcnHost "github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/mcnerror"
	"github.com/docker/machine/libmachine/state"
)

const (
	DEFAULT_NAME   = utils.PROGRAM_NAME
	DEFAULT_CPU    = 2
	DEFAULT_MEMORY = 2048
)

type Host struct {
	Host   *mcnHost.Host
	Client *libmachine.Client
}

func NewHost(create bool, force bool) (*Host, error) {
	// Create a Virtualbox host if an existing host doesn't already exist.
	client := libmachine.NewClient(mcndirs.GetBaseDir())

	existing, err := client.Load(DEFAULT_NAME)
	if err != nil {
		switch err.(type) {
		case mcnerror.ErrHostDoesNotExist:
		default:
			log.Fatal(err)
		}
	}

	if existing != nil && create && !force {
		return nil, errors.New("Host already exists.")
	} else if existing != nil && !create {
		log.Debug("Existing host found.")
		return &Host{
			Host:   existing,
			Client: client,
		}, nil
	}

	if create {
		driver := virtualbox.NewDriver(DEFAULT_NAME, mcndirs.GetBaseDir())
		driver.CPU = DEFAULT_CPU
		driver.Memory = DEFAULT_MEMORY
		// Disable the '/Users' mount, we handle that ourselves via
		// `docker-unisync`.
		driver.NoShare = true

		data, err := json.Marshal(driver)
		if err != nil {
			log.Fatal(err)
		}

		pluginDriver, err := client.NewPluginDriver("virtualbox", data)
		if err != nil {
			log.Fatal(err)
		}

		h, err := client.NewHost(pluginDriver)
		if err != nil {
			log.Fatal(err)
		}

		if err := client.Create(h); err != nil {
			log.Fatal(err)
		}

		return &Host{
			Host:   h,
			Client: client,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("The %v host doesn't exist and create was not specified.", DEFAULT_NAME))
}

func (h *Host) Start() error {
	s, err := h.Host.Driver.GetState()
	if err != nil {
		return err
	}

	if s == state.Running {
		return errors.New("Host is already running.")
	} else if s == state.Starting {
		return errors.New("Host is starting.")
	}

	err = h.Host.Start()
	if err != nil {
		return err
	}
	return nil
}

func (h *Host) Stop() error {
	s, err := h.Host.Driver.GetState()
	if err != nil {
		return err
	}

	if s == state.Stopped {
		return errors.New("Host is already stopped.")
	} else if s == state.Stopping {
		return errors.New("Host is stopping.")
	}

	err = h.Host.Stop()
	if err != nil {
		return err
	}
	return nil
}

func (h *Host) Restart() error {
	err := h.Host.Restart()
	if err != nil {
		return err
	}
	return nil
}

func (h *Host) Destroy() error {
	err := h.Stop()
	if err != nil {
		return err
	}
	return h.Client.Remove(h.Host.Name)
}
