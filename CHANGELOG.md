# Changelog

All notable changes to `go-agent-launch` are documented in this file. Per-release notes are also published as GitHub Releases.

## Unreleased

## v0.3.0 — 2026-05-17

### Added — `PlanFromLaunch`: the LaunchSpec → LaunchPlan bridge

The S4 parameterized-launch model (`LaunchSpec` + `LaunchBag` + var
resolution + `Render`) produced a `RenderResult` and stopped there — the
execution pipeline (`launcher.Compile → Prepare → Plant`) consumes a
`LaunchPlan`, and nothing assembled one from the new model. `PlanFromLaunch`
is that missing seam, so both front-ends drive the same shipped pipeline.

- **`agentlaunch.PlanFromLaunch(PlanFromLaunchInput) (LaunchPlan, error)`** —
  a pure, deterministic, I/O-free transform. The caller owns resolution
  (vars, `runner`→`RuntimeBinding`, agent); the bridge owns assembly. The
  returned plan is `Validate()`-clean and ready for `launcher.Compile` — the
  shipped Compile → Prepare → Plant pipeline is reused unchanged.
- `runner` is the sole runtime source: the bridge never reads the spec's
  embedded `BootSpec.Runtime`. An empty/invalid `RuntimeBinding` is a hard
  error — there is no spec-baked fallback.
- The agent identity is a bridge input (`PlanFromLaunchInput.Agent`), never
  read from the LaunchSpec — agent content stays consumer-owned.
- `MCP` / `Injection` / `Metadata` on the returned plan are a base the
  consumer overlays. The raw `isolation` token, the `RuntimeBinding.Timeout`,
  and the spec/bag ids are surfaced on `Metadata.Annotations`.

### Changed

- `Version` is now `v0.3.0` — it had been left at the `v0.1.0-dev`
  placeholder through v0.2.0.

### Notes

- The `runner` → `RuntimeBinding` registry resolver helper
  (`ResolveRuntimeBinding`) is a v0.3.1 fast-follow. It is a caller
  convenience — `PlanFromLaunch` is pure and each consumer resolves `runner`
  its own way — so it does NOT gate the S5 dispatch-wiring step; the bridge
  does, and the bridge ships here.

## v0.2.0 — 2026-05-17

First tagged release of the parameterized-launch engine: the S1 platform
contracts, the S3 directory registry, and the S4 Boot Assembly Spec /
var-resolution / materialization / launch-spec / parity layers — plus the
dependency bump below so the engine builds against the fixed provider and
session libraries the S5 cutover consumers require.

### Changed — dependency bump: go-providers v0.20.0, go-agent-sessions v0.9.5

Bumps the provider/runtime libraries (go-providers v0.17.1 → v0.20.0,
go-agent-sessions v0.9.4 → v0.9.5) so the engine builds and materializes
against the libraries that carry the agent-execution fixes the S5 cutover
depends on:

- **go-providers v0.20.0** — `CodexAdapter.ApprovalPolicy` / `SandboxMode`;
  the planted codex `config.toml` always carries a non-interactive
  `approval_policy` / `sandbox_mode` (fixes the headless-codex approval
  deadlock). Includes v0.18.0–v0.19.0: current-schema
  `.claude/settings.json` and first-class `ClaudeAdapter.PermissionMode`.
- **go-agent-sessions v0.9.5** — the `jsonrpc-stdio` runtime answers
  server-initiated requests instead of dropping them (fixes the codex
  `app-server` approval-elicitation deadlock); adds `JsonRpcRequestHook`.

No go-agent-launch API change — the bump is additive on the dependency
side; the full build + test suite is green against the new versions.

### Added — old-vs-new launch-plan parity harness (S4.5, CW-20260517-0030)

The cutover gate for the Tether platform reshape: before S5 flips live
consumers onto the new launch model, the new model must provably
reproduce the old model's launch-identity resolution. This adds the
harness that proves it.

- **`agentlaunch/parity/`** — a new subpackage.
  - `RunParity` resolves each catalog-corpus launch BOTH ways — the
    legacy static-file path (`catalog.LoadGlobal` + `GlobalCatalog.Resolve`
    → `LaunchPlan`) and the new S4.4 path (`LoadLaunchSpec` +
    `LoadLaunchBag` + S4.2 var resolution + `LaunchSpec.Render`) — and
    diffs the resulting `NormalizedPlan` (project, work_dir, runner,
    isolation): the launch identity both models share.
  - `expectedDiffs` / `expectedOldErrors` are the annotated registry of
    intentional divergences. A non-zero diff passes only if it is
    documented there; an unexplained diff fails the harness.
  - The catalog corpus is the 11 S4.4 launch bags that re-express a real
    legacy launch (the synthetic `tether-minimum` bag has no legacy
    counterpart and is excluded).
  - `NormalizeRuntimeKinds` is a documented in-memory pre-processing step:
    the live `~/.tether/catalog/` providers mostly omit `runtime_kind`, so
    the old-side `Resolve` would fail; this bridges `bootstrap.mode` →
    `runtime_kind`. It never writes the catalog.
  - The harness is read-only on `~/.tether/catalog/` and ships a
    self-contained fixture catalog so CI runs deterministically with no
    Tether install; `TestParity_LiveCatalog` opportunistically checks the
    live catalog when present.
- **`.github/workflows/check.yml`** — CI workflow (fmt / vet /
  golangci-lint / `go test -race` / govulncheck) with an explicit
  `parity harness (S4.5 cutover gate)` step that runs on every push/PR.

Parity result over the corpus: 5 launches identical, 6 explained
divergences (5 × the `agent-mux` project-repoint, 1 × the
`hollislabs-web-writer` dangling-agent legacy defect). Zero unexplained
diffs.

### Added — launch Specs + templates (S4.4, CW-20260517-0029)

Collapses the legacy `~/.tether/catalog/launches/` (64 files) and
`boot-profiles/` directories into one parameterized `LaunchSpec` plus
per-launch input bags. Builds on the S4.1 `AssemblySpec`, S4.2 var
resolution, and S4.3 materializer; adds no new field shapes to those
frozen types.

- **`agentlaunch/launchspec.go`** — the S4.4 wrapper layer.
  - `LaunchSpec` embeds `AssemblySpec` and pins the minimum-config input
    contract: a launch is valid with only `work_dir` + `runner`
    (required); `isolation` + `bus` are the only other knobs (optional).
  - `LaunchBag` is one concrete launch invocation — the input bag that
    replaces a single legacy `launches/*.yaml`. The `.worktree` twin is
    gone: a worktree launch is the same bag with `isolation: worktree`.
  - `ValidateMinimumConfig` enforces the two-knob minimum and rejects
    unknown bag keys.
  - `LoadLaunchSpec` / `LoadLaunchBag` parse + validate the artifacts.
  - Provider is the `runner` input feeding the runtime-binding, never
    blueprint identity (D3).
- **`agentlaunch/testdata/specs/`** — the re-expressed catalog: one
  canonical `launch-assembly.yaml` LaunchSpec (boot-profile slots folded
  into vars + a merge-tag template), a starter set of 8 common-setup
  templates, and 12 representative launch bags spanning the
  agent/project/provider/mode cross-product. Loadable by the S4.5 parity
  harness.

## v0.1.0 — 2026-05-15

Initial scaffold. See sprint `SP-20260514-0003`.

### Added — provider bootdir planting (CW-20260515-0106)

`go-agent-launch` now owns provider bootdir planting end to end: consumers
no longer reimplement the provider-specific planting layer per app.

- **`agentlaunch/providerplant`** — go-providers integration subpackage.
  - `Plant` materializes a provider's go-providers `BootDirSpec` files into
    an already-`Prepare`d launch's bootdir, then rewires `PreparedLaunch`
    `Env` / `Argv` / `Workdir` from the spec.
  - `PrepareAndPlant` is the one-call API (`launcher.Prepare` + `Plant`).
  - `DefaultResolver` maps the matrix `provider×runtime` pair to the right
    go-providers adapter (claude, codex exec vs. app-server, opencode).
  - `PlantContextFor` translates `PreparedLaunch` → `provider.PlantContext`.
- **`agentlaunch/sessionshim`** — go-agent-sessions integration subpackage.
  `ToSessionLaunch` converts a `PreparedLaunch` into the binary +
  `agentsessions.StartOptions` the runtime consumes (`AutoPlantBootDir`
  left off — providerplant already planted the dir).
- **`InjectionSpec.NativeFiles`** + the new `NativeFile` type — first-class
  provider-native extra files. `NativeFileSkill` resolves to the provider
  skill convention (`.claude/skills/<id>.md`, `.opencode/skills/<id>.md`);
  `NativeFileRaw` plants a caller-placed file verbatim.
- **`agentlaunch.ValidateBootDirRelPath`** — exported bootdir-relative
  path-safety check (the rule `LaunchPlan.Validate` applies to
  `BootDirOverlay` keys), re-run by the planter as defence in depth.
- Planting order is fixed and deterministic: provider files →
  `NativeFiles` → `BootDirOverlay` (overlay wins on path collision).
- New runnable example under `examples/providerplant`.
