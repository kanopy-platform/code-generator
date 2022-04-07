package snippets

import (
	"io"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func GenerateMarshalJSON(c *generator.Context, t *types.Type, w io.Writer) error {
	return nil
}
