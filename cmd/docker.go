package cmd

import (
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

			host, err := tools.DockerAutodetectHost()

			if err != nil {
				return cli.NewExitError(err, 1)
			}

			if host != "" {
				goclitools.Log("Found docker host at:", host)
			} else {
				goclitools.Log("Host docker not found but current configuration works without specifying DOCKER_HOST")
			}

			return nil
		},
	}

}
