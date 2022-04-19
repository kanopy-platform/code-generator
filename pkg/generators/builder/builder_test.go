package builder

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

const tagName = "kanopy:builder"

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

func TestBuilderPatternGenerator_TypeDoesNotNeedGeneration(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "a", "NoGeneration")
	g := b.NewBuilder(pkg, "")
	buf := &bytes.Buffer{}
	assert.NoError(t, g.GenerateType(&generator.Context{}, typeToGenerate, buf))
	assert.Empty(t, buf)
}

func TestBuilderPatternGenerator_TypeNeedsGeneration(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "a", "AStruct")
	g := b.NewBuilder(pkg, tagName)
	buf := &bytes.Buffer{}
	err := g.GenerateType(&generator.Context{}, typeToGenerate, buf)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
}

func TestBuilderPatternGenerator_AllPackageTypesNeedGeneration(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "b", "BStruct")
	g := b.NewBuilder(pkg, tagName)
	buf := &bytes.Buffer{}
	c := &generator.Context{}
	c.Namers = g.Namers(c)
	err := g.GenerateType(c, typeToGenerate, buf)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
}

func TestBuilderPatternGenerator_AllTypesInPackageWithTypeOptout(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "b", "OptOutStruct")
	g := b.NewBuilder(pkg, tagName)
	buf := &bytes.Buffer{}
	c := &generator.Context{}
	c.Namers = g.Namers(c)
	err := g.GenerateType(c, typeToGenerate, buf)
	assert.NoError(t, err)
	assert.Empty(t, buf)
}

func TestBuilderPattern_TypeContainTypeMeta(t *testing.T) {
	_, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	g := &BuilderPatternGenerator{}
	assert.True(t, g.hasTypeMetaEmbedded(typeToGenerate))
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
	g := b.NewBuilder(pkg, tagName)
	buf := &bytes.Buffer{}
	c := &generator.Context{}
	c.Namers = g.Namers(c)
	err := g.GenerateType(c, typeToGenerate, buf)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "func (in *CDeployment) DeepCopy() *CDeployment")
	assert.Contains(t, buf.String(), "func (in *CDeployment) DeepCopyInto(out *CDeployment)")
	assert.Contains(t, buf.String(), "func (o CDeployment) MarshalJSON()")
	assert.Contains(t, buf.String(), "func (o CDeployment) MarshalJSON()")
	assert.Contains(t, buf.String(), "d.SchemeGroupVersion")
}

func TestBuilderPattern_GenerateInit(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, _ := newTestGeneratorType(t, "c", "CDeployment")
	g := b.NewBuilder(pkg, tagName)
	buf := &bytes.Buffer{}
	c := &generator.Context{}
	c.Namers = g.Namers(c)
	assert.NoError(t, g.Init(c, buf))
	assert.Contains(t, buf.String(), "mergeMapStringString")
}
