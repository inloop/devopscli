package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/inloop/goclitools"
	"github.com/urfave/cli"
)

func GitlabGoCmd() cli.Command {
	return cli.Command{
		Name: "go",
		Subcommands: []cli.Command{
			GitlabGoTestCmd(),
			GitlabGoBuildCmd(),
		},
	}
}

func GitlabGoTestCmd() cli.Command {
	return cli.Command{
		Name: "test",
		Action: func(c *cli.Context) error {
			projectUrl := os.Getenv("CI_PROJECT_URL")
			goPath := os.Getenv("GOPATH")

			if projectUrl == "" {
				return cli.NewExitError("missing CI_PROJECT_URL environment variable", 1)
			}

			if goPath == "" {
				return cli.NewExitError("missing GOPATH environment variable", 1)
			}

			if err := gitlabGoTest(projectUrl, goPath); err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func gitlabGoTest(projectUrl, goPath string) error {
	projectPath := strings.Replace(projectUrl, "https://", "http://", 1)
	projectPath = strings.Replace(projectUrl, "http://", goPath, 1)
	parent := path.Dir(projectUrl)

	if err := goclitools.RunInteractive(fmt.Sprintf("mkdir -p %s", parent)); err != nil {
		return err
	}
	if err := goclitools.RunInteractive(fmt.Sprintf("cp -R . %s", projectUrl)); err != nil {
		return err
	}

	if err := goclitools.RunInteractiveInDir("go get -d -v ./... && go test ./...", projectPath); err != nil {
		return err
	}
	return nil
}

func GitlabGoBuildCmd() cli.Command {
	return cli.Command{
		Name: "build",
		Action: func(c *cli.Context) error {
			projectUrl := os.Getenv("CI_PROJECT_URL")
			goPath := os.Getenv("GOPATH")

			if projectUrl == "" {
				return cli.NewExitError("missing CI_PROJECT_URL environment variable", 1)
			}

			if goPath == "" {
				return cli.NewExitError("missing GOPATH environment variable", 1)
			}

			if err := gitlabGoBuild(projectUrl, goPath); err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func gitlabGoBuild(projectUrl, goPath string) error {

	// - mkdir -p $GOPATH/src/git.inloop.eu/inloop-ci
	// - cp -R . $GOPATH/src/git.inloop.eu/inloop-ci/ios-provisioning-cli
	// - cd $GOPATH/src/git.inloop.eu/inloop-ci/ios-provisioning-cli
	// - go get -d -v ./...
	// - go get github.com/mitchellh/gox
	// - gox -osarch="darwin/amd64" ./...
	// - cp ./ios-provisioning-cli_darwin_amd64 $CI_PROJECT_DIR/ios-provisioning

	projectPath := strings.Replace(projectUrl, "https://", "http://", 1)
	projectPath = strings.Replace(projectUrl, "http://", goPath, 1)
	parent := path.Dir(projectUrl)

	if err := goclitools.RunInteractive(fmt.Sprintf("mkdir -p %s", parent)); err != nil {
		return err
	}
	if err := goclitools.RunInteractive(fmt.Sprintf("cp -R . %s", projectUrl)); err != nil {
		return err
	}

	if err := goclitools.RunInteractiveInDir("go get -d -v ./... && go get github.com/mitchellh/gox", projectPath); err != nil {
		return err
	}

	if err := goclitools.RunInteractiveInDir("gox ./... && cp ./* $CI_PROJECT_DIR/bin", projectPath); err != nil {
		return err
	}

	return nil
}
