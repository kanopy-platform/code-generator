package builder

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kanopy-platform/code-generator/pkg/generators"
	"github.com/kanopy-platform/code-generator/pkg/generators/index"
	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

var defaultIndex = generators.NewPackageTypeIndex()

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

	defaultIndex.TypesByTypePath = index.BuildPackageIndex(defaultIndex.TypesByTypePath, pkg)

	n := pkg.Types[selector]
	assert.NotNil(t, n)
	return pkg, n
}

func newGeneratorContext(g generator.Generator) *generator.Context {
	c := &generator.Context{}
	c.Namers = g.Namers(c)
	return c
}

func TestBuilderPatternGenerator_Filter(t *testing.T) {
	tests := []struct {
		description string
		dir         string
		structName  string
		wantGen     bool
	}{
		{
			description: "Types do not need generation",
			dir:         "a",
			structName:  "NoGeneration",
		},
		{
			description: "Types need generation",
			dir:         "a",
			structName:  "AStruct",
			wantGen:     true,
		},
		{
			description: "All package types need generation",
			dir:         "b",
			structName:  "BStruct",
			wantGen:     true,
		},
		{
			description: "Type generation opt-out",
			dir:         "b",
			structName:  "OptOutStruct",
		},
	}

	for _, test := range tests {
		b := &BuilderPatternGeneratorFactory{}
		pkg, typeToGenerate := newTestGeneratorType(t, test.dir, test.structName)
		g := b.NewBuilder(pkg, defaultIndex)
		c := newGeneratorContext(g)
		assert.Equal(t, test.wantGen, g.Filter(c, typeToGenerate), test.description)
	}
}

func TestBuilderPattern_ImportTrackerToAliasNames(t *testing.T) {
	tracker := newImportTracker(generators.NewPackageTypeIndex())
	_, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	assert.Equal(t, "testdatac", golangNameToImportAlias(tracker, typeToGenerate.Name))

	_, typeToGenerate = newTestGeneratorType(t, "c/d", "MockDeployment")
	assert.Equal(t, "cd", golangNameToImportAlias(tracker, typeToGenerate.Name))
}

func TestBuilderPattern_ObjectMetaGeneratesSnippets(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	_, specTypeToGenerate := newTestGeneratorType(t, "c", "MockSpec")
	g := b.NewBuilder(pkg, defaultIndex)
	buf := &bytes.Buffer{}
	c := newGeneratorContext(g)
	assert.True(t, g.Filter(c, typeToGenerate))
	assert.True(t, g.Filter(c, specTypeToGenerate))
	assert.NoError(t, g.GenerateType(c, typeToGenerate, buf))

	// constructor
	assert.Contains(t, buf.String(), "func NewCDeployment(name string) *CDeployment")
	// setters
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithName(in string) *CDeployment")
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithSpec(in *MockSpec) *CDeployment")
	// deepcopy
	assert.Contains(t, buf.String(), "func (in *CDeployment) DeepCopy() *CDeployment")
	assert.Contains(t, buf.String(), "func (in *CDeployment) DeepCopyInto(out *CDeployment)")
}

func TestBuilderPattern_NonObjectMetaGeneratesSnippets(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "d", "DPolicyRule")
	g := b.NewBuilder(pkg, defaultIndex)
	buf := &bytes.Buffer{}
	c := newGeneratorContext(g)
	assert.True(t, g.Filter(c, typeToGenerate))
	assert.NoError(t, g.GenerateType(c, typeToGenerate, buf))

	// constructor
	assert.Contains(t, buf.String(), "func NewDPolicyRule() *DPolicyRule")
	// setters
	assert.Contains(t, buf.String(), "func (o *DPolicyRule) AppendVerbs(in ...string) *DPolicyRule")
	assert.Contains(t, buf.String(), "func (o *DPolicyRule) AppendListOfInts(in ...int) *DPolicyRule")
	// no deepcopy
	assert.NotContains(t, buf.String(), "DeepCopy")
	assert.NotContains(t, buf.String(), "DeepCopyInto")
}

func TestBuilderAliasPrimitiveType(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "d", "DPolicyRule")
	_, aliasToGenerate := newTestGeneratorType(t, "d", "AliasType")
	g := b.NewBuilder(pkg, defaultIndex)
	buf := &bytes.Buffer{}
	c := newGeneratorContext(g)
	assert.True(t, g.Filter(c, typeToGenerate))
	assert.True(t, g.Filter(c, aliasToGenerate))
	assert.NoError(t, g.GenerateType(c, typeToGenerate, buf))
	t.Log(buf.String())
	assert.Contains(t, buf.String(), "func (o *DPolicyRule) WithAliasType(in AliasType) *DPolicyRule")
}

func TestBuilderAliasPrimitiveTypeNotGenerated(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "d", "DPolicyRule")
	g := b.NewBuilder(pkg, defaultIndex)
	buf := &bytes.Buffer{}
	c := newGeneratorContext(g)
	assert.True(t, g.Filter(c, typeToGenerate))
	assert.NoError(t, g.GenerateType(c, typeToGenerate, buf))
	assert.NotContains(t, buf.String(), "func (o *DPolicyRule) WithToggleAliasWithoutRef(in AnotherAlias) *DPolicyRule")
}

func TestBuilderPattern_GenerateSettersForType(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	_, specTypeToGenerate := newTestGeneratorType(t, "c", "MockSpec")
	g := b.NewBuilder(pkg, defaultIndex)
	buf := &bytes.Buffer{}
	c := newGeneratorContext(g)
	assert.True(t, g.Filter(c, typeToGenerate))
	assert.True(t, g.Filter(c, specTypeToGenerate))
	assert.NoError(t, g.GenerateType(c, typeToGenerate, buf))

	// ObjectMeta setters
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithName(in string) *CDeployment")
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithLabels(in map[string]string) *CDeployment")
	assert.NotContains(t, buf.String(), "AppendFinalizers")
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithIntPtr(in int) *CDeployment")
	// Spec setters
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithSpec(in *MockSpec) *CDeployment")
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithPointerSpec(in *MockSpec) *CDeployment")
	assert.Contains(t, buf.String(), "func (o *CDeployment) AppendSpecs(in ...*MockSpec) *CDeployment")
	assert.NotContains(t, buf.String(), "SpecNoGen")
	assert.NotContains(t, buf.String(), "PointerSpecNoGen")
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithPrimitive(in int) *CDeployment")
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithBool(in ...bool) *CDeployment")
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithPointerBool(in ...bool) *CDeployment")
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithMapStringByteSlice(in map[string][]byte) *CDeployment")
}

func TestBuilderPattern_ObjectMetaGeneratesImportLines(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	g := b.NewBuilder(pkg, defaultIndex)
	c := newGeneratorContext(g)
	assert.NoError(t, g.GenerateType(c, typeToGenerate, &bytes.Buffer{}))

	imports := g.Imports(c)
	assert.Len(t, imports, 4) // 4 types are tagged for importing
	assert.Contains(t, strings.Join(imports, ""), "cmeta")
	assert.Contains(t, strings.Join(imports, ""), "cd")

}

func TestBuilderPattern_GenerateInit(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, _ := newTestGeneratorType(t, "c", "CDeployment")
	g := b.NewBuilder(pkg, defaultIndex)
	buf := &bytes.Buffer{}
	c := newGeneratorContext(g)
	assert.NoError(t, g.Init(c, buf))
	assert.Contains(t, buf.String(), "mergeMapStringString")
}
