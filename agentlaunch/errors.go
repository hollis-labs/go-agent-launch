package agentlaunch

import "errors"

// Sentinel errors returned by the Validate methods on LaunchPlan,
// CompiledLaunch, and PreparedLaunch. All sentinels are
// errors.Is-comparable; callers branch on them with errors.Is rather
// than string match.
//
// The errors document which field failed; they do not carry the value
// itself. Validation that needs to report a specific offending value
// wraps the sentinel with fmt.Errorf("%w: %s", sentinel, value).
var (
	// ErrMissingProjectID is returned when LaunchPlan.Project.ID is empty.
	ErrMissingProjectID = errors.New("agentlaunch: missing project id")

	// ErrMissingAgentID is returned when LaunchPlan.Agent.ID is empty.
	ErrMissingAgentID = errors.New("agentlaunch: missing agent id")

	// ErrMissingProviderID is returned when LaunchPlan.Provider.ID is empty.
	ErrMissingProviderID = errors.New("agentlaunch: missing provider id")

	// ErrUnsupportedWorkspaceMode is returned when LaunchPlan.Workspace.Mode
	// is not one of the four declared WorkspaceMode values.
	ErrUnsupportedWorkspaceMode = errors.New("agentlaunch: unsupported workspace mode")

	// ErrUnsupportedLaunchMode is returned when LaunchPlan.Mode is not one
	// of the three declared LaunchMode values.
	ErrUnsupportedLaunchMode = errors.New("agentlaunch: unsupported launch mode")

	// ErrUnknownRuntime is returned when LaunchPlan.Runtime is not one of
	// the four declared RuntimeKind values. This is a value-level check;
	// the provider × runtime support matrix (which legal (provider, runtime)
	// pairs exist) lives in a sibling package and is NOT enforced here.
	ErrUnknownRuntime = errors.New("agentlaunch: unknown runtime kind")

	// ErrMissingBootProfile is returned when neither
	// LaunchPlan.BootProfile.CatalogPath nor LaunchPlan.BootProfile.Inline
	// is set.
	ErrMissingBootProfile = errors.New("agentlaunch: boot profile must set CatalogPath or Inline")

	// ErrUnsupportedBootMode is returned when an inline boot profile's
	// BootMode is not one of "none", "stdin", or "planted".
	ErrUnsupportedBootMode = errors.New("agentlaunch: unsupported inline boot mode")

	// ErrUnsafeInjectionTarget is returned when InjectionSpec.BootDirOverlay
	// contains a key that escapes the bootdir: an absolute path, a path
	// containing ".." segments, or a path that targets a reserved name
	// like ".git/". Path-safety is enforced defensively at validation
	// time so consumers cannot trick the preparer into writing outside
	// the planted bootdir.
	ErrUnsafeInjectionTarget = errors.New("agentlaunch: unsafe injection overlay target")

	// ErrCompiledMissingPlan is returned by CompiledLaunch.Validate when
	// the embedded Plan has not been populated.
	ErrCompiledMissingPlan = errors.New("agentlaunch: compiled launch has no source plan")

	// ErrPreparedMissingBootDir is returned by PreparedLaunch.Validate
	// when PlantedBootDir is empty — the preparer must have materialized
	// a bootdir before a PreparedLaunch is considered valid.
	ErrPreparedMissingBootDir = errors.New("agentlaunch: prepared launch missing planted bootdir")

	// ErrPreparedMissingWorkspaceDir is returned by PreparedLaunch.Validate
	// when WorkspaceDir is empty.
	ErrPreparedMissingWorkspaceDir = errors.New("agentlaunch: prepared launch missing workspace dir")

	// ErrPreparedMissingArgv is returned by PreparedLaunch.Validate when
	// Argv has length zero — at minimum the spawned binary must be named.
	ErrPreparedMissingArgv = errors.New("agentlaunch: prepared launch missing argv")

	// ErrUnknownNativeFileKind is returned by NativeFile.Validate when
	// NativeFile.Kind is not one of the declared NativeFileKind values.
	ErrUnknownNativeFileKind = errors.New("agentlaunch: unknown native file kind")

	// ErrNativeFileMissingID is returned by NativeFile.Validate when a
	// NativeFileSkill entry has an empty ID — the ID is required because
	// it derives the planted filename (.claude/skills/<ID>.md etc.).
	ErrNativeFileMissingID = errors.New("agentlaunch: native skill file missing id")

	// ErrNativeFileUnsafeID is returned by NativeFile.Validate when a
	// NativeFileSkill entry's ID is not a safe single path segment (it
	// contains a separator, a "." / ".." traversal token, or a rune
	// outside [A-Za-z0-9._-]).
	ErrNativeFileUnsafeID = errors.New("agentlaunch: native skill file id is not a safe path segment")

	// ErrNativeFileMissingRelPath is returned by NativeFile.Validate when
	// a NativeFileRaw entry has an empty RelPath.
	ErrNativeFileMissingRelPath = errors.New("agentlaunch: native raw file missing relpath")
)
