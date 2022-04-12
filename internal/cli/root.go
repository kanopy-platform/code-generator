package cli

import (
	"github.com/kanopy-platform/code-generator/pkg/generators"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		GeneratorArgs: gengoargs.Default(),
	}

	for _, opt := range opts {
		opt(rootCommand)
	}

	rootCommand.GeneratorArgs.OutputFileBaseName = "zz_generated_builders"

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

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	return setupGlobalLogLevel()
}

func (r *rootCommand) runE(cmd *cobra.Command, args []string) error {

	return r.GeneratorArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
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
