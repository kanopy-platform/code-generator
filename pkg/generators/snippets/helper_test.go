/*
 Helper functions for _test.go files
*/

package snippets

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/gengo/v2/generator"
	"k8s.io/gengo/v2/namer"
	"k8s.io/gengo/v2/parser"
	"k8s.io/gengo/v2/types"
)

const testPackage string = "./testdata"

// testPackagePath is the canonical import path of testPackage. gengo/v2 keys
// the universe and the raw namer by import path rather than by directory.
const testPackagePath string = "github.com/kanopy-platform/code-generator/pkg/generators/snippets/testdata"

func nameSystem() namer.NameSystems {
	return namer.NameSystems{
		"public": namer.NewPublicNamer(1),
		"raw":    namer.NewRawNamer(testPackagePath, nil),
	}
}

func defaultNameSystem() string {
	return "public"
}

func newTestGeneratorContext() (*generator.Context, error) {
	p := parser.NewWithOptions(parser.Options{})

	c, err := generator.NewContext(p, nameSystem(), defaultNameSystem())
	if err != nil {
		return nil, err
	}

	return c, nil
}

func newTestType(t *testing.T, selector string) *types.Type {
	dir := testPackage
	p := parser.NewWithOptions(parser.Options{})
	// Go tooling excludes "testdata" from "..." expansion, so list the
	// fixture packages explicitly.
	paths, err := p.FindPackages(dir, dir+"/a", dir+"/b")
	assert.NoError(t, err)
	assert.NoError(t, p.LoadPackages(paths...))

	findTypes, err := p.NewUniverse()
	assert.NoError(t, err)
	dir = paths[0]

	n := findTypes[dir].Types[selector]
	assert.NotNil(t, n)

	return n
}

// Helper function to select a nested member from a type
func getMemberFromType(t *testing.T, in *types.Type, selector ...string) types.Member {
	require.NotNil(t, in)

	currType := in
	var member types.Member

	for _, s := range selector {
		foundMember := false

		for _, m := range currType.Members {
			if m.Name == s {
				foundMember = true
				member = m
				currType = m.Type
				break
			}
		}
		require.True(t, foundMember, fmt.Sprintf("Member %q not found", s))
	}

	require.NotEqual(t, types.Member{}, member)

	return member
}
