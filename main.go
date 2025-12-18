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

// gitCommandResult holds the output of a git command execution.
type gitCommandResult struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

// executeGitCommand runs a git command and handles verbose logging.
func executeGitCommand(ctx context.Context, verbose bool, args ...string) (*gitCommandResult, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	result := &gitCommandResult{}

	cmd.Stdout = &result.stdout
	cmd.Stderr = &result.stderr

	if verbose {
		log.Infof("execute `%s`", strings.Join(cmd.Args, " "))
	}

	err := cmd.Run()
	if err != nil {
		if verbose && result.stderr.Len() > 0 {
			errOut := colorable.NewColorableStderr()

			_, writeErr := errOut.Write(result.stderr.Bytes())
			if writeErr != nil {
				log.Warnf("failed to write stderr: %v", writeErr)
			}

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

func history(ctx context.Context, verbose bool, args []string) error {
	cfg, err := parseConfig()
	if err != nil {
		return err
	}

	for _, v := range cfg.Verifiers {
		log.Infof("%s -->", v)

		result, err := executeGitCommand(ctx, verbose, "log", "--color", args[0], args[1], "--", v)
		if err != nil {
			return err
		}

		log.Info(result.stdout.String())
		log.Infof("<-- %s", v)
	}

	return nil
}

func verify(ctx context.Context, verbose bool, args []string) error {
	cfg, err := parseConfig()
	if err != nil {
		return err
	}

	status := map[string]bool{}

	for k, v := range cfg.Verifiers {
		result, err := executeGitCommand(ctx, verbose, "diff", "--color", args[0], args[1], "--", v)
		if err != nil {
			return err
		}

		if verbose && result.stdout.Len() > 0 {
			errOut := colorable.NewColorableStderr()

			_, writeErr := errOut.Write(result.stdout.Bytes())
			if writeErr != nil {
				log.Warnf("failed to write stdout: %v", writeErr)
			}

			fmt.Fprintln(errOut)
		}

		status[k] = result.stdout.Len() > 0
	}

	b, err := json.Marshal(status)
	if err != nil {
		return errors.Wrap(err, "failed to marshal status")
	}

	log.Info(string(b))

	return nil
}

func parseConfig() (*Config, error) {
	yml, err := os.ReadFile(".verifydog.yml")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	out := &Config{}

	err = yaml.Unmarshal(yml, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
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

	return history(ctx, verbose, cmd.Args().Slice())
}

func verifyAction(ctx context.Context, cmd *cli.Command) error {
	verbose := cmd.Bool("verbose")
	if cmd.Args().Len() != 2 {
		return errors.New("require 2 only commit reference")
	}

	if verbose {
		log.Info("start verifydog with verbose mode")
	}

	return verify(ctx, verbose, cmd.Args().Slice())
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

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		return errors.Wrap(err, "command execution failed")
	}

	return nil
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
