package tools

import (
	"errors"
	"fmt"
	"os"

	"github.com/inloop/goclitools"
)

// DockerAutodetectHost detect docker host by iterating over possible options and checking if `docker info` works
func DockerAutodetectHost() (string, error) {
	checks := []string{"tcp://docker:2375", "tcp://0.0.0.0:2375", "unix:///var/run/docker.sock", ""}

	dockerPort := os.Getenv("DOCKER_PORT")
	if dockerPort != "" {
		checks = append(checks, dockerPort)
	}

	checks = append(checks, "")

	for _, check := range checks {
		if DockerCheckHost(check) {
			return check, nil
		}
	}

	return "", errors.New("could not find docker host")
}

// CheckDockerHost ...
func DockerCheckHost(host string) bool {
	cmd := fmt.Sprintf("DOCKER_HOST=%s docker info", host)
	if host == "" {
		cmd = "docker info"
	}
	_, err := goclitools.Run(cmd)
	return err == nil
}
