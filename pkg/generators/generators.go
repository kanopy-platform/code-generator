package generators

import (
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
)

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	const prependPackageNames = 1
	return namer.NameSystems{
		"public": namer.NewPublicNamer(prependPackageNames),
		"raw":    namer.NewRawNamer("", nil),
	}
}

func DefaultNameSystem() string {
	return "public"
}

func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	packages := generator.Packages{}
	// TODO - fill in
	return packages
}
