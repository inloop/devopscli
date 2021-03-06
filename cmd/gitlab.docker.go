package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/inloop/goclitools"
	"github.com/urfave/cli"
)

func GitlabDockerCmd() cli.Command {
	return cli.Command{
		Name: "docker",
		Subcommands: []cli.Command{
			GitlabDockerBuildCmd(),
			GitlabDockerLoginCmd(),
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
				Name:  "tag-suffix",
				Value: "",
				Usage: "Tag suffix",
			},
			cli.BoolFlag{
				Name:  "no-cache",
				Usage: "Disable build cache",
			},
			cli.BoolFlag{
				Name:  "specific-tag",
				Usage: "Generate tag names bound to specific version/commit",
			},
			cli.BoolFlag{
				Name:  "push-latest",
				Usage: "Push currently build ALSO as latest tag",
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
				EnvVar: "CI_JOB_TOKEN",
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
			cli.StringFlag{
				Name:  "file",
				Usage: "Name of the Dockerfile (Default is 'PATH/Dockerfile')",
			},
		},
		Action: func(c *cli.Context) error {

			image := c.String("image")
			buildPath := c.String("path")
			file := c.String("file")
			tag := c.String("tag")

			if tag == "" {
				tag = tagForRefName(os.Getenv("CI_COMMIT_REF_NAME"), c.Bool("specific-tag"))
			}

			tag = c.String("tag-prefix") + tag + c.String("tag-suffix")

			if tag == "" {
				tag = "latest"
			}

			if err := DockerDetectHostAndUpdateEnv(); err != nil {
				return cli.NewExitError(err, 1)
			}

			buildParams := []string{}

			if file != "" {
				buildParams = append(buildParams, "--file "+file)
			}
			if c.Bool("no-cache") {
				buildParams = append(buildParams, "--no-cache")
			}

			loginCmd := fmt.Sprintf("docker login -u %s -p %s %s", c.String("username"), c.String("password"), c.String("registry"))
			buildCmd := fmt.Sprintf("docker build -t %s:%s %s %s", image, tag, strings.Join(buildParams, " "), buildPath)
			pushCmd := fmt.Sprintf("docker push %s:%s", image, tag)

			cmds := []string{loginCmd, buildCmd, pushCmd}

			if c.Bool("push-latest") {
				latestTagCmd := fmt.Sprintf("docker tag %s:%s %s:latest", image, tag, image)
				pushLatestCmd := fmt.Sprintf("docker push %s:latest", image)
				cmds = append(cmds, latestTagCmd, pushLatestCmd)
			}

			if err := goclitools.RunInteractive(strings.Join(cmds, " && ")); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
}

func GitlabDockerLoginCmd() cli.Command {
	return cli.Command{
		Name:  "login",
		Usage: "Login using CI job token",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "username, u",
				Value: "gitlab-ci-token",
				Usage: "Docker registry username",
			},
			cli.StringFlag{
				Name:   "password, p",
				Value:  "",
				Usage:  "Docker registry password",
				EnvVar: "CI_JOB_TOKEN",
			},
			cli.StringFlag{
				Name:   "registry, r",
				Value:  "",
				Usage:  "Docker registry url",
				EnvVar: "CI_REGISTRY",
			},
		},
		Action: func(c *cli.Context) error {
			if err := DockerDetectHostAndUpdateEnv(); err != nil {
				return cli.NewExitError(err, 1)
			}

			loginCmd := fmt.Sprintf("docker login -u %s -p %s %s", c.String("username"), c.String("password"), c.String("registry"))
			if err := goclitools.RunInteractive(loginCmd); err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

// tagForRefName maps reference name to tagName
// master -> latest
// develop -> unstable
// release/x -> x
// other -> other
func tagForRefName(ref string, specific bool) string {
	tag := ref
	if specific {
		tag = fmt.Sprintf("%s", os.Getenv("CI_BUILD_REF")[0:8])
	} else if tag == "master" {
		tag = "latest"
	} else if tag == "develop" {
		tag = "unstable"
	} else if strings.HasPrefix(tag, "release/") {
		tag = strings.Replace(tag, "release/", "", 1)
	}
	return strings.Replace(tag, "/", "-", -1)
}
