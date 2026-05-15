// Package agentlaunch provides a portfolio-shared launch substrate for
// agent runtimes. It defines the LaunchPlan → CompiledLaunch →
// PreparedLaunch pipeline that turns a declarative catalog entry plus a
// runtime selection into a materialized boot directory and a
// ready-to-Start config the caller hands to go-agent-sessions.
//
// The library sits above github.com/hollis-labs/go-agent-sessions and
// github.com/hollis-labs/go-providers, and below the app-specific
// orchestrators (Tether, Torque, Nanite) that each previously grew their
// own near-identical launch pipelines. It owns:
//
//   - the LaunchPlan / CompiledLaunch / PreparedLaunch types and the
//     Compile / Prepare entry points that move between them,
//   - the provider × runtime support matrix that validates a plan against
//     the capabilities of the selected go-providers adapter and the
//     go-agent-sessions lifecycle shape,
//   - the Tether-compatible catalog schema so a single catalog entry can
//     drive every consumer in the portfolio.
//
// The library is intentionally app-neutral. It imports go-agent-sessions
// and go-providers but no app-specific repository (Tether, Torque,
// Nanite). Consumers configure the pipeline through caller-supplied
// types and sinks rather than direct dependencies on any orchestrator.
//
// The current scaffold ships this package documentation only. The
// LaunchPlan / CompiledLaunch / PreparedLaunch surface, the Compile and
// Prepare entry points, and the provider × runtime matrix land in
// subsequent Phase 1 subagents under sprint SP-20260514-0003.
package agentlaunch
