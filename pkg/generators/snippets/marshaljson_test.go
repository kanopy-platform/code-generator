package snippets

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func TestGenerateMarshalJSON(t *testing.T) {
	t.Parallel()

	ctx, err := newTestGeneratorContext()
	require.NoError(t, err)

	tests := []struct {
		ctx                            *generator.Context
		t                              *types.Type
		schemeGroupVersionPackageAlias string
		wantErr                        error
		want                           string
	}{
		{
			// nil pointer ctx
			ctx:     nil,
			t:       newTestNamespaceType(),
			wantErr: errors.New("nil pointer"),
		},
		{
			// nil pointer t
			ctx:     ctx,
			t:       nil,
			wantErr: errors.New("nil pointer"),
		},
		{
			ctx:                            ctx,
			t:                              newTestNamespaceType(),
			schemeGroupVersionPackageAlias: "corev1",
			wantErr:                        nil,
			want: `// MarshalJSON is an autogenerated marshaling function, setting the TypeMeta.
func (o Namespace) MarshalJSON() ([]byte, error) {
	o.Namespace.TypeMeta = metav1.TypeMeta{Kind: "Namespace", APIVersion: corev1.SchemeGroupVersion.String()}
	return json.Marshal(o.Namespace)
}

`,
		},
	}

	for _, test := range tests {
		var b bytes.Buffer
		err := GenerateMarshalJSON(&b, test.ctx, test.t, test.schemeGroupVersionPackageAlias)
		if test.wantErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, b.String())
		}
	}
}
