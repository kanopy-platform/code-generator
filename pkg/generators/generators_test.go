package generators

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func TestExtractTag_ValueFromComments(t *testing.T) {
	tag := "test:tag"
	comments := []string{fmt.Sprintf("+%s=value", tag)}
	expected := "value"
	assert.Equal(t, expected, extractTag(tag, comments))
}

func TestExtracTag_ValueFromEmptyComments(t *testing.T) {
	tag := "test:tag"
	expected := ""
	assert.Equal(t, expected, extractTag(tag, []string{""}))
}

func TestExtractTag_ReturnFirstValueWithMultipleValues(t *testing.T) {
	tag := "test:tag"
	comments := []string{fmt.Sprintf("+%s=value,value2", tag)}
	assert.Equal(t, "value", extractTag(tag, comments))
}

func TestNewGenerators(t *testing.T) {
	g := New(&MockBuilderFactory{})
	assert.NotNil(t, g)
}

func TestNewGenerator_WithBoilerplate(t *testing.T) {
	expected := "// some basic boilerplate"
	g := New(&MockBuilderFactory{}, WithBoilerplate(expected))
	assert.NotNil(t, g)
	assert.Equal(t, expected, g.Boilerplate)
}

func TestPackages_NoPackagesNeedGeneration(t *testing.T) {
	ctx := generator.Context{}
	args := args.GeneratorArgs{}
	g := New(&MockBuilderFactory{})
	packages := g.Packages(&ctx, &args)
	assert.Empty(t, packages)
}

func TestPackages_NeedsGeneration(t *testing.T) {
	a, ctx := testDataGeneratorSetup(t, "./testdata/a/...")
	g := New(&MockBuilderFactory{})
	packages := g.Packages(ctx, a)
	assert.Len(t, packages, 1)
}

func TestPackages_SingleTypeNeedsGeneration(t *testing.T) {
	a, ctx := testDataGeneratorSetup(t, "./testdata/c/...")
	g := New(&MockBuilderFactory{})
	packages := g.Packages(ctx, a)
	assert.Len(t, packages, 1)
}

func TestPackages_EmptyDir(t *testing.T) {
	a, ctx := testDataGeneratorSetup(t, "./testdata/b/...")
	g := New(&MockBuilderFactory{})
	packages := g.Packages(ctx, a)
	assert.Len(t, packages, 0)
}

func TestPackages_TypesDoNotNeedGeneration(t *testing.T) {
	a, ctx := testDataGeneratorSetup(t, "./testdata/d/...")
	g := New(&MockBuilderFactory{})
	packages := g.Packages(ctx, a)
	assert.Len(t, packages, 0)
}

func TestPackage_FilterPackageIncluded(t *testing.T) {
	_, ctx := testDataGeneratorSetup(t, "./testdata/d/...")
	p := &types.Package{
		Path: "testme",
	}

	assert.True(t, filterFuncByPackagePath(p)(ctx, &types.Type{
		Name: types.Name{
			Package: "testme",
		},
	}))
}

func TestPackage_GeneratorFuncForPackage(t *testing.T) {
	_, ctx := testDataGeneratorSetup(t, "./testdata/d/...")
	p := &types.Package{
		Path: "testme",
	}
	g := New(&MockBuilderFactory{})
	assert.Len(t, g.generatorFuncForPackage("base", p, "true")(ctx), 1)
}

func testDataGeneratorSetup(t *testing.T, dir string) (*args.GeneratorArgs, *generator.Context) {
	a := args.Default()

	a.InputDirs = []string{dir}

	b, err := a.NewBuilder()
	assert.NoError(t, err)

	ctx, err := generator.NewContext(b, NameSystems(), DefaultNameSystem())
	assert.NoError(t, err)
	return a, ctx
}

type MockBuilderFactory struct {
	generator.DefaultGen
}

func (m *MockBuilderFactory) NewBuilder(outputFileBaseName string, pkg *types.Package, tagValue string) generator.Generator {
	return m
}
