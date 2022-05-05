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
		g := b.NewBuilder(pkg)
		c := newGeneratorContext(g)
		assert.Equal(t, test.wantGen, g.Filter(c, typeToGenerate), test.description)
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
	_, specTypeToGenerate := newTestGeneratorType(t, "c", "MockSpec")
	g := b.NewBuilder(pkg)
	buf := &bytes.Buffer{}
	c := newGeneratorContext(g)
	assert.True(t, g.Filter(c, typeToGenerate))
	assert.True(t, g.Filter(c, specTypeToGenerate))
	assert.NoError(t, g.GenerateType(c, typeToGenerate, buf))

	// constructor
	assert.Contains(t, buf.String(), "func NewCDeployment(name string) *CDeployment")
	// deepcopy
	assert.Contains(t, buf.String(), "func (in *CDeployment) DeepCopy() *CDeployment")
	assert.Contains(t, buf.String(), "func (in *CDeployment) DeepCopyInto(out *CDeployment)")
	// setters
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithName(in string) *CDeployment")
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithSpec(in *MockSpec) *CDeployment")
}

func TestBuilderPattern_NonTypeMetaGeneratesSnippets(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "d", "DPolicyRule")
	g := b.NewBuilder(pkg)
	buf := &bytes.Buffer{}
	c := newGeneratorContext(g)
	assert.True(t, g.Filter(c, typeToGenerate))
	assert.NoError(t, g.GenerateType(c, typeToGenerate, buf))

	// constructor
	assert.Contains(t, buf.String(), "func NewDPolicyRule() *DPolicyRule")
	// no deepcopy
	assert.NotContains(t, buf.String(), "DeepCopy()")
	assert.NotContains(t, buf.String(), "DeepCopyInto")
	// setters
	assert.Contains(t, buf.String(), "func (o *DPolicyRule) AppendVerbs(in ...string) *DPolicyRule")
	assert.Contains(t, buf.String(), "func (o *DPolicyRule) AppendListOfInts(in ...int) *DPolicyRule")
}

func TestBuilderPattern_GenerateSettersForType(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	_, specTypeToGenerate := newTestGeneratorType(t, "c", "MockSpec")
	g := b.NewBuilder(pkg)
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
	assert.Contains(t, buf.String(), "func (o *CDeployment) AppendSpecs(in ...*MockSpec) *CDeployment")
	assert.NotContains(t, buf.String(), "SpecNoGen")
	assert.Contains(t, buf.String(), "func (o *CDeployment) WithPrimitive(in bool) *CDeployment")
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
	c := newGeneratorContext(g)
	assert.NoError(t, g.Init(c, buf))
	assert.Contains(t, buf.String(), "mergeMapStringString")
}
