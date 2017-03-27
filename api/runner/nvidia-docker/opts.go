// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package nvidia_docker

import (
	//"fmt"

	//"github.com/NVIDIA/nvidia-docker/src/docker"
	//"github.com/NVIDIA/nvidia-docker/src/nvidia"
)

//func GenerateLocalArgs(image string, vols []string) ([]string, error) {
//	cv, err := nvidia.GetCUDAVersion()
//	if err != nil {
//		return nil, err
//	}
//	if err := cudaSupported(image, cv); err != nil {
//		return nil, err
//	}
//
//	d, err := devicesArgs()
//	if err != nil {
//		return nil, err
//	}
//	v, err := volumesArgs(vols)
//	if err != nil {
//		return nil, err
//	}
//	return append(d, v...), nil
//}

//func DevicesArgs() ([]string, error) {
//	var args []string
//
//	cdevs, err := nvidia.GetControlDevicePaths()
//	if err != nil {
//		return nil, err
//	}
//	for i := range cdevs {
//		args = append(args, fmt.Sprintf("--device=%s", cdevs[i]))
//	}
//
//	devs, err := nvidia.LookupDevices(nvidia.LookupMinimal)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(GPU) == 0 {
//		for i := range devs {
//			args = append(args, fmt.Sprintf("--device=%s", devs[i].Path))
//		}
//	} else {
//		devs, err := nvidia.FilterDevices(devs, GPU)
//		if err != nil {
//			return nil, err
//		}
//		for i := range devs {
//			args = append(args, fmt.Sprintf("--device=%s", devs[i].Path))
//		}
//	}
//	return args, nil
//}
//
//func VolumesArgs(vols []string) ([]string, error) {
//	args := make([]string, 0, len(vols))
//
//	drv, err := nvidia.GetDriverVersion()
//	if err != nil {
//		return nil, err
//	}
//	for _, vol := range nvidia.Volumes {
//		for _, v := range vols {
//			if v == vol.Name {
//				// Check if the volume exists locally otherwise fallback to using the plugin
//				n := fmt.Sprintf("%s_%s", vol.Name, drv)
//				if _, err := docker.VolumeInspect(n); err == nil {
//					args = append(args, fmt.Sprintf("--volume=%s:%s:%s", n, vol.Mountpoint, vol.MountOptions))
//				} else {
//					args = append(args, fmt.Sprintf("--volume-driver=%s", nvidia.DockerPlugin))
//					args = append(args, fmt.Sprintf("--volume=%s:%s:%s", n, vol.Mountpoint, vol.MountOptions))
//				}
//				break
//			}
//		}
//	}
//	return args, nil
//}
//
//func enableCUDA(image string) (error) {
//	cv, err := nvidia.GetCUDAVersion()
//	if err != nil {
//		return err
//	}
//	if err := cudaSupported(image, cv); err != nil {
//		return err
//	}
//
//	return err
//}
//
//func cudaSupported(image, version string) error {
//	var vmaj, vmin int
//	var lmaj, lmin int
//
//	label, err := docker.Label(image, labelCUDAVersion)
//	if err != nil {
//		return err
//	}
//	if label == "" {
//		return nil
//	}
//	if _, err := fmt.Sscanf(version, "%d.%d", &vmaj, &vmin); err != nil {
//		return err
//	}
//	if _, err := fmt.Sscanf(label, "%d.%d", &lmaj, &lmin); err != nil {
//		return err
//	}
//	if lmaj > vmaj || (lmaj == vmaj && lmin > vmin) {
//		return fmt.Errorf("unsupported CUDA version: driver %s < image %s", version, label)
//	}
//	return nil
//}
