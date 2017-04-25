package server

import (
	"fmt"
	"context"
	"net/http"
	"io"
	"os"
	"strings"
	"errors"
	"path/filepath"
	"text/template"
	"bytes"
	"archive/tar"
	"compress/gzip"

	driverscommon "github.com/cmdhema/runner/drivers"
	"github.com/cmdhema/runner/drivers/docker"
	"github.com/iron-io/functions/api/models"
	"github.com/cmdhema/runner/common"
	"github.com/gin-gonic/gin"
)
var Path = "/home/taejoon/kjwook/iron/function/"

const tplDockerfile = `FROM {{ .BaseImage }}
WORKDIR /function
ADD . /function/
{{ if ne .Entrypoint "" }} ENTRYPOINT [{{ .Entrypoint }}] {{ end }}
{{ if ne .Cmd "" }} CMD [{{ .Cmd }}] {{ end }}
`

var acceptableDeeplearningImages = map[string]string {
	"tensorflow":		"tensorflow/tensorflow:latest-devel-gpu",
	"theano":		"cmdhema/cuda-theano:8.0",
	"torch":		"cmdhema/cuda-torch:8.0",
}

func (s *Server) handleBuild(c *gin.Context) {

	fmt.Println("API server handleBuild")
	ctx := c.MustGet("ctx").(context.Context)
	log := common.Logger(ctx)

	var build models.BuildWrapper
	err := c.BindJSON(&build)
	if err != nil {
		log.WithError(err).Debug(models.ErrInvalidJSON)
		c.JSON(http.StatusBadRequest, simpleError(models.ErrInvalidJSON))
		return
	}

	if build.Build == nil {
		log.Debug(models.ErrAppsMissingNew)
		//c.JSON(http.StatusBadRequest, simpleError(models.ErrAppsMissingNew))
		return
	}

	writeTmpDockerfile(Path, build.Build)
	env := common.NewEnvironment(func(e *common.Environment) {})
	driver := docker.NewDocker(env, *(&driverscommon.Config{}))
	err = driver.Build(build.Build.Name)
	//fmt.Println(driver.Exec())
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, buildResponse{"Successfully remote build", build.Build})

}

func (s *Server) handleBuilds(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	log := common.Logger(ctx)

	fmt.Println("handling req...")
	title := c.PostForm("title")

	file, header, err  := c.Request.FormFile("file")
	fmt.Println(title)
	fmt.Println(header.Filename)
	fmt.Println(file)
	out, err := os.Create(Path + header.Filename)
	if err != nil {
		return
	}
	fmt.Println(Path + header.Filename + " Create success!")
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Debug(err)
		return
	}
	fmt.Println("Upload success")

	err = unTarFuncFiles(Path + header.Filename,Path)
	if err != nil {
		log.Debug(err)
	}
	var buildParams = models.Build {
		Name:   	c.PostForm("Name"),
		Deeplearning:	c.PostForm("Deeplearning"),
		Entrypoint:	c.PostForm("Entrypoint"),
		Runtime:	c.PostForm("Runtime"),
	}

	writeTmpDockerfile(Path, &buildParams)
	env := common.NewEnvironment(func(e *common.Environment) {})
	driver := docker.NewDocker(env, *(&driverscommon.Config{}))
	err = driver.Build(buildParams.Name)
	//fmt.Println(driver.Exec())
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, buildResponse{"Successfully remote build", &buildParams})
}

func unTarFuncFiles(src string, dest string) error{
	f, _ := os.Open(src)
	gzr, err := gzip.NewReader(f)
	defer gzr.Close()
	if err != nil {
		return err
	}

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

			// return any other error
		case err != nil:
			return err

			// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}
		// the target location where the dir/file should be created
		target := filepath.Join(dest, header.Name)
		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

			// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}

func writeTmpDockerfile(dir string, build *models.Build) error {

	fmt.Println("Write dockerfile")
	filename := strings.Split(build.Entrypoint, " ")[1]

	if build.Code != "" {
		writeFuncCodeToFile(build.Code, filename)
	}

	if build.Entrypoint == "" && build.Cmd == "" {
		return errors.New("entrypoint and cmd are missing, you must provide one or the other")
	}

	deepLearning := build.Deeplearning
	rt, ok := acceptableDeeplearningImages[deepLearning]
	if !ok {
		return fmt.Errorf("cannot use deeplearning framework %s", deepLearning)
	}

	fd, err := os.Create(filepath.Join(dir, "Dockerfile"))
	if err != nil {
		return err
	}
	defer fd.Close()

	// convert entrypoint string to slice
	bufferEp := stringToSlice(build.Entrypoint)
	bufferCmd := stringToSlice(build.Cmd)

	t := template.Must(template.New("Dockerfile").Parse(tplDockerfile))
	err = t.Execute(fd, struct {
		BaseImage, Entrypoint, Cmd string
	}{rt, bufferEp.String(), bufferCmd.String()})

	return nil
}

func writeFuncCodeToFile(code string, filename string) error {
	fo, err := os.Create(Path + filename)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(code))
	if err != nil {
		return err
	}

	return nil
}

func stringToSlice(in string) bytes.Buffer {
	epvals := strings.Fields(in)
	var buffer bytes.Buffer
	for i, s := range epvals {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString("\"")
		buffer.WriteString(s)
		buffer.WriteString("\"")
	}
	return buffer
}