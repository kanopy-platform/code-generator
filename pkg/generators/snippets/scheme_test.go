package snippets

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/gengo/generator"
)

func TestGenerateAddToScheme_Snippet(t *testing.T) {
	ctx, err := newTestGeneratorContext()
	require.NoError(t, err)

	tests := []struct {
		description string
		aliases     []string
	}{
		{
			description: "empty function",
		},
		{
			description: "single alias",
			aliases:     []string{"a"},
		},
		{
			description: "multiple alias",
			aliases:     []string{"a", "b"},
		},
	}

	for _, test := range tests {
		var b bytes.Buffer
		sw := generator.NewSnippetWriter(&b, ctx, "$", "$")
		sw.Do(GenerateAddToScheme(test.aliases))
		assert.Equal(t, strings.Count(b.String(), "SchemeBuilder"), len(test.aliases), test.description)
		for _, a := range test.aliases {
			assert.Equal(t, strings.Count(b.String(), fmt.Sprintf("%s.SchemeBuilder", a)), 1, test.description)
		}
	}
}
