package tools

import (
	"errors"
	"fmt"

	"github.com/inloop/goclitools"
)

// DockerAutodetectHost detect docker host by iterating over possible options and checking if `docker info` works
func DockerAutodetectHost() (string, error) {
	checks := []string{"tcp://docker:2375", "tcp://0.0.0.0:2375", "unix:///var/run/docker.sock",""}

	for _, check := range checks {
		cmd := fmt.Sprintf("DOCKER_HOST=%s docker info")
		if check == "" {
			cmd = "docker info"
		}
		_, err := goclitools.Run(cmd)
		if err == nil {
			return check, nil
		}
	}

	return "", errors.New("could not find docker host")
}
