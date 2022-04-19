package builder

import (
	"io"
	"strings"

	"k8s.io/gengo/types"

	"k8s.io/gengo/generator"
)

// TODO: this will be refactored out in another PR
type DefaultBuilderFactory struct {
	generator.DefaultGen
}

func (d *DefaultBuilderFactory) NewBuilder(outputFileBaseName string, pkg *types.Package, tagName string) generator.Generator {
	return d
}

type BuilderPatternGenerator struct {
	generator.DefaultGen
	tagName    string
	pkgToBuild *types.Package
	allTypes   bool
}

type BuilderPatternGeneratorFactory struct{}

func (d *BuilderPatternGeneratorFactory) NewBuilder(outputFileBaseName string, pkg *types.Package, tagName string) generator.Generator {
	return &BuilderPatternGenerator{tagName: tagName, pkgToBuild: pkg, allTypes: allTypes(tagName, pkg)}
}

func allTypes(tagName string, pkg *types.Package) bool {
	return doesPackageNeedGeneration(extractTag(tagName, pkg.Comments))
}

func doesPackageNeedGeneration(tag string) bool {
	return tag == "package" // package type TODO refactor
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

func (b *BuilderPatternGenerator) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	if !b.needsGeneration(t) {
		return nil
	}

	sw := generator.NewSnippetWriter(w, c, "$", "$")

	sw.Do("// auto generated", nil)

	return nil
}

func (b *BuilderPatternGenerator) needsGeneration(t *types.Type) bool {
	if !b.doesTypeNeedGeneration(t) && !b.allTypes {
		return false
	}
	return true
}

func (b *BuilderPatternGenerator) doesTypeNeedGeneration(t *types.Type) bool {
	tag := extractTag(b.tagName, t.CommentLines)
	return tag == "true"
}
