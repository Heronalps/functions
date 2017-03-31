package main

import (
	"errors"
	"fmt"
	"os"

	"strings"

	"github.com/iron-io/functions/fn/langs"
	"github.com/urfave/cli"
)

func init() {
	for rt := range fileExtToRuntime {
		fnInitRuntimes = append(fnInitRuntimes, rt)
	}
}

type createFnCmd struct {
	name           string
	force          bool
	runtime        string
	entrypoint     string
	format         string
	maxConcurrency int
	cmd		string
	deeplearning 	string
}

func createFn() cli.Command {
	a := createFnCmd{}

	return cli.Command {
		Name: "create",
		Usage: "create a local func.yaml file",
		Description: "Creates a func.yaml file in the current directory.",
		Action: a.create,
		Flags: []cli.Flag {
			cli.BoolFlag{
				Name:        "force, f",
				Usage:       "overwrite existing func.yaml",
				Destination: &a.force,
			},
			cli.StringFlag{
				Name:        "runtime",
				Usage:       "choose an existing runtime - " + strings.Join(fnInitRuntimes, ", "),
				Destination: &a.runtime,
			},
			cli.StringFlag{
				Name:        "entrypoint",
				Usage:       "entrypoint is the command to run to start this function - equivalent to Dockerfile ENTRYPOINT.",
				Destination: &a.entrypoint,
			},
			cli.StringFlag{
				Name:        "format",
				Usage:       "hot function IO format - json or http",
				Destination: &a.format,
				Value:       "",
			},
			cli.IntFlag{
				Name:        "max-concurrency",
				Usage:       "maximum concurrency for hot function",
				Destination: &a.maxConcurrency,
				Value:       1,
			},
			cli.StringFlag{
				Name:		"deeplearning, dl",
				Usage:		"docker image for deeplearning framework",
				Destination:	&a.deeplearning,
			},
			cli.StringFlag{
				Name:		"name, n",
				Usage:		"function name",
				Destination:	&a.name,
			},
		},
	}
}

func (a *createFnCmd) create(c *cli.Context) error {

	if a.name == "" {
		return errors.New("name can't be null")
	}

	if  a.deeplearning == ""  {
		return errors.New("deeplearning can't be null")
	}

	if !a.force {
		ff, err := loadFuncfile()
		if _, ok := err.(*notFoundError); !ok && err != nil {
			return err
		}
		if ff != nil {
			return errors.New("function file already exists")
		}
	}

	err := a.buildFuncFile(c)
	if err != nil {
		return err
	}

	var ffmt *string
	if a.format != "" {
		ffmt = &a.format
	}

	ff := &funcfile{
		Name:           a.name,
		Runtime:        &a.runtime,
		Version:        initialVersion,
		Entrypoint:     a.entrypoint,
		Cmd:            a.cmd,
		Format:         ffmt,
		MaxConcurrency: &a.maxConcurrency,
		Deeplearning:	&a.deeplearning,
	}
	path := "/" + a.name
	ff.Path = &path

	if err := encodeFuncfileYAML("func.yaml", ff); err != nil {
		return err
	}

	fmt.Println("func.yaml created.")

	return nil
}

func (a *createFnCmd) buildFuncFile(c *cli.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error detecting current working directory: %s\n", err)
	}

	if exists("Dockerfile") {
		fmt.Println("Dockerfile found, will use that to build.")
		return nil
	}

	var rt string
	if a.runtime == "" {
		rt, err = detectRuntime(pwd)
		if err != nil {
			return err
		}
		a.runtime = rt
		fmt.Printf("assuming %v runtime\n", rt)
	}

	if _, ok := acceptableFnRuntimes[a.runtime]; !ok {
		return fmt.Errorf("init does not support the %s runtime, you'll have to create your own Dockerfile for this function", a.runtime)
	}

	helper := langs.GetLangHelper(a.runtime)
	if helper == nil {
		fmt.Printf("No helper found for %s runtime, you'll have to pass in the appropriate flags or use a Dockerfile.", a.runtime)
	}

	if a.entrypoint == "" {
		if helper != nil {
			a.entrypoint = helper.Entrypoint()
		}
	}
	if a.cmd == "" {
		if helper != nil {
			a.cmd = helper.Cmd()
		}
	}
	if a.entrypoint == "" && a.cmd == "" {
		return fmt.Errorf("could not detect entrypoint or cmd for %v, use --entrypoint and/or --cmd to set them explicitly", a.runtime)
	}

	return nil
}