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
	jsonpatch "github.com/evanphx/json-patch/v5"

	injectorConsts "mosn.io/layotto/pkg/injector/consts"
)

// injectRequired check if the sidecar should be injected
func (c *SidecarConfig) injectRequired() bool {
	return c.SidecarInject && !c.podContainsSidecarContainer()
}

// GetPatch returns the patch to apply to a Pod to inject the Layotto sidecar
func (c *SidecarConfig) GetPatch() (patchOps jsonpatch.Patch, err error) {
	// If Layotto is not enabled, or if the layotto container is already present, return
	if !c.injectRequired() {
		return nil, nil
	}

	patchOps = jsonpatch.Patch{}

	// Get volume mounts
	volumeMounts := c.getVolumeMounts()

	// Get the sidecar container
	sidecarContainer, err := c.getSidecarContainer(getSidecarContainerOpts{
		VolumeMounts: volumeMounts,
	})
	if err != nil {
		return nil, err
	}

	patchOps = append(patchOps,
		NewPatchOperation("add", PatchPathContainers+"/-", sidecarContainer),
	)

	return patchOps, nil
}

// podContainsSidecarContainer returns true if the pod contains a sidecar container (i.e. a container named "layotto").
func (c *SidecarConfig) podContainsSidecarContainer() bool {
	for _, c := range c.pod.Spec.Containers {
		if c.Name == injectorConsts.SidecarContainerName {
			return true
		}
	}
	return false
}
