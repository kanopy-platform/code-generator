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

func TestGenerateSetter(t *testing.T) {
	t.Parallel()

	ctx, err := newTestGeneratorContext()
	require.NoError(t, err)

	tests := []struct {
		ctx       *generator.Context
		root      *types.Type
		member    types.Member
		hasParent bool
		embedded  *types.Type
		wantErr   error
		want      string
	}{
		{
			// nil pointer t
			ctx:     ctx,
			root:    nil,
			wantErr: errors.New("nil pointer"),
		},
	}

	for _, test := range tests {
		var b bytes.Buffer
		err := GenerateSetter(&b, test.ctx, test.root, test.member, test.hasParent, test.embedded)
		if test.wantErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, b.String())
		}
	}
}
