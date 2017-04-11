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

func inject() cli.Command {
	cmd := injectcmd{}
	return cli.Command{
		Name:   "inject",
		Usage:  "inject function to running container",
		Action: cmd.inject,
	}
}

type injectcmd struct {
	verbose bool
	remote bool

	client *fnclient.Functions

}

// build will take the found valid function and build it
func (b *injectcmd) inject(c *cli.Context) error {

	b.client = apiClient()

	path, err := os.Getwd()
	if err != nil {
		return err
	}
	fn, err := findFuncfile(path)

	if err != nil {
		return err
	}

	funcfile, err := parsefuncfile(fn)

	if err != nil {
		return err
	}

	body := &models.BuildWrapper{Build: &models.Build{
		Name:   	funcfile.FullName(),
		Code:		getFuncCode(path),
		Deeplearning:	*funcfile.Deeplearning,
		Entrypoint:	funcfile.Entrypoint,
		Runtime:	*funcfile.Runtime,
		FileName:	*funcfile.FileName,
	}}

	resp, err := b.client.Inject.PostInject(&apibuild.PostBuildParams{
		Context: context.Background(),
		Body:    body,
	})

	if err != nil {
		fmt.Println("Inject function error")
		return err
	}

	if resp != nil {
		fmt.Println("Inject function success!!!")
	}

	return nil
}

