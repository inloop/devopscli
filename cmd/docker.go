package cmd

import (
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

// DockerDetectHostAndUpdateEnv ...
func DockerDetectHostAndUpdateEnv() error {

	dockerHost := os.Getenv("DOCKER_HOST")

	if !tools.DockerCheckHost(dockerHost) {
		goclitools.Log("Current DOCKER_HOST(" + dockerHost + ") seems to not work properly!")
		dockerHost = ""
	}

	if dockerHost == "" {
		goclitools.Log("Trying to find another docker host...")
		host, err := tools.DockerAutodetectHost()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		goclitools.Log("New working docker host found:", dockerHost)
		dockerHost = host
	}

	if dockerHost != "" {
		os.Setenv("DOCKER_HOST", dockerHost)
	}
	return nil
}
