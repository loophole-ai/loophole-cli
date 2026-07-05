package version

import "runtime/debug"

// Build-time parameters set via -ldflags
var Version = "1.7.6"

// A user may install loophole using `go install github.com/loophole-ai/loophole-cli@latest`.
// without -ldflags, in which case the version above is unset. As a workaround
// we use the embedded build version that *is* set when using `go install` (and
// is only set for `go install` and not for `go build`).
func init() {
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
	// bin built using `go install`
	Version = mainVersion
}
