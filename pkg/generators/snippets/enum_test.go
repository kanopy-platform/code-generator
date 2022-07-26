package snippets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func TestGenerateEnumSetter(t *testing.T) {
	t.Parallel()

	ctx, err := newTestGeneratorContext()
	assert.NoError(t, err)

	tests := []struct {
		description string
		enumVals    []string
		want        string
	}{
		{
			description: "single value",
			enumVals:    []string{"val1"},
			want:        "const MyEnumVal1 MyEnum = \"val1\"\n",
		},
		{
			description: "multi value",
			enumVals:    []string{"val1", "val2"},
			want:        "const MyEnumVal1 MyEnum = \"val1\"\nconst MyEnumVal2 MyEnum = \"val2\"\n",
		},
		{
			description: "all caps",
			enumVals:    []string{"VAL"},
			want:        "const MyEnumVal MyEnum = \"VAL\"\n",
		},
	}

	tt := enumTestType()

	for _, test := range tests {
		var b bytes.Buffer
		sw := generator.NewSnippetWriter(&b, ctx, "$", "$")
		sw.Do(GenerateEnumSetter(&tt, test.enumVals))
		assert.NoError(t, sw.Error(), test.description)
		assert.Equal(t, test.want, b.String(), test.description)
	}
}

func enumTestType() types.Type {
	tt := types.Type{
		Name: types.Name{
			Package: "./",
			Name:    "MyEnum",
			Path:    "",
		},
	}
	return tt
}
