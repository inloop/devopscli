package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/inloop/devopscli/tools"

	"github.com/inloop/goclitools"
	"github.com/urfave/cli"
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
				Name:  "tag-prefix",
				Value: "",
				Usage: "Tag prefix",
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
			cli.StringFlag{
				Name:  "path",
				Value: ".",
				Usage: "Docker build path used as root folder for building image",
			},
		},
		Action: func(c *cli.Context) error {

			image := c.String("image")
			buildPath := c.String("path")
			tag := c.String("tag")

			if tag == "" {
				tag = os.Getenv("CI_COMMIT_REF_NAME")
				if tag == "develop" {
					tag = "latest"
				}
			}

			tag = c.String("tag-prefix") + tag

			if tag != "" {
				image += ":" + tag
			}

			dockerHost := os.Getenv("DOCKER_HOST")

			if !tools.DockerCheckHost(dockerHost) {
				goclitools.Log("Current DOCKER_HOST is not working")
				dockerHost = ""
			}

			if dockerHost == "" {
				goclitools.Log("Detecting docker host...")
				host, err := tools.DockerAutodetectHost()
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				dockerHost = host
			}

			loginCmd := fmt.Sprintf("docker login -u %s -p %s %s", c.String("username"), c.String("password"), c.String("registry"))
			buildCmd := fmt.Sprintf("docker build -t %s %s", image, buildPath)
			pushCmd := fmt.Sprintf("docker push %s", image)

			cmds := []string{loginCmd, buildCmd, pushCmd}

			if dockerHost != "" {
				cmd := fmt.Sprintf("export DOCKER_HOST=%s", dockerHost)
				cmds = append([]string{cmd}, cmds...)
			}

			if err := goclitools.RunInteractive(strings.Join(cmds, " && ")); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}

}
