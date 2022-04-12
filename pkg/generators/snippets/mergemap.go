package snippets

import (
	"errors"
	"io"

	"k8s.io/gengo/generator"
)

func GenerateMergeMapStringString(w io.Writer, ctx *generator.Context) error {
	if ctx == nil {
		return errors.New("nil pointer")
	}

	raw := `// mergeMapStringString creates a new map and loads it from map args
// This function takes a least 2 args. Later map args take precedence.
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

`
	return writeSnippet(w, ctx, raw)
}
