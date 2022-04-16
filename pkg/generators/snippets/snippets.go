package snippets

import (
	"errors"
	"io"

	"k8s.io/gengo/generator"
)

// writeSnippet uses NewSnippetWriter to parse the input text and runs args through it.
//
// w is the destination
// text is the templated text
// args is used as a lookup for the template directives in text, typically a map or struct.
//   Set args to nil if there are no template directives within text.
func writeSnippet(w io.Writer, ctx *generator.Context, text string, args generator.Args) error {
	if ctx == nil {
		return errors.New("nil pointer")
	}

	sw := generator.NewSnippetWriter(w, ctx, "$", "$")
	sw.Do(text, args)

	return sw.Error()
}
