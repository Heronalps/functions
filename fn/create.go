package main

import (
	"errors"
	"fmt"
	//"os"
	//"path/filepath"

	"strings"

	//"github.com/iron-io/functions/fn/langs"
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

	deeplearning string
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
				Name:		"deeplearning",
				Usage:		"docker image for deeplearning framework",
				Destination:	&a.deeplearning,
			},
			cli.StringFlag{
				Name:		"name",
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
	fmt.Println("func.yaml created.")

	return nil
}