package cmd

import (
	"fmt"
	"os"

	"github.com/inloop/devopscli/tools"

	"github.com/inloop/goclitools"
	"github.com/urfave/cli"
)

func DockerCmd() cli.Command {
	return cli.Command{
		Name: "docker",
		Subcommands: []cli.Command{
			DockerDetectHostCmd(),
			DockerGetHostCmd(),
		},
	}
}

// DockerDetectHostCmd ...
func DockerDetectHostCmd() cli.Command {
	return cli.Command{
		Name:  "detect-host",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {

			err := DockerDetectHostAndUpdateEnv()

			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

// DockerDetectHostCmd ...
func DockerGetHostCmd() cli.Command {
	return cli.Command{
		Name:  "get-host",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {

			host, err := DockerDetectHost()

			if err != nil {
				return cli.NewExitError(err, 1)
			} else {
				fmt.Print(host)
			}

			return nil
		},
	}
}

// DockerDetectHostAndUpdateEnv ...
func DockerDetectHost() (string, error) {

	dockerHost := os.Getenv("DOCKER_HOST")

	if !tools.DockerCheckHost(dockerHost) {
		dockerHost = ""
	}

	if dockerHost == "" {
		host, err := tools.DockerAutodetectHost()
		if err != nil {
			return dockerHost, err
		}
		dockerHost = host
	}

	return dockerHost, nil
}

// DockerDetectHostAndUpdateEnv ...
func DockerDetectHostAndUpdateEnv() error {

	dockerHost, err := DockerDetectHost()

	if err != nil {
		return err
	}

	if dockerHost != "" {
		goclitools.Log("Working docker host found:", dockerHost)
		os.Setenv("DOCKER_HOST", dockerHost)
	}
	return nil
}
