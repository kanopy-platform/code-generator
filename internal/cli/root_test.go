package cli

import (
	"testing"

	"github.com/kanopy-platform/code-generator/pkg/generators"
	"github.com/stretchr/testify/assert"
	gengoargs "k8s.io/gengo/args"
)

func TestRootCommandGeneratorArgs(t *testing.T) {
	tests := []struct {
		args []string
		want *gengoargs.GeneratorArgs
	}{
		{
			args: []string{"--bounding-dirs=dir", "--input-dirs=test", "--output-base=./src", "--output-package=pkg", "--output-file-base=zz-gen", "--go-header-file=myfile", "--verify-only", "--build-tag=abc", "--trim-path-prefix=src"},
			want: func() *gengoargs.GeneratorArgs {
				g := gengoargs.Default()

				g.CustomArgs = &generators.CustomArgs{BoundingDirs: []string{"dir"}}
				g.InputDirs = []string{"test"}
				g.OutputBase = "./src"
				g.OutputPackagePath = "pkg"
				g.OutputFileBaseName = "zz-gen"
				g.GoHeaderFilePath = "myfile"
				g.VerifyOnly = true
				g.GeneratedBuildTag = "abc"
				g.TrimPathPrefix = "src"

				return g
			}(),
		},
	}

	for _, test := range tests {
		g := gengoargs.Default()
		root := NewRootCommand(WithGeneratorArgs(g))

		assert.NoError(t, root.ParseFlags(test.args))
		assert.NoError(t, root.PersistentPreRunE(root, test.args))

		assert.Equal(t, test.want, g)
	}

}
