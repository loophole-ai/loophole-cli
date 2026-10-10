package version

import (
	"runtime/debug"
	"strings"
)

// unset is the placeholder compiled in. A release build replaces it with
// -X, so a value that is still this placeholder means no version was injected.
//
// The placeholder has to be distinguishable from every real version, including
// the one compiled in below. An earlier attempt used the release number itself
// as the marker, which cannot tell "not injected" apart from "injected with
// the same number", so a release binary fell through and reported a pseudo
// version instead of the tag it was built from.
const unset = "unset-at-build-time"

// defaultVersion is what a build reports when nothing was injected at build
// time. Release builds overwrite it.
var Version = unset

// version is resolved once, at startup, in this order:
//
//  1. the value injected with -X, which is how every release is built
//  2. a real tag recorded by `go install module@version`
//  3. the compiled in fallback below
func init() {
	if Version != unset {
		return
	}

	Version = defaultVersion

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	switch recorded := info.Main.Version; {
	case recorded == "" || recorded == "(devel)" || recorded == "unknown":
	case isPseudoVersion(recorded):
	default:
		Version = recorded
	}
}

// isPseudoVersion reports whether a recorded version describes a checkout
// rather than a release.
//
// Go stamps binaries built from a local repository with a pseudo version, which
// carries a build timestamp and a commit hash. Those two forms are the ones
// seen here: v0.0.0-20060102150405-abcdefabcdef when the tree has no tag, and
// v1.0.1-0.20060102150405-abcdefabcdef when it is one commit past a tag.
func isPseudoVersion(v string) bool {
	return strings.HasPrefix(v, "v0.0.0-") || strings.Contains(v, "-0.") || strings.Contains(v, "+")
}

// defaultVersion is reported when the build recorded no usable version.
const defaultVersion = "1.0.4"
