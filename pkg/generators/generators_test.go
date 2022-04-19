package generators

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func TestExtractCommentTag(t *testing.T) {
	tag := "test:tag"
	fmtTag := "+%s=%s"
	tests := []struct {
		description string
		comments    []string
		want        string
	}{
		{
			description: "Empty value from no comments",
			comments:    []string{},
			want:        "",
		},
		{
			description: "Empty value from empty comments",
			comments:    []string{""},
			want:        "",
		},
		{
			description: "Value from comments",
			comments:    []string{fmt.Sprintf(fmtTag, tag, "value")},
			want:        "value",
		},
		{
			description: "Return first value with multiple values",
			comments:    []string{fmt.Sprintf(fmtTag, tag, "value,value2")},
			want:        "value",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, extractTag(tag, test.comments), test.description)
	}
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

func TestPackages_Generation(t *testing.T) {
	tests := []struct {
		description  string
		testInputDir string
		want         int
	}{
		{
			description:  "All defaults",
			testInputDir: "",
			want:         0,
		},
		{
			description:  "One package needs generation",
			testInputDir: "./testdata/a/...",
			want:         1,
		},
		{
			description:  "Empty dir",
			testInputDir: "./testdata/b/...",
			want:         0,
		},
		{
			description:  "Single type in package needs generation",
			testInputDir: "./testdata/c/...",
			want:         1,
		},
		{
			description:  "Types do not need generation",
			testInputDir: "./testdata/d/...",
			want:         0,
		},
	}

	for _, test := range tests {
		ctx := &generator.Context{}
		a := &args.GeneratorArgs{}
		g := New(&MockBuilderFactory{})
		if test.testInputDir != "" {
			a, ctx = testDataGeneratorSetup(t, test.testInputDir)
		}
		packages := g.Packages(ctx, a)
		assert.Len(t, packages, test.want, test.description)
	}
}

func TestPackage_FilterPackage(t *testing.T) {
	_, ctx := testDataGeneratorSetup(t, "./testdata/d/...")

	tests := []struct {
		description string
		pkg         *types.Package
		name        types.Name
		want        bool
	}{
		{
			description: "package included",
			pkg:         &types.Package{Path: "testme"},
			name:        types.Name{Package: "testme"},
			want:        true,
		},
		{
			description: "package excluded",
			pkg:         &types.Package{Path: "excludeme"},
			name:        types.Name{Package: "testme"},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, filterFuncByPackagePath(test.pkg)(ctx, &types.Type{Name: test.name}), test.description)
	}
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

	ctx, err := generator.NewContext(b, NameSystems(), DefaultNameSystem)
	assert.NoError(t, err)
	return a, ctx
}

type MockBuilderFactory struct {
	generator.DefaultGen
}

func (m *MockBuilderFactory) NewBuilder(outputFileBaseName string, pkg *types.Package, tagValue string) generator.Generator {
	return m
}
