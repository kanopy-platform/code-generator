package generators

import (
	"path"
	"strings"

	"github.com/kanopy-platform/code-generator/pkg/generators/index"
	"github.com/kanopy-platform/code-generator/pkg/generators/tags"
	log "github.com/sirupsen/logrus"
	"k8s.io/gengo/v2/generator"
	"k8s.io/gengo/v2/namer"
	"k8s.io/gengo/v2/types"
)

type BuilderFactory interface {
	NewBuilder(pkg *types.Package, packageIndex *PackageTypeIndex) generator.Generator
}

type PackageTypeIndex struct {
	TypesByTypePath map[string]*types.Type
	// PackageRoot is the directory the packages being generated live under.
	PackageRoot string
}

func NewPackageTypeIndex() *PackageTypeIndex {
	return &PackageTypeIndex{
		TypesByTypePath: map[string]*types.Type{},
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

func (g *Generators) Targets(context *generator.Context) []generator.Target {
	gp := []generator.Target{}
	for _, v := range context.Inputs {
		pkg := context.Universe[v]
		if !tags.IsPackageTagged(pkg.Comments) && !doPackageTypesNeedGeneration(pkg) {
			continue
		}

		log.Infof("Package: %s marked for generation.", pkg.Name)
		buildPackageIndex(g.Index, pkg)
		g.Index.PackageRoot = commonDir(g.Index.PackageRoot, path.Dir(pkg.Path))

		gp = append(gp, &generator.SimpleTarget{
			PkgName:        pkg.Name,
			PkgPath:        pkg.Path,
			PkgDir:         pkg.Dir,
			HeaderComment:  []byte(g.Boilerplate),
			FilterFunc:     filterFuncByPackagePath(pkg),
			GeneratorsFunc: g.generatorFuncForPackage(pkg),
		})
	}

	return gp
}

// commonDir returns the deepest directory a and b share, or b if a is empty.
func commonDir(a, b string) string {
	if a == "" {
		return b
	}

	as := strings.Split(a, namer.GoSeparator)
	bs := strings.Split(b, namer.GoSeparator)

	shared := 0
	for shared < len(as) && shared < len(bs) && as[shared] == bs[shared] {
		shared++
	}

	return strings.Join(as[:shared], namer.GoSeparator)
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
	for _, t := range pkg.Types {
		if tags.IsTypeEnabled(t) {
			return true
		}
	}
	return false
}
