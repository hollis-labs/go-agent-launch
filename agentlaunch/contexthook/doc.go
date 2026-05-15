// Package contexthook bridges go-agent-launch's ContextHook extension
// point to the go-agent-context assembly pipeline.
//
// # Why a subpackage
//
// agentlaunch declares the hook TYPE (ContextHook) but ships no concrete
// implementation — by design. Per the boundary documented in
// agentlaunch/hooks.go, the default behaviour is "no hook registered →
// empty boot prompt." A concrete implementation belongs OUTSIDE the
// top-level agentlaunch package so the launch contract does not gain a
// reverse dependency on context assembly.
//
// This subpackage supplies that implementation, backed by
// github.com/hollis-labs/go-agent-context. The adapter is a one-way
// bridge: contexthook IMPORTS agentcontext but agentcontext NEVER
// imports any agentlaunch type. Operators verify this with:
//
//	go list -m all                 # in go-agent-launch (require)
//	grep agent-launch go.mod       # in go-agent-context (no match)
//
// # Dependency state
//
// As of CW-20260515-0010, go-agent-context is unpublished — it lives on
// a feature branch with no tag. This module pins it via a replace
// directive pointing at the local checkout (../go-agent-context). The
// replace is the documented bridge for SP-20260514-0004 Phase 2; once
// go-agent-context tags v0.1.0 (or similar), the replace can be dropped
// in favour of a real require version.
//
// # Usage
//
// The common case is a one-liner over the default eight-resolver set:
//
//	resolverMap := resolvers.WithSkillIndex(resolvers.Default())
//	provider, _ := agentcontext.NewProvider(resolverMap, agentcontext.DefaultRenderer{})
//	hook := contexthook.New(provider, contexthook.Config{
//	    SlotExtractor: myExtractor,    // see Config.SlotExtractor
//	    PlantArtifacts: true,          // optional, for debugging
//	})
//	prepared, _ := launcher.Prepare(ctx, compiled, launcher.WithContextHook(hook))
//
// See contexthook.New for the contract between hook → provider → bootdir.
package contexthook
