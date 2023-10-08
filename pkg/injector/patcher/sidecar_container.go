// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package patcher

import (
	corev1 "k8s.io/api/core/v1"

	injectorConsts "mosn.io/layotto/pkg/injector/consts"
)

type getSidecarContainerOpts struct {
	VolumeMounts []corev1.VolumeMount
}

// getSidecarContainer returns the Container object for the sidecar.
func (c *SidecarConfig) getSidecarContainer(opts getSidecarContainerOpts) (*corev1.Container, error) {
	ports := []corev1.ContainerPort{
		{
			ContainerPort: c.SidecarAPIGRPCPort,
			Name:          injectorConsts.SidecarGRPCPortName,
		},
	}
	// Get the command (/layotto)
	cmd := []string{"/runtime/layotto", "start"}
	args := []string{"-c", "configs/config.json"}

	// Create the container object
	container := &corev1.Container{
		Name:            injectorConsts.SidecarContainerName,
		Image:           c.SidecarImage,
		ImagePullPolicy: c.ImagePullPolicy,
		Ports:           ports,
		Command:         cmd,
		Args:            args,
		Env: []corev1.EnvVar{
			{
				Name:  "NAMESPACE",
				Value: c.Namespace,
			},
			{
				Name: "POD_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
		},
		VolumeMounts: opts.VolumeMounts,
	}

	return container, nil
}
