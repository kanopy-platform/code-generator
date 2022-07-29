package snippets

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

const allValue = "*"
const allSuffix = "All"

func GenerateEnumSetter(inputType *types.Type, enumOptions []string) (string, generator.Args) {
	args := generator.Args{
		"type": inputType,
		"name": inputType.Name.Name,
	}

	raw := ""

	for _, val := range enumOptions {
		raw += fmt.Sprintf(`const $.name$%s $.name$ = "%s"`, toSuffix(val), val)
		raw += "\n"
	}

	return raw, args
}

func toSuffix(v string) string {
	suffix := v
	if suffix == allValue {
		suffix = allSuffix
	}
	caser := cases.Title(language.English)
	return caser.String(suffix)
}
