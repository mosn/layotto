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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch/v5"
	log "github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (i *injector) handleRequest(w http.ResponseWriter, r *http.Request) {
	if err := validateRequest(r); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := readRequestBody(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var patchOps jsonpatch.Patch
	patchedSuccessfully := false

	ar := admissionv1.AdmissionReview{}
	_, gvk, err := i.deserializer.Decode(body, nil, &ar)
	if err != nil {
		log.Errorf("Can't decode body: %v", err)
	} else {
		patchOps, err = i.getPodPatchOperations(r.Context(), &ar)
		if err == nil {
			patchedSuccessfully = true
		}
	}

	var admissionResponse *admissionv1.AdmissionResponse
	if err != nil {
		admissionResponse = errorToAdmissionResponse(err)
		log.Errorf("Sidecar layotto-injector failed to inject. Error: %s", err)
	} else if len(patchOps) == 0 {
		admissionResponse = &admissionv1.AdmissionResponse{
			Allowed: true,
		}
	} else {
		var patchBytes []byte
		patchBytes, err = json.Marshal(patchOps)
		if err != nil {
			admissionResponse = errorToAdmissionResponse(err)
		} else {
			admissionResponse = &admissionv1.AdmissionResponse{
				Allowed: true,
				Patch:   patchBytes,
				PatchType: func() *admissionv1.PatchType {
					pt := admissionv1.PatchTypeJSONPatch
					return &pt
				}(),
			}
		}
	}

	admissionReview := admissionv1.AdmissionReview{
		Response: admissionResponse,
	}
	if admissionResponse != nil && ar.Request != nil {
		admissionReview.Response.UID = ar.Request.UID
		admissionReview.SetGroupVersionKind(*gvk)
	}

	respBytes, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("Sidecar layotto-injector failed to inject. Can't serialize response: %s", err)
		return
	}
	w.Header().Set("Content-Type", runtime.ContentTypeJSON)
	_, err = w.Write(respBytes)
	if err != nil {
		log.Errorf("Sidecar layotto-injector failed to inject. Failed to write response: %v", err)
		return
	}

	if patchedSuccessfully {
		log.Infof("Sidecar layotto-injector succeeded injection.")
	} else {
		log.Errorf("Admission succeeded, but pod was not patched. No sidecar injected.")
	}
}

// errorToAdmissionResponse is a helper function to create an AdmissionResponse
// with an embedded error.
func errorToAdmissionResponse(err error) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

func validateRequest(req *http.Request) error {
	if req.Method != http.MethodPost {
		return fmt.Errorf("wrong http verb. got %s", req.Method)
	}
	if req.Body == nil {
		return errors.New("empty body")
	}
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("wrong content type. expected 'application/json', got: '%s'", contentType)
	}
	return nil
}

func readRequestBody(req *http.Request) ([]byte, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read Request Body: %v", err)
	}
	return body, nil
}
