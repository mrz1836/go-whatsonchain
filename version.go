package whatsonchain

import (
	"runtime/debug"
	"sync"
)

const (
	// modulePath is the import path of this module, used to find its version in
	// the build info of the program that embeds it
	modulePath = "github.com/mrz1836/go-whatsonchain"

	// develVersion is reported when the module version is unknown: a local
	// build, a test binary, or a module replaced by a directory
	develVersion = "(devel)"

	// userAgentPrefix starts the default user agent, followed by the module version
	userAgentPrefix = "go-whatsonchain: "
)

// defaultUserAgent returns the default user agent for all requests, such as
// "go-whatsonchain: v1.4.0".
//
// The version is read once from the build info of the running program, so it
// always matches the release in use. A library cannot be stamped with -ldflags
// and releases are created from tags alone, so a hard-coded version constant
// would drift from the real release.
var defaultUserAgent = sync.OnceValue(func() string { //nolint:gochecknoglobals // computed once, read-only
	info, _ := debug.ReadBuildInfo()
	return userAgentPrefix + moduleVersion(info)
})

// moduleVersion returns the version of this module recorded in the build info:
// the main module's version when this module is being built directly, or the
// version of the dependency (or of its replacement) when it is imported. It
// returns develVersion when the version is unknown.
func moduleVersion(info *debug.BuildInfo) string {
	if info == nil {
		return develVersion
	}

	if info.Main.Path == modulePath {
		return knownVersion(info.Main.Version)
	}

	for _, dep := range info.Deps {
		if dep == nil || dep.Path != modulePath {
			continue
		}
		if dep.Replace != nil {
			// The replacement is the code that runs; a directory replacement has no version
			return knownVersion(dep.Replace.Version)
		}
		return knownVersion(dep.Version)
	}

	return develVersion
}

// knownVersion returns version, or develVersion when it is empty or already "(devel)"
func knownVersion(version string) string {
	if version == "" {
		return develVersion
	}
	return version
}
