package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
	"novacloud.cz/udi-cli/utils"
)

func GitlabDockerCmd() cli.Command {
	return cli.Command{
		Name: "docker",
		Subcommands: []cli.Command{
			GitlabDockerBuildCmd(),
		},
	}
}

func GitlabDockerBuildCmd() cli.Command {
	return cli.Command{
		Name: "build",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "tag, t",
				Value: "",
				Usage: "Tag for image (deduced using $CI_COMMIT_REF_NAME)",
			},
			cli.StringFlag{
				Name:   "image, i",
				Value:  "",
				Usage:  "Image name",
				EnvVar: "CI_REGISTRY_IMAGE",
			},
			cli.StringFlag{
				Name:  "email, e",
				Value: "john.doe@example.com",
				Usage: "Docker registry email",
			},
			cli.StringFlag{
				Name:  "username, u",
				Value: "gitlab-ci-token",
				Usage: "Docker registry username",
			},
			cli.StringFlag{
				Name:   "password, p",
				Value:  "",
				Usage:  "Docker registry password",
				EnvVar: "CI_BUILD_TOKEN",
			},
			cli.StringFlag{
				Name:   "registry, r",
				Value:  "",
				Usage:  "Docker registry url",
				EnvVar: "CI_REGISTRY",
			},
		},
		Action: func(c *cli.Context) error {

			image := c.String("image")
			tag := c.String("tag")

			if tag == "" {
				tag = utils.Getenv("CI_COMMIT_REF_NAME", "")
				if tag == "develop" {
					tag = "latest"
				}
			}

			if tag != "" {
				image += ":" + tag
			}

			loginCmd := fmt.Sprintf("docker login -e %s -u %s -p %s %s", c.String("email"), c.String("username"), c.String("password"), c.String("registry"))
			buildCmd := fmt.Sprintf("docker build -t %s .", image)
			pushCmd := fmt.Sprintf("docker push %s", image)

			cmds := []string{loginCmd, buildCmd, pushCmd}

			if err := utils.RunInteractive(strings.Join(cmds, " && ")); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}

}
