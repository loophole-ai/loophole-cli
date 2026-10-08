package version

import (
	"runtime/debug"
	"strings"
)

// Build-time parameters set via -ldflags
var Version = "1.0.0"

// defaultVersion marks the value compiled in rather than injected at build
// time. Release builds overwrite it with -X, and only then is the build
// metadata consulted.
const defaultVersion = "1.0.0"

// A user may install loophole using `go install github.com/loophole-ai/loophole-cli@latest`.
// without -ldflags, in which case the version above is unset. As a workaround
// we use the embedded build version that *is* set when using `go install` (and
// is only set for `go install` and not for `go build`).
//
// The build metadata is only consulted when no version was injected. A plain
// `go build` inside a checkout reports a pseudo version such as
// v0.0.0-20260101000000-abcdef+dirty, which was overriding the real version
// and made a release binary report a different version than the tag it was
// built from.
func init() {
	if Version != defaultVersion {
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		// < go v1.18
		return
	}
	mainVersion := info.Main.Version
	if mainVersion == "" || mainVersion == "(devel)" || mainVersion == "unknown" {
		// bin not built using `go install` or version not embedded
		return
	}
	// A pseudo version describes the checkout rather than a release, so it is
	// less useful than the compiled in number.
	if strings.HasPrefix(mainVersion, "v0.0.0-") {
		return
	}
	// bin built using `go install`
	Version = mainVersion
}
