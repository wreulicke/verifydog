package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	cli "gopkg.in/urfave/cli.v1"
	yaml "gopkg.in/yaml.v2"
)

func verify() error {
	cfg, err := ParseConfig()
	if err != nil {
		return err

	}

	status := map[string]bool{}
	for k, v := range cfg.Verifiers {
		cmd := exec.Command("git", "diff", os.Args[1], os.Args[2], "--", v)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			return err
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
