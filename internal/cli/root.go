package cli

import (
	"github.com/kanopy-platform/code-generator/pkg/generators"
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
		Use:   "kanopy-codegen",
		Short: "Kanopy Builder code generator",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			return nil
		},
		RunE: rootCommand.runE,
	}

	flags := cmd.PersistentFlags()
	flagLogLevel(flags)
	flags.StringSliceP("bounding-dirs", "b", []string{}, "specify directories to bound the generation")
	flagGeneratorArgs(flags, rootCommand.GeneratorArgs)

	return cmd
}

func (r *rootCommand) runE(cmd *cobra.Command, args []string) error {
	// TODO add argument to override outputfilebasename from args

	return r.GeneratorArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	)
}
