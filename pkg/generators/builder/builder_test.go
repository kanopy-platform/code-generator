package builder

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/kanopy-platform/code-generator/pkg/generators"
	"github.com/kanopy-platform/code-generator/pkg/generators/index"
	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/v2/generator"
	"k8s.io/gengo/v2/parser"
	"k8s.io/gengo/v2/types"
)

var defaultIndex = generators.NewPackageTypeIndex()

func newTestGeneratorType(t *testing.T, dir string, selector string) (*types.Package, *types.Type) {
	testDir := fmt.Sprintf("./testdata/%s", dir)
	p := parser.NewWithOptions(parser.Options{})
	paths, err := p.FindPackages(testDir)
	assert.NoError(t, err)
	assert.NoError(t, p.LoadPackages(testDir))
	findTypes, err := p.NewUniverse()
	assert.NoError(t, err)
	testDir = paths[0]
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
	const root = "github.com/10gen/kanopy/pkg/builder"
	tracker := newImportTracker(root+"/argo", root)

	alias := func(pkg string) string {
		return golangNameToImportAlias(tracker, root, types.Name{Package: pkg})
	}

	// named by parent directory and leaf
	assert.Equal(t, "corev1", alias("k8s.io/api/core/v1"))
	assert.Equal(t, "metav1", alias("k8s.io/apimachinery/pkg/apis/meta/v1"))
	assert.Equal(t, "workflowv1alpha1", alias("github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"))

	// the package root is dropped
	assert.Equal(t, "k8s", alias(root+"/k8s"))
	assert.Equal(t, "crossplane", alias(root+"/crossplane"))
	assert.Equal(t, "crossplaneprovideraws", alias(root+"/crossplane/provideraws"))

	// a collision walks further up the path
	tracker.AddType(&types.Type{Name: types.Name{Package: "k8s.io/api/core/v1", Name: "Pod"}})
	assert.Equal(t, "apicorev1", alias("k8s.io/api/core/v1"))
}

func TestBuilderPattern_ImportTrackerWithoutPackageRoot(t *testing.T) {
	tracker := newImportTracker("github.com/10gen/kanopy/pkg/builder/argo", "")
	assert.Equal(t, "buildercrossplane", golangNameToImportAlias(tracker, "",
		types.Name{Package: "github.com/10gen/kanopy/pkg/builder/crossplane"}))
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
	// Index package c first so the wrapper for the d.MockSpec members resolves
	// to c.MockSpec, a type outside of the package being generated.
	newTestGeneratorType(t, "c", "CDeployment")
	newTestGeneratorType(t, "c", "MockSpec")
	pkg, typeToGenerate := newTestGeneratorType(t, "e", "EDeployment")
	defaultIndex.PackageRoot = "github.com/kanopy-platform/code-generator/pkg/generators/builder/testdata"
	g := b.NewBuilder(pkg, defaultIndex)
	buf := &bytes.Buffer{}
	c := newGeneratorContext(g)
	assert.True(t, g.Filter(c, typeToGenerate))
	assert.NoError(t, g.GenerateType(c, typeToGenerate, buf))

	// The setters reference the wrapper type from package c.
	assert.Contains(t, buf.String(), "func (o *EDeployment) WithSpec(in *c.MockSpec) *EDeployment")

	imports := g.Imports(c)
	assert.Len(t, imports, 1)
	// c sits directly under the package root, so it is named by its leaf
	assert.Equal(t, `c "github.com/kanopy-platform/code-generator/pkg/generators/builder/testdata/c"`, imports[0])

	// Only packages the generated body references are imported, and never the
	// package being generated.
	for _, line := range imports {
		alias, path, found := strings.Cut(line, " ")
		assert.True(t, found, line)
		assert.NotEmpty(t, alias)
		assert.NotContains(t, path, pkg.Path+"\"", "must not self-import")
	}
}

func TestBuilderPattern_NoImportLinesWhenBodyIsLocal(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "c", "CDeployment")
	g := b.NewBuilder(pkg, defaultIndex)
	c := newGeneratorContext(g)
	assert.NoError(t, g.GenerateType(c, typeToGenerate, &bytes.Buffer{}))

	// Everything CDeployment references is wrapped in its own package.
	assert.Empty(t, g.Imports(c))
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
