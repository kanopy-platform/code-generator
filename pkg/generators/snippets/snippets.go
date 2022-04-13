package snippets

import (
	"errors"
	"io"

	"k8s.io/gengo/generator"
)

func writeSnippet(w io.Writer, ctx *generator.Context, text string, args generator.Args) error {
	if ctx == nil {
		return errors.New("nil pointer")
	}

	sw := generator.NewSnippetWriter(w, ctx, "$", "$")
	sw.Do(text, args)

	return sw.Error()
}
