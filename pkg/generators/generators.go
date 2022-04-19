package generators

import (
	"strings"

	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

type BuilderFactory interface {
	NewBuilder(outfileBaseName string, pkg *types.Package, tagName string) generator.Generator
}

const (
	DefaultNameSystem = "public"
	tagName           = "kanopy:builder"
	tagValuePackage   = "package"
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
		tagValue := extractTag(tagName, pkg.Comments)
		if doesPackageNeedGeneration(tagValue) || doPackageTypesNeedGeneration(pkg) {
			packages = append(packages, &generator.DefaultPackage{
				PackageName:   pkg.Name,
				PackagePath:   pkg.Path,
				HeaderText:    []byte(g.Boilerplate),
				FilterFunc:    filterFuncByPackagePath(pkg),
				GeneratorFunc: g.generatorFuncForPackage(arguments.OutputFileBaseName, pkg, tagValue),
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

func (g *Generators) generatorFuncForPackage(baseName string, pkg *types.Package, tagValue string) func(c *generator.Context) []generator.Generator {
	return func(c *generator.Context) []generator.Generator {
		return []generator.Generator{
			g.Builder.NewBuilder(baseName, pkg, tagValue),
		}
	}
}

func doPackageTypesNeedGeneration(pkg *types.Package) bool {
	for _, t := range pkg.Types {
		if doesTypeNeedGeneration(t) {
			return true
		}
	}
	return false
}

func doesTypeNeedGeneration(t *types.Type) bool {
	tag := extractTag(tagName, t.CommentLines)
	return tag == "true"
}

func doesPackageNeedGeneration(tag string) bool {
	return tag == tagValuePackage
}

func extractTag(tag string, comments []string) string {
	vals := types.ExtractCommentTags("+", comments)[tag]
	if len(vals) == 0 {
		return ""
	}

	return getFirstTagValue(vals...)
}

func getFirstTagValue(values ...string) string {
	return strings.Split(values[0], ",")[0]
}
