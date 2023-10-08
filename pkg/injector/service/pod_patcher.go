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

package service

import (
	"context"
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"
	log "github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"

	"mosn.io/layotto/pkg/injector/patcher"
)

func (i *injector) getPodPatchOperations(ctx context.Context, ar *admissionv1.AdmissionReview) (patchOps jsonpatch.Patch, err error) {
	pod := &corev1.Pod{}
	err = json.Unmarshal(ar.Request.Object.Raw, pod)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal raw object: %w", err)
	}
	log.Infof(
		"AdmissionReview for Kind=%v, Namespace=%s Name=%s (%s) UID=%v patchOperation=%v UserInfo=%v",
		ar.Request.Kind, ar.Request.Namespace, ar.Request.Name, pod.Name, ar.Request.UID, ar.Request.Operation, ar.Request.UserInfo,
	)

	// Create the sidecar configuration object from the pod
	sidecar := patcher.NewSidecarConfig(pod)
	sidecar.Namespace = ar.Request.Namespace
	sidecar.ImagePullPolicy = i.config.GetPullPolicy()

	// Default value for the sidecar image, which can be overridden by annotations
	sidecar.SidecarImage = i.config.SidecarImage

	// Set the configuration from annotations
	sidecar.SetFromPodAnnotations()

	// Get the patch to apply to the pod
	// Patch may be empty if there's nothing that needs to be done
	patchOps, err = sidecar.GetPatch()
	if err != nil {
		return nil, err
	}

	return patchOps, nil
}
