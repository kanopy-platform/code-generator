/*
 Helper functions for _test.go files
*/

package snippets

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

// Set packageName to match package of test structs to avoid package prefix
const packageName string = "./testdata"

func nameSystem() namer.NameSystems {
	return namer.NameSystems{
		"public": namer.NewPublicNamer(1),
		"raw":    namer.NewRawNamer(packageName, nil),
	}
}

func defaultNameSystem() string {
	return "public"
}

func newTestGeneratorContext() (*generator.Context, error) {
	args := args.Default()

	b, err := args.NewBuilder()
	if err != nil {
		return nil, err
	}

	c, err := generator.NewContext(b, nameSystem(), defaultNameSystem())
	if err != nil {
		return nil, err
	}

	return c, nil
}

func newSomeStructType(t *testing.T, selector string) *types.Type {
	dir := "./testdata"
	d := args.Default()
	d.IncludeTestFiles = true
	d.InputDirs = []string{dir + "/..."}

	b, err := d.NewBuilder()
	assert.NoError(t, err)

	findTypes, err := b.FindTypes()
	assert.NoError(t, err)

	n := findTypes[dir].Types[selector]
	fmt.Println(findTypes, n.Members)
	assert.NotNil(t, n)

	return n
}
