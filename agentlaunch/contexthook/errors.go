package contexthook

import "errors"

// Sentinel errors emitted by the contexthook adapter. Callers branch on
// these with errors.Is.
var (
	// ErrProviderNil is returned when New is invoked with a nil
	// agentcontext.ContextProvider. We deliberately defer the check
	// until the hook fires (not at construction) so dependency-injection
	// patterns that nil-check late still surface a clean sentinel.
	ErrProviderNil = errors.New("provider is nil")

	// ErrCompiledNil is returned when the hook is invoked with a nil
	// *agentlaunch.CompiledLaunch. The preparer should never do this in
	// practice — the sentinel is defensive against test stubs.
	ErrCompiledNil = errors.New("compiled is nil")

	// ErrProviderResultNil is returned when a custom ContextProvider
	// returns (nil, nil) from Assemble. Defensive — DefaultProvider
	// never does this, but a hand-rolled provider might.
	ErrProviderResultNil = errors.New("provider returned nil result")

	// ErrPlantArtifacts wraps any filesystem error encountered while
	// planting per-slot artifacts under <bootDir>/context/.
	ErrPlantArtifacts = errors.New("plant artifacts failed")
)
