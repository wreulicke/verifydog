package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	yaml "gopkg.in/yaml.v2"
)

var version string

type Config struct {
	Verifiers map[string]string
}

// gitCommandResult holds the output of a git command execution
type gitCommandResult struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

// executeGitCommand runs a git command and handles verbose logging
func executeGitCommand(verbose bool, args ...string) (*gitCommandResult, error) {
	cmd := exec.Command("git", args...)
	result := &gitCommandResult{}

	cmd.Stdout = &result.stdout
	cmd.Stderr = &result.stderr

	if verbose {
		log.Infof("execute `%s`", strings.Join(cmd.Args, " "))
	}

	if err := cmd.Run(); err != nil {
		if verbose && result.stderr.Len() > 0 {
			errOut := colorable.NewColorableStderr()
			errOut.Write(result.stderr.Bytes())
			fmt.Fprintln(errOut)
		}

		line, _, readErr := bufio.NewReader(&result.stderr).ReadLine()
		if readErr != nil {
			return nil, errors.Wrap(err, "git command failed. cannot get details.")
		}

		return nil, errors.Errorf("git command failed. message: `%s`", string(line))
	}

	return result, nil
}

func history(verbose bool, args []string) error {
	cfg, err := parseConfig()
	if err != nil {
		return err
	}

	for _, v := range cfg.Verifiers {
		fmt.Printf("%s -->\r\n", v)

		result, err := executeGitCommand(verbose, "log", "--color", args[0], args[1], "--", v)
		if err != nil {
			return err
		}

		fmt.Println(result.stdout.String())
		fmt.Printf("<-- %s\r\n", v)
	}

	return nil
}

func verify(verbose bool, args []string) error {
	cfg, err := parseConfig()
	if err != nil {
		return err
	}

	status := map[string]bool{}

	for k, v := range cfg.Verifiers {
		result, err := executeGitCommand(verbose, "diff", "--color", args[0], args[1], "--", v)
		if err != nil {
			return err
		}

		if verbose && result.stdout.Len() > 0 {
			errOut := colorable.NewColorableStderr()
			errOut.Write(result.stdout.Bytes())
			fmt.Fprintln(errOut)
		}

		status[k] = result.stdout.Len() > 0
	}

	b, _ := json.Marshal(status)
	fmt.Println(string(b))

	return nil
}

func parseConfig() (*Config, error) {
	yml, err := os.ReadFile(".verifydog.yml")
	if err != nil {
		return nil, err
	}

	out := &Config{}
	if err := yaml.Unmarshal(yml, out); err != nil {
		return nil, err
	}

	return out, nil
}

func historyAction(ctx context.Context, cmd *cli.Command) error {
	verbose := cmd.Bool("verbose")
	if cmd.Args().Len() != 2 {
		return errors.New("requires 2 only commit reference")
	}

	if verbose {
		log.Info("start verifydog with verbose mode")
	}

	return history(verbose, cmd.Args().Slice())
}

func verifyAction(ctx context.Context, cmd *cli.Command) error {
	verbose := cmd.Bool("verbose")
	if cmd.Args().Len() != 2 {
		return errors.New("require 2 only commit reference")
	}

	if verbose {
		log.Info("start verifydog with verbose mode")
	}

	return verify(verbose, cmd.Args().Slice())
}

func mainInternal() error {
	cmd := &cli.Command{
		Name:    "verifydog",
		Usage:   "verify diff between versions",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "verbose mode",
			},
		},
		Action: verifyAction,
		Commands: []*cli.Command{
			{
				Name:        "history",
				Description: "show history",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "verbose",
						Usage: "verbose mode",
					},
				},
				Action: historyAction,
			},
		},
	}

	return cmd.Run(context.Background(), os.Args)
}

func main() {
	log.SetOutput(colorable.NewColorableStderr())
	f := log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}
	log.SetFormatter(&f)

	err := mainInternal()
	if err != nil {
		log.Fatal(err)
	}
}
