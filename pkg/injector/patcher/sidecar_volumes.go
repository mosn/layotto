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
	"strings"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"

	"mosn.io/layotto/pkg/injector/consts"
)

// getVolumeMounts returns the list of VolumeMount's for the sidecar container.
func (c *SidecarConfig) getVolumeMounts() []corev1.VolumeMount {
	vs := append(
		parseVolumeMountsString(c.VolumeMounts, true),
		parseVolumeMountsString(c.VolumeMountsRW, false)...,
	)

	vs = append(vs, getConfigVolumeMount(c.ConfigVolume, true))

	volumeMounts := make([]corev1.VolumeMount, 0)
	for _, v := range vs {
		if podContainsVolume(c.pod, v.Name) {
			volumeMounts = append(volumeMounts, v)
		} else {
			log.Warnf("Volume %s is not present in pod %s, skipping", v.Name, c.pod.GetName())
		}
	}

	return volumeMounts
}

func podContainsVolume(pod *corev1.Pod, name string) bool {
	for _, volume := range pod.Spec.Volumes {
		if volume.Name == name {
			return true
		}
	}
	return false
}

// parseVolumeMountsString parses the annotation and returns volume mounts.
// The format of the annotation is: "mountPath1:hostPath1,mountPath2:hostPath2"
// The readOnly parameter applies to all mounts.
func parseVolumeMountsString(volumeMountStr string, readOnly bool) []corev1.VolumeMount {
	vs := strings.Split(volumeMountStr, ",")
	volumeMounts := make([]corev1.VolumeMount, 0, len(vs))
	for _, v := range vs {
		vmount := strings.Split(strings.TrimSpace(v), ":")
		if len(vmount) != 2 {
			continue
		}
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      vmount[0],
			MountPath: vmount[1],
			ReadOnly:  readOnly,
		})
	}
	return volumeMounts
}

// getConfigVolumeMount returns the layotto config volume mount.
// Currently, the path of the Layotto configuration file is "/runtime/configs/config.json"
func getConfigVolumeMount(configVolumeName string, readOnly bool) corev1.VolumeMount {
	volumeMount := corev1.VolumeMount{
		Name:      configVolumeName,
		MountPath: consts.LayottoConfigPath,
		ReadOnly:  readOnly,
	}
	return volumeMount
}
