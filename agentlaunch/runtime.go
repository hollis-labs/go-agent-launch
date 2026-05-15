package agentlaunch

// RuntimeKind is the typed string enum that names the runtime lifecycle
// shape a launched session uses. The set mirrors the four runtime kinds
// recognised by go-agent-sessions (subprocess-per-turn, PTY,
// streaming-stdio, JSON-RPC-stdio) so a CompiledLaunch can be translated
// to a go-agent-sessions StartOptions / Runtime selection without an
// intermediate string protocol.
//
// Valid returns true for the four declared constants only; all other
// values — including the zero value — are rejected by
// LaunchPlan.Validate via ErrUnknownRuntime.
//
// The provider × runtime support matrix (which provider IDs legally
// pair with which RuntimeKind values) is owned by a sibling subpackage
// and intentionally NOT enforced here. RuntimeKind.Valid only tells you
// the value is one of the recognised constants.
type RuntimeKind string

const (
	// RuntimeSubprocess names the turn-based subprocess-per-turn runtime:
	// each turn spawns a fresh child, exits when the turn is done, and
	// the next turn re-spawns. Maps to go-agent-sessions's adapter runtime.
	RuntimeSubprocess RuntimeKind = "subprocess"

	// RuntimePTY names the long-lived PTY runtime: a single child runs
	// for the whole session behind a pseudo-terminal master; SendInput
	// writes bytes to the PTY. Maps to go-agent-sessions's PTY runtime.
	RuntimePTY RuntimeKind = "pty"

	// RuntimeStreamingStdio names the long-lived non-PTY runtime that
	// speaks NDJSON over stdin/stdout. Maps to go-agent-sessions's
	// streaming-stdio runtime.
	RuntimeStreamingStdio RuntimeKind = "streaming-stdio"

	// RuntimeJsonRpcStdio names the long-lived non-PTY runtime that
	// speaks JSON-RPC 2.0 over stdin/stdout. Maps to go-agent-sessions's
	// jsonrpc-stdio runtime.
	RuntimeJsonRpcStdio RuntimeKind = "jsonrpc-stdio"
)

// Valid reports whether the receiver is one of the four declared
// RuntimeKind constants. The zero value ("") is not valid.
func (r RuntimeKind) Valid() bool {
	switch r {
	case RuntimeSubprocess, RuntimePTY, RuntimeStreamingStdio, RuntimeJsonRpcStdio:
		return true
	default:
		return false
	}
}

// String returns the underlying token (e.g. "subprocess"). The zero
// value returns "".
func (r RuntimeKind) String() string {
	return string(r)
}
