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
	"testing"

	corev1 "k8s.io/api/core/v1"

	injectorConsts "mosn.io/layotto/pkg/injector/consts"
)

func TestInjectRequired(t *testing.T) {
	t.Run("returns true when sidecar injection is enabled and pod does not contain sidecar", func(t *testing.T) {
		config := &SidecarConfig{
			SidecarInject: true,
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "not-sidecar"},
					},
				},
			},
		}

		if !config.injectRequired() {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("returns false when sidecar injection is disabled", func(t *testing.T) {
		config := &SidecarConfig{
			SidecarInject: false,
		}

		if config.injectRequired() {
			t.Errorf("Expected false, got true")
		}
	})

	t.Run("returns false when pod already contains sidecar", func(t *testing.T) {
		config := &SidecarConfig{
			SidecarInject: true,
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: injectorConsts.SidecarContainerName},
					},
				},
			},
		}

		if config.injectRequired() {
			t.Errorf("Expected false, got true")
		}
	})
}

func TestGetPatch(t *testing.T) {
	t.Run("returns nil when sidecar injection is not required", func(t *testing.T) {
		config := &SidecarConfig{
			SidecarInject: false,
		}

		patch, err := config.GetPatch()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if patch != nil {
			t.Errorf("Expected nil, got %v", patch)
		}
	})

	t.Run("returns patch when sidecar injection is required", func(t *testing.T) {
		config := &SidecarConfig{
			SidecarInject: true,
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "not-sidecar"},
					},
				},
			},
		}

		patch, err := config.GetPatch()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(patch) == 0 {
			t.Errorf("Expected patch, got nil or empty")
		}
	})
}

func TestPodContainsSidecarContainer(t *testing.T) {
	t.Run("returns true when pod contains sidecar container", func(t *testing.T) {
		config := &SidecarConfig{
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: injectorConsts.SidecarContainerName},
					},
				},
			},
		}

		if !config.podContainsSidecarContainer() {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("returns false when pod does not contain sidecar container", func(t *testing.T) {
		config := &SidecarConfig{
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "not-sidecar"},
					},
				},
			},
		}

		if config.podContainsSidecarContainer() {
			t.Errorf("Expected false, got true")
		}
	})
}
