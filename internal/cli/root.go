package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kanopy-platform/code-generator/pkg/generators"
	"github.com/kanopy-platform/code-generator/pkg/generators/builder"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/mod/modfile"
	gengoargs "k8s.io/gengo/args"
)

type rootCommand struct {
	GeneratorArgs *gengoargs.GeneratorArgs
}

func WithGeneratorArgs(g *gengoargs.GeneratorArgs) func(*rootCommand) {
	return func(r *rootCommand) {
		r.GeneratorArgs = g
	}
}

func NewRootCommand(opts ...func(*rootCommand)) *cobra.Command {
	rootCommand := &rootCommand{
		GeneratorArgs: gengoargs.Default().WithoutDefaultFlagParsing(),
	}

	for _, opt := range opts {
		opt(rootCommand)
	}

	rootCommand.GeneratorArgs.OutputFileBaseName = "zz_generated_builders"
	// Settings for optimization:
	// - Do not include _test.go files (default)
	// - Add a header to generated files so gengo will ignore them to reduce parse time
	rootCommand.GeneratorArgs.IncludeTestFiles = false
	rootCommand.GeneratorArgs.GeneratedBuildTag = "ignore_autogenerated"

	cmd := &cobra.Command{
		Use:               "kanopy-codegen",
		Short:             "Kanopy Builder code generator",
		PersistentPreRunE: rootCommand.prerun,
		RunE:              rootCommand.runE,
	}

	rootCommand.setupFlags(cmd)

	return cmd
}

func (r *rootCommand) setupFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flagLogLevel(flags)
	flagGeneratorArgs(flags, r.GeneratorArgs)

	customArgs := &generators.CustomArgs{}
	flagCustomGeneratorArgs(flags, customArgs)
	r.GeneratorArgs.CustomArgs = customArgs
}

func (r *rootCommand) prerun(cmd *cobra.Command, args []string) error {

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}

	return setupGlobalLogLevel()
}

func (r *rootCommand) runE(cmd *cobra.Command, args []string) error {
	headerLines := []string{
		fmt.Sprintf("//go:build !%s\n", r.GeneratorArgs.GeneratedBuildTag),
		"/* DO NOT EDIT */",
		"/* autogenerated by kanopy-platform/code-generator */",
		"\n",
	}

	mod, err := packageRoot()
	if err != nil {
		return err
	}

	g := generators.New(&builder.BuilderPatternGeneratorFactory{OutputFileBaseName: r.GeneratorArgs.OutputFileBaseName},
		generators.WithBoilerplate(strings.Join(headerLines, "\n")), generators.WithPackageRoot(mod))
	return r.GeneratorArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem,
		g.Packages,
	)
}

func setupGlobalLogLevel() error {
	// set log level
	logLevel, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		return err
	}

	log.SetLevel(logLevel)
	return nil
}

func packageRoot() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	goModPath, err := findFile(path, "go.mod")
	if err != nil {
		return "", err
	}

	goMod := goModPath + "/go.mod"
	goModF, err := os.ReadFile(goMod)
	if err != nil {
		return "", err
	}

	modFile, err := modfile.Parse(goMod, goModF, nil)
	if err != nil {
		return "", err
	}

	pkgRoot := strings.Replace(path, goModPath, "", 1)
	return modFile.Module.Mod.Path + pkgRoot + "/", nil
}

func findFile(path, file string) (string, error) {
	current, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	fp := filepath.Join(current, file)

	if _, err = os.Stat(fp); err == nil {
		return current, nil
	} else if path == "/" {
		return "", fmt.Errorf("Could not find file %s in path", file)
	}
	return findFile(filepath.Dir(current), file)
}
