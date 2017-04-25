package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"

	"context"
	fnclient "github.com/cmdhema/functions_go/client"
	apibuild "github.com/cmdhema/functions_go/client/build"
	"github.com/cmdhema/functions_go/models"
	"net/http"
	"bytes"
	"mime/multipart"
	"path/filepath"
	"io"
	"log"
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
	all string
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
		cli.StringFlag{
			Name:		"all, a",
			Usage:		"build all files",
			Destination:	&b.all,
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
	} else if b.remote == true && b.all == "" {
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

	} else if b.remote == true && b.all != "" {
		funcfile, err := parsefuncfile(fn)
		extraParams := map[string]string{
			"Name":		funcfile.FullName(),
			"Deeplearning":	*funcfile.Deeplearning,
			"Entrypoint":	funcfile.Entrypoint,
			"Runtime":	*funcfile.Runtime,
		}

		request, err := newfileUploadRequest("http://192.168.0.11:8080/v1/builds", extraParams, "file", b.all)
		if err != nil {
			log.Fatal(err)
		}
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		} else {
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Header)

			fmt.Println(body)
		}
	}

	return nil
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}