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
	pkgToBuild   *types.Package
	allTypes     bool
	imports      namer.ImportTracker
	enabledTypes map[string]*types.Type
}

type BuilderPatternGeneratorFactory struct {
	OutputFileBaseName string
}

func (d *BuilderPatternGeneratorFactory) NewBuilder(pkg *types.Package) generator.Generator {
	return &BuilderPatternGenerator{
		DefaultGen: generator.DefaultGen{
			OptionalName: d.OutputFileBaseName,
		},
		pkgToBuild:   pkg,
		allTypes:     isAllTypes(pkg),
		imports:      newImportTracker(),
		enabledTypes: map[string]*types.Type{},
	}
}

func isAllTypes(pkg *types.Package) bool {
	return tags.IsPackageTagged(pkg.Comments)
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
	log.Infof("Generating type: %s", t.Name.Name)

	sw := generator.NewSnippetWriter(w, c, "$", "$")

	// generate constructor and setters
	if hasObjectMetaEmbedded(t) {
		parentTypeOfObjectMeta := getParentOfEmbeddedType(t, ObjectMeta)
		objectMetaType := getMemberTypeFromType(parentTypeOfObjectMeta, ObjectMeta)
		b.imports.AddType(parentTypeOfObjectMeta)
		b.imports.AddType(objectMetaType)
		sw.Do(snippets.GenerateConstructorForObjectMeta(t))
		b.generateSettersForType(sw, t, objectMetaType)
	} else {
		sw.Do(snippets.GenerateEmptyConstructor(t, true))
	}

	for _, member := range t.Members {
		b.generateSettersForType(sw, t, member.Type)
	}

	// generate deepcopy
	if hasTypeMetaEmbedded(t) {
		parentTypeOfTypeMeta := getParentOfEmbeddedType(t, TypeMeta)
		b.imports.AddType(parentTypeOfTypeMeta)
		b.imports.AddType(getMemberTypeFromType(parentTypeOfTypeMeta, TypeMeta))
		sw.Do(snippets.GenerateDeepCopy(t))
	}

	return sw.Error()
}

func (b *BuilderPatternGenerator) generateSettersForType(sw *generator.SnippetWriter, root *types.Type, parent *types.Type) {
	setter := snippets.NewSetter(root, parent, true)

	for _, m := range parent.Members {
		if m.Embedded || !includeMember(parent, m) {
			continue
		}

		switch {
		case m.Type.Kind == types.Map:
			sw.Do(setter.GenerateSetterForMap(m))
		case m.Type.Kind == types.Slice:
			sliceType := m.Type.Elem
			switch sliceType.Kind {
			case types.Struct:
				if b.isTypeEnabled(sliceType) {
					sw.Do(setter.GenerateSetterForEmbeddedSlice(m, b.getWrapperType(sliceType)))
				}
			default:
				sw.Do(setter.GenerateSetterForMemberSlice(m))
			}
		case m.Type.Kind == types.Struct:
			if b.isTypeEnabled(m.Type) {
				sw.Do(setter.GenerateSetterForEmbeddedStruct(m, b.getWrapperType(m.Type)))
			}
		case m.Type.Kind == types.Pointer && m.Type.Elem.Kind == types.Builtin:
			sw.Do(setter.GenerateSetterForPointerToBuiltinType(m))
		default:
			sw.Do(setter.GenerateSetterForPrimitiveType(m))
		}
	}
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

func (b *BuilderPatternGenerator) isTypeEnabled(t *types.Type) bool {
	typeName := t.Name.String()
	_, exists := b.enabledTypes[typeName]
	return exists
}

func (b *BuilderPatternGenerator) getWrapperType(t *types.Type) *types.Type {
	typeName := t.Name.String()
	if parent, ok := b.enabledTypes[typeName]; ok {
		return parent
	}
	return nil
}

func hasTypeMetaEmbedded(t *types.Type) bool {
	if p := getParentOfEmbeddedType(t, TypeMeta); p != nil {
		return true
	}
	return false
}

func hasObjectMetaEmbedded(t *types.Type) bool {
	if p := getParentOfEmbeddedType(t, ObjectMeta); p != nil {
		return true
	}
	return false
}

func getParentOfEmbeddedType(t *types.Type, name string) *types.Type {
	for _, m := range t.Members {
		if m.Embedded {
			if mm := getMemberTypeFromType(m.Type, name); mm != nil {
				return m.Type
			}
		}
	}
	return nil
}

func getMemberFromType(t *types.Type, name string) types.Member {
	for _, mm := range t.Members {
		if mm.Name == name {
			return mm
		}
	}
	return types.Member{}
}

func getMemberTypeFromType(t *types.Type, name string) *types.Type {
	for _, mm := range t.Members {
		if mm.Name == name {
			return mm.Type
		}
	}
	return nil
}

func (b *BuilderPatternGenerator) Filter(c *generator.Context, t *types.Type) bool {
	log.Debugf("Checking type: %s", t.Name.Name)
	if !b.needsGeneration(t) {
		return false
	}

	for _, m := range t.Members {
		if m.Embedded {
			childType := m.Type.Name.String()
			b.enabledTypes[childType] = t
		}
	}

	return true
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
