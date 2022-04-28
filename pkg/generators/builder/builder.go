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
	pkgToBuild      *types.Package
	allTypes        bool
	imports         namer.ImportTracker
	typesToGenerate map[string]*types.Type
}

type BuilderPatternGeneratorFactory struct {
	OutputFileBaseName string
}

func (d *BuilderPatternGeneratorFactory) NewBuilder(pkg *types.Package) generator.Generator {
	return &BuilderPatternGenerator{
		DefaultGen: generator.DefaultGen{
			OptionalName: d.OutputFileBaseName,
		},
		pkgToBuild:      pkg,
		allTypes:        isAllTypes(pkg),
		imports:         newImportTracker(),
		typesToGenerate: map[string]*types.Type{},
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
	log.Infof("Generating type: %s", t.Name.Name)

	sw := generator.NewSnippetWriter(w, c, "$", "$")

	// generate the constructor
	if hasObjectMetaEmbedded(t) {
		parentTypeOfObjectMeta := getParentOfEmbeddedType(t, "ObjectMeta")
		b.imports.AddType(parentTypeOfObjectMeta)
		b.imports.AddType(getMemberFromType(parentTypeOfObjectMeta, "ObjectMeta"))
		sw.Do(snippets.GenerateConstructorForObjectMeta(t))
	} else {
		sw.Do(snippets.GenerateEmptyConstructor(t))
	}

	if hasTypeMetaEmbedded(t) {
		parentTypeOfTypeMeta := getParentOfEmbeddedType(t, "TypeMeta")
		b.imports.AddType(parentTypeOfTypeMeta)
		b.imports.AddType(getMemberFromType(parentTypeOfTypeMeta, "TypeMeta"))
		sw.Do(snippets.GenerateDeepCopy(t))
	}

	// generate setters for struct
	if hasObjectMetaEmbedded(t) {
		parentTypeOfObjectMeta := getParentOfEmbeddedType(t, "ObjectMeta")
		objectMetaType := getMemberFromType(parentTypeOfObjectMeta, "ObjectMeta")
		b.generateSettersForType(sw, t, objectMetaType)
		b.generateSettersForType(sw, t, parentTypeOfObjectMeta)
	} else {
		b.generateSettersForType(sw, t, t.Members[0].Type)
	}

	return sw.Error()
}

func (b *BuilderPatternGenerator) generateSettersForType(sw *generator.SnippetWriter, root *types.Type, parent *types.Type) {
	for _, m := range parent.Members {
		if m.Embedded {
			continue
		}

		if m.Type.Kind == types.Map {
			sw.Do(snippets.GenerateSetterForMap(root, parent, m))
		} else if m.Type.Kind == types.Slice {
			sliceType := m.Type.Elem

			switch sliceType.Kind {
			case types.Struct:
				// skip adding setters for un-enabled structs
				inputType := b.getTypeEnabledForGeneration(sliceType)
				if inputType != nil {
					sw.Do(snippets.GenerateSetterForEmbeddedSlice(root, parent, m, inputType))
				}
			default:
				sw.Do(snippets.GenerateSetterForMemberSlice(root, parent, m))
			}
		} else if m.Type.Kind == types.Struct {
			// skip adding setters for un-enabled structs
			inputType := b.getTypeEnabledForGeneration(m.Type)
			if inputType != nil {
				sw.Do(snippets.GenerateSetterForEmbeddedStruct(root, parent, m, inputType))
			}
		} else if m.Type.Kind == types.Pointer && m.Type.Elem.Kind == types.Builtin {
			sw.Do(snippets.GenerateSetterForPointerToBuiltinType(root, parent, m))
		} else {
			sw.Do(snippets.GenerateSetterForPrimitiveType(root, parent, m))
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

func (b *BuilderPatternGenerator) getTypeEnabledForGeneration(t *types.Type) *types.Type {
	typeName := t.Name.String()
	if wrapper, ok := b.typesToGenerate[typeName]; ok {
		return wrapper
	}
	return nil
}

func hasTypeMetaEmbedded(t *types.Type) bool {
	if p := getParentOfEmbeddedType(t, "TypeMeta"); p != nil {
		return true
	}
	return false
}

func hasObjectMetaEmbedded(t *types.Type) bool {
	if p := getParentOfEmbeddedType(t, "ObjectMeta"); p != nil {
		return true
	}
	return false
}

func getParentOfEmbeddedType(t *types.Type, name string) *types.Type {
	for _, m := range t.Members {
		if m.Embedded {
			if mm := getMemberFromType(m.Type, name); mm != nil {
				return m.Type
			}
		}
	}
	return nil
}

func getMemberFromType(t *types.Type, name string) *types.Type {
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
			typeName := m.Type.Name.String()
			b.typesToGenerate[typeName] = t
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
