/*
 * Copyright 2026 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package hosted manages API servers, native services and their dependencies.
// It has no dependency on a transport, product SDK or MOSN host. Business
// interfaces remain with the caller; Runtime owns only their lifecycle.
package hosted

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Module methods must honor context cancellation. Status may run concurrently
// with lifecycle methods and must return a safe snapshot without side effects.
// Start returns after establishing service readiness, not after serving ends.
// Its context bounds the startup operation, not the service lifetime; the module
// owns and cancels background work when it closes.
// Close must support partial initialization and forcibly release owned resources
// after a failed Drain. Neither Drain nor Close may close borrowed dependencies.
type Module interface {
	Init(context.Context) error
	Start(context.Context) error
	Status(context.Context) (Health, error)
	Drain(context.Context) error
	Close(context.Context) error
}

type Health struct {
	Ready  bool
	Reason string
}

// Registration describes one owned instance. New is called at most once.
// Factories must clean up resources they cannot return to the runtime. A non-nil
// module returned together with an error is still owned and closed by Runtime.
type Registration struct {
	Name      string
	DependsOn []string
	Optional  bool
	New       func(context.Context) (Module, error)
}

type Options struct {
	// Rollback uses a fresh context because startup may already be cancelled.
	RollbackTimeout time.Duration
	DrainTimeout    time.Duration
	CloseTimeout    time.Duration
}

type State string

const (
	Created      State = "created"
	Initializing State = "initializing"
	Starting     State = "starting"
	Running      State = "running"
	Failed       State = "failed"
	Draining     State = "draining"
	Closing      State = "closing"
	Closed       State = "closed"
)

type ModuleStatus struct {
	Name           string
	State          State
	Optional       bool
	Ready          bool
	Reason         string
	LastError      string
	CloseAttempted bool // failures remain visible; an attempt is not proof of release
}

type Snapshot struct {
	State   State
	Ready   bool
	Modules []ModuleStatus // stable dependency order
}

// ModuleError retains the module and operation for diagnostics and errors.Is.
type ModuleError struct {
	Name      string
	Operation string
	Err       error
}

func (e *ModuleError) Error() string { return fmt.Sprintf("%s %s: %v", e.Name, e.Operation, e.Err) }
func (e *ModuleError) Unwrap() error { return e.Err }

// Errors preserves all failures, including cleanup errors, on Go 1.18.
type Errors []error

func (e Errors) Error() string {
	parts := make([]string, len(e))
	for i, err := range e {
		parts[i] = err.Error()
	}
	return strings.Join(parts, "; ")
}

func (e Errors) Is(target error) bool {
	for _, err := range e {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
