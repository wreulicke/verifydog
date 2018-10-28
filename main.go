package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"

	cli "gopkg.in/urfave/cli.v1"
	yaml "gopkg.in/yaml.v2"
)

var version string

func history(verbose bool, args []string) error {
	cfg, err := parseConfig()
	if err != nil {
		return err
	}

	for _, v := range cfg.Verifiers {
		fmt.Printf("%s -->\r\n", v)
		cmd := exec.Command("git", "log", "--color", args[0], args[1], "--", v)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if verbose {
			log.Infof("execute `%s`", strings.Join(cmd.Args, " "))
		}
		err := cmd.Run()
		if err != nil {
			b := stderr.Bytes()
			if verbose && len(b) != 0 {
				errOut := colorable.NewColorableStderr()
				errOut.Write(b)
				fmt.Fprintln(errOut)
			}
			line, _, readErr := bufio.NewReader(bytes.NewReader(b)).ReadLine()
			if readErr != nil {
				return errors.Wrap(err, "git diff is failed. cannot get deitals.")
			}
			return errors.Errorf("git diff is failed. message: `%s`", string(line))
		}
		fmt.Println(stdout.String())
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
		cmd := exec.Command("git", "diff", "--color", args[0], args[1], "--", v)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if verbose {
			log.Infof("execute `%s`", strings.Join(cmd.Args, " "))
		}
		err := cmd.Run()
		if err != nil {
			b := stderr.Bytes()
			if verbose && len(b) != 0 {
				errOut := colorable.NewColorableStderr()
				errOut.Write(b)
				fmt.Fprintln(errOut)
			}
			line, _, readErr := bufio.NewReader(bytes.NewReader(b)).ReadLine()
			if readErr != nil {
				return errors.Wrap(err, "git diff is failed. cannot get deitals.")
			}
			return errors.Errorf("git diff is failed. message: `%s`", string(line))
		}
		b := stdout.Bytes()
		if verbose && len(b) != 0 {
			errOut := colorable.NewColorableStderr()
			errOut.Write(stdout.Bytes())
			fmt.Fprintln(errOut)
		}
		status[k] = len(b) != 0
	}
	b, _ := json.Marshal(status)
	fmt.Println(string(b))
	return nil
}

type Config struct {
	Verifiers map[string]string
}

func parseConfig() (*Config, error) {
	yml, err := ioutil.ReadFile(".verifydog.yml")
	if err != nil {
		return nil, err
	}
	out := &Config{}
	if err := yaml.Unmarshal(yml, out); err != nil {
		return nil, err
	}
	return out, nil
}

func historyAction(c *cli.Context) error {
	verbose := c.Bool("verbose")
	if len(c.Args()) != 2 {
		return errors.New("requires 2 only commit reference")
	}
	if verbose {
		log.Info("start verifydog with verbose mode")
	}
	return history(verbose, c.Args())
}

func verifyAction(c *cli.Context) error {
	verbose := c.Bool("verbose")
	if len(c.Args()) != 2 {
		return errors.New("require 2 only commit reference")
	}
	if verbose {
		log.Info("start verifydog with verbose mode")
	}
	return verify(verbose, c.Args())
}

func mainInternal() error {
	app := cli.NewApp()
	app.Name = "verifydog"
	app.Usage = "verify diff between versions"
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "verbose mode",
		},
	}
	app.Action = verifyAction
	app.Commands = []cli.Command{
		cli.Command{
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "verbose",
					Usage: "verbose mode",
				},
			},
			Name:        "history",
			Description: "show history",
			Action:      historyAction,
		},
	}
	return app.Run(os.Args)
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
