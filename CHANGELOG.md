# Changelog

All notable changes to `go-agent-launch` are documented in this file. Per-release notes are also published as GitHub Releases.

## Unreleased

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
