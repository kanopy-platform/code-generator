package snippets

import (
	"fmt"
	"strings"
	"unicode"

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

	suffix = strings.ReplaceAll(suffix, "kubernetes.io/", "")

	if strings.Contains(suffix, "-") {
		var sb strings.Builder
		up := false
		for i, r := range suffix {
			if i == 0 {
				sb.WriteRune(unicode.ToUpper(r))
				continue
			}
			if r == '-' {
				up = true
				continue
			}
			if up {
				up = false
				sb.WriteRune(unicode.ToUpper(r))
			} else {
				sb.WriteRune(r)
			}
		}
		suffix = sb.String()
	}

	if strings.ToLower(suffix) == suffix || strings.ToUpper(suffix) == suffix {
		caser := cases.Title(language.English)
		return caser.String(suffix)
	}
	return suffix
}
