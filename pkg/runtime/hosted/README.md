# Hosted Runtime

`hosted` manages API servers, native services and shared dependencies through
one lifecycle. Its core imports only the Go standard library. It does not create
an API server, application callback connection or MOSN process implicitly.

This is the foundation for the mosn3 runtime migration. The existing
`runtime.MosnRuntime.Run` entry point is unchanged; switching production startup
to this core and adapting ZDAS/ZCache SDKs are subsequent integration steps.

## Registration and ownership

Register each owned instance once using a unique name and declare its direct
dependencies. Factories execute in deterministic dependency order. Modules are
required unless `Optional` is set. Invalid graphs fail before factories run.

Each module implements:

```go
type Module interface {
    Init(context.Context) error
    Start(context.Context) error
    Status(context.Context) (Health, error)
    Drain(context.Context) error
    Close(context.Context) error
}
```

`Init` prepares resources; `Start` returns when service readiness is established.
The runtime checks initial health before starting dependents. A module that
needs asynchronous warmup should wait within its startup context. Contexts bound
individual operations; modules own their background work until Close.

Business interfaces remain typed and outside the core. Assembly code injects
the same component or engine into its users, as shown in `cmd/layotto_hosted`.
Borrowers do not close shared resources. A factory returning the same pointer
under a second name is rejected to prevent double ownership.

## Failure, state and shutdown

- Factory results belong to Runtime immediately, including a non-nil module
  returned with an error. Objects that fail Init or Start are still closed.
- Required failures roll back all created modules. Optional failures are
  recorded and cleaned up; their dependents cannot start. Cleanup failure is
  fatal even for an optional module.
- Status reads live health. Dependency failures propagate to dependent
  readiness. Lifecycle errors and close attempts remain visible in snapshots.
  Call Status with a bounded context; module status methods must be concurrent-safe.
- Shutdown withdraws readiness immediately and cancels pending startup.
  Drain and Close execute in reverse dependency order. All close attempts run
  even after an earlier error. Errors preserve module, operation and errors.Is.
- Shutdown uses one overall deadline, with time reserved for Close; a shorter
  caller deadline wins. Rollback has its own bounded context because the startup
  context may already have been cancelled.
- Methods must honor cancellation. The core does not run uncooperative SDK
  operations in abandoned goroutines or claim that a timed-out operation freed
  its resources. A product adapter must supply its actual force-stop behavior.
- Start/Shutdown are idempotent and serialized. Runtime is single-use; create a
  new instance after startup failure or shutdown. Final close calls are not
  retried automatically; failures remain in the shutdown result.

## Server adapters

`grpcserver` owns a raw gRPC server and listener. Its registration callback can
initialize and register existing Layotto API implementations. Its optional
Release callback cleans up API-owned callback connections or other resources,
including partial registration failures. Business components remain separate
dependencies. GracefulStop falls back to Stop when the drain deadline expires.

`tcpserver` owns a listener and accepted connections. It stops accepting during
drain and waits for sessions; Close cancels handlers and closes remaining sockets.
Handlers must return on cancellation or connection closure. Product adapters
remain responsible for protocol-aware request rejection, transactions and ACKs.
An existing native Server can alternatively implement Module directly without
using this helper or replacing its networking implementation.

## Run without Envoy

From the Layotto repository root:

```sh
go run ./cmd/layotto_hosted -mode api
go run ./cmd/layotto_hosted -mode native
go run ./cmd/layotto_hosted -mode both
```

The example exposes the real Layotto SayHello API at `127.0.0.1:34904` and a
TCP echo service at `127.0.0.1:16379`. The latter demonstrates native hosting;
it is not a Redis or MySQL implementation. Addresses can be overridden with
`-api-address` and `-native-address`; `127.0.0.1:0` allocates a random port.
SIGINT/SIGTERM invoke the same runtime shutdown flow.

```sh
go test -race ./pkg/runtime/hosted/... ./cmd/layotto_hosted
```

Tests cover graph validation, failure rollback, live health propagation, shared
ownership, concurrent status/shutdown, shutdown during startup, cleanup errors,
real gRPC calls, TCP sessions, deadline-driven force stop, and all three hosting
modes. No external middleware or Envoy is required.
