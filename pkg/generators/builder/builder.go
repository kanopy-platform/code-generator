package builder

import (
	"go/token"
	"io"
	"strings"

	"github.com/kanopy-platform/code-generator/pkg/generators"
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
	//enabledTypes map[string]*types.Type
	packageIndex *generators.PackageTypeIndex
}

type BuilderPatternGeneratorFactory struct {
	OutputFileBaseName string
}

func (d *BuilderPatternGeneratorFactory) NewBuilder(pkg *types.Package, packageIndex *generators.PackageTypeIndex) generator.Generator {
	return &BuilderPatternGenerator{
		DefaultGen: generator.DefaultGen{
			OptionalName: d.OutputFileBaseName,
		},
		pkgToBuild: pkg,
		allTypes:   isAllTypes(pkg),
		imports:    newImportTracker(packageIndex),
		//enabledTypes: map[string]*types.Type{},
		packageIndex: packageIndex,
	}
}

func isAllTypes(pkg *types.Package) bool {
	return tags.IsPackageTagged(pkg.Comments)
}

func newImportTracker(index *generators.PackageTypeIndex) namer.ImportTracker {
	tracker := namer.NewDefaultImportTracker(types.Name{})
	tracker.IsInvalidType = func(*types.Type) bool { return false }
	tracker.LocalName = func(name types.Name) string { return golangNameToImportAlias(&tracker, name) }
	tracker.PrintImport = func(path, name string) string {
		path = strings.Replace(path, "./", index.Gomod, 1)
		return name + " \"" + path + "\""
	}
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
	sw.Do(snippets.GenerateVariadicBool(), nil)
	sw.Do(snippets.GenerateBoolPointer(), nil)
	return sw.Error()
}

func (b *BuilderPatternGenerator) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	log.Infof("Generating type: %s", t.Name.Name)

	sw := generator.NewSnippetWriter(w, c, "$", "$")

	if t.IsPrimitive() {
		sw.Do(snippets.GenerateEnumSetter(t, tags.GetEnumOptions(t)))
		return sw.Error()
	}

	if hasObjectMetaEmbedded(t) {
		parentTypeOfObjectMeta := getParentOfEmbeddedType(t, ObjectMeta)
		objectMetaType := getMemberTypeFromType(parentTypeOfObjectMeta, ObjectMeta)
		b.imports.AddType(parentTypeOfObjectMeta)
		b.imports.AddType(objectMetaType)
		sw.Do(snippets.GenerateConstructorForObjectMeta(t))
		sw.Do(snippets.GenerateDeepCopy(t))
		b.generateSettersForType(sw, t, objectMetaType)
	} else {
		sw.Do(snippets.GenerateEmptyConstructor(t, true))
	}

	for _, member := range t.Members {
		log.Debugf("generateSettersForType %v - Type : %v", member.Name, member.Type)
		b.generateSettersForType(sw, t, member.Type)
	}

	return sw.Error()
}

func (b *BuilderPatternGenerator) generateSettersForType(sw *generator.SnippetWriter, root *types.Type, parent *types.Type) {
	setter := snippets.NewSetter(root, parent, true)

	for _, m := range parent.Members {
		if m.Embedded || !includeMember(parent, m) {
			continue
		}

		log.Debugf("parentMember %v - Type : %v -- Kind: %s", m.Name, m.Type, m.Type.Kind)

		switch {
		case m.Type.Kind == types.Map:
			keyType := m.Type.Key
			elemType := m.Type.Elem
			switch {
			case keyType == types.String && elemType == types.String:
				sw.Do(setter.GenerateSetterForMapStringString(m))
			default:
				sw.Do(setter.GenerateSetterForMap(m))
			}
		case m.Type.Kind == types.Slice:
			sliceType := m.Type.Elem
			switch sliceType.Kind {
			case types.Struct, types.Pointer:
				log.Debugf("generateSettersForType - Slice -> Struct : %v - Type : %v", m.Name, m.Type)
				if b.isTypeEnabled(m.Type) {

					if sliceType.Kind == types.Pointer {
						log.Debugf("\t %v is enabled -> GenerateSetterForEmbeddedSlicePointer", m.Type)
						sw.Do(setter.GenerateSetterForEmbeddedSlicePointer(m, b.getWrapperType(sliceType)))
					} else {
						log.Debugf("\t %v is enabled -> GenerateSetterForEmbeddedSlice", m.Type)
						sw.Do(setter.GenerateSetterForEmbeddedSlice(m, b.getWrapperType(sliceType)))
					}
				}
			default:
				if b.isTypeEnabled(m.Type) || sliceType.Kind == types.Builtin {
					log.Debugf("\t %v is default   (kind - %s)", m.Type, sliceType.Kind)
					sw.Do(setter.GenerateSetterForMemberSlice(m))
				}
			}
		case m.Type.Kind == types.Struct:
			log.Debugf("generateSettersForType - Struct : %v", m.Type)
			if b.isTypeEnabled(m.Type) {
				log.Debugf("\t %v is enabled", m.Type)
				sw.Do(setter.GenerateSetterForEmbeddedStruct(m, b.getWrapperType(m.Type)))
			}
		case m.Type.Kind == types.Pointer:
			pointerType := m.Type.Elem
			switch pointerType.Kind {
			case types.Builtin:
				if pointerType == types.Bool {
					sw.Do(setter.GenerateSetterForPointerToBool(m))
				} else {
					sw.Do(setter.GenerateSetterForPointerToBuiltinType(m))
				}
			case types.Struct:
				log.Debugf("generateSettersForType - Pointer -> Struct : %v", pointerType)
				if b.isTypeEnabled(pointerType) {
					log.Debugf("\t %v is enabled", pointerType)
					sw.Do(setter.GenerateSetterForEmbeddedPointer(m, b.getWrapperType(pointerType)))
				}
			case types.Alias:
				log.Debugf("generateSettersForType - Alias : %v", m.Type)
				if b.isTypeEnabled(m.Type) {
					wrap := b.getWrapperType(m.Type)
					log.Debugf("WRAPAP %#v", wrap)
					sw.Do(setter.GenerateSetterForAliasPointerPrimitive(m, wrap))
				}
			default:
				sw.Do(setter.GenerateSetterForType(m))
			}
		case m.Type == types.Bool:
			sw.Do(setter.GenerateSetterForBool(m))
		case m.Type.Kind == types.Alias:
			if m.Type.Underlying.Kind == types.Builtin && b.isTypeEnabled(m.Type) {
				log.Debugf("Kind Alias - generateSetterForTypeEnum - enhanced : %v - Type: %v", m.Name, m.Type.Name)
				wrap := b.getWrapperType(m.Type)
				sw.Do(setter.GenerateSetterForTypeEnum(m, wrap))
			}
		default:
			log.Debugf("generateSettersForType - Default : %v - Type: %v", m.Name, m.Type.Name)
			if b.isTypeEnabled(m.Type) {
				log.Debugf("\t GenerateSetterForType : %v - Type: %v", m.Name, m.Type.Name)
				sw.Do(setter.GenerateSetterForType(m))
			}
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
	// _, exists := b.enabledTypes[typeName]
	// return exists

	log.Debugf("isTypeEnabled for %s", typeName)
	typeName = strings.Replace(typeName, "[]", "", 1)
	typeName = strings.Replace(typeName, "*", "", 1)
	_, exists := b.packageIndex.TypesByTypePath[typeName]
	return exists || t.Kind == types.Builtin
}

// func (b *BuilderPatternGenerator) isTypePrimitiveEnabled(t types.Member) bool {
// 	typeName := t.Name
// 	_, exists := b.enabledTypes[typeName]
// 	return exists
// }

// func (b *BuilderPatternGenerator) getEnabledPrimitiveType(t types.Member) *types.Type {
// 	typeName := t.Name
// 	return b.enabledTypes[typeName]
// }

func (b *BuilderPatternGenerator) getWrapperType(t *types.Type) *types.Type {
	typeName := t.Name.String()

	typeName = strings.Replace(typeName, "[]", "", 1)
	typeName = strings.Replace(typeName, "*", "", 1)
	log.Debugf("getWrapperType for %s", typeName)

	// if parent, ok := b.enabledTypes[typeName]; ok {
	// 	return parent
	// }

	return b.packageIndex.TypesByTypePath[typeName]
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
	// if !b.needsGeneration(t) {
	// 	return false
	// }

	//b.enablePrimitiveType(t)

	// for _, m := range t.Members {
	// 	b.enableEmbededMemberType(m, t)
	// }

	return b.needsGeneration(t)
}

// func (b *BuilderPatternGenerator) enableEmbededMemberType(m types.Member, t *types.Type) {
// 	if m.Embedded {
// 		childType := m.Type.Name.String()
// 		b.enabledTypes[childType] = t
// 	}
// }

// func (b *BuilderPatternGenerator) enablePrimitiveType(t *types.Type) {
// 	if t.IsPrimitive() {
// 		b.enabledTypes[t.Name.Name] = t
// 	}
// }

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

	eval := strings.HasSuffix(pkg, b.pkgToBuild.Path[2:]+"\"")

	log.Debugf("Suffix Test |%s| has suffix |%s| == %t", pkg, b.pkgToBuild.Path[2:]+"\"", eval)
	return !eval
}
