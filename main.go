package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"

	cli "gopkg.in/urfave/cli.v1"
	yaml "gopkg.in/yaml.v2"
)

func verify(verbose bool, args []string) error {
	cfg, err := ParseConfig()
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
		err := cmd.Run()
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, stderr.String())
			}
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, stdout.String())
		}
		status[k] = len(stdout.String()) != 0
	}
	b, _ := json.Marshal(status)
	fmt.Println(string(b))
	return nil
}

type Config struct {
	Verifiers map[string]string
}

func ParseConfig() (*Config, error) {
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

func mainInternal() error {
	app := cli.NewApp()
	app.Name = "verifydog"
	app.Usage = "verify diff between versions"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "show git command output",
		},
	}
	app.Action = func(c *cli.Context) error {
		verbose := c.Bool("verbose")
		if len(c.Args()) != 2 {
			return errors.New("requires 2 commit reference")
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "start verifydog with verbose mode")
		}
		return verify(verbose, c.Args())
	}
	return app.Run(os.Args)
}

func main() {
	err := mainInternal()
	if err != nil {
		log.Fatal(err)
	}
}
