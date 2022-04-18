package snippets

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/gengo/generator"
)

func TestWriteSnippetWithArgs(t *testing.T) {
	t.Parallel()

	ctx, err := newTestGeneratorContext()
	require.NoError(t, err)

	tests := []struct {
		ctx     *generator.Context
		text    string
		args    generator.Args
		wantErr error
		want    string
	}{
		{
			// nil pointer ctx
			ctx:     nil,
			wantErr: errors.New("nil pointer"),
		},
		{
			ctx: ctx,
			text: `// Test Comment
func (o $.type|raw$) MarshalJSON(in $.inputType$) $.returnType$ {
	o.$.type|raw$.TypeMeta = metav1.TypeMeta{Kind: "$.type|raw$", APIVersion: $.alias$.SchemeGroupVersion.String()}
	return json.Marshal(o.$.type|raw$)
}

`,
			args: generator.Args{
				"type":       newSampleTestType(t, "TestStruct"),
				"alias":      "corev1",
				"inputType":  "string",
				"returnType": "error",
			},
			wantErr: nil,
			want: `// Test Comment
func (o TestStruct) MarshalJSON(in string) error {
	o.TestStruct.TypeMeta = metav1.TypeMeta{Kind: "TestStruct", APIVersion: corev1.SchemeGroupVersion.String()}
	return json.Marshal(o.TestStruct)
}

`,
		},
	}

	for _, test := range tests {
		var b bytes.Buffer
		err := writeSnippetWithArgs(&b, test.ctx, test.text, test.args)
		if test.wantErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, b.String())
		}
	}
}
