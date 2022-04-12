package snippets

import (
	"io"

	"k8s.io/gengo/generator"
)

func writeSnippet(w io.Writer, ctx *generator.Context, text string) error {
	return writeSnippetWithArgs(w, ctx, text, nil)
}

func writeSnippetWithArgs(w io.Writer, ctx *generator.Context, text string, args generator.Args) error {
	sw := generator.NewSnippetWriter(w, ctx, "$", "$")
	sw.Do(text, args)

	return sw.Error()
}
