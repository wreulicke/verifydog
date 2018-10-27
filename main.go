package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	cli "gopkg.in/urfave/cli.v1"
)

func verify() error {
	fmt.Println(os.Args)
	path, err := exec.LookPath(os.Args[1])
	if err != nil {
		return err
	}
	return syscall.Exec(path, os.Args[1:], os.Environ())
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
