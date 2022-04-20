package builder

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func newTestGeneratorType(t *testing.T, dir string, selector string) (*types.Package, *types.Type) {
	testDir := fmt.Sprintf("./testdata/%s", dir)
	d := args.Default()
	d.IncludeTestFiles = true
	d.InputDirs = []string{testDir + ""}
	d.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree())
	b, err := d.NewBuilder()
	assert.NoError(t, err)
	findTypes, err := b.FindTypes()
	assert.NoError(t, err)

	pkg := findTypes[testDir]
	assert.NotNil(t, pkg)
	n := pkg.Types[selector]
	assert.NotNil(t, n)
	return pkg, n
}

func newGeneratorContext(g generator.Generator) *generator.Context {
	c := &generator.Context{}
	c.Namers = g.Namers(c)
	return c
}

func TestBuilderPatternGenerator_NeedsGeneration(t *testing.T) {
	tests := []struct {
		description string
		dir         string
		structName  string
		wantEmpty   bool
	}{
		{
			description: "Types do not need generation",
			dir:         "a",
			structName:  "NoGeneration",
			wantEmpty:   true,
		},
		{
			description: "Types need generation",
			dir:         "a",
			structName:  "AStruct",
		},
		{
			description: "All package types need generation",
			dir:         "b",
			structName:  "BStruct",
		},
		{
			description: "Type generation opt-out",
			dir:         "b",
			structName:  "OptOutStruct",
			wantEmpty:   true,
		},
	}

	for _, test := range tests {
		b := &BuilderPatternGeneratorFactory{}
		pkg, typeToGenerate := newTestGeneratorType(t, test.dir, test.structName)
		g := b.NewBuilder(pkg)
		buf := &bytes.Buffer{}
		c := newGeneratorContext(g)
		assert.NoError(t, g.GenerateType(c, typeToGenerate, buf), test.description)
		if test.wantEmpty {
			assert.Empty(t, buf, test.description)
		} else {
			assert.NotEmpty(t, buf, test.description)
		}
	}
}
func TestBuilderPattern_TypeContainTypeMeta(t *testing.T) {
	_, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	assert.True(t, hasTypeMetaEmbedded(typeToGenerate))
}

func TestBuilderPattern_ImportTrackerToAliasNames(t *testing.T) {
	tracker := newImportTracker()
	_, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	assert.Equal(t, "testdatac", golangNameToImportAlias(tracker, typeToGenerate.Name))

	_, typeToGenerate = newTestGeneratorType(t, "c/d", "MockDeployment")
	assert.Equal(t, "cd", golangNameToImportAlias(tracker, typeToGenerate.Name))
}

func TestBuilderPattern_TypeMetaGeneratesSnippets(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	g := b.NewBuilder(pkg)
	buf := &bytes.Buffer{}
	c := newGeneratorContext(g)
	assert.NoError(t, g.GenerateType(c, typeToGenerate, buf))
	assert.Contains(t, buf.String(), "func (in *CDeployment) DeepCopy() *CDeployment")
	assert.Contains(t, buf.String(), "func (in *CDeployment) DeepCopyInto(out *CDeployment)")
	assert.Contains(t, buf.String(), "func (o CDeployment) MarshalJSON()")
	assert.Contains(t, buf.String(), "d.SchemeGroupVersion")
}

func TestBuilderPattern_TypeMetaGeneratesImportLines(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	g := b.NewBuilder(pkg)
	c := newGeneratorContext(g)
	assert.NoError(t, g.GenerateType(c, typeToGenerate, &bytes.Buffer{}))

	imports := g.Imports(c)
	assert.Len(t, imports, 2)
	assert.Contains(t, strings.Join(imports, ""), "cmeta")
	assert.Contains(t, strings.Join(imports, ""), "cd")

}

func TestBuilderPattern_GenerateInit(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, _ := newTestGeneratorType(t, "c", "CDeployment")
	g := b.NewBuilder(pkg)
	buf := &bytes.Buffer{}
	c := &generator.Context{}
	c.Namers = g.Namers(c)
	assert.NoError(t, g.Init(c, buf))
	assert.Contains(t, buf.String(), "mergeMapStringString")
}
