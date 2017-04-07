package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"

	"context"
	fnclient "github.com/cmdhema/functions_go/client"
	apibuild "github.com/cmdhema/functions_go/client/build"
	"github.com/cmdhema/functions_go/models"
)

func build() cli.Command {
	cmd := buildcmd{}
	flags := append([]cli.Flag{}, cmd.flags()...)
	return cli.Command{
		Name:   "build",
		Usage:  "build function version",
		Flags:  flags,
		Action: cmd.build,
	}
}

type buildcmd struct {
	verbose bool
	remote bool

	client *fnclient.Functions

}

func (b *buildcmd) flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:        "v",
			Usage:       "verbose mode",
			Destination: &b.verbose,
		},
		cli.BoolFlag{
			Name :		"remote, r",
			Usage :		"remote mode",
			Destination:	&b.remote,
		},
	}
}

// build will take the found valid function and build it
func (b *buildcmd) build(c *cli.Context) error {

	b.client = apiClient()

	verbwriter := verbwriter(b.verbose)

	path, err := os.Getwd()
	if err != nil {
		return err
	}
	fn, err := findFuncfile(path)

	if err != nil {
		return err
	}

	if b.remote == false {
		fmt.Println("Local build")
		ff, err := buildfunc(verbwriter, fn)

		if err != nil {
			return err
		}

		fmt.Printf("Function %v local built successfully.\n", ff.FullName())
	} else {
		funcfile, err := parsefuncfile(fn)

		fmt.Println("Remote build : " + funcfile.FullName())
		if err != nil {
			return err
		}

		body := &models.BuildWrapper{Build: &models.Build{
			Name:   	funcfile.FullName(),
			Code:		getFuncCode(path),
			Deeplearning:	*funcfile.Deeplearning,
			Entrypoint:	funcfile.Entrypoint,
			Runtime:	*funcfile.Runtime,
		}}

		resp, err := b.client.Build.PostBuild(&apibuild.PostBuildParams{
			Context: context.Background(),
			Body:    body,
		})

		if err != nil {
			fmt.Println("Error Remote Build!!!!!")

			return err
		}

		if resp != nil {
			fmt.Println("Remote build success")
		} else {
			fmt.Println("Error Remote Build")
		}

		fmt.Println("Remote build success!!!")

	}

	return nil
}

