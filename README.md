# go-agent-launch

`go-agent-launch` is a small, standalone Go library that provides a portfolio-shared **launch substrate** for agent runtimes — the layer that takes a declarative `LaunchPlan` (catalog entry + profile + runtime selection), compiles it into a `CompiledLaunch` (resolved provider + boot-dir spec + capability matrix), and prepares it into a `PreparedLaunch` (materialized bootdir + ready-to-`Start` config) that downstream consumers feed straight into [`go-agent-sessions`](https://github.com/hollis-labs/go-agent-sessions).

It sits **above** [`go-agent-sessions`](https://github.com/hollis-labs/go-agent-sessions) and [`go-providers`](https://github.com/hollis-labs/go-providers), and **below** the app-specific orchestrators (Tether, Torque, Nanite) that each previously grew their own near-identical launch pipelines. The library defines the shared `LaunchPlan` → `CompiledLaunch` → `PreparedLaunch` flow, the provider × runtime support matrix, and a Tether-compatible catalog schema so a single catalog entry can drive every consumer.

## Status

v0 — foundation. Public API in flux; expect breaking changes before the v0.1.0 tag. This repository currently ships the package scaffold only; the launch pipeline types, compiler, and preparer land in subsequent Phase 1 subagents under sprint `SP-20260514-0003`.

## Install

```bash
go get github.com/hollis-labs/go-agent-launch
```

## Usage

The launch pipeline is three stages — `Compile` → `Prepare` → `Plant` —
plus a conversion shim into [`go-agent-sessions`](https://github.com/hollis-labs/go-agent-sessions):

```go
import (
    "github.com/hollis-labs/go-agent-launch/agentlaunch/launcher"
    "github.com/hollis-labs/go-agent-launch/agentlaunch/providerplant"
    "github.com/hollis-labs/go-agent-launch/agentlaunch/sessionshim"
)

// 1. Compile the declarative LaunchPlan.
compiled, err := launcher.Compile(ctx, plan)

// 2. Prepare (workspace + bootdir) AND Plant (provider boot files,
//    native skills/files, injection overlay) in one call.
prepared, err := providerplant.PrepareAndPlant(ctx, compiled)

// 3. Convert into the go-agent-sessions launch handoff.
launch, err := sessionshim.ToSessionLaunch(prepared)
// launch.Binary + launch.Options → agentsessions Manager.Start
```

### Provider bootdir planting

`providerplant` owns the provider-specific planting layer so consumers
(Tether, Torque, Nanite) do not each reimplement it. `Plant` resolves the
go-providers adapter for the launch's `provider×runtime` pair, renders its
`BootDirSpec`, and writes — in a fixed, deterministic order — the provider
files, then `InjectionSpec.NativeFiles`, then `InjectionSpec.BootDirOverlay`
(overlay wins on a path collision). It then rewires the `PreparedLaunch`
`Env` / `Argv` / `Workdir` from the spec.

`NativeFile` gives first-class support for provider-native extra files
without app-specific callbacks: a `NativeFileSkill` entry lands at the
provider's skill path (`.claude/skills/<id>.md`, `.opencode/skills/<id>.md`),
a `NativeFileRaw` entry is planted verbatim at a path-validated location.

The integration subpackages (`providerplant`, `sessionshim`, `contexthook`)
each carry exactly one external dependency, so importing the core
`agentlaunch` package never pulls in `go-providers` or `go-agent-sessions`.

A runnable end-to-end example lives under
[`examples/providerplant`](./examples/providerplant) — run it with
`go run ./examples/providerplant`.

## License

MIT — see [LICENSE](./LICENSE).
