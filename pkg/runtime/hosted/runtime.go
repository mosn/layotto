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

package hosted

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrStopped = errors.New("runtime has been stopped")

type entry struct {
	registration Registration
	module       Module
	state        State
	err          error
	started      bool // Start attempted; Drain may be needed even if it failed.
	closed       bool
	retainsDeps  bool
}

// Runtime is single-use. Start and Shutdown are idempotent; a failed startup
// cannot be retried on partially closed objects. Create a new Runtime to retry.
type Runtime struct {
	mu       sync.RWMutex
	gate     chan struct{}
	entries  []*entry
	options  Options
	state    State
	stop     bool
	cancel   context.CancelFunc
	startErr error
	stopErr  error
	stopped  bool
}

// New validates the complete graph before any factory or other side effect.
func New(registrations []Registration, options Options) (*Runtime, error) {
	if options.RollbackTimeout < 0 || options.DrainTimeout < 0 || options.CloseTimeout < 0 {
		return nil, errors.New("runtime timeouts must not be negative")
	}
	if options.RollbackTimeout == 0 {
		options.RollbackTimeout = 5 * time.Second
	}
	if options.DrainTimeout == 0 {
		options.DrainTimeout = 30 * time.Second
	}
	if options.CloseTimeout == 0 {
		options.CloseTimeout = 5 * time.Second
	}
	byName := make(map[string]Registration, len(registrations))
	for _, reg := range registrations {
		if reg.Name == "" || reg.Name != strings.TrimSpace(reg.Name) || reg.New == nil {
			return nil, fmt.Errorf("invalid registration %q", reg.Name)
		}
		if _, exists := byName[reg.Name]; exists {
			return nil, fmt.Errorf("duplicate module %q", reg.Name)
		}
		reg.DependsOn = append([]string(nil), reg.DependsOn...)
		sort.Strings(reg.DependsOn)
		byName[reg.Name] = reg
	}
	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)
	visited := make(map[string]int)
	var ordered []*entry
	var visit func(string) error
	visit = func(name string) error {
		reg, exists := byName[name]
		if !exists {
			return fmt.Errorf("missing dependency %q", name)
		}
		if visited[name] == 1 {
			return fmt.Errorf("dependency cycle at %q", name)
		}
		if visited[name] == 2 {
			return nil
		}
		visited[name] = 1
		for _, dep := range reg.DependsOn {
			if err := visit(dep); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
		visited[name] = 2
		ordered = append(ordered, &entry{registration: reg, state: Created})
		return nil
	}
	for _, name := range names {
		if err := visit(name); err != nil {
			return nil, err
		}
	}
	return &Runtime{gate: make(chan struct{}, 1), entries: ordered, options: options, state: Created}, nil
}

func (r *Runtime) acquire(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case r.gate <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Runtime) set(e *entry, state State, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e.state = state
	if err != nil {
		e.err = err
	}
}

func (r *Runtime) Start(ctx context.Context) error {
	if err := r.acquire(ctx); err != nil {
		return err
	}
	defer func() { <-r.gate }()
	r.mu.Lock()
	if r.stop {
		r.mu.Unlock()
		return ErrStopped
	}
	if r.state == Running {
		r.mu.Unlock()
		return nil
	}
	if r.state == Failed {
		err := r.startErr
		r.mu.Unlock()
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	r.cancel, r.state = cancel, Starting
	r.mu.Unlock()
	defer func() { cancel(); r.mu.Lock(); r.cancel = nil; r.mu.Unlock() }()

	available := make(map[string]bool)
	for _, e := range r.entries {
		err := ctx.Err()
		if err == nil {
			for _, dep := range e.registration.DependsOn {
				if !available[dep] {
					err = fmt.Errorf("dependency %s unavailable", dep)
					break
				}
			}
		}
		if err == nil {
			err = r.startOne(ctx, e)
		} else {
			err = &ModuleError{e.registration.Name, "dependencies", err}
		}
		if err == nil {
			available[e.registration.Name] = true
			continue
		}
		r.set(e, Failed, err)
		if e.registration.Optional && ctx.Err() == nil {
			cleanup, cancelCleanup := context.WithTimeout(context.Background(), r.options.RollbackTimeout)
			cleanupErr := r.cleanup(cleanup, []*entry{e})
			cancelCleanup()
			if cleanupErr == nil {
				continue
			}
			return r.failStart(Errors{err, cleanupErr}, cleanupErr)
		}
		return r.failStart(err, nil)
	}
	r.mu.Lock()
	// Shutdown may have cancelled startup while the last module returned.
	if ctx.Err() != nil || r.stop {
		err := ctx.Err()
		if err == nil {
			err = context.Canceled
		}
		r.mu.Unlock()
		return r.failStart(err, nil)
	}
	r.state = Running
	r.mu.Unlock()
	return nil
}

func nilModule(m Module) bool {
	if m == nil {
		return true
	}
	v := reflect.ValueOf(m)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Func, reflect.Interface, reflect.Slice, reflect.Chan:
		return v.IsNil()
	}
	return false
}

func (r *Runtime) startOne(ctx context.Context, e *entry) error {
	fail := func(op string, err error) error { return &ModuleError{e.registration.Name, op, err} }
	m, err := e.registration.New(ctx)
	if !nilModule(m) {
		// A shared pointer must be registered once and borrowed by its users.
		// Reject aliases before taking ownership a second time.
		v := reflect.ValueOf(m)
		if v.Kind() == reflect.Ptr {
			for _, previous := range r.entries {
				if previous.module != nil && reflect.TypeOf(previous.module) == v.Type() && reflect.ValueOf(previous.module).Pointer() == v.Pointer() {
					return fail("create", fmt.Errorf("instance already owned by %s", previous.registration.Name))
				}
			}
		}
		r.mu.Lock()
		e.module = m
		r.mu.Unlock()
	}
	if err != nil {
		return fail("create", err)
	}
	if nilModule(m) {
		return fail("create", errors.New("factory returned nil module"))
	}
	if err := ctx.Err(); err != nil {
		return fail("create", err)
	}
	r.set(e, Initializing, nil)
	if err := m.Init(ctx); err != nil {
		return fail("init", err)
	}
	if err := ctx.Err(); err != nil {
		return fail("init", err)
	}
	r.mu.Lock()
	e.started = true
	e.state = Starting
	r.mu.Unlock()
	if err := m.Start(ctx); err != nil {
		return fail("start", err)
	}
	if err := ctx.Err(); err != nil {
		return fail("start", err)
	}
	h, err := m.Status(ctx)
	if err != nil {
		return fail("status", err)
	}
	if err := ctx.Err(); err != nil {
		return fail("status", err)
	}
	if !h.Ready {
		return fail("status", fmt.Errorf("not ready: %s", h.Reason))
	}
	r.set(e, Running, nil)
	return nil
}

func (r *Runtime) failStart(err, previousCleanup error) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.options.RollbackTimeout)
	defer cancel()
	cleanupErr := r.cleanup(ctx, r.entries)
	if cleanupErr != nil {
		err = Errors{err, cleanupErr}
	}
	r.mu.Lock()
	r.state, r.startErr = Failed, err
	// Rollback has already attempted cleanup. Preserve its result across Shutdown.
	r.stopped, r.stopErr = true, cleanupErr
	if previousCleanup != nil {
		r.stopErr = previousCleanup
		if cleanupErr != nil {
			r.stopErr = Errors{previousCleanup, cleanupErr}
		}
	}
	r.mu.Unlock()
	return err
}

// Shutdown immediately withdraws readiness and cancels in-progress startup.
// The context bounds both waiting for startup and shutdown itself. Modules must
// cooperate with cancellation; Runtime does not abandon calls in goroutines.
func (r *Runtime) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	r.stop = true
	if r.cancel != nil {
		r.cancel()
	}
	if !r.stopped {
		r.state = Draining
	}
	r.mu.Unlock()
	// Even an expired caller should force-close resources when no operation is
	// in progress. If startup/another shutdown owns the gate, waiting is bounded.
	select {
	case r.gate <- struct{}{}:
	default:
		select {
		case r.gate <- struct{}{}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	defer func() { <-r.gate }()
	r.mu.RLock()
	stopped, result := r.stopped, r.stopErr
	r.mu.RUnlock()
	if stopped {
		return result
	}
	ctx, cancel := context.WithTimeout(ctx, r.options.DrainTimeout+r.options.CloseTimeout)
	defer cancel()
	err := r.cleanup(ctx, r.entries)
	r.mu.Lock()
	r.stopped, r.stopErr = true, err
	if err != nil {
		r.state = Failed
	} else {
		r.state = Closed
	}
	r.mu.Unlock()
	return err
}

// cleanup uses one deadline, reserving time for Close. Even with an expired
// context it invokes Close to allow immediate force-stop, except on dependencies
// still borrowed by modules returning ErrResourcesInUse.
func (r *Runtime) cleanup(ctx context.Context, entries []*entry) error {
	deadline, _ := ctx.Deadline()
	remaining := time.Until(deadline)
	reserve := r.options.CloseTimeout
	if reserve > remaining/2 {
		reserve = remaining / 2
	}
	if reserve < 0 {
		reserve = 0
	}
	drainDeadline := time.Now().Add(r.options.DrainTimeout)
	if d := deadline.Add(-reserve); d.Before(drainDeadline) {
		drainDeadline = d
	}
	drainCtx, cancel := context.WithDeadline(ctx, drainDeadline)
	defer cancel()
	var failures Errors
	retained := make(map[string]bool)
	for _, e := range entries {
		if e.retainsDeps {
			r.retainDependencies(e, retained)
		}
	}
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		if e.module == nil || e.closed || !e.started {
			continue
		}
		r.set(e, Draining, nil)
		if err := e.module.Drain(drainCtx); err != nil {
			err = &ModuleError{e.registration.Name, "drain", err}
			failures = append(failures, err)
			r.set(e, Draining, err)
		}
	}
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		if e.module == nil || e.closed {
			continue
		}
		if retained[e.registration.Name] {
			err := &ModuleError{e.registration.Name, "close skipped", ErrResourcesInUse}
			failures = append(failures, err)
			r.set(e, Failed, err)
			continue
		}
		r.set(e, Closing, nil)
		err := e.module.Close(ctx)
		if errors.Is(err, ErrResourcesInUse) {
			r.mu.Lock()
			e.retainsDeps = true
			r.mu.Unlock()
			r.retainDependencies(e, retained)
		}
		if err != nil {
			err = &ModuleError{e.registration.Name, "close", err}
			failures = append(failures, err)
		}
		r.mu.Lock()
		e.closed = true
		if err != nil {
			e.err = err
		}
		if e.err != nil {
			e.state = Failed
		} else {
			e.state = Closed
		}
		r.mu.Unlock()
	}
	if len(failures) != 0 {
		return failures
	}
	return nil
}

// Dependencies of a module with surviving work must remain allocated. Other
// independent branches still close, and retained modules report no Close attempt.
func (r *Runtime) retainDependencies(e *entry, retained map[string]bool) {
	for _, name := range e.registration.DependsOn {
		if retained[name] {
			continue
		}
		retained[name] = true
		for _, dep := range r.entries {
			if dep.registration.Name == name {
				r.retainDependencies(dep, retained)
				break
			}
		}
	}
}

// Status queries live health in dependency order. An optional dependency that
// becomes unhealthy also makes its required dependents unready. No module call
// runs under the runtime lock. The caller must provide a bounded context.
func (r *Runtime) Status(ctx context.Context) Snapshot {
	r.mu.RLock()
	snapshot := Snapshot{State: r.state, Ready: r.state == Running && !r.stop}
	entries := make([]entry, len(r.entries))
	for i, e := range r.entries {
		entries[i] = *e
	}
	r.mu.RUnlock()
	ready := make(map[string]bool)
	for _, e := range entries {
		s := ModuleStatus{Name: e.registration.Name, State: e.state, Optional: e.registration.Optional, CloseAttempted: e.closed}
		if e.err != nil {
			s.LastError = e.err.Error()
		}
		if e.state == Running && !e.closed {
			if err := ctx.Err(); err != nil {
				s.Reason = err.Error()
			} else {
				h, err := e.module.Status(ctx)
				if err == nil {
					err = ctx.Err()
				}
				if err != nil {
					s.Reason = err.Error()
				} else {
					s.Ready, s.Reason = h.Ready, h.Reason
				}
			}
			for _, dep := range e.registration.DependsOn {
				if !ready[dep] {
					s.Ready = false
					s.Reason = "dependency " + dep + " unavailable"
					break
				}
			}
		} else {
			s.Reason = string(e.state)
		}
		ready[s.Name] = s.Ready
		if !s.Optional && !s.Ready {
			snapshot.Ready = false
		}
		snapshot.Modules = append(snapshot.Modules, s)
	}
	// Do not return overall readiness from a snapshot taken before drain began.
	r.mu.RLock()
	snapshot.Ready = snapshot.Ready && r.state == Running && !r.stop
	r.mu.RUnlock()
	return snapshot
}
