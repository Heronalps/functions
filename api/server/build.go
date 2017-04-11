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

	driverscommon "github.com/cmdhema/runner/drivers"
	"github.com/cmdhema/runner/drivers/docker"
	"github.com/iron-io/functions/api/models"
	"github.com/cmdhema/runner/common"
	"github.com/gin-gonic/gin"
	"bytes"
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
	"theano":		"kaixhin/cuda-theano:8.0",
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

	c.JSON(http.StatusOK, buildResponse{"Successfully remote build", build})

}

func writeTmpDockerfile(dir string, build *models.Build) error {

	fmt.Println("Write dockerfile")
	filename := strings.Split(build.Entrypoint, " ")[1]
	writeFuncCodeToFile(build.Code, filename)

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