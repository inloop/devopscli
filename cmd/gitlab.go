package cmd

import "github.com/urfave/cli"

func GitlabCmd() cli.Command {
	return cli.Command{
		Name: "gitlab",
		Subcommands: []cli.Command{
			GitlabDockerCmd(),
			GitlabGoCmd(),
		},
	}
}
