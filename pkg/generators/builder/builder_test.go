package builder

import (
	"bytes"
	"fmt"
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
	d.InputDirs = []string{testDir + "/..."}
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
	g := b.NewBuilder("zz_test", pkg, "")
	buf := &bytes.Buffer{}
	assert.NoError(t, g.GenerateType(&generator.Context{}, typeToGenerate, buf))
	assert.Empty(t, buf)
}

func TestBuilderPatternGenerator_TypeNeedsGeneration(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "a", "AStruct")
	g := b.NewBuilder("zz_test", pkg, tagName)
	buf := &bytes.Buffer{}
	err := g.GenerateType(&generator.Context{}, typeToGenerate, buf)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
}

func TestBuilderPatternGenerator_AllPackageTypesNeedGeneration(t *testing.T) {
	b := &BuilderPatternGeneratorFactory{}
	pkg, typeToGenerate := newTestGeneratorType(t, "b", "BStruct")
	g := b.NewBuilder("zz_test", pkg, tagName)
	buf := &bytes.Buffer{}
	err := g.GenerateType(&generator.Context{}, typeToGenerate, buf)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
}
