package providerplant

import (
	"context"
	"errors"
	"testing"

	"github.com/hollis-labs/go-providers/provider"

	"github.com/hollis-labs/go-agent-launch/agentlaunch"
	"github.com/hollis-labs/go-agent-launch/agentlaunch/launcher"
)

func TestDefaultResolver_Claude(t *testing.T) {
	a, err := DefaultResolver(compiledFor(t, "claude", agentlaunch.RuntimePTY))
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if _, ok := a.(*provider.ClaudeAdapter); !ok {
		t.Errorf("got %T, want *provider.ClaudeAdapter", a)
	}
}

func TestDefaultResolver_CodexExecVsAppServer(t *testing.T) {
	exec, err := DefaultResolver(compiledFor(t, "codex", agentlaunch.RuntimeSubprocess))
	if err != nil {
		t.Fatalf("resolve exec: %v", err)
	}
	if cx, ok := exec.(*provider.CodexAdapter); !ok || cx.Mode == "app-server" {
		t.Errorf("subprocess runtime: got %T mode=%q, want exec-mode CodexAdapter", exec, modeOf(exec))
	}

	app, err := DefaultResolver(compiledFor(t, "codex", agentlaunch.RuntimeJsonRpcStdio))
	if err != nil {
		t.Fatalf("resolve app-server: %v", err)
	}
	if cx, ok := app.(*provider.CodexAdapter); !ok || cx.Mode != "app-server" {
		t.Errorf("jsonrpc runtime: got %T mode=%q, want app-server CodexAdapter", app, modeOf(app))
	}
}

func TestDefaultResolver_Opencode(t *testing.T) {
	a, err := DefaultResolver(compiledFor(t, "opencode", agentlaunch.RuntimeSubprocess))
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	oc, ok := a.(*provider.OpencodeAdapter)
	if !ok {
		t.Fatalf("got %T, want *provider.OpencodeAdapter", a)
	}
	if oc.Agent != "agent-name" {
		t.Errorf("OpencodeAdapter.Agent = %q, want agent-name (from AgentSpec.Name)", oc.Agent)
	}
}

func TestDefaultResolver_NilCompiled(t *testing.T) {
	if _, err := DefaultResolver(nil); !errors.Is(err, ErrNilCompiled) {
		t.Fatalf("DefaultResolver(nil) err = %v, want ErrNilCompiled", err)
	}
}

// TestPlant_WithAdapterOverride proves WithAdapter bypasses resolver
// lookup — here a plain codex adapter planted for a claude launch.
func TestPlant_WithAdapterOverride(t *testing.T) {
	isolateHome(t)
	prepared, err := launcher.Prepare(context.Background(), compiledFor(t, "claude", agentlaunch.RuntimePTY))
	if err != nil {
		t.Fatalf("prepare: %v", err)
	}
	if err := Plant(context.Background(), prepared, WithAdapter(provider.NewCodexAdapter())); err != nil {
		t.Fatalf("plant: %v", err)
	}
	// The codex adapter's BootDirSpec was planted despite the claude plan.
	assertExists(t, prepared.PlantedBootDir, "config.toml")
}

func modeOf(a provider.BootDirProvider) string {
	if cx, ok := a.(*provider.CodexAdapter); ok {
		return cx.Mode
	}
	return ""
}
