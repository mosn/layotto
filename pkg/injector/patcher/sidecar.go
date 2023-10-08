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
	"reflect"
	"strconv"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"

	"mosn.io/layotto/pkg/common"
)

// SidecarConfig contains the configuration for the sidecar container.
// Its parameters can be read from annotations on a pod.
type SidecarConfig struct {
	SidecarAPIGRPCPort int32 `default:"34904"`

	Namespace       string
	ImagePullPolicy corev1.PullPolicy

	SidecarInject  bool   `annotation:"layotto/sidecar-inject"`
	SidecarImage   string `annotation:"layotto/sidecar-image"`
	ConfigVolume   string `annotation:"layotto/config-volume"`
	VolumeMounts   string `annotation:"layotto/volume-mounts"`
	VolumeMountsRW string `annotation:"layotto/volume-mounts-rw"`

	pod *corev1.Pod
}

// NewSidecarConfig returns a ContainerConfig object for a pod.
func NewSidecarConfig(pod *corev1.Pod) *SidecarConfig {
	c := &SidecarConfig{
		pod: pod,
	}
	c.setDefaultValues()
	return c
}

func (c *SidecarConfig) setDefaultValues() {
	// Iterate through the fields using reflection
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldT := val.Type().Field(i)
		fieldV := val.Field(i)
		def := fieldT.Tag.Get("default")
		if !fieldV.CanSet() || def == "" {
			continue
		}

		// Assign the default value
		setValueFromString(fieldT.Type, fieldV, def, "")
	}
}

func (c *SidecarConfig) SetFromPodAnnotations() {
	c.setFromAnnotations(c.pod.Annotations)
}

// setFromAnnotations updates the object with properties from an annotation map.
func (c *SidecarConfig) setFromAnnotations(an map[string]string) {
	// Iterate through the fields using reflection
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldV := val.Field(i)
		fieldT := val.Type().Field(i)
		key := fieldT.Tag.Get("annotation")
		if !fieldV.CanSet() || key == "" {
			continue
		}

		// Skip annotations that are not defined or which have an empty value
		if an[key] == "" {
			continue
		}

		// Assign the value
		setValueFromString(fieldT.Type, fieldV, an[key], key)
	}
}

func setValueFromString(rt reflect.Type, rv reflect.Value, val string, key string) bool {
	switch rt.Kind() {
	case reflect.Pointer:
		pt := rt.Elem()
		pv := reflect.New(rt.Elem()).Elem()
		if setValueFromString(pt, pv, val, key) {
			rv.Set(pv.Addr())
		}
	case reflect.String:
		rv.SetString(val)
	case reflect.Bool:
		rv.SetBool(common.StringToBool(val))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			rv.SetInt(v)
		} else {
			log.Warnf("Failed to parse int value from annotation %s (annotation will be ignored): %v", key, err)
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(val, 10, 64)
		if err == nil {
			rv.SetUint(v)
		} else {
			log.Warnf("Failed to parse uint value from annotation %s (annotation will be ignored): %v", key, err)
			return false
		}
	}

	return true
}
