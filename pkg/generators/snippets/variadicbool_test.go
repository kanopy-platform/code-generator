package snippets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/gengo/generator"
)

func TestGenerateVariadicBool(t *testing.T) {
	t.Parallel()

	ctx, err := newTestGeneratorContext()
	require.NoError(t, err)

	want := `// variadicBool selects the first element in the passed in list if non-empty. Otherwise the default return is "true".
func variadicBool(in ...bool) bool {
	if len(in) > 0 {
		return in[0]
	}
	return true
}

`
	var b bytes.Buffer
	sw := generator.NewSnippetWriter(&b, ctx, "$", "$")
	sw.Do(GenerateVariadicBool(), nil)
	assert.NoError(t, sw.Error())
	assert.Equal(t, want, b.String())
}

func TestGenerateBoolPointer(t *testing.T) {
	t.Parallel()

	ctx, err := newTestGeneratorContext()
	require.NoError(t, err)

	want := `// boolPointer returns a pointer to a bool.
func boolPointer(in bool) *bool {
	return &in
}

`
	var b bytes.Buffer
	sw := generator.NewSnippetWriter(&b, ctx, "$", "$")
	sw.Do(GenerateBoolPointer(), nil)
	assert.NoError(t, sw.Error())
	assert.Equal(t, want, b.String())
}
