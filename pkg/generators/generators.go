package generators

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/tags"
	log "github.com/sirupsen/logrus"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

type BuilderFactory interface {
	NewBuilder(pkg *types.Package) generator.Generator
}

const (
	DefaultNameSystem = "public"
)

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	const prependPackageNames = 1
	return namer.NameSystems{
		"public": namer.NewPublicNamer(prependPackageNames),
		"raw":    namer.NewRawNamer("", nil),
	}
}

type Generators struct {
	Boilerplate string
	Builder     BuilderFactory
}

func WithBoilerplate(boilerplate string) func(g *Generators) {
	return func(g *Generators) {
		g.Boilerplate = boilerplate
	}
}

func New(builderFactory BuilderFactory, opts ...func(g *Generators)) *Generators {
	g := &Generators{Builder: builderFactory}
	for _, o := range opts {
		o(g)
	}
	return g
}

func (g *Generators) Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	packages := generator.Packages{}
	for _, v := range context.Inputs {
		pkg := context.Universe[v]
		tagValue := tags.Extract(pkg.Comments)
		if tags.IsPackageTagged(tagValue) || doPackageTypesNeedGeneration(pkg) {
			log.Infof("Package: %s marked for generation.", pkg.Name)
			packages = append(packages, &generator.DefaultPackage{
				PackageName:   pkg.Name,
				PackagePath:   pkg.Path,
				HeaderText:    []byte(g.Boilerplate),
				FilterFunc:    filterFuncByPackagePath(pkg),
				GeneratorFunc: g.generatorFuncForPackage(pkg),
			})
		}
	}

	return packages
}

func filterFuncByPackagePath(pkg *types.Package) func(c *generator.Context, t *types.Type) bool {
	return func(c *generator.Context, t *types.Type) bool {
		return t.Name.Package == pkg.Path
	}
}

func (g *Generators) generatorFuncForPackage(pkg *types.Package) func(c *generator.Context) []generator.Generator {
	return func(c *generator.Context) []generator.Generator {
		return []generator.Generator{
			g.Builder.NewBuilder(pkg),
		}
	}
}

func doPackageTypesNeedGeneration(pkg *types.Package) bool {
	for _, t := range pkg.Types {
		if tags.IsTypeEnabled(t) {
			return true
		}
	}
	return false
}
