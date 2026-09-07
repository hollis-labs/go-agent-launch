# go-agent-launch

The launch substrate: it compiles a declarative `LaunchPlan` into a
`CompiledLaunch` (resolved provider, boot-dir spec, capability matrix) and
prepares it into a `PreparedLaunch` (materialized bootdir, ready-to-start
config) that `go-agent-sessions` consumes. It plans and materializes launches;
it does not start, supervise or tear down a session, and it holds no
app-specific orchestration logic.

## Start Here

- `README.md` explains the `Compile` → `Prepare` → `Plant` pipeline and the
  `RuntimeBinding` / `BootSpec` split.
- `agentlaunch/launcher/compile.go` turns a `LaunchPlan` into a `CompiledLaunch`.
- `agentlaunch/providerplant/plant.go` materializes the bootdir and applies the
  consumer overlay.
- `agentlaunch/matrix/matrix.go` owns which provider × runtime pairs are
  supported.
- `agentlaunch/catalog/loader.go` reads the Tether-compatible catalog schema.
- `agentlaunch/sessionshim/shim.go` is the handoff into `go-agent-sessions`.
- `agentlaunch/contexthook/hook.go` is where `go-agent-context` plugs in.
- `agentlaunch/parity/doc.go` states the parity harness contract.

## Commands

```bash
go vet ./...
go test -race -count=1 ./...
go test -race -count=1 -v ./agentlaunch/parity/...
```

CI runs all three. `TestParity_LiveCatalog` resolves against the operator's
live `~/.tether/catalog`, so it skips where no catalog exists and can go red
from machine-local catalog edits alone — see `## Boundaries` for how to read a
red result.

## Boundaries

This module was absorbed into `agentkit` as `agentkit/agentlaunch` at agentkit
v0.1.0, and this repo is maintenance-only. New work belongs in `agentkit`.

`CHANGELOG.md` and the git tags are the authority for what has shipped here.

`matrix/` is the single source of truth for supported provider × runtime pairs.
An unsupported pair must be refused, never silently defaulted to a working one —
`TestLookupUnsupportedPairs`, `TestResolve_UnsupportedRuntime` and
`TestCompileMatrixErrorWraps` guard that.

Consumer overlays win on runtime-critical launch fields, but the overlay key
set is validated rather than trusted: `TestPlant_UnsafeOverlayRejected` and
`TestValidateOverlayKey` exist so an overlay cannot reach outside the bootdir.

`Compile` is deterministic for a given plan (`TestCompileDeterministic`), and
the parity harness is read-only against whatever catalog it inspects
(`TestParity_ReadOnlyContract`). Neither is incidental — both are what let a
single catalog entry drive every consumer.

`agentlaunch/parity/expected_diffs.go` is the harness's built-in registry of
explained divergences, and it is rot-guarded: an entry that stops being
observed fails the run. A red parity result usually means the live catalog
moved, not that this library regressed — read the diff before changing code.
