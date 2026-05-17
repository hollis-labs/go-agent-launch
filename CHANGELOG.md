# Changelog

All notable changes to `go-agent-launch` are documented in this file. Per-release notes are also published as GitHub Releases.

## Unreleased

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
