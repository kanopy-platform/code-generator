package snippets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/gengo/generator"
)

func TestGenerateMergeMapStringString(t *testing.T) {
	t.Parallel()

	ctx, err := newTestGeneratorContext()
	require.NoError(t, err)

	tests := []struct {
		ctx     *generator.Context
		wantErr error
		want    string
	}{
		{
			ctx:     ctx,
			wantErr: nil,
			want: `// mergeMapStringString creates a new map and loads it from map args
// This function takes at least 2 args. Later map args take precedence.
func mergeMapStringString(m1 map[string]string, mapArgs ...map[string]string) map[string]string {
	outMap := map[string]string{}
	for k, v := range m1 {
		outMap[k] = v
	}

	for _, m := range mapArgs {
		for k, v := range m {
			outMap[k] = v
		}
	}
	return outMap
}

`,
		},
	}

	for _, test := range tests {
		var b bytes.Buffer
		err := GenerateMergeMapStringString(&b, test.ctx)
		if test.wantErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, b.String())
		}
	}
}
