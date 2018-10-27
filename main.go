package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"

	cli "gopkg.in/urfave/cli.v1"
	yaml "gopkg.in/yaml.v2"
)

func verify() error {
	cfg, err := ParseConfig()
	if err != nil {
		return err

	}

	for _, v := range cfg.Verfiers {
		path, err := exec.LookPath("git")
		if err != nil {
			return err
		}
		args := []string{
			"git",
			"diff",
			os.Args[1],
			os.Args[2],
			"--",
			v,
		}
		err = syscall.Exec(path, args, os.Environ())
		if err != nil {
			return err
		}
	}
	return nil
}

type Config struct {
	Verfiers map[string]string
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
	app.Action = func(c *cli.Context) error {
		return verify()
	}
	return app.Run(os.Args)
}

func main() {
	err := mainInternal()
	if err != nil {
		log.Fatal(err)
	}
}
