package runner

import (
	"context"
	"io"
	"time"
	"fmt"
	"log"
	"github.com/fsouza/go-dockerclient"
	"github.com/iron-io/functions/api/runner/task"
	"github.com/cmdhema/runner/drivers"
	"github.com/cmdhema/nvidia-docker/src/nvidia"
)

var (
	VolumesPath string = "/var/lib/nvidia-docker/volumes"
)
type containerTask struct {
	ctx    context.Context
	cfg    *task.Config
	canRun chan bool
}

func assert(err error) {
	if err != nil {
		log.Panicln("Error:", err)
	}
}

func (t *containerTask) Volumes() [][2]string {

	log.Println("Loading NVIDIA unified memory")
	assert(nvidia.LoadUVM())

	log.Println("Loading NVIDIA management library")
	assert(nvidia.Init())
	defer func() { assert(nvidia.Shutdown()) }()

	Volumes, err := nvidia.LookupVolumes(VolumesPath)
	assert(err)

	volumes := [][2]string{}
	for _, vol := range Volumes {
		hostDir := fmt.Sprintf("%s_%s", vol.VolumeInfo.Name, vol.Version)
		containerDir := fmt.Sprintf("%s:%s", vol.VolumeInfo.Mountpoint, vol.VolumeInfo.MountOptions)
		volumes = append(volumes, [2]string{hostDir, containerDir})

	}

	return volumes
}

func (t * containerTask) Devices() [][3]string {

	log.Println("Loading NVIDIA unified memory")
	assert(nvidia.LoadUVM())

	log.Println("Loading NVIDIA management library")
	assert(nvidia.Init())
	defer func() { assert(nvidia.Shutdown()) }()

	Devices, err := nvidia.LookupDevices()
	assert(err)

	devices := [][3]string{}
	for _, device := range Devices {
		if device.NVMLDevice != nil {
			devices = append(devices, [3]string{device.NVMLDevice.Path,device.NVMLDevice.Path,"rwm"})
		} else {
			fmt.Println("CUDA")
		}
	}
	ControlDevices, err := nvidia.GetControlDevicePaths()
	for i := range ControlDevices {
		fmt.Println(ControlDevices[i])
		devices = append(devices, [3]string{ControlDevices[i], ControlDevices[i],"rwm"})
	}

	return devices
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
