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
			GitlabGoCompileCmd(),
			GitlabGoBuildCmd(), // deprecated
		},
	}
}

// GitlabGoTestCmd
func GitlabGoTestCmd() cli.Command {
	return cli.Command{
		Name: "test",
		Action: func(c *cli.Context) error {
			projectURL := os.Getenv("CI_PROJECT_URL")
			projectDir := os.Getenv("CI_PROJECT_DIR")
			goSourcePath := path.Join(os.Getenv("GOPATH"), "src")

			if projectURL == "" {
				return cli.NewExitError("missing CI_PROJECT_URL environment variable", 1)
			}
			if projectDir == "" {
				return cli.NewExitError("missing CI_PROJECT_DIR environment variable", 1)
			}
			if goSourcePath == "" {
				return cli.NewExitError("missing GOPATH environment variable", 1)
			}

			if err := gitlabGoTest(projectDir, projectURL, goSourcePath); err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func gitlabGoTest(projectDir, projectUrl, goSourcePath string) error {
	projectPath := strings.Replace(projectUrl, "https://", "http://", 1)
	projectPath = strings.Replace(projectPath, "http:/", goSourcePath, 1)
	parent := path.Dir(projectPath)

	if err := goclitools.RunInteractive(fmt.Sprintf("mkdir -p %s", parent)); err != nil {
		return err
	}
	if err := goclitools.RunInteractiveInDir(fmt.Sprintf("cp -R . %s", projectPath), projectDir); err != nil {
		return err
	}

	return goclitools.RunInteractiveInDir("go get -d -v ./... && go test ./...", projectPath)
}

// GitlabGoCompileCmd
func GitlabGoCompileCmd() cli.Command {
	return cli.Command{
		Name: "compile",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "osarch",
				Value: "",
				Usage: "Space-separated list of os/arch pairs to build for",
			},
		},
		Action: func(c *cli.Context) error {
			projectURL := os.Getenv("CI_PROJECT_URL")
			projectDir := os.Getenv("CI_PROJECT_DIR")
			goPath := os.Getenv("GOPATH")
			osarch := c.String("osarch")

			if projectURL == "" {
				return cli.NewExitError("missing CI_PROJECT_URL environment variable", 1)
			}
			if projectDir == "" {
				return cli.NewExitError("missing CI_PROJECT_DIR environment variable", 1)
			}
			if goPath == "" {
				return cli.NewExitError("missing GOPATH environment variable", 1)
			}

			if err := gitlabGoCompile(projectDir, projectURL, goPath, osarch); err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

// GitlabGoBuildCmd
func GitlabGoBuildCmd() cli.Command {
	return cli.Command{
		Name: "build",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "osarch",
				Value: "",
				Usage: "Space-separated list of os/arch pairs to build for",
			},
		},
		Action: func(c *cli.Context) error {
			return cli.NewExitError("deprecated, please use `compile` command instead", 1)
		},
	}
}

func gitlabGoCompile(projectDir, projectURL, goSourcePath, osarch string) error {

	// - mkdir -p $GOPATH/src/git.inloop.eu/inloop-ci
	// - cp -R . $GOPATH/src/git.inloop.eu/inloop-ci/ios-provisioning-cli
	// - cd $GOPATH/src/git.inloop.eu/inloop-ci/ios-provisioning-cli
	// - go get -d -v ./...
	// - go get github.com/mitchellh/gox
	// - gox -osarch="darwin/amd64" ./...
	// - cp ./ios-provisioning-cli_darwin_amd64 $CI_PROJECT_DIR/ios-provisioning

	projectPath := strings.Replace(projectURL, "https://", "http://", 1)
	projectPath = strings.Replace(projectPath, "http:/", goSourcePath, 1)
	projectDirBin := path.Join(projectDir, "bin")
	parent := path.Dir(projectPath)

	if err := goclitools.RunInteractive(fmt.Sprintf("mkdir -p %s", parent)); err != nil {
		return err
	}
	if err := goclitools.RunInteractive(fmt.Sprintf("mkdir -p %s", projectDirBin)); err != nil {
		return err
	}
	if err := goclitools.RunInteractive(fmt.Sprintf("cp -R . %s", projectPath)); err != nil {
		return err
	}

	if err := goclitools.RunInteractiveInDir("go get -d -v ./... && go get github.com/mitchellh/gox", projectPath); err != nil {
		return err
	}

	cmd := fmt.Sprintf("gox -output=\"bin/{{.Dir}}_{{.OS}}_{{.Arch}}\" -osarch=\"%s\" ./... && cp ./bin/* %s", osarch, projectDirBin)
	return goclitools.RunInteractiveInDir(cmd, projectPath)
}
