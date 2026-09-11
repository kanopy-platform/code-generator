package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommandGeneratorArgs(t *testing.T) {
	tests := []struct {
		args []string
		want *Args
	}{
		{
			args: []string{"--input-dirs=test", "--output-file-base=zz-gen", "--build-tag=abc"},
			want: &Args{
				InputDirs:          []string{"test"},
				OutputFileBaseName: "zz-gen",
				GeneratedBuildTag:  "abc",
			},
		},
	}

	for _, test := range tests {
		g := &Args{}
		root := NewRootCommand(WithGeneratorArgs(g))

		assert.NoError(t, root.ParseFlags(test.args))
		assert.NoError(t, root.PersistentPreRunE(root, test.args))

		assert.Equal(t, test.want, g)
	}
}

// The flags below were removed in the gengo/v2 migration because nothing
// consumed them. Assert they are rejected rather than silently ignored.
func TestRootCommandRemovedFlags(t *testing.T) {
	for _, flag := range []string{
		"--verify-only",
		"--include-test-files",
		"--output-package=pkg",
		"--trim-path-prefix=src",
		"--go-header-file=myfile",
		"--bounding-dirs=dir",
		"--output-base=./src",
	} {
		t.Run(flag, func(t *testing.T) {
			root := NewRootCommand(WithGeneratorArgs(&Args{}))
			assert.Error(t, root.ParseFlags([]string{flag}))
		})
	}
}
