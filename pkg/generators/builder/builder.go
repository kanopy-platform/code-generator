package builder

import (
	"go/token"
	"io"
	"strings"

	"github.com/kanopy-platform/code-generator/pkg/generators/snippets"
	"github.com/kanopy-platform/code-generator/pkg/generators/tags"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"

	log "github.com/sirupsen/logrus"
	"k8s.io/gengo/generator"
)

type BuilderPatternGenerator struct {
	generator.DefaultGen
	pkgToBuild *types.Package
	allTypes   bool
	imports    namer.ImportTracker
}

type BuilderPatternGeneratorFactory struct {
	OutputFileBaseName string
}

func (d *BuilderPatternGeneratorFactory) NewBuilder(pkg *types.Package) generator.Generator {
	return &BuilderPatternGenerator{
		DefaultGen: generator.DefaultGen{
			OptionalName: d.OutputFileBaseName,
		},
		pkgToBuild: pkg,
		allTypes:   isAllTypes(pkg),
		imports:    newImportTracker(),
	}
}

func isAllTypes(pkg *types.Package) bool {
	return tags.IsPackageTagged(tags.Extract(pkg.Comments))
}

func newImportTracker() namer.ImportTracker {
	tracker := namer.NewDefaultImportTracker(types.Name{})
	tracker.IsInvalidType = func(*types.Type) bool { return false }
	tracker.LocalName = func(name types.Name) string { return golangNameToImportAlias(&tracker, name) }
	tracker.PrintImport = func(path, name string) string { return name + " \"" + path + "\"" }
	return &tracker
}

func golangNameToImportAlias(tracker namer.ImportTracker, t types.Name) string {
	path := t.Package
	dirs := strings.Split(path, namer.GoSeperator)
	const immediateParentPosition = 2

	for n := len(dirs) - immediateParentPosition; n >= 0; n-- {
		name := sanitizeGoImportDir(sliceFromParent(dirs, n))

		if isGolangNameImportTracked(tracker, name) {
			continue
		}

		return prefixGoKeywordsWithUnderscore(name)
	}
	return ""
}

func sliceFromParent(in []string, parent int) []string {
	return in[parent:]
}

func isGolangNameImportTracked(tracker namer.ImportTracker, name string) bool {
	_, found := tracker.PathOf(name)
	return found
}

func prefixGoKeywordsWithUnderscore(name string) string {
	out := name
	if token.Lookup(name).IsKeyword() {
		out = "_" + name
	}
	return out
}

func sanitizeGoImportDir(dirs []string) string {
	name := strings.Join(dirs, "")
	return pathToLegalGoName(name)
}

func pathToLegalGoName(in string) string {
	out := strings.Replace(in, "_", "", -1)
	out = strings.Replace(out, ".", "", -1)
	return strings.Replace(out, "-", "", -1)
}

func (b *BuilderPatternGenerator) Init(c *generator.Context, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")
	sw.Do(snippets.GenerateMergeMapStringString(), nil)
	return sw.Error()
}

func (b *BuilderPatternGenerator) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	log.Debugf("Checking type: %s", t.Name.Name)
	if !b.needsGeneration(t) {
		return nil
	}
	log.Infof("Generating type: %s", t.Name.Name)

	sw := generator.NewSnippetWriter(w, c, "$", "$")

	sw.Do("//auto generated\n", nil) // TODO can be removed once setters are constructed for types
	if hasTypeMetaEmbedded(t) {
		parentTypeOfTypeMeta := getParentOfTypeMeta(t)
		b.imports.AddType(parentTypeOfTypeMeta)
		b.imports.AddType(getTypeMetaFromType(parentTypeOfTypeMeta))
		sw.Do(snippets.GenerateDeepCopy(t))
		sw.Do(snippets.GenerateMarshalJSON(t, b.imports.LocalNameOf(parentTypeOfTypeMeta.Name.Package)))
	}

	// TODO generate setters for struct

	return sw.Error()
}

func (b *BuilderPatternGenerator) needsGeneration(t *types.Type) bool {
	if b.doesTypeOptout(t) || (!b.doesTypeNeedGeneration(t) && !b.allTypes) {
		return false
	}

	return true
}

func (b *BuilderPatternGenerator) doesTypeOptout(t *types.Type) bool {
	return tags.IsTypeOptedOut(t)
}

func (b *BuilderPatternGenerator) doesTypeNeedGeneration(t *types.Type) bool {
	v := tags.IsTypeEnabled(t)
	log.Debugf("Type: %s, Tag: %v", t.Name, v)
	return v
}

func hasTypeMetaEmbedded(t *types.Type) bool {
	if p := getParentOfTypeMeta(t); p != nil {
		return true
	}
	return false
}

func getParentOfTypeMeta(t *types.Type) *types.Type {
	for _, m := range t.Members {
		if m.Embedded {
			if mm := getTypeMetaFromType(m.Type); mm != nil {
				return m.Type
			}
		}
	}
	return nil
}

func getTypeMetaFromType(t *types.Type) *types.Type {
	for _, mm := range t.Members {
		if mm.Name == "TypeMeta" {
			return mm.Type
		}
	}
	return nil
}

func (b *BuilderPatternGenerator) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(b.pkgToBuild.Path, b.imports),
	}
}

func (b *BuilderPatternGenerator) Imports(c *generator.Context) (imports []string) {
	importLines := []string{}
	for _, singleImport := range b.imports.ImportLines() {
		if b.isNotTargetPackage(singleImport) {
			importLines = append(importLines, singleImport)
		}
	}

	return importLines
}

func (b *BuilderPatternGenerator) isNotTargetPackage(pkg string) bool {
	if pkg == b.pkgToBuild.Path {
		return false
	}
	if strings.HasSuffix(pkg, "\""+b.pkgToBuild.Path+"\"") {
		return false
	}
	return true
}
