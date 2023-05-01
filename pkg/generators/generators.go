package generators

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/index"
	"github.com/kanopy-platform/code-generator/pkg/generators/tags"
	log "github.com/sirupsen/logrus"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

type BuilderFactory interface {
	NewBuilder(pkg *types.Package, packageIndex *PackageTypeIndex) generator.Generator
}

type PackageTypeIndex struct {
	TypesByTypePath map[string]*types.Type
	Gomod           string
}

func NewPackageTypeIndex() *PackageTypeIndex {
	return &PackageTypeIndex{
		TypesByTypePath: map[string]*types.Type{},
		Gomod:           "github.com/10gen/kanopy/pkg/builder/",
	}
}

const (
	DefaultNameSystem = "public"
)

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
	Index       *PackageTypeIndex
}

func WithBoilerplate(boilerplate string) func(g *Generators) {
	return func(g *Generators) {
		g.Boilerplate = boilerplate
	}
}

func New(builderFactory BuilderFactory, opts ...func(g *Generators)) *Generators {
	g := &Generators{
		Boilerplate: "",
		Builder:     builderFactory,
		Index:       NewPackageTypeIndex(),
	}
	for _, o := range opts {
		o(g)
	}
	return g
}

func (g *Generators) Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	packages := []*types.Package{}
	for _, v := range context.Inputs {
		pkg := context.Universe[v]
		if tags.IsPackageTagged(pkg.Comments) || doPackageTypesNeedGeneration(pkg) {
			log.Infof("Package: %s marked for generation.", pkg.Name)
			packages = append(packages, pkg)
			buildPackageIndex(g.Index, pkg)
		}
	}

	gp := generator.Packages{}
	for _, pkg := range packages {
		if tags.IsPackageTagged(pkg.Comments) || doPackageTypesNeedGeneration(pkg) {
			log.Infof("Package: %s marked for generation.", pkg.Name)
			gp = append(gp, &generator.DefaultPackage{
				PackageName:   pkg.Name,
				PackagePath:   pkg.Path,
				HeaderText:    []byte(g.Boilerplate),
				FilterFunc:    filterFuncByPackagePath(pkg),
				GeneratorFunc: g.generatorFuncForPackage(pkg),
			})
		}
	}

	return gp
}

func filterFuncByPackagePath(pkg *types.Package) func(c *generator.Context, t *types.Type) bool {
	return func(c *generator.Context, t *types.Type) bool {
		return t.Name.Package == pkg.Path
	}
}

func (g *Generators) generatorFuncForPackage(pkg *types.Package) func(c *generator.Context) []generator.Generator {
	return func(c *generator.Context) []generator.Generator {
		return []generator.Generator{
			g.Builder.NewBuilder(pkg, g.Index),
		}
	}
}

func buildPackageIndex(packageIndex *PackageTypeIndex, pkg *types.Package) {
	packageIndex.TypesByTypePath = index.BuildPackageIndex(packageIndex.TypesByTypePath, pkg)
}

func doPackageTypesNeedGeneration(pkg *types.Package) bool {
	cnt := 0
	for _, t := range pkg.Types {
		if tags.IsTypeEnabled(t) {
			cnt++
			log.Debugf("Type marked for generation: Name: %s \t \n %#v", t.Name, t.Members)
		}
	}
	return cnt > 0
}
