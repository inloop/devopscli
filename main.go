package main

import (
	"os"

	"github.com/inloop/devopscli/cmd"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "DevOps cli"
	app.Usage = "..."
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		cmd.GitlabCmd(),
		cmd.DockerCmd(),
	}

	app.Run(os.Args)
}
