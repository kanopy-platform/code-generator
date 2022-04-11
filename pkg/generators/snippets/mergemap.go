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

	sw := generator.NewSnippetWriter(w, ctx, "$", "$")

	sw.Do("// mergeMapStringString creates a new map and loads it from map args\n", nil)
	sw.Do("// this function takes a least 2 args and later map args take precedence.\n", nil)
	sw.Do("func mergeMapStringString(m1 map[string]string, mapArgs ...map[string]string) map[string]string {\n", nil)
	sw.Do("	// populate initial map\n", nil)
	sw.Do("	outMap := map[string]string{}\n", nil)
	sw.Do("	for k, v := range m1 {\n", nil)
	sw.Do(" 	outMap[k] = v\n", nil)
	sw.Do("	}\n", nil)
	sw.Do("	// iterate all args\n", nil)
	sw.Do("	for _, m := range mapArgs {\n", nil)
	sw.Do(" 	for k, v := range m {\n", nil)
	sw.Do("	 		outMap[k] = v\n", nil)
	sw.Do("		}\n", nil)
	sw.Do("	}\n", nil)
	sw.Do("	return outMap\n", nil)
	sw.Do("}\n", nil)

	return sw.Error()
}
