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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
)

func TestHandleRequest(t *testing.T) {
	i, err := NewInjector(Config{
		TLSCertFile:  "test-cert",
		TLSKeyFile:   "test-key",
		SidecarImage: "test-image",
		Namespace:    "test-ns",
	})
	assert.NoError(t, err)

	injector := i.(*injector)

	podBytes, _ := json.Marshal(corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test-app",
			Namespace:         "default",
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Annotations: map[string]string{
				"layotto/sidecar-inject": "true",
				"layotto/config-volume":  "layotto-config-vol",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
					Ports: []corev1.ContainerPort{
						{ContainerPort: 80},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "layotto-config-vol",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "layotto-config",
							},
						},
					},
				},
			},
		},
	})

	testCases := []struct {
		testName         string
		request          admissionv1.AdmissionReview
		contentType      string
		expectStatusCode int
		expectPatched    bool
	}{
		{
			"TestSidecarInjectSuccess",
			admissionv1.AdmissionReview{
				Request: &admissionv1.AdmissionRequest{
					UID:       uuid.NewUUID(),
					Kind:      metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
					Name:      "test-app",
					Namespace: "test-ns",
					Operation: "CREATE",
					Object:    runtime.RawExtension{Raw: podBytes},
				},
			},
			runtime.ContentTypeJSON,
			http.StatusOK,
			true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(injector.handleRequest))
	defer ts.Close()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			requestBytes, err := json.Marshal(tc.request)
			assert.NoError(t, err)

			resp, err := http.Post(ts.URL, tc.contentType, bytes.NewBuffer(requestBytes))
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectStatusCode, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var ar admissionv1.AdmissionReview
				err = json.Unmarshal(body, &ar)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectPatched, len(ar.Response.Patch) > 0)
			}
		})
	}
}
