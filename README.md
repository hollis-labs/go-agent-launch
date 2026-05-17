# go-agent-launch

`go-agent-launch` is a small, standalone Go library that provides a portfolio-shared **launch substrate** for agent runtimes — the layer that takes a declarative `LaunchPlan` (catalog entry + profile + runtime selection), compiles it into a `CompiledLaunch` (resolved provider + boot-dir spec + capability matrix), and prepares it into a `PreparedLaunch` (materialized bootdir + ready-to-`Start` config) that downstream consumers feed straight into [`go-agent-sessions`](https://github.com/hollis-labs/go-agent-sessions).

It sits **above** [`go-agent-sessions`](https://github.com/hollis-labs/go-agent-sessions) and [`go-providers`](https://github.com/hollis-labs/go-providers), and **below** the app-specific orchestrators (Tether, Torque, Nanite) that each previously grew their own near-identical launch pipelines. The library defines the shared `LaunchPlan` → `CompiledLaunch` → `PreparedLaunch` flow, the provider × runtime support matrix, and a Tether-compatible catalog schema so a single catalog entry can drive every consumer.

## Status

`v0.1.x` shared-lib contract. The existing `LaunchPlan -> CompiledLaunch -> PreparedLaunch` flow remains the stable launch handoff, and the boot-assembly API is now split explicitly into:

- `agentlaunch.RuntimeBinding`: the hot-path-readable `provider / model / runtime_kind / args / timeout` contract.
- `agentlaunch.BootSpec`: the parameterized blueprint carrying `inputs[]`, `files[]`, `vars[]`, `injections[]`, and the associated runtime binding.

`agentlaunch.RuntimeKind` is canonical across those surfaces. Consumer overlays still win on runtime-critical launch fields when plans are compiled and prepared.

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

## Frozen Boot Contract

The shared boot contract published by this library is:

- `BootInput`: Hadron-style typed inputs with `name`, `type`, `required`, `default`, and `description`.
- `VarSpec`: derived vars with `source`, `freshness`, `fallback`, `on_error`, `phase`, and `secret`.
- `BootFileSpec` / `BootInjectionSpec`: path-safe bootdir materialization targets.
- `Materializer`: idempotent populate-against-existing-dir plus partial replant by file ID, injection ID, or slot ref.

Var source kinds are `literal`, `file`, `call`, and `cmd`. `call` supports `mcp` and `http` transports. Secret vars may not persist inline literal/fallback values at rest, and `call` / `cmd` vars require an explicit trust or authorization gate in schema.

## Directory Registry Contract

The shared directory-registry contract published by this library is local-first and file-backed by design:

- `RegistrationRecord` registers a stable `kind / owner / namespace / name` handle against a local contract file, with optional directory enrichment metadata.
- `RegistryEnvelope` defines the `register`, `deregister`, `query`, and `health` envelope with `resolution=local-first`.
- Published kind contracts include `agent-source`, `skill-source`, `mcp-server`, `runtime-binding`, `boot-spec`, `execution-template`, and `contract-object`.

`agent-source` and `skill-source` carry resolver handles only, not profile/skill bodies. `execution-template` composes references to distinct `runtime-binding` and `boot-spec` kinds rather than inlining them back together.

## License

MIT — see [LICENSE](./LICENSE).
