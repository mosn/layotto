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
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const (
	port = 8443
)

// Injector is the interface for the layotto runtime sidecar injection component.
type Injector interface {
	Run() error
}

type injector struct {
	config       Config
	deserializer runtime.Decoder
	server       *http.Server
}

// Run implements Injector.
func (i *injector) Run() error {
	log.Info("Server started on http://localhost:8443")
	return i.server.ListenAndServeTLS(i.config.TLSCertFile, i.config.TLSKeyFile)
}

// NewInjector returns a new Injector instance.
func NewInjector(config Config) (Injector, error) {
	mux := http.NewServeMux()

	i := &injector{
		config: config,
		deserializer: serializer.NewCodecFactory(
			runtime.NewScheme(),
		).UniversalDeserializer(),
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			ReadHeaderTimeout: 10 * time.Second,
		},
	}

	mux.HandleFunc("/mutate", i.handleRequest)
	return i, nil
}
