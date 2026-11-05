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
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type events struct {
	sync.Mutex
	values []string
}

func (e *events) add(s string) { e.Lock(); defer e.Unlock(); e.values = append(e.values, s) }
func (e *events) list() []string {
	e.Lock()
	defer e.Unlock()
	return append([]string(nil), e.values...)
}

type fake struct {
	name      string
	events    *events
	fail      string
	err       error
	unhealthy int32
	start     func(context.Context) error
	drain     func(context.Context) error
}

func (f *fake) op(op string) error {
	f.events.add(f.name + "." + op)
	if f.fail == op {
		return f.err
	}
	return nil
}
func (f *fake) Init(context.Context) error { return f.op("init") }
func (f *fake) Start(ctx context.Context) error {
	if err := f.op("start"); err != nil {
		return err
	}
	if f.start != nil {
		return f.start(ctx)
	}
	return nil
}
func (f *fake) Status(context.Context) (Health, error) {
	return Health{Ready: atomic.LoadInt32(&f.unhealthy) == 0}, nil
}
func (f *fake) Drain(ctx context.Context) error {
	if err := f.op("drain"); err != nil {
		return err
	}
	if f.drain != nil {
		return f.drain(ctx)
	}
	return nil
}
func (f *fake) Close(context.Context) error { return f.op("close") }
func reg(f *fake, deps ...string) Registration {
	return Registration{Name: f.name, DependsOn: deps, New: func(context.Context) (Module, error) { f.events.add(f.name + ".new"); return f, nil }}
}
func mustRuntime(t *testing.T, regs []Registration, opts Options) *Runtime {
	t.Helper()
	r, err := New(regs, opts)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
func bounded(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)
	return ctx
}

func TestDependencyOrderAndSingleOwnership(t *testing.T) {
	e := &events{}
	dep, api, native := &fake{name: "dep", events: e}, &fake{name: "api", events: e}, &fake{name: "native", events: e}
	r := mustRuntime(t, []Registration{reg(native, "dep"), reg(api, "dep"), reg(dep)}, Options{})
	ctx := bounded(t)
	for i := 0; i < 2; i++ {
		if err := r.Start(ctx); err != nil {
			t.Fatal(err)
		}
	}
	if !r.Status(ctx).Ready {
		t.Fatal("runtime not ready")
	}
	for i := 0; i < 2; i++ {
		if err := r.Shutdown(ctx); err != nil {
			t.Fatal(err)
		}
	}
	want := []string{"dep.new", "dep.init", "dep.start", "api.new", "api.init", "api.start", "native.new", "native.init", "native.start", "native.drain", "api.drain", "dep.drain", "native.close", "api.close", "dep.close"}
	if got := e.list(); !reflect.DeepEqual(got, want) {
		t.Fatalf("events = %v", got)
	}
	if r.Status(ctx).Ready {
		t.Fatal("ready after shutdown")
	}
	if !errors.Is(r.Start(ctx), ErrStopped) {
		t.Fatal("restart allowed")
	}
}

func TestGraphValidationHasNoSideEffects(t *testing.T) {
	called := false
	newModule := func(context.Context) (Module, error) { called = true; return nil, nil }
	for _, tc := range []struct {
		name string
		regs []Registration
	}{
		{"duplicate", []Registration{{Name: "a", New: newModule}, {Name: "a", New: newModule}}},
		{"missing", []Registration{{Name: "a", New: newModule, DependsOn: []string{"b"}}}},
		{"cycle", []Registration{{Name: "a", New: newModule, DependsOn: []string{"b"}}, {Name: "b", New: newModule, DependsOn: []string{"a"}}}},
		{"self", []Registration{{Name: "a", New: newModule, DependsOn: []string{"a"}}}},
		{"empty", []Registration{{New: newModule}}},
		{"nilFactory", []Registration{{Name: "a"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := New(tc.regs, Options{}); err == nil {
				t.Fatal("accepted invalid graph")
			}
		})
	}
	if called {
		t.Fatal("factory called during validation")
	}
}

func TestStartupFailureRollsBackPartialModule(t *testing.T) {
	for _, op := range []string{"init", "start", "create"} {
		t.Run(op, func(t *testing.T) {
			e := &events{}
			cause := errors.New("boom")
			a, b := &fake{name: "a", events: e}, &fake{name: "b", events: e, fail: op, err: cause}
			bReg := reg(b, "a")
			if op == "create" {
				bReg.New = func(context.Context) (Module, error) { return b, cause }
			}
			r := mustRuntime(t, []Registration{reg(a), bReg}, Options{})
			ctx := bounded(t)
			if err := r.Start(ctx); !errors.Is(err, cause) {
				t.Fatalf("error = %v", err)
			}
			before := e.list()
			if !reflect.DeepEqual(before[len(before)-2:], []string{"b.close", "a.close"}) {
				t.Fatal(before)
			}
			if op == "init" && strings.Contains(strings.Join(before, ","), "b.drain") {
				t.Fatal("drained unstarted module")
			}
			if err := r.Start(ctx); !errors.Is(err, cause) {
				t.Fatal(err)
			}
			_ = r.Shutdown(ctx)
			if !reflect.DeepEqual(before, e.list()) {
				t.Fatal("cleanup replayed")
			}
		})
	}
}

func TestOptionalFailureAndRequiredDependent(t *testing.T) {
	for _, dependent := range []bool{false, true} {
		t.Run(map[bool]string{false: "independent", true: "dependent"}[dependent], func(t *testing.T) {
			e := &events{}
			bad := &fake{name: "bad", events: e, fail: "init", err: errors.New("bad")}
			good := &fake{name: "good", events: e}
			badReg := reg(bad)
			badReg.Optional = true
			goodReg := reg(good)
			if dependent {
				goodReg.DependsOn = []string{"bad"}
			}
			r := mustRuntime(t, []Registration{goodReg, badReg}, Options{})
			ctx := bounded(t)
			err := r.Start(ctx)
			if dependent {
				if err == nil {
					t.Fatal("dependent started")
				}
				if strings.Contains(strings.Join(e.list(), ","), "good.new") {
					t.Fatal("dependent factory called")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				s := r.Status(ctx)
				if !s.Ready || s.Modules[0].State != Failed || !s.Modules[0].CloseAttempted || s.Modules[0].LastError == "" {
					t.Fatalf("snapshot=%+v", s)
				}
			}
			_ = r.Shutdown(ctx)
		})
	}
}

func TestLiveDependencyHealth(t *testing.T) {
	e := &events{}
	dep := &fake{name: "dep", events: e}
	app := &fake{name: "app", events: e}
	d := reg(dep)
	d.Optional = true
	r := mustRuntime(t, []Registration{d, reg(app, "dep")}, Options{})
	ctx := bounded(t)
	if err := r.Start(ctx); err != nil {
		t.Fatal(err)
	}
	atomic.StoreInt32(&dep.unhealthy, 1)
	s := r.Status(ctx)
	if s.Ready || s.Modules[1].Ready {
		t.Fatal("dependent ready despite dependency failure")
	}
	atomic.StoreInt32(&dep.unhealthy, 0)
	if !r.Status(ctx).Ready {
		t.Fatal("health did not recover")
	}
	_ = r.Shutdown(ctx)
}

func TestShutdownCancelsStartup(t *testing.T) {
	e := &events{}
	entered := make(chan struct{})
	f := &fake{name: "a", events: e, start: func(ctx context.Context) error { close(entered); <-ctx.Done(); return ctx.Err() }}
	r := mustRuntime(t, []Registration{reg(f)}, Options{})
	started := make(chan error, 1)
	go func() { started <- r.Start(context.Background()) }()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("start blocked")
	}
	if err := r.Shutdown(bounded(t)); err != nil {
		t.Fatal(err)
	}
	if err := <-started; !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if got := strings.Count(strings.Join(e.list(), ","), "a.close"); got != 1 {
		t.Fatalf("close calls=%d", got)
	}
}

func TestShutdownBudgetAndErrorAggregation(t *testing.T) {
	e := &events{}
	cause := errors.New("close failed")
	a := &fake{name: "a", events: e, fail: "close", err: cause}
	b := &fake{name: "b", events: e, drain: func(ctx context.Context) error { <-ctx.Done(); return ctx.Err() }}
	r := mustRuntime(t, []Registration{reg(a), reg(b, "a")}, Options{DrainTimeout: 20 * time.Millisecond, CloseTimeout: 100 * time.Millisecond})
	ctx := bounded(t)
	if err := r.Start(ctx); err != nil {
		t.Fatal(err)
	}
	err := r.Shutdown(ctx)
	if !errors.Is(err, cause) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
	if s := r.Status(ctx); s.Ready || s.State != Failed {
		t.Fatalf("snapshot=%+v", s)
	}
	before := e.list()
	if err2 := r.Shutdown(ctx); err2.Error() != err.Error() {
		t.Fatal("shutdown result changed")
	}
	if !reflect.DeepEqual(before, e.list()) {
		t.Fatal("cleanup repeated")
	}
}

func TestRollbackUsesFreshContext(t *testing.T) {
	e := &events{}
	ctx, cancel := context.WithCancel(context.Background())
	f := &fake{name: "a", events: e, start: func(context.Context) error { cancel(); return context.Canceled }, drain: func(ctx context.Context) error {
		if ctx.Err() != nil {
			t.Error("rollback used cancelled startup context")
		}
		return nil
	}}
	r := mustRuntime(t, []Registration{reg(f)}, Options{})
	if err := r.Start(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestConcurrentStatusAndShutdown(t *testing.T) {
	e := &events{}
	f := &fake{name: "a", events: e}
	r := mustRuntime(t, []Registration{reg(f)}, Options{})
	ctx := bounded(t)
	if err := r.Start(ctx); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				r.Status(ctx)
			}
		}()
	}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.Shutdown(ctx); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if n := strings.Count(strings.Join(e.list(), ","), "a.close"); n != 1 {
		t.Fatalf("close calls=%d", n)
	}
}

func TestNilFactoryResultAndStopBeforeStart(t *testing.T) {
	for _, typedNil := range []bool{false, true} {
		r := mustRuntime(t, []Registration{{Name: "nil", New: func(context.Context) (Module, error) {
			if typedNil {
				var f *fake
				return f, nil
			}
			return nil, nil
		}}}, Options{})
		if err := r.Start(bounded(t)); err == nil {
			t.Fatal("accepted nil module")
		}
	}
	e := &events{}
	r := mustRuntime(t, []Registration{reg(&fake{name: "a", events: e})}, Options{})
	if err := r.Shutdown(bounded(t)); err != nil {
		t.Fatal(err)
	}
	if err := r.Start(bounded(t)); !errors.Is(err, ErrStopped) {
		t.Fatal(err)
	}
	if len(e.list()) != 0 {
		t.Fatal("created module after shutdown")
	}
}

func TestSharedPointerCannotHaveTwoOwners(t *testing.T) {
	e := &events{}
	f := &fake{name: "shared", events: e}
	a, b := reg(f), reg(f)
	a.Name, b.Name = "a", "b"
	r := mustRuntime(t, []Registration{a, b}, Options{})
	if err := r.Start(bounded(t)); err == nil || !strings.Contains(err.Error(), "already owned by a") {
		t.Fatalf("startup=%v", err)
	}
	if n := strings.Count(strings.Join(e.list(), ","), "shared.close"); n != 1 {
		t.Fatalf("close calls=%d", n)
	}
}

func TestExpiredShutdownStillAttemptsForceClose(t *testing.T) {
	e := &events{}
	f := &fake{name: "a", events: e, drain: func(ctx context.Context) error { return ctx.Err() }}
	r := mustRuntime(t, []Registration{reg(f)}, Options{})
	if err := r.Start(bounded(t)); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := r.Shutdown(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if n := strings.Count(strings.Join(e.list(), ","), "a.close"); n != 1 {
		t.Fatalf("close calls=%d", n)
	}
}

func TestUnreadyStartupAndCleanupErrorArePreserved(t *testing.T) {
	e := &events{}
	cause := errors.New("cannot close")
	f := &fake{name: "a", events: e, unhealthy: 1, fail: "close", err: cause}
	r := mustRuntime(t, []Registration{reg(f)}, Options{})
	err := r.Start(bounded(t))
	if !errors.Is(err, cause) || !strings.Contains(err.Error(), "not ready") {
		t.Fatal(err)
	}
	s := r.Status(bounded(t))
	if s.Ready || s.Modules[0].State != Failed || !s.Modules[0].CloseAttempted {
		t.Fatalf("snapshot=%+v", s)
	}
}
