package providerplant

import (
	"fmt"

	"github.com/hollis-labs/go-providers/provider"

	"github.com/hollis-labs/go-agent-launch/agentlaunch"
	"github.com/hollis-labs/go-agent-launch/agentlaunch/matrix"
)

// AdapterResolver maps a CompiledLaunch to the go-providers adapter whose
// BootDirSpec should be planted. Callers override the built-in
// DefaultResolver via WithResolver — e.g. to select bare-mode Claude, a
// pinned CLI variant, or a custom provider not in the matrix.
//
// A resolver MUST return an adapter that also implements
// provider.BootDirProvider; Plant surfaces ErrNoBootDirSpec otherwise.
type AdapterResolver func(*agentlaunch.CompiledLaunch) (provider.BootDirProvider, error)

// DefaultResolver resolves the standard adapter for a launch's
// provider×runtime pair via the matrix:
//
//   - claude   → provider.NewClaudeAdapter()
//   - codex    → provider.NewCodexAdapter(), or NewCodexAdapterAppServer()
//     for the jsonrpc-stdio runtime (the app-server daemon rejects the
//     --cd flag, so its BootDirSpec suppresses ProjectDirArg)
//   - opencode → &provider.OpencodeAdapter{Agent: <agent name>}
//
// The agent name fed to the opencode adapter is AgentSpec.Name, falling
// back to AgentSpec.ID — the same precedence Prepare uses for
// PreparedPlantContext.AgentName.
//
// DefaultResolver returns adapters in their plain (non-bare) shape: the
// planted BootDirSpec files are identical across the PTY / streaming /
// subprocess variants of a provider, so the runtime variant only needs
// to be distinguished where it changes the spec (the codex app-server
// case above). Consumers that need bare-mode Claude flag injection wire
// that at the go-agent-sessions boundary.
func DefaultResolver(compiled *agentlaunch.CompiledLaunch) (provider.BootDirProvider, error) {
	if compiled == nil || compiled.Plan == nil {
		return nil, ErrNilCompiled
	}
	plan := compiled.Plan
	desc, err := matrix.Lookup(plan.Provider, plan.Runtime)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAdapterResolution, err)
	}
	switch desc.BootDirRenderer {
	case matrix.BootDirRendererClaude:
		return provider.NewClaudeAdapter(), nil
	case matrix.BootDirRendererCodex:
		if plan.Runtime == agentlaunch.RuntimeJsonRpcStdio {
			return provider.NewCodexAdapterAppServer(), nil
		}
		return provider.NewCodexAdapter(), nil
	case matrix.BootDirRendererOpencode:
		return &provider.OpencodeAdapter{Agent: agentName(plan)}, nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownRenderer, desc.BootDirRenderer)
	}
}

// agentName returns the display name the opencode adapter and the
// PlantContext renderers reference: AgentSpec.Name, then AgentSpec.ID.
func agentName(plan *agentlaunch.LaunchPlan) string {
	if plan.Agent.Name != "" {
		return plan.Agent.Name
	}
	return plan.Agent.ID
}
