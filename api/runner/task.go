package runner

import (
	"context"
	"io"
	"time"
	//"fmt"
	"log"
	"github.com/fsouza/go-dockerclient"
	"github.com/iron-io/functions/api/runner/task"
	"github.com/iron-io/runner/drivers"
	"github.com/NVIDIA/nvidia-docker/src/nvidia"
)

type containerTask struct {
	ctx    context.Context
	cfg    *task.Config
	canRun chan bool
}

func (t *containerTask) Volumes() [][2]string {
	//vols, err := nvidia_docker.VolumesNeeded(t.cfg.Image)
	//if  err != nil {
	//	volumes := nvidia_docker.VolumesArgs(vols)
	//	fmt.Println(volumes)
	//	devices := nvidia_docker.DevicesArgs()
	//	fmt.Println(devices)
	//}

	//log.Println("Discovering GPU devices")
	//Devices, err = nvidia.LookupDevices()
	//assert(err)
	//
	//log.Println("Provisioning volumes at", "/var/lib/nvidia-docker/volumes")
	//Volumes, err = nvidia.LookupVolumes("/var/lib/nvidia-docker/volumes")
	//assert(err)
	//
	//fmt.Println("Device")
	//fmt.Println(Devices)
	//fmt.Println("Volumes")
	//fmt.Println(Volumes)

	return [][2]string{}
}


func (t *containerTask) Command() string { return "" }

func (t *containerTask) EnvVars() map[string]string {
	return t.cfg.Env
}
func (t *containerTask) Input() io.Reader {
	return t.cfg.Stdin
}

func (t *containerTask) Labels() map[string]string {
	return map[string]string{
		"LogName": t.cfg.AppName,
	}
}

func (t *containerTask) Id() string                         { return t.cfg.ID }
func (t *containerTask) Route() string                      { return "" }
func (t *containerTask) Image() string                      { return t.cfg.Image }
func (t *containerTask) Timeout() time.Duration             { return t.cfg.Timeout }
func (t *containerTask) Logger() (stdout, stderr io.Writer) { return t.cfg.Stdout, t.cfg.Stderr }
//func (t *containerTask) Volumes() [][2]string               { return [][2]string{} }
func (t *containerTask) WorkDir() string                    { return "" }

func (t *containerTask) Close()                 {}
func (t *containerTask) WriteStat(drivers.Stat) {}

// FIXME: for now just use empty creds => public docker hub image
func (t *containerTask) DockerAuth() docker.AuthConfiguration { return docker.AuthConfiguration{} }
