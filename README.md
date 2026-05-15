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

```go
// TODO(phase-1): example after Compile/Prepare lands
```

A runnable end-to-end example will live under [`examples/`](./examples) once
the `LaunchPlan` / `CompiledLaunch` / `PreparedLaunch` surface and the
`Compile` / `Prepare` entry points are committed.

## License

MIT — see [LICENSE](./LICENSE).
