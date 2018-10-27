package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	cli "gopkg.in/urfave/cli.v1"
)

func verify() error {
	path, err := exec.LookPath("git")
	if err != nil {
		return err
	}
	args := []string{
		"git",
		"diff",
		os.Args[1],
		os.Args[2],
	}
	return syscall.Exec(path, args, os.Environ())
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
