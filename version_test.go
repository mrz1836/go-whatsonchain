package whatsonchain

import (
	"bufio"
	"os"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestModuleVersion covers every way this module can appear in a program's build info
func TestModuleVersion(t *testing.T) {
	t.Parallel()

	const otherModule = "example.com/app"

	tests := []struct {
		name string
		info *debug.BuildInfo
		want string
	}{
		{"no build info", nil, develVersion},
		{
			"main module with a release version",
			&debug.BuildInfo{Main: debug.Module{Path: modulePath, Version: "v1.4.0"}},
			"v1.4.0",
		},
		{
			"main module with a VCS-stamped pseudo-version",
			&debug.BuildInfo{Main: debug.Module{Path: modulePath, Version: "v1.4.1-0.20261119120000-abcdef123456+dirty"}},
			"v1.4.1-0.20261119120000-abcdef123456+dirty",
		},
		{
			"main module in development",
			&debug.BuildInfo{Main: debug.Module{Path: modulePath, Version: "(devel)"}},
			develVersion,
		},
		{
			"main module without a version",
			&debug.BuildInfo{Main: debug.Module{Path: modulePath}},
			develVersion,
		},
		{
			"dependency with a release version",
			&debug.BuildInfo{
				Main: debug.Module{Path: otherModule, Version: "v0.1.0"},
				Deps: []*debug.Module{
					{Path: "github.com/stretchr/testify", Version: "v1.11.1"},
					{Path: modulePath, Version: "v1.4.0"},
				},
			},
			"v1.4.0",
		},
		{
			"dependency replaced by another version",
			&debug.BuildInfo{
				Main: debug.Module{Path: otherModule},
				Deps: []*debug.Module{
					{Path: modulePath, Version: "v1.4.0", Replace: &debug.Module{Path: "example.com/fork", Version: "v1.5.0-fork"}},
				},
			},
			"v1.5.0-fork",
		},
		{
			"dependency replaced by a local directory",
			&debug.BuildInfo{
				Main: debug.Module{Path: otherModule},
				Deps: []*debug.Module{
					{Path: modulePath, Version: "v1.4.0", Replace: &debug.Module{Path: "../go-whatsonchain"}},
				},
			},
			develVersion,
		},
		{
			"dependency without a version",
			&debug.BuildInfo{
				Main: debug.Module{Path: otherModule},
				Deps: []*debug.Module{{Path: modulePath}},
			},
			develVersion,
		},
		{
			"nil dependency entries are skipped",
			&debug.BuildInfo{
				Main: debug.Module{Path: otherModule},
				Deps: []*debug.Module{nil, {Path: modulePath, Version: "v1.3.1"}},
			},
			"v1.3.1",
		},
		{
			"a different major version is not this module",
			&debug.BuildInfo{
				Main: debug.Module{Path: otherModule},
				Deps: []*debug.Module{{Path: modulePath + "/v2", Version: "v2.0.0"}},
			},
			develVersion,
		},
		{
			"module missing from the build",
			&debug.BuildInfo{
				Main: debug.Module{Path: otherModule, Version: "v0.1.0"},
				Deps: []*debug.Module{{Path: "github.com/stretchr/testify", Version: "v1.11.1"}},
			},
			develVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, moduleVersion(tt.info))
		})
	}
}

// TestModulePathMatchesGoMod guards against modulePath drifting from go.mod (for
// example after a /v2 move), which would silently turn the version into "(devel)".
func TestModulePathMatchesGoMod(t *testing.T) {
	t.Parallel()

	f, err := os.Open("go.mod")
	require.NoError(t, err)
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if path, ok := strings.CutPrefix(strings.TrimSpace(scanner.Text()), "module "); ok {
			assert.Equal(t, modulePath, strings.TrimSpace(path))
			return
		}
	}
	require.NoError(t, scanner.Err())
	t.Fatal("no module directive in go.mod")
}
