package server

import (
	"fmt"
	"context"
	"net/http"

	driverscommon "github.com/cmdhema/runner/drivers"
	"github.com/cmdhema/runner/drivers/docker"
	"github.com/iron-io/functions/api/models"
	"github.com/cmdhema/runner/common"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleInject(c *gin.Context) {

	fmt.Println("API server handleInject")
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
		return
	}

	fmt.Println(build.Build.Name)

	env := common.NewEnvironment(func(e *common.Environment) {})
	driver := docker.NewDocker(env, *(&driverscommon.Config{}))
	err = driver.Upload(build.Build.Code)

	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, injectResponse{"Successfully injected function to container", build})

}
