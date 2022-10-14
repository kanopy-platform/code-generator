package cli

import (
	"github.com/kanopy-platform/code-generator/pkg/generators"
	"github.com/spf13/pflag"
	gengoargs "k8s.io/gengo/args"
)

func flagLogLevel(flags *pflag.FlagSet) {
	flags.String("log-level", "info", "Configure log level")
}

// This replaces the https://github.com/kubernetes/gengo/blob/master/args/args.go#L102 AddFlags method.  The Cobra framework
// uses `-h` shortflag for the "help" method.
// Notable change: 'go-header-file' shortflag was changed from 'h' to 'e'.
func flagGeneratorArgs(fs *pflag.FlagSet, g *gengoargs.GeneratorArgs) {
	fs.StringSliceVarP(&g.InputDirs, "input-dirs", "i", g.InputDirs, "Comma-separated list of import paths to get input types from.")
	fs.StringVarP(&g.OutputBase, "output-base", "o", g.OutputBase, "Output base; defaults to $GOPATH/src/ or ./ if $GOPATH is not set.")
	fs.StringVarP(&g.OutputPackagePath, "output-package", "p", g.OutputPackagePath, "Base package path.")
	fs.StringVarP(&g.OutputFileBaseName, "output-file-base", "O", g.OutputFileBaseName, "Base name (without .go suffix) for output files.")
	fs.StringVarP(&g.GoHeaderFilePath, "go-header-file", "e", g.GoHeaderFilePath, "File containing boilerplate header text. The string YEAR will be replaced with the current 4-digit year.")
	fs.BoolVar(&g.VerifyOnly, "verify-only", g.VerifyOnly, "If true, only verify existing output, do not write anything.")
	fs.BoolVar(&g.IncludeTestFiles, "include-test-files", g.IncludeTestFiles, "If true, includes _test.go files.")
	fs.StringVar(&g.GeneratedBuildTag, "build-tag", g.GeneratedBuildTag, "A Go build tag to use to identify files generated by this command. Should be unique.")
	fs.StringVar(&g.TrimPathPrefix, "trim-path-prefix", g.TrimPathPrefix, "If set, trim the specified prefix from --output-package when generating files.")
}

func flagCustomGeneratorArgs(fs *pflag.FlagSet, customArgs *generators.CustomArgs) {
	fs.StringSliceVar(&customArgs.BoundingDirs, "bounding-dirs", customArgs.BoundingDirs, "specify directories to bound the generation")
}
